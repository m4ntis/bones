package cpu

import (
	"github.com/m4ntis/bones/controller"
	"github.com/m4ntis/bones/ines"
	"github.com/m4ntis/bones/ppu"
)

// TODO: Worker should probably moved away from the cpu package, as it is more
// on the whole NES level

// Worker is used to run the NES.
//
// The worker is initialized with a parsed NES ROM, running it's programme,
// reading from the controller and displaying to the given display.
type Worker struct {
	c *CPU
	p *ppu.PPU

	nmi chan bool
}

// NewWorker initializes an instance of a worker, returning the instance.
//
// The controller passed to the worker is only read fby the worker, and expected
// to be controlled by the caller.
func NewWorker(rom *ines.ROM, disp ppu.Displayer, ctrl *controller.Controller) *Worker {
	nmi := make(chan bool)

	p := ppu.New(rom.Header.Mirroring, nmi, disp)
	p.LoadROM(rom)

	ram := RAM{}
	c := New(&ram)
	c.LoadROM(rom)

	ram.CPU = c
	ram.PPU = p
	ram.Ctrl = ctrl

	return &Worker{
		c: c,
		p: p,

		nmi: nmi,
	}
}

// Start starts running the NES.
func (w *Worker) Start() {
	go w.handleNmi()

	for {
		w.execNext()
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
