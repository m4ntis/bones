package cpu

import (
	"github.com/m4ntis/bones/models"
	"github.com/m4ntis/bones/ppu"
)

type Worker struct {
	c *CPU
	p *ppu.PPU

	nmi chan bool
}

func NewWorker(rom *models.ROM, disp ppu.Displayer, ctrl *models.Controller) *Worker {
	nmi := make(chan bool)

	p := ppu.New(nmi, disp)
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
	w.c.HandleInterupts()
	cycles := w.c.ExecNext()
	for i := 0; i < cycles*3; i++ {
		w.p.Cycle()
	}
}
