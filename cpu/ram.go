package cpu

import "github.com/m4ntis/bones/ppu"

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
)

type RAM struct {
	data [RamSize]byte

	CPU *CPU
	PPU *ppu.PPU
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
	default:
		d = r.data[addr]
	}

	r.data[addr] = d
	return d
}

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
	}

	// We write it anyway, even if mapped i/o, so RAM.Observe can see the value
	r.data[addr] = d
	return
}

func (r *RAM) Observe(addr int) byte {
	addr = getAddr(addr)
	return r.data[addr]
}
