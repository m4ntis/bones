package cpu

import (
	"github.com/m4ntis/bones/models"
)

type CPU struct {
	RAM *RAM

	Reg *Registers
}

func NewCPU() *CPU {
	return &CPU{
		Reg: &Registers{
			PC: 0x8000,
			SP: 0xFF,
		},
	}
}

func (cpu *CPU) LoadROM(rom *models.ROM) {
	// Load first 2 pages of PrgROM (not supporting mappers as of yet)
	copy(rom.PrgROM[0][:], cpu.RAM.data[:models.PRG_ROM_PAGE_SIZE])
	copy(rom.PrgROM[1][:],
		cpu.RAM.data[models.PRG_ROM_PAGE_SIZE:2*models.PRG_ROM_PAGE_SIZE])
}

// Interrupt Handling
func (cpu *CPU) IRQ() {
	if cpu.Reg.I == CLEAR {
		cpu.interrupt(0xfffe)
	}
}

func (cpu *CPU) NMI() {
	if int(*cpu.RAM.Fetch(0x2000)&1<<7) == CLEAR {
		cpu.interrupt(0xfffa)
	}
}

func (cpu *CPU) Reset() {
	cpu.interrupt(0xfffc)
}

func (cpu *CPU) interrupt(handlerAddr int) {
	// push PCH
	cpu.push(byte(cpu.Reg.PC >> 8))
	// push PCL
	cpu.push(byte(cpu.Reg.PC & 0xff))
	// push P
	cpu.push(cpu.Reg.GetP())

	cpu.Reg.I = 1

	// fetch PCL from $fffe and PCH from $ffff
	cpu.Reg.PC = int(*cpu.RAM.Fetch(handlerAddr)) |
		int(*cpu.RAM.Fetch(handlerAddr + 1))<<8
	return
}

// Stack operations
func (cpu *CPU) push(b byte) {
	*cpu.RAM.Fetch(cpu.getStackAddr()) = b
	cpu.Reg.SP--
}

func (cpu *CPU) pull() byte {
	cpu.Reg.SP++
	return *cpu.RAM.Fetch(cpu.getStackAddr())
}

func (cpu *CPU) getStackAddr() int {
	return int(cpu.Reg.SP) | (1 << 8)
}
