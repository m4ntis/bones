// Package bones implements a worker that runs the NES and provides an API for
// breaking and debugging it's CPU.
package bones

import (
	"github.com/m4ntis/bones/asm"
	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/ines"
	"github.com/m4ntis/bones/io"
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

	Code  asm.Code
	PCIdx int

	Err error
}

// NES runs the CPU and PPU, providing a simple debugging API.
type NES struct {
	// Breaks are published on this channel when run in ModeDebug.
	Breaks chan BreakState

	c *cpu.CPU
	p *ppu.PPU

	running bool
	stopc   chan struct{}

	mode Mode

	/* Mode debug related NES state */
	bps   breakPoints
	instQ asm.Code

	continuec chan struct{}
	nextc     chan struct{}
}

// Mode represents CPU running type (run/debug)
type Mode int

const (
	ModeRun Mode = iota
	ModeDebug
)

// New creates a runnable instance of an NES.
//
// mode determines whether the NES will publish breaks and errors (ModeDebug)
// or just just run the CPU and panic on error (ModeRun).
func New(disp ppu.Displayer, ctrl *io.Controller, mode Mode) *NES {
	p := ppu.New(disp)
	c := cpu.New(p, ctrl)

	return &NES{
		c: c,
		p: p,

		running: false,
		mode:    mode,

		Breaks: make(chan BreakState),

		continuec: make(chan struct{}),
		nextc:     make(chan struct{}),
	}
}

func (n *NES) Load(rom *ines.ROM) {
	n.p.Load(rom)
	n.c.Load(rom)
}

// Start starts running the NES until Stop is called.
//
// Start is blocking and should be run in a goroutine of it's own.
func (n *NES) Start() {
	n.stopc = make(chan struct{})
	n.running = true

	go n.handleNmi()

	if n.mode == ModeRun {
		n.startRun()
		return
	}

	n.startDebug()
}

// startRun runs the CPU without checking breakpoints or errors.
func (n *NES) startRun() {
	for {
		select {
		case <-n.stopc:
			return
		default:
			panicOnErr(n.execNext())
		}
	}
}

// startDebug checks for breakpoints and publishes errors after each cycle.
func (n *NES) startDebug() {
	n.bps = breakPoints{
		n.c.Reg.PC: true,
	}

	for {
		select {
		case <-n.stopc:
			return
		default:
			n.handleBps()
			n.handleError(n.execNext())
		}
	}
}

// Stop sends a signal to stop the cpu on the next cycle.
func (n *NES) Stop() {
	if n.running {
		close(n.stopc)
		n.running = false
	}
}

// Continue resumes the programme's execution until the next breakpoint is hit.
func (n *NES) Continue() {
	n.continuec <- struct{}{}
}

// Next executes the next opcode in the programme and breaks.
func (n *NES) Next() {
	n.nextc <- struct{}{}
}

// Break adds a breakpoint at addr.
//
// Note that breakpoint hitting won't be checked if the NES was run in ModeRun.
func (n *NES) Break(addr int) {
	n.bps[addr] = true
}

// Delete attempts to delete an existing breakpoint at an address, returning
// whether there was a breakpoint in that address or not.
func (n *NES) Delete(addr int) (ok bool) {
	_, ok = n.bps[addr]
	if ok {
		delete(n.bps, addr)
	}
	return ok
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

func (n *NES) Vectors() [3]int {
	return n.c.Vectors()
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
				asm.DisassembleRAM(n.c.RAM, n.c.Reg.PC, instFutureSize+1)...),
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
	for {
		select {
		case <-n.p.NMI:
			n.c.NMI()
		case <-n.stopc:
			return
		}
	}
}

func (n *NES) handleError(err error) {
	if err != nil {
		n.breakOper(err)
	}
}

func (n *NES) addInstToQ() {
	// TODO: This is extremly inefficient
	n.instQ = append(n.instQ, asm.DisassembleRAM(n.c.RAM, n.c.Reg.PC, 1)[0])

	if len(n.instQ) > instHistorySize {
		n.instQ = n.instQ[1:]
	}
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
