package cpu

import (
	"image"

	"github.com/m4ntis/bones/models"
	"github.com/m4ntis/bones/ppu"
)

type Displayer interface {
	Display(image.Image)
}

type Worker struct {
	c *CPU
	p *ppu.PPU

	disp  Displayer
	frame *models.Frame

	nmi chan bool

	pixelc chan models.Pixel
}

func NewWorker(rom *models.ROM, d Displayer, ctrl *models.Controller) *Worker {
	nmi := make(chan bool)
	pixelc := make(chan models.Pixel)

	p := ppu.New(nmi, pixelc)
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

		disp:  d,
		frame: &models.Frame{},

		nmi: nmi,

		pixelc: pixelc,
	}
}

func (w *Worker) Start() {
	go w.handleNmi()
	go w.handlePixel()

	for {
		w.execNext()
	}
}

func (w *Worker) handleNmi() {
	for <-w.nmi {
		w.c.NMI()
	}
}

func (w *Worker) handlePixel() {
	for pix := range w.pixelc {
		w.frame.Push(pix)
		if pix.X == 255 && pix.Y == 239 {
			w.disp.Display(w.frame.Create())
		}
	}
}

func (w *Worker) execNext() {
	w.c.HandleInterupts()
	cycles := w.c.ExecNext()
	for i := 0; i < cycles*3; i++ {
		w.p.Cycle()
	}
}
