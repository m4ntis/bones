package dbg

import (
	"sync"

	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/disass"
	"github.com/m4ntis/bones/drawer"
	"github.com/m4ntis/bones/models"
	"github.com/m4ntis/bones/ppu"
)

type breakPoints map[int]bool

type BreakData struct {
	RAM *cpu.RAM
	Reg *cpu.Regs

	Disass disass.Disassembly
}

type Worker struct {
	c *cpu.CPU
	p *ppu.PPU

	drawer *drawer.Drawer
	frame  *models.Frame

	d      disass.Disassembly
	bps    breakPoints
	bpsMux *sync.Mutex

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
func NewWorker(rom *models.ROM, vals chan<- BreakData, d *drawer.Drawer) *Worker {
	nmi := make(chan bool)
	p := ppu.New(nmi)
	p.LoadROM(rom)

	ram := cpu.RAM{}
	c := cpu.New(&ram)
	c.LoadROM(rom)

	ram.CPU = c
	ram.PPU = p

	return &Worker{
		c: c,
		p: p,

		drawer: d,
		frame:  &models.Frame{},

		d: disass.Disassemble(rom.PrgROM),
		bps: breakPoints{
			c.Reg.PC: true,
		},
		bpsMux: &sync.Mutex{},

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

	w.bpsMux.Lock()
	defer w.bpsMux.Unlock()

	w.bps[addr] = true
	return true
}

func (w *Worker) Delete(addr int) (success bool) {
	w.bpsMux.Lock()
	defer w.bpsMux.Unlock()

	_, success = w.bps[addr]
	if success {
		delete(w.bps, addr)
	}
	return success
}

func (w *Worker) DeleteAll() {
	w.bpsMux.Lock()
	defer w.bpsMux.Unlock()

	for addr, _ := range w.bps {
		delete(w.bps, addr)
	}
}

func (w *Worker) List() (breaks []int) {
	w.bpsMux.Lock()
	defer w.bpsMux.Unlock()

	for addr, _ := range w.bps {
		breaks = append(breaks, addr)
	}
	return breaks
}

func (w *Worker) handleBps() {
	w.bpsMux.Lock()
	_, ok := w.bps[w.c.Reg.PC]
	w.bpsMux.Unlock()

	if ok {
		w.breakOper()
	}
}

func (w *Worker) breakOper() {
	for {
		w.vals <- BreakData{
			Disass: w.d,
			RAM:    w.c.RAM,
			Reg:    w.c.Reg,
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
