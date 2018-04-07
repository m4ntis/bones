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
	SprPaletteIdx    = 0x3f10
	PaletteMirrorIdx = 0x3f20
	RamMirrorIdx     = 0x4000
	RamSize          = 0x10000
)

type VRAM struct {
	data [RamSize]byte
}

func getAddr(addr int) int {
	addr %= RamMirrorIdx
	if addr >= TablesMirrorIdx && addr < BgrPaletteIdx {
		return addr - 0x1000
	}
	if addr >= PaletteMirrorIdx && addr < RamMirrorIdx {
		if addr == 0x3f10 || addr == 0x3f14 || addr == 0x3f18 || addr == 0x3f0c {
			return addr - 0x10
		}
		return (addr-BgrPaletteIdx)%0x20 + BgrPaletteIdx
	}
	return addr
}

func (v *VRAM) Read(addr int) byte {
	return v.data[getAddr(addr)]
}

func (v *VRAM) Write(addr int, data byte) {
	v.data[getAddr(addr)] = data
}
