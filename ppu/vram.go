package ppu

import "github.com/m4ntis/bones/ines/mapper"

const (
	RAMSize = 0x10000

	ptSize = 0x1000

	pt0Addr           = 0x0
	pt1Addr           = 0x1000
	nt0Addr           = 0x2000
	at0Addr           = 0x23c0
	nt1Addr           = 0x2400
	at1Addr           = 0x27c0
	nt2Addr           = 0x2800
	at2Addr           = 0x2bc0
	nt3Addr           = 0x2c00
	at3Addr           = 0x2fc0
	tablesMirrorAddr  = 0x3000
	bgrPaletteAddr    = 0x3f00
	sprPaletteAddr    = 0x3f10
	paletteMirrorAddr = 0x3f20
	ramMirrorAddr     = 0x4000
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
	addr %= ramMirrorAddr

	if addr >= tablesMirrorAddr && addr < bgrPaletteAddr {
		return addr - 0x1000
	}

	if addr >= paletteMirrorAddr && addr < ramMirrorAddr {
		if addr == 0x3f10 || addr == 0x3f14 || addr == 0x3f18 || addr == 0x3f0c {
			return addr - 0x10
		}
		return (addr-bgrPaletteAddr)%0x20 + bgrPaletteAddr
	}

	return addr
}

// Read returns the a value in a specified address.
func (v *VRAM) Read(addr int) byte {
	addr = stripMirror(addr)

	// addr it a PT address
	if addr < nt0Addr {
		d, _ := v.Mapper.Read(addr)
		return d
	}

	return v.data[addr]
}

// Write sets a specified address with a given value.
func (v *VRAM) Write(addr int, d byte) {
	addr = stripMirror(addr)

	// addr it a PT address
	if addr < nt0Addr {
		v.Mapper.Write(addr, d)
		return
	}

	v.data[addr] = d
}
