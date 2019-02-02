package cpu

import (
	"github.com/m4ntis/bones/ines"
	"github.com/m4ntis/bones/io"
	"github.com/m4ntis/bones/ppu"
	"github.com/pkg/errors"
)

const (
	RAMSize = 0x10000

	zeroPageAddr       = 0x0
	stackAddr          = 0x100
	ramAddr            = 0x200
	ramMirrorAddr      = 0x800
	ppuRegAddr         = 0x2000
	ppuRegMirrorAddr   = 0x2008
	ioRegAddr          = 0x4000
	cartridgeSpaceAddr = 0x4020
)

const (
	ppuCtrlAddr   = 0x2000
	ppuMaskAddr   = 0x2001
	ppuStatusAddr = 0x2002
	oamAddrAddr   = 0x2003
	oamDataAddr   = 0x2004
	ppuScrollAddr = 0x2005
	ppuAddrAddr   = 0x2006
	ppuDataAddr   = 0x2007
	oamDMAAddr    = 0x4014
	ctrl1Addr     = 0x4016
)

// RAM holds the mos 6502's 16k (64 when mirrored) of on chip memory.
//
// The RAM contains some memory mapped i/o and therefore should be initialized
// and passed to the CPU as well as the mapped components.
//
// All RAM accessing methods contain logic for mirrored address translation.
type RAM struct {
	data [RAMSize]byte

	Mapper ines.Mapper

	CPU  *CPU
	PPU  *ppu.PPU
	Ctrl *io.Controller
}

// stripMirror returns the underlying address after mirroring.
func stripMirror(addr int) int {
	addr &= 0xffff

	// Internal RAM mirroring
	if addr >= ramMirrorAddr && addr < ppuRegAddr {
		return addr % 0x800
	}

	if addr >= ppuRegMirrorAddr && addr < ioRegAddr {
		return (addr-ppuRegAddr)%0x8 + ppuRegAddr
	}

	return addr
}

func (r *RAM) readMMIO(addr int) (d byte, err error) {
	switch addr {
	case ppuCtrlAddr:
		return 0, errors.New("Invalid read from PPUCtrl")
	case ppuMaskAddr:
		return 0, errors.New("Invalid read from PPUMask")
	case ppuStatusAddr:
		d = r.PPU.Regs.PPUStatusRead()
	case oamAddrAddr:
		return 0, errors.New("Invalid read from OAMAddr")
	case oamDataAddr:
		d = r.PPU.Regs.OAMDataRead()
	case ppuScrollAddr:
		return 0, errors.New("Invalid read from PPUScroll")
	case ppuAddrAddr:
		return 0, errors.New("Invalid read from PPUAddr")
	case ppuDataAddr:
		d = r.PPU.Regs.PPUDataRead()
	case oamDMAAddr:
		return 0, errors.New("Invalid read from OAMDMA")
	case ctrl1Addr:
		d = r.Ctrl.Read()
	default:
		// Read from PPU i/o register mirroring

		return 0, nil
	}

	// We keep the read at the ram location for observe to be able to see the
	// last i/o operation
	r.data[addr] = d
	return d, nil
}

func (r *RAM) writeMMIO(addr int, d byte) (cycles int, err error) {
	switch addr {
	case ppuCtrlAddr:
		r.PPU.Regs.PPUCtrlWrite(d)
	case ppuMaskAddr:
		r.PPU.Regs.PPUMaskWrite(d)
	case ppuStatusAddr:
		return 0, nil
	case oamAddrAddr:
		r.PPU.Regs.OAMAddrWrite(d)
	case oamDataAddr:
		r.PPU.Regs.OAMDataWrite(d)
	case ppuScrollAddr:
		r.PPU.Regs.PPUScrollWrite(d)
	case ppuAddrAddr:
		r.PPU.Regs.PPUAddrWrite(d)
	case ppuDataAddr:
		r.PPU.Regs.PPUDataWrite(d)
	case oamDMAAddr:
		// TODO: Move DMA to CPU struct. Incrementing the cycles in that method
		// will allow to remove the cycles from ram and operand api

		var oamData [256]byte
		copy(oamData[:], r.data[int(d)<<8:int(d+1)<<8])
		r.PPU.DMA(oamData)

		cycles += 513
		// extra cycle on odd cycles
		if r.CPU.oddCycle {
			cycles++
		}
	case ctrl1Addr:
		r.Ctrl.Strobe(d & 1)
	}

	// r.data is updated regardless of i/o reg write in order to be able to
	// "observe" the value later
	r.data[addr] = d

	return cycles, nil
}

// Read fetches a byte from memory, cartridge or i/o register, specified by addr.
func (r *RAM) Read(addr int) (d byte, err error) {
	addr = stripMirror(addr)

	// Read from cartridge
	if addr >= cartridgeSpaceAddr {
		return r.Mapper.Read(addr)
	}

	// Read from MMIO
	if addr >= ppuRegAddr {
		return r.readMMIO(addr)
	}

	// Read from internal RAM
	return r.data[addr], nil
}

// TODO: Consider inlining the Must fucntions

// MustRead calls Read but panics instead of returning an error.
func (r *RAM) MustRead(addr int) byte {
	d, err := r.Read(addr)
	if err != nil {
		panic(err)
	}

	return d
}

// Write puts a value to memory, PRG-RAM or i/o register, specified by addr.
//
// Write returns a cycle count the memory access took, as writing to some i/o
// registers may block the cpu and take up cycles, such as DMA.
func (r *RAM) Write(addr int, d byte) (cycles int, err error) {
	addr = stripMirror(addr)

	// Write to cartridge
	if addr >= cartridgeSpaceAddr {
		return 0, r.Mapper.Write(addr, d)
	}

	// Write to MMIO
	if addr >= ppuRegAddr {
		return r.writeMMIO(addr, d)
	}

	// Write to internal RAM
	r.data[addr] = d
	return 0, nil
}

// MustWrite calls Write but panics instead of returning an error.
func (r *RAM) MustWrite(addr int, d byte) (cycles int) {
	cycles, err := r.Write(addr, d)
	if err != nil {
		panic(err)
	}

	return cycles
}

// Observe is used as an api for debuggers, letting the caller read the value in
// RAM without triggering memory mapped i/o operations.
//
// Reading from memory mapped i/o locations will return the last value written
// to them.
func (r *RAM) Observe(addr int) (d byte, err error) {
	if addr < 0 || addr > RAMSize {
		return 0, errors.Errorf("Invalid RAM observing addr $%04x",
			addr)
	}

	addr = stripMirror(addr)

	// Read from cartridge
	if addr >= cartridgeSpaceAddr {
		return r.Mapper.Observe(addr)
	}

	// Read from internal RAM
	return r.data[addr], nil
}

// TODO: Consider changing api to a single "access" method (maybe internally)
// to avoid all code duplication between read and write methods
