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

func getIndex(index int) int {
	if index < 0 || index > RamSize {
		panic("RAM accessing index out of range")
	}

	if index >= TablesMirrorIdx && index < BgrPaletteIdx {
		return index - 0x1000
	}
	if index >= PaletteMirrorIdx && index < RamMirrorIdx {
		return (index-BgrPaletteIdx)%0x20 + BgrPaletteIdx
	}
	if index >= RamMirrorIdx {
		return index % 0x4000
	}
	return index
}

func (r *RAM) Read(index int) byte {
	return r.data[getIndex(index)]
}

func (r *RAM) Write(index int, data byte) {
	r.data[getIndex(index)] = data
}
