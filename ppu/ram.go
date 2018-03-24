package ppu

const (
	PT0Idx           = 0x0
	PT1Idx           = 0x1000
	NT0Idx           = 0x2000
	AT0Idx           = 0x23c0
	NT1Idx           = 0x2400
	AT1Idx           = 0x27c0
	NT2Idx           = 0x2800
	AT2Idx           = 0x2bc0
	NT3Idx           = 0x2c00
	AT3Idx           = 0x2fc0
	TablesMirrorIdx  = 0x3000
	BgrPaletteIdx    = 0x3f00
	SpritePaletteIdx = 0x3f10
	PaletteMirrorIdx = 0x3f20
	RamMirrorIdx     = 0x4000
	RamSize          = 0x10000
)

type RAM struct {
	data [RamSize]byte
}

func getAddr(addr int) int {
	if addr < 0 || addr > RamSize {
		panic("RAM accessing addr out of range")
	}

	if addr >= TablesMirrorIdx && addr < BgrPaletteIdx {
		return addr - 0x1000
	}
	if addr >= PaletteMirrorIdx && addr < RamMirrorIdx {
		return (addr-BgrPaletteIdx)%0x20 + BgrPaletteIdx
	}
	if addr >= RamMirrorIdx {
		return addr % 0x4000
	}
	return addr
}

func (r *RAM) Read(addr int) byte {
	return r.data[getAddr(addr)]
}

func (r *RAM) Write(addr int, data byte) {
	r.data[getAddr(addr)] = data
}
