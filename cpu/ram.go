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

func getIndex(index int) int {
	if index < 0 || index > RamSize {
		panic("RAM accessing index out of range")
	}

	if index >= RamMirrorBeginIdx && index < LowerIORegBeginIdx {
		return index % 0x800
	}
	if index >= LowerIORegMirrorBeginIdx && index < UpperIORegBeginIdx {
		return (index-LowerIORegBeginIdx)%0x8 + LowerIORegBeginIdx
	}
	return index
}

func (r *RAM) Fetch(index int) *byte {
	return &r.data[getIndex(index)]
}
