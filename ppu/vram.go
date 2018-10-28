package ppu

import "github.com/m4ntis/bones/ines/mapper"

const (
	PT0Addr           = 0x0
	PT1Addr           = 0x1000
	NT0Addr           = 0x2000
	AT0Addr           = 0x23c0
	NT1Addr           = 0x2400
	AT1Addr           = 0x27c0
	NT2Addr           = 0x2800
	AT2Addr           = 0x2bc0
	NT3Addr           = 0x2c00
	AT3Addr           = 0x2fc0
	TablesMirrorAddr  = 0x3000
	BgrPaletteAddr    = 0x3f00
	SprPaletteAddr    = 0x3f10
	PaletteMirrorAddr = 0x3f20
	RAMMirrorAddr     = 0x4000

	PTSize  = 0x1000
	RAMSize = 0x10000
)

// VRAM holds the Ricoh 2A03's 16kb (64 when mirrored) of on board memory.
//
// All VRAM accessing methods contain logic for mirrored address translation.
type VRAM struct {
	data [RAMSize]byte

	Mapper mapper.Mapper
}

// stripMirror returns the underlying address after mirroring.
func stripMirror(addr int) int {
	// VRAM is mirrored entirely every `RAMMirrorAddr` bytes
	addr %= RAMMirrorAddr

	if addr >= TablesMirrorAddr && addr < BgrPaletteAddr {
		return addr - 0x1000
	}

	if addr >= PaletteMirrorAddr && addr < RAMMirrorAddr {
		if addr == 0x3f10 || addr == 0x3f14 || addr == 0x3f18 || addr == 0x3f0c {
			return addr - 0x10
		}
		return (addr-BgrPaletteAddr)%0x20 + BgrPaletteAddr
	}

	return addr
}

// Read returns the a value in a specified address.
func (v *VRAM) Read(addr int) byte {
	addr = stripMirror(addr)

	// addr it a PT address
	if addr < NT0Addr {
		d, _ := v.Mapper.Read(addr)
		return d
	}

	return v.data[addr]
}

// Write sets a specified address with a given value.
func (v *VRAM) Write(addr int, d byte) {
	addr = stripMirror(addr)

	// addr it a PT address
	if addr < NT0Addr {
		v.Mapper.Write(addr, d)
		return
	}

	v.data[addr] = d
}
