package cpu

import (
	"image"

	"github.com/m4ntis/bones/models"
	"github.com/m4ntis/bones/ppu"
)

type Drawer interface {
	Draw(image.Image)
}

type Worker struct {
	c *CPU
	p *ppu.PPU

	drawer Drawer
	frame  *models.Frame

	nmi chan bool
}

func NewWorker(rom *models.ROM, d Drawer) *Worker {
	nmi := make(chan bool)
	p := ppu.New(nmi)
	p.LoadROM(rom)

	ram := RAM{}
	c := New(&ram)
	c.LoadROM(rom)

	ram.CPU = c
	ram.PPU = p

	return &Worker{
		c: c,
		p: p,

		drawer: d,
		frame:  &models.Frame{},

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
		// This isn't the correct place for this
		w.drawer.Draw(w.frame.Create())
	}
}

func (w *Worker) execNext() {
	w.c.HandleInterupts()
	cycles := w.c.ExecNext()
	for i := 0; i < cycles*3; i++ {
		w.frame.Push(w.p.Cycle())
	}
}
