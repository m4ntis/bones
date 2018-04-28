// Package cpu provides an api for the mos 6502
package cpu

import "github.com/m4ntis/bones/ines"

// CPU represents the mos 6502 and implements its functionality.
//
// CPU should be loaded with a ROM, and can then execute the programme opcode by
// opcode.
type CPU struct {
	RAM *RAM
	Reg *Regs

	cycles int

	irq   bool
	nmi   bool
	reset bool
}

// New creates an instance of the CPU struct.
//
// ram is passed to the CPU instead of initialized within, as it is shared with
// other components via memory mapped i/o. It is the caller's responsibility to
// initialize and pass it to the other parts of the NES.
func New(ram *RAM) *CPU {
	return &CPU{
		RAM: ram,
		Reg: &Regs{},

		irq:   false,
		nmi:   false,
		reset: false,
	}
}

// LoadROM loads a parsed ROM to CPU memory, and initialized the PC register to
// the reset vector.
func (cpu *CPU) LoadROM(rom *ines.ROM) {
	if len(rom.PrgROM) > 1 {
		// Load first 2 pages of PrgROM (not supporting mappers as of yet)
		copy(cpu.RAM.data[0x8000:0x8000+ines.PrgROMPageSize], rom.PrgROM[0][:])
		copy(cpu.RAM.data[0x8000+ines.PrgROMPageSize:0x8000+2*ines.PrgROMPageSize],
			rom.PrgROM[1][:])
	} else {
		// If there is only one page of prg rom, load it to $c000 ~ $ffff
		copy(cpu.RAM.data[0x8000+ines.PrgROMPageSize:0x8000+2*ines.PrgROMPageSize],
			rom.PrgROM[0][:])
	}

	// Init pc to the reset handler addr
	cpu.Reg.PC = int(cpu.RAM.Read(0xfffc)) | int(cpu.RAM.Read(0xfffd))<<8
}

// ExecNext reads the next opcode from RAM, executes it and returns the cycle
// count.
func (cpu *CPU) ExecNext() (cycles int) {
	defer cpu.handleInterrupts()

	op := OpCodes[cpu.RAM.Read(cpu.Reg.PC)]

	// We are doing this manually cus there are only 3 posibilities and writing
	// logic to describe this would be ugly IMO
	if op.Mode.OpsLen == 1 {
		cycles = op.Exec(cpu, cpu.RAM.Read(cpu.Reg.PC+1))
	} else if op.Mode.OpsLen == 2 {
		cycles = op.Exec(cpu, cpu.RAM.Read(cpu.Reg.PC+1), cpu.RAM.Read(cpu.Reg.PC+2))
	} else {
		cycles = op.Exec(cpu)
	}

	// TODO: decrement cycles after 1786830 cycles
	cpu.cycles += cycles
	return cycles
}

func (cpu *CPU) handleInterrupts() {
	if cpu.reset {
		cpu.interrupt(0xfffc)
		cpu.reset = false
	} else if cpu.nmi {
		cpu.interrupt(0xfffa)
		cpu.nmi = false
	} else if cpu.irq {
		cpu.interrupt(0xfffe)
		cpu.irq = false
	}
}

func (cpu *CPU) IRQ() {
	if cpu.Reg.I == Clear {
		cpu.irq = true
	}
}

func (cpu *CPU) NMI() {
	cpu.nmi = true
}

func (cpu *CPU) Reset() {
	cpu.reset = true
}

func (cpu *CPU) interrupt(handlerAddr int) {
	// push PCH
	cpu.push(byte(cpu.Reg.PC >> 8))
	// push PCL
	cpu.push(byte(cpu.Reg.PC & 0xff))
	// push P
	cpu.push(cpu.Reg.GetP())

	cpu.Reg.I = 1

	cpu.Reg.PC = int(cpu.RAM.Read(handlerAddr)) |
		int(cpu.RAM.Read(handlerAddr+1))<<8
	return
}

// Stack operations
func (cpu *CPU) push(d byte) {
	cpu.RAM.Write(cpu.getStackAddr(), d)
	cpu.Reg.SP--
}

func (cpu *CPU) pull() byte {
	cpu.Reg.SP++
	return cpu.RAM.Read(cpu.getStackAddr())
}

func (cpu *CPU) getStackAddr() int {
	return int(cpu.Reg.SP) | (1 << 8)
}
