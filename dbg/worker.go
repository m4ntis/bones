// Package dbg provides a worker that runs the NES with an api for breaking the
// cpu and handling breaks
package dbg

import (
	"github.com/m4ntis/bones/controller"
	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/disass"
	"github.com/m4ntis/bones/ines"
	"github.com/m4ntis/bones/ppu"
)

type breakPoints map[int]bool

// BreakState describes the state of the NES and the running programme when the
// cpu breaks.
type BreakState struct {
	Reg  *cpu.Regs
	RAM  *cpu.RAM
	VRAM *ppu.VRAM

	Disass disass.Disassembly
}

// Worker runs the NES and provides an api for all basic debugging
// functionality.
type Worker struct {
	c *cpu.CPU
	p *ppu.PPU

	d   disass.Disassembly
	bps breakPoints

	nmi chan bool

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
func NewWorker(rom *ines.ROM, vals chan<- BreakState, disp ppu.Displayer,
	ctrl *controller.Controller) *Worker {
	nmi := make(chan bool)

	p := ppu.New(rom.Header.Mirroring, rom.Mapper, nmi, disp)

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

		d: disass.Disassemble(rom.PrgROM),
		bps: breakPoints{
			c.Reg.PC: true,
		},

		nmi: nmi,

		continuec: make(chan bool),
		nextc:     make(chan bool),
		vals:      vals,
	}
}

// Start makes the worker start running the NES.
//
// Start should be run in a goroutine.
func (w *Worker) Start() {
	go w.handleNmi()

	for {
		w.handleBps()
		w.execNext()
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

// Break adds a breakpoint at an address, and returns whether the address is a
// valid breaking address (the start of a new instruction).
func (w *Worker) Break(addr int) (success bool) {
	if w.d.IndexOf(addr) == -1 {
		return false
	}

	w.bps[addr] = true
	return true
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
		w.breakOper()
	}
}

func (w *Worker) breakOper() {
	for {
		w.vals <- BreakState{
			Reg:  w.c.Reg,
			RAM:  w.c.RAM,
			VRAM: w.p.VRAM,

			Disass: w.d,
		}

		select {
		case <-w.continuec:
			return
		case <-w.nextc:
			w.execNext()
			continue
		}
	}
}

func (w *Worker) handleNmi() {
	for <-w.nmi {
		w.c.NMI()
	}
}

func (w *Worker) execNext() {
	cycles := w.c.ExecNext()
	for i := 0; i < cycles*3; i++ {
		w.p.Cycle()
	}
}
