package cpu

import (
	"fmt"

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
	ReadAccess = iota
	WriteAccess
)

type AccessMode int

// RAM represents the cpu's 64k of onboard ram.
//
// The api for accessing RAM addresses is done with the RAM.Fetch and RAM.Commit
// methods. See respective documentation for more detail.
//
// The Fetch/Commit api is more flexible than the more obious Read/Write api by
// letting the caller delegate ram access by passing a "handle" it gets using
// Fetch rather than passing a reference to ram and an addr. In addition, it
// lets the caller aggregate multible read/write accesses into one logical
// access, without invoking the Read or Write multiple times. This is usefull
// when say reading a value multiple times, each time testing a different bit.
// Another example would be changing the value in RAM relative to itself. With
// a Read/Write api, that would mean you would have to invoke a Read at least
// once, despite needing to classify the whole access as one write action.
//
// The aparrent drawback would be that the responsibility of specifiying that a
// read/write has occured is handed to the caller. But if you consider the
// flexibility of aggregating multiple accesses into one logical access, the
// case with the NES, there wouldn't be a way of avoiding that delegation of
// responsibility.
type RAM struct {
	data [RamSize]byte

	iom IOMapper
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

// Fetch returns a "handle" to a byte in RAM, specified by addr.
//
// Fetch lets you access RAM values, but it is REQUIRED of the caller to call
// Commit after a logical read/write has occured. See description of RAM for
// the reasoning of this api.
func (r *RAM) Fetch(addr int) *byte {
	return &r.data[getAddr(addr)]
}

// Commit signifies that a read/write access occured on a certain address.
//
// This lets the RAM handle read/write events, such as updating memory mapped
// i/o.
func (r *RAM) Commit(addr int, mode AccessMode) {
	addr = getAddr(addr)

	if addr >= LowerIORegBeginIdx && addr < LowerIORegMirrorBeginIdx ||
		addr >= UpperIORegBeginIdx && addr < ExpansionRomBeginIdx {
		switch mode {
		case ReadAccess:
			r.iom.OnRead(addr)
		case WriteAccess:
			r.iom.OnWrite(addr, r.data[addr])
		default:
			panic(fmt.Sprintf("Invalid access mode %d", mode))
		}
	}
}

type IOMapper struct {
	ppu *ppu.PPU
}

func (iom IOMapper) OnRead(addr int) {

}

func (iom IOMapper) OnWrite(addr int, d byte) {

}
