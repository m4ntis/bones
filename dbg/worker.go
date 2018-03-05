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

type DbgWorker struct {
	c *cpu.CPU
	d disass.Disassembly

	bps    breakPoints
	bpsMux *sync.Mutex

	continuec chan bool
	nextc     chan bool
	vals      chan<- BreakData
}

// NewDbgWorker creates a dbg worker that will start a cpu that will run on the
// given ROM.
//
// The vals channel is the channel containing the data returned each time the
// cpu breaks, describing the current cpu state.
func NewDbgWorker(rom *models.ROM, vals chan<- BreakData) *DbgWorker {
	c := cpu.NewCPU()
	c.LoadROM(rom)

	return &DbgWorker{
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
func (dw *DbgWorker) Start() {
	for {
		dw.handleBps()
		dw.c.ExecNext()
		dw.c.HandleInterupts()
	}
}

func (dw *DbgWorker) Continue() {
	dw.continuec <- true
}

func (dw *DbgWorker) Next() {
	dw.nextc <- true
}

func (dw *DbgWorker) Break(addr int) {
	dw.bpsMux.Lock()
	defer dw.bpsMux.Unlock()

	dw.bps[addr] = true
}

func (dw *DbgWorker) Clear(addr int) {
	dw.bpsMux.Lock()
	defer dw.bpsMux.Unlock()

	delete(dw.bps, addr)
}

func (dw *DbgWorker) ClearAll() {
	dw.bpsMux.Lock()
	defer dw.bpsMux.Unlock()

	for addr, _ := range dw.bps {
		delete(dw.bps, addr)
	}
}

func (dw *DbgWorker) handleBps() {
	dw.bpsMux.Lock()
	_, ok := dw.bps[dw.c.Reg.PC-0x8000]
	dw.bpsMux.Unlock()

	if ok {
		dw.breakOper()
	}
}

func (dw *DbgWorker) breakOper() {
	for {
		dw.vals <- BreakData{
			d:   dw.d,
			ram: dw.c.RAM,
			reg: dw.c.Reg,
		}

		select {
		case <-dw.continuec:
			return
		case <-dw.nextc:
			continue
		}
	}
}
