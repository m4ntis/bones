package cpu

import (
	"github.com/m4ntis/bones/controller"
	"github.com/m4ntis/bones/ppu"
)

const (
	ZeroPageBeginIdx         = 0x0
	StackBeginIdx            = 0x100
	RamBeginIdx              = 0x200
	RamMirrorBeginIdx        = 0x800
	LowerIORegBeginIdx       = 0x2000
	LowerIORegMirrorBeginIdx = 0x2008
	UpperIORegBeginIdx       = 0x4000
	ExpansionRomBeginIdx     = 0x4020
	SramBeginIdx             = 0x4000
	PrgRomLowerBeginIdx      = 0x8000
	PrgRomUpperBeginIdx      = 0xc000
	RamSize                  = 0x10000
)

const (
	PPUCtrl   = 0x2000
	PPUMask   = 0x2001
	PPUStatus = 0x2002
	OAMAddr   = 0x2003
	OAMData   = 0x2004
	PPUScroll = 0x2005
	PPUAddr   = 0x2006
	PPUData   = 0x2007
	OAMDMA    = 0x4014
	Ctrl1     = 0x4016
)

// TODO: should consider whethere the ram should know all of it's memory
// mappings or where there should be some sort of part to manage it (motherboard
// or something of the sort).

// RAM holds the mos 6502's 64k of on chip memory.
//
// The RAM contains some memory mapped i/o and therefore should be initialized
// and passed to the CPU as well as the mapped components.
//
// All RAM accessing methods contain logic for mirrored address translation.
type RAM struct {
	data [RamSize]byte

	CPU  *CPU
	PPU  *ppu.PPU
	Ctrl *controller.Controller
}

// getAddr returns the underlying address after mapping.
func getAddr(addr int) int {
	if addr < 0 || addr > RamSize {
		panic("RAM accessing addr out of range")
	}

	if addr >= RamMirrorBeginIdx && addr < LowerIORegBeginIdx {
		return addr % 0x800
	}
	if addr >= LowerIORegMirrorBeginIdx && addr < UpperIORegBeginIdx {
		return (addr-LowerIORegBeginIdx)%0x8 + LowerIORegBeginIdx
	}
	return addr
}

// Read returns the byte in the address specified.
func (r *RAM) Read(addr int) byte {
	addr = getAddr(addr)

	var d byte

	switch addr {
	case PPUCtrl:
		panic("Invalid read from PPUCtrl")
	case PPUMask:
		panic("Invalid read from PPUMask")
	case PPUStatus:
		d = r.PPU.PPUStatusRead()
	case OAMAddr:
		panic("Invalid read from OAMAddr")
	case OAMData:
		d = r.PPU.OAMDataRead()
	case PPUScroll:
		panic("Invalid read from PPUScroll")
	case PPUAddr:
		panic("Invalid read from PPUAddr")
	case PPUData:
		d = r.PPU.PPUDataRead()
	case OAMDMA:
		panic("Invalid read from OAMDMA")
	case Ctrl1:
		d = r.Ctrl.Read()
	default:
		d = r.data[addr]
	}

	// We keep the read at the ram location for observe to be able to see the
	// last i/o operation
	r.data[addr] = d
	return d
}

// Write writes a value to the specified address.
//
// Write returns a cycle count the memory access took, as some memory mapped i/o
// operations may block the cpu and take up cycles, such as DMA.
func (r *RAM) Write(addr int, d byte) (cycles int) {
	addr = getAddr(addr)

	switch addr {
	case PPUCtrl:
		r.PPU.PPUCtrlWrite(d)
	case PPUMask:
		r.PPU.PPUMaskWrite(d)
	case PPUStatus:
		panic("Invalid write to PPUStatus")
	case OAMAddr:
		r.PPU.OAMAddrWrite(d)
	case OAMData:
		r.PPU.OAMDataWrite(d)
	case PPUScroll:
		r.PPU.PPUScrollWrite(d)
	case PPUAddr:
		r.PPU.PPUAddrWrite(d)
	case PPUData:
		r.PPU.PPUDataWrite(d)
	case OAMDMA:
		var oamData [256]byte
		copy(oamData[:], r.data[int(d)<<8:int(d+1)<<8])
		r.PPU.DMA(oamData)

		cycles += 513
		// extra cycle on odd cycles
		if r.CPU.cycles%2 == 1 {
			cycles++
		}
	case Ctrl1:
		r.Ctrl.Strobe(d & 1)
	}

	// We write it anyway, even if mapped i/o, so RAM.Observe can see the value
	r.data[addr] = d
	return
}

// Observe is used as an api for debuggers, letting the caller read the value in
// RAM without triggering memory mapped i/o operations.
//
// Reading from memory mapped i/o locations will return the last value written
// to these locations.
func (r *RAM) Observe(addr int) byte {
	addr = getAddr(addr)
	return r.data[addr]
}
