// Package cpu provides an API for executing MOS 6502's instructions
package cpu

import (
	"github.com/m4ntis/bones/ines"
	"github.com/m4ntis/bones/io"
	"github.com/m4ntis/bones/ppu"
	"github.com/pkg/errors"
)

const (
	NMIVector   = 0xfffa
	ResetVector = 0xfffc
	IRQVector   = 0xfffe
)

// CPU implements the mos 6502.
//
// CPU should be loaded with a ROM, and can then execute the programme opcode by
// opcode.
//
// CPU exports its RAM and registers which can both be read and written to.
type CPU struct {
	RAM *RAM
	Reg *Registers

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
func New(p *ppu.PPU, ctrl *io.Controller) *CPU {
	ram := &RAM{}

	c := &CPU{
		RAM: ram,
		Reg: &Registers{},

		irq:   false,
		nmi:   false,
		reset: false,
	}

	ram.CPU = c
	ram.PPU = p
	ram.Ctrl = ctrl

	return c
}

// Load connects a CPU's RAM to a ROM mapper and inits the CPU's PC to the
// ROM's reset vector.
func (cpu *CPU) Load(rom *ines.ROM) {
	cpu.RAM.Mapper = rom.Mapper
	cpu.resetPC()
}

func (cpu *CPU) Vectors() [3]int {
	return [3]int{
		int(cpu.RAM.MustRead(NMIVector)) |
			int(cpu.RAM.MustRead(NMIVector+1))<<8,
		int(cpu.RAM.MustRead(ResetVector)) |
			int(cpu.RAM.MustRead(ResetVector+1))<<8,
		int(cpu.RAM.MustRead(IRQVector)) |
			int(cpu.RAM.MustRead(IRQVector+1))<<8,
	}
}

// TODO: You'd think that cycles should be internal to the CPU... I should
// probably think about it.

// ExecNext fetches the next opcode from RAM and executes it.
//
// ExecNext returns cycle count the whole operation took and an error if one
// occured.
func (cpu *CPU) ExecNext() (cycles int, err error) {
	defer cpu.handleInterrupts()

	code, err := cpu.RAM.Read(cpu.Reg.PC)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to read opcode from memory")
	}

	op, ok := OpCodes[code]
	if !ok {
		return 0, errors.Errorf("Invalid opcode to execute: %02x, PC: %04x",
			code, cpu.Reg.PC)
	}

	// This is switched instead of iterated because generalizing operand
	// handling to all 3 cases and iterating would probably turn out uglier.
	switch op.Mode.OpsLen {
	case 1:
		op1, err := cpu.RAM.Read(cpu.Reg.PC + 1)
		if err != nil {
			return 0, errors.Wrap(err, "Failed to read operand from memory")
		}

		cycles, err = op.Exec(cpu, op1)
		if err != nil {
			return 0, errors.Wrap(err, "Failed to execute opcode")
		}
	case 2:
		op1, err := cpu.RAM.Read(cpu.Reg.PC + 1)
		op2, err := cpu.RAM.Read(cpu.Reg.PC + 2)
		if err != nil {
			return 0, errors.Wrap(err, "Failed to read operands from memory")
		}

		cycles, err = op.Exec(cpu, op1, op2)
		if err != nil {
			return 0, errors.Wrap(err, "Failed to execute opcode")
		}
	default:
		cycles, err = op.Exec(cpu)
		if err != nil {
			return 0, errors.Wrap(err, "Failed to execute opcode")
		}
	}

	// TODO: decrement cycles after 1786830 cycles, as of now there's a
	// potential overflow
	cpu.cycles += cycles
	return cycles, nil
}

// resetPC sets the cpu's PC to the reset vector.
func (cpu *CPU) resetPC() {
	cpu.Reg.PC = int(cpu.RAM.MustRead(ResetVector)) |
		int(cpu.RAM.MustRead(ResetVector+1))<<8
}

// TODO: Make interrupt handlers constant
func (cpu *CPU) handleInterrupts() {
	if cpu.reset {
		cpu.interrupt(ResetVector)
		cpu.reset = false
	} else if cpu.nmi {
		cpu.interrupt(NMIVector)
		cpu.nmi = false
	} else if cpu.irq {
		cpu.interrupt(IRQVector)
		cpu.irq = false
	}
}

// IRQ causes the IRQ handler to be executed right after the current opcode
// finished execution, as long as the 'I' bit of the status register is reset.
func (cpu *CPU) IRQ() {
	if cpu.Reg.I == clear {
		cpu.irq = true
	}
}

// NMi causes the NMI handler to be executed right after the current opcode
// finished execution.
func (cpu *CPU) NMI() {
	cpu.nmi = true
}

// Reset causes the reset handler to be executed right after the current opcode
// finished execution.
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

	// TODO: Find appropriate place for this outside the interrupt function so
	// the CPU can call it when initiating PC in the beginning.
	cpu.Reg.I = set

	cpu.Reg.PC = int(cpu.RAM.MustRead(handlerAddr)) |
		int(cpu.RAM.MustRead(handlerAddr+1))<<8
	return
}

// Stack operations
func (cpu *CPU) push(d byte) {
	cpu.RAM.Write(cpu.getStackAddr(), d)
	cpu.Reg.SP--
}

func (cpu *CPU) pull() byte {
	cpu.Reg.SP++
	return cpu.RAM.MustRead(cpu.getStackAddr())
}

func (cpu *CPU) getStackAddr() int {
	return int(cpu.Reg.SP) | (1 << 8)
}
