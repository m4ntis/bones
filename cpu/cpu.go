// Package cpu provides an api for the mos 6502
package cpu

import "fmt"

// CPU implements the mos 6502.
//
// CPU should be loaded with a ROM, and can then execute the programme opcode by
// opcode.
//
// CPU exports its RAM and registers which can both be read and written to.
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

// ResetPC sets the pc to the reset handler addr
//
// TODO: Make sure ram is initialized before reseting PC
func (cpu *CPU) ResetPC() {
	cpu.Reg.PC = int(cpu.RAM.Read(0xfffc)) | int(cpu.RAM.Read(0xfffd))<<8
}

// ExecNext reads the next opcode from RAM, executes it and returns the cycle
// count.
func (cpu *CPU) ExecNext() (cycles int) {
	defer cpu.handleInterrupts()

	code := cpu.RAM.Read(cpu.Reg.PC)
	op, ok := OpCodes[code]
	if !ok {
		// TODO: Don't panic
		panic(fmt.Sprintf("Invalid opcode to execute: %02x, PC: %04x", code,
			cpu.Reg.PC))
	}

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

// IRQ causes the IRQ handler to be executed right after the current opcode
// finished execution, as long as the 'I' bit of the status register is reset.
func (cpu *CPU) IRQ() {
	if cpu.Reg.I == Clear {
		cpu.irq = true
	}
}

// NMi causes the IRQ handler to be executed right after the current opcode
// finished execution.
func (cpu *CPU) NMI() {
	cpu.nmi = true
}

// Reset causes the IRQ handler to be executed right after the current opcode
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

	cpu.Reg.I = Set

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
