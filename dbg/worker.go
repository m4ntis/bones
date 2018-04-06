package dbg

import (
	"image"

	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/disass"
	"github.com/m4ntis/bones/models"
	"github.com/m4ntis/bones/ppu"
)

type breakPoints map[int]bool

type BreakData struct {
	Reg  *cpu.Regs
	CRAM *cpu.RAM
	VRAM *ppu.VRAM

	Disass disass.Disassembly
}

type Displayer interface {
	Display(image.Image)
}

type Worker struct {
	c *cpu.CPU
	p *ppu.PPU

	d   disass.Disassembly
	bps breakPoints

	nmi chan bool

	continuec chan bool
	nextc     chan bool
	vals      chan<- BreakData
}

// NewWorker creates a dbg worker that will start a cpu that will run on the
// given ROM.
//
// The vals channel is the channel containing the data returned each time the
// cpu breaks, describing the current cpu state.
func NewWorker(rom *models.ROM, vals chan<- BreakData, disp ppu.Displayer, ctrl *models.Controller) *Worker {
	nmi := make(chan bool)

	p := ppu.New(nmi, disp)
	p.LoadROM(rom)

	ram := cpu.RAM{}
	c := cpu.New(&ram)
	c.LoadROM(rom)

	ram.CPU = c
	ram.PPU = p
	ram.Ctrl = ctrl

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

// Start starts the debug worker.
//
// Runs in a loop, should be run in a goroutine
func (w *Worker) Start() {
	go w.handleNmi()

	for {
		w.handleBps()
		w.execNext()
	}
}

func (w *Worker) Continue() {
	w.continuec <- true
}

func (w *Worker) Next() {
	w.nextc <- true
}

func (w *Worker) Break(addr int) (success bool) {
	if w.d.IndexOf(addr) == -1 {
		return false
	}

	w.bps[addr] = true
	return true
}

func (w *Worker) Delete(addr int) (success bool) {
	_, success = w.bps[addr]
	if success {
		delete(w.bps, addr)
	}
	return success
}

func (w *Worker) DeleteAll() {
	for addr, _ := range w.bps {
		delete(w.bps, addr)
	}
}

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
		w.vals <- BreakData{
			Reg:  w.c.Reg,
			CRAM: w.c.RAM,
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
	w.c.HandleInterupts()
	for i := 0; i < cycles*3; i++ {
		w.p.Cycle()
	}
}
