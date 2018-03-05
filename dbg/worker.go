package dbg

import (
	"sync"

	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/disass"
	"github.com/m4ntis/bones/models"
)

type breakPoints map[int]bool

type BreakData struct {
	ram *cpu.RAM
	reg *cpu.Registers

	d disass.Disassembly
}

type Worker struct {
	c *cpu.CPU
	d disass.Disassembly

	bps    breakPoints
	bpsMux *sync.Mutex

	continuec chan bool
	nextc     chan bool
	vals      chan<- BreakData
}

// NewWorker creates a dbg worker that will start a cpu that will run on the
// given ROM.
//
// The vals channel is the channel containing the data returned each time the
// cpu breaks, describing the current cpu state.
func NewWorker(rom *models.ROM, vals chan<- BreakData) *Worker {
	c := cpu.NewCPU()
	c.LoadROM(rom)

	return &Worker{
		c: c,
		d: disass.Disassemble(rom.PrgROM),

		bps: breakPoints{
			0x0000: true,
		},
		bpsMux: &sync.Mutex{},

		continuec: make(chan bool),
		nextc:     make(chan bool),
		vals:      vals,
	}
}

// Start starts the debug worker.
//
// Runs in a loop, should be run in a goroutine
func (w *Worker) Start() {
	for {
		w.handleBps()
		w.c.ExecNext()
		w.c.HandleInterupts()
	}
}

func (w *Worker) Continue() {
	w.continuec <- true
}

func (w *Worker) Next() {
	w.nextc <- true
}

func (w *Worker) Break(addr int) {
	w.bpsMux.Lock()
	defer w.bpsMux.Unlock()

	w.bps[addr] = true
}

func (w *Worker) Clear(addr int) {
	w.bpsMux.Lock()
	defer w.bpsMux.Unlock()

	delete(w.bps, addr)
}

func (w *Worker) ClearAll() {
	w.bpsMux.Lock()
	defer w.bpsMux.Unlock()

	for addr, _ := range w.bps {
		delete(w.bps, addr)
	}
}

func (w *Worker) handleBps() {
	w.bpsMux.Lock()
	_, ok := w.bps[w.c.Reg.PC-0x8000]
	w.bpsMux.Unlock()

	if ok {
		w.breakOper()
	}
}

func (w *Worker) breakOper() {
	for {
		w.vals <- BreakData{
			d:   w.d,
			ram: w.c.RAM,
			reg: w.c.Reg,
		}

		select {
		case <-w.continuec:
			return
		case <-w.nextc:
			continue
		}
	}
}
