package cpu

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

type RAM struct {
	data [RamSize]byte
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
	return r.data[getAddr(addr)]
}

func (r *RAM) Write(addr int, d byte) {
	r.data[getAddr(addr)] = d
}
