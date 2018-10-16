// Package dbg provides a worker that runs the NES with an api for breaking the
// cpu and handling breaks
package dbg

import (
	"github.com/m4ntis/bones/controller"
	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/disass"
	"github.com/m4ntis/bones/ines"
	"github.com/m4ntis/bones/ppu"
	"github.com/pkg/errors"
)

const (
	instHistorySize = 5
	instFutureSize  = 5
)

type breakPoints map[int]bool

// BreakState describes the state of the NES and the running programme when the
// cpu breaks.
type BreakState struct {
	Reg  *cpu.Registers
	RAM  *cpu.RAM
	VRAM *ppu.VRAM

	Code  disass.Code
	PCIdx int

	Err error
}

// Worker runs the NES and provides an api for all basic debugging
// functionality.
type Worker struct {
	c *cpu.CPU
	p *ppu.PPU

	bps   breakPoints
	instQ disass.Code

	continuec chan bool
	nextc     chan bool
	vals      chan<- BreakState
}

// NewWorker creates an instance of a dbg worker.
//
// vals is a channel being populated by the worker each time it breaks,
// containing information about the current state of the NES.
//
// ctrl will be read by the cpu the worker runs and is expected to be controlled
// by the caller.
func NewWorker(rom *ines.ROM, disp ppu.Displayer, ctrl *controller.Controller,
	vals chan<- BreakState) *Worker {

	p := ppu.New(rom.Header.Mirroring, rom.Mapper, disp)

	ram := cpu.RAM{}
	c := cpu.New(&ram)

	ram.CPU = c
	ram.PPU = p
	ram.Ctrl = ctrl
	ram.Mapper = rom.Mapper

	c.ResetPC()

	return &Worker{
		c: c,
		p: p,

		bps: breakPoints{
			c.Reg.PC: true,
		},

		continuec: make(chan bool),
		nextc:     make(chan bool),
		vals:      vals,
	}
}

// Start starts running the NES.
//
// Start is blocking and should be run in a goroutine of it's own.
func (w *Worker) Start() {
	go w.handleNmi()

	for {
		w.handleBps()
		w.handleError(w.execNext())
	}
}

// Continue resumes the programme's execution until the next breakpoint is hit.
func (w *Worker) Continue() {
	w.continuec <- true
}

// Next executes the next opcode in the programme and breaks.
func (w *Worker) Next() {
	w.nextc <- true
}

// Break adds a breakpoint at addr.
func (w *Worker) Break(addr int) {
	w.bps[addr] = true
}

// Delete attempts to delete an existing breakpoint at an address, returning
// whether there was a breakpoint in that address or not.
func (w *Worker) Delete(addr int) (success bool) {
	_, success = w.bps[addr]
	if success {
		delete(w.bps, addr)
	}
	return success
}

// DeleteAll removes all set breakpoints.
func (w *Worker) DeleteAll() {
	for addr, _ := range w.bps {
		delete(w.bps, addr)
	}
}

// List returns the list of breakpoints set.
func (w *Worker) List() (breaks []int) {
	for addr, _ := range w.bps {
		breaks = append(breaks, addr)
	}
	return breaks
}

func (w *Worker) handleBps() {
	_, ok := w.bps[w.c.Reg.PC]

	if ok {
		w.breakOper(nil)
	}
}

func (w *Worker) breakOper(err error) {
	for {
		w.vals <- BreakState{
			Reg:  w.c.Reg,
			RAM:  w.c.RAM,
			VRAM: w.p.VRAM,

			Code: append(w.instQ,
				disass.DisassembleRAM(w.c.RAM, w.c.Reg.PC, instFutureSize+1)...),
			PCIdx: len(w.instQ),

			Err: err,
		}

		select {
		case <-w.continuec:
			return
		case <-w.nextc:
			w.handleError(w.execNext())
			continue
		}
	}
}

func (w *Worker) handleNmi() {
	for <-w.p.NMI {
		w.c.NMI()
	}
}

func (w *Worker) execNext() error {
	// Add the next instruction to be executed immediately to queue. This is so
	// the queue will be updated before PC is incremented to next instruction.
	//
	// TODO: This might be an issue if the instruction fails to execute.
	w.addInstToQ()

	cycles, err := w.c.ExecNext()
	if err != nil {
		return errors.Wrap(err, "Failed to execute next opcode")
	}

	for i := 0; i < cycles*3; i++ {
		w.p.Cycle()
	}

	return nil
}

func (w *Worker) handleError(err error) {
	if err != nil {
		w.breakOper(err)
	}
}

func (w *Worker) addInstToQ() {
	// TODO: This is extremly inefficient
	w.instQ = append(w.instQ, disass.DisassembleRAM(w.c.RAM, w.c.Reg.PC, 1)[0])

	if len(w.instQ) > instHistorySize {
		w.instQ = w.instQ[1:]
	}
}
