package cpu

import (
	"github.com/m4ntis/bones/controller"
	"github.com/m4ntis/bones/ines/mapper"
	"github.com/m4ntis/bones/ppu"
	"github.com/pkg/errors"
)

const (
	ZeroPageAddr       = 0x0
	StackAddr          = 0x100
	RAMAddr            = 0x200
	RAMMirrorAddr      = 0x800
	PPURegAddr         = 0x2000
	PPURegMirrorAddr   = 0x2008
	IORegAddr          = 0x4000
	CartridgeSpaceAddr = 0x4020
	RAMSize            = 0x10000
)

const (
	PPUCtrlAddr   = 0x2000
	PPUMaskAddr   = 0x2001
	PPUStatusAddr = 0x2002
	OAMAddrAddr   = 0x2003
	OAMDataAddr   = 0x2004
	PPUScrollAddr = 0x2005
	PPUAddrAddr   = 0x2006
	PPUDataAddr   = 0x2007
	OAMDMAAddr    = 0x4014
	Ctrl1Addr     = 0x4016
)

// RAM holds the mos 6502's 16k (64 when mirrored) of on chip memory.
//
// The RAM contains some memory mapped i/o and therefore should be initialized
// and passed to the CPU as well as the mapped components.
//
// All RAM accessing methods contain logic for mirrored address translation.
type RAM struct {
	data [RAMSize]byte

	Mapper mapper.Mapper

	CPU  *CPU
	PPU  *ppu.PPU
	Ctrl *controller.Controller
}

// stripMirror returns the underlying address after mirroring.
func stripMirror(addr int) int {
	// Internal RAM mirroring
	if addr >= RAMMirrorAddr && addr < PPURegAddr {
		return addr % 0x800
	}

	// PPU i/o register mirroring
	if addr >= PPURegMirrorAddr && addr < IORegAddr {
		return (addr-PPURegAddr)%0x8 + PPURegAddr
	}

	return addr
}

func (r *RAM) readMMIO(addr int) (d byte, err error) {
	switch addr {
	case PPUCtrlAddr:
		return 0, errors.New("Invalid read from PPUCtrl")
	case PPUMaskAddr:
		return 0, errors.New("Invalid read from PPUMask")
	case PPUStatusAddr:
		d = r.PPU.PPUStatusRead()
	case OAMAddrAddr:
		return 0, errors.New("Invalid read from OAMAddr")
	case OAMDataAddr:
		d = r.PPU.OAMDataRead()
	case PPUScrollAddr:
		return 0, errors.New("Invalid read from PPUScroll")
	case PPUAddrAddr:
		return 0, errors.New("Invalid read from PPUAddr")
	case PPUDataAddr:
		d = r.PPU.PPUDataRead()
	case OAMDMAAddr:
		return 0, errors.New("Invalid read from OAMDMA")
	case Ctrl1Addr:
		d = r.Ctrl.Read()
	default:
		// TODO: Consider returning an error of a not implemented mmio
		return 0, nil
	}

	// We keep the read at the ram location for observe to be able to see the
	// last i/o operation
	r.data[addr] = d
	return d, nil
}

func (r *RAM) writeMMIO(addr int, d byte) (cycles int, err error) {
	switch addr {
	case PPUCtrlAddr:
		r.PPU.PPUCtrlWrite(d)
	case PPUMaskAddr:
		r.PPU.PPUMaskWrite(d)
	case PPUStatusAddr:
		return 0, errors.New("Invalid write to PPUStatus")
	case OAMAddrAddr:
		r.PPU.OAMAddrWrite(d)
	case OAMDataAddr:
		r.PPU.OAMDataWrite(d)
	case PPUScrollAddr:
		r.PPU.PPUScrollWrite(d)
	case PPUAddrAddr:
		r.PPU.PPUAddrWrite(d)
	case PPUDataAddr:
		r.PPU.PPUDataWrite(d)
	case OAMDMAAddr:
		// TODO: Move DMA to CPU struct. Incrementing the cycles in that method
		// will allow to remove the cycles from ram and operand api

		var oamData [256]byte
		copy(oamData[:], r.data[int(d)<<8:int(d+1)<<8])
		r.PPU.DMA(oamData)

		cycles += 513
		// extra cycle on odd cycles
		if r.CPU.cycles%2 == 1 {
			cycles++
		}
	case Ctrl1Addr:
		r.Ctrl.Strobe(d & 1)
	}

	// r.data is updated regardless of i/o reg write in order to be able to
	// "observe" the value later
	r.data[addr] = d

	return cycles, nil
}

// Read fetches a byte from memory, cartridge or i/o register, specified by addr.
func (r *RAM) Read(addr int) (d byte, err error) {
	if addr < 0 || addr > RAMSize {
		return 0, errors.Errorf("Invalid RAM reading addr $%04x",
			addr)
	}

	addr = stripMirror(addr)

	// Read from cartridge
	if addr >= CartridgeSpaceAddr {
		return r.Mapper.Read(addr)
	}

	// Read from MMIO
	if addr >= PPURegAddr {
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
	if addr < 0 || addr > RAMSize {
		return 0, errors.Errorf("Invalid RAM writing addr $%04x",
			addr)
	}

	addr = stripMirror(addr)

	// Write to cartridge
	if addr >= CartridgeSpaceAddr {
		return 0, r.Mapper.Write(addr, d)
	}

	// Write to MMIO
	if addr >= PPURegAddr {
		return r.writeMMIO(addr, d)
	}

	// Write to internal RAM
	r.data[addr] = d
	return 0, nil
}

// MustWrite calls Write but panics instead of returning an error.
func (r *RAM) MustWrite(addr int, d byte) (cycles int) {
	cycles, err := r.write(addr, d)
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
	if addr >= CartridgeSpaceAddr {
		return r.Mapper.Observe(addr)
	}

	// Read from internal RAM
	return r.data[addr]
}

// TODO: Consider changing api to a single "access" method (maybe internally)
// to avoid all code duplication between read and write methods
