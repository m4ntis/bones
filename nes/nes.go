// package nes implements a worker that runs the NES and provdes an API for
// breaking and debugging.
package nes

import (
	"github.com/m4ntis/bones/controller"
	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/disass"
	"github.com/m4ntis/bones/ines"
	"github.com/m4ntis/bones/ines/mapper"
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

// NES runs the CPU and PPU, providing a basic debugging API.
type NES struct {
	Breaks chan BreakState

	c *cpu.CPU
	p *ppu.PPU

	bps   breakPoints
	instQ disass.Code

	continuec chan bool
	nextc     chan bool
}

type Mode int

const (
	ModeRun Mode = iota
	ModeDebug
)

// New creates a runnable instance of an NES.
func New(rom *ines.ROM,
	disp ppu.Displayer,
	ctrl *controller.Controller,
	mode Mode) *NES {

	p := ppu.New(rom.Header.Mirroring, rom.Mapper, disp)
	var ram *cpu.RAM
	c := cpu.New(ram)

	initRAM(ram, c, p, ctrl, rom.Mapper)
	c.ResetPC()

	var bps breakPoints
	if mode == ModeDebug {
		bps = breakPoints{
			c.Reg.PC: true,
		}
	}

	return &NES{
		Breaks: make(chan BreakState),

		c: c,
		p: p,

		bps: bps,

		continuec: make(chan bool),
		nextc:     make(chan bool),
	}
}

// Start starts running the NES.
//
// Start is blocking and should be run in a goroutine of it's own.
func (n *NES) Start() {
	go n.handleNmi()

	for {
		n.handleBps()
		n.handleError(n.execNext())
	}
}

// Continue resumes the programme's execution until the next breakpoint is hit.
func (n *NES) Continue() {
	n.continuec <- true
}

// Next executes the next opcode in the programme and breaks.
func (n *NES) Next() {
	n.nextc <- true
}

// Break adds a breakpoint at addr.
func (n *NES) Break(addr int) {
	n.bps[addr] = true
}

// Delete attempts to delete an existing breakpoint at an address, returning
// whether there was a breakpoint in that address or not.
func (n *NES) Delete(addr int) (success bool) {
	_, success = n.bps[addr]
	if success {
		delete(n.bps, addr)
	}
	return success
}

// DeleteAll removes all set breakpoints.
func (n *NES) DeleteAll() {
	for addr, _ := range n.bps {
		delete(n.bps, addr)
	}
}

// List returns the list of breakpoints set.
func (n *NES) List() (breaks []int) {
	for addr, _ := range n.bps {
		breaks = append(breaks, addr)
	}
	return breaks
}

func (n *NES) handleBps() {
	_, ok := n.bps[n.c.Reg.PC]

	if ok {
		n.breakOper(nil)
	}
}

func (n *NES) breakOper(err error) {
	for {
		n.Breaks <- BreakState{
			Reg:  n.c.Reg,
			RAM:  n.c.RAM,
			VRAM: n.p.VRAM,

			Code: append(n.instQ,
				disass.DisassembleRAM(n.c.RAM, n.c.Reg.PC, instFutureSize+1)...),
			PCIdx: len(n.instQ),

			Err: err,
		}

		select {
		case <-n.continuec:
			return
		case <-n.nextc:
			n.handleError(n.execNext())
			continue
		}
	}
}

func (n *NES) handleNmi() {
	for <-n.p.NMI {
		n.c.NMI()
	}
}

func (n *NES) execNext() error {
	// Add the next instruction to be executed immediately to queue. This is so
	// the queue will be updated before PC is incremented to next instruction.
	//
	// TODO: This might be an issue if the instruction fails to execute.
	n.addInstToQ()

	cycles, err := n.c.ExecNext()
	if err != nil {
		return errors.Wrap(err, "Failed to execute next opcode")
	}

	for i := 0; i < cycles*3; i++ {
		n.p.Cycle()
	}

	return nil
}

func (n *NES) handleError(err error) {
	if err != nil {
		n.breakOper(err)
	}
}

func (n *NES) addInstToQ() {
	// TODO: This is extremly inefficient
	n.instQ = append(n.instQ, disass.DisassembleRAM(n.c.RAM, n.c.Reg.PC, 1)[0])

	if len(n.instQ) > instHistorySize {
		n.instQ = n.instQ[1:]
	}
}

func initRAM(ram *cpu.RAM, cpu *cpu.CPU, ppu *ppu.PPU, ctrl *controller.Controller, mapper mapper.Mapper) {
	ram.CPU = cpu
	ram.PPU = ppu
	ram.Ctrl = ctrl
	ram.Mapper = mapper
}
