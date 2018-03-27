package cpu

import "fmt"

// AddressingMode defines one of the 2a03's ways of addressing the operands of
// the different opcodes.
//
// Each addressing mode is responsible of fetching the operands in it's way,
// and calling the operation with them.
//
// The bool tells whether there needs to be a page boundry check
//
// The addressing mode returns the amount of extra cycles caused by page boundry
// crossing, if any.
type AddressingMode struct {
	Name   string
	OpsLen int
	Format func([]byte) string

	address func(*CPU, Operation, bool, ...byte) int
}

var (
	Implied = AddressingMode{
		Name:   "Implied",
		OpsLen: 0,
		Format: func(ops []byte) string { return "" },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			cpu.Reg.PC++
			op(cpu, NilOperand{})
			return
		},
	}

	Accumulator = AddressingMode{
		Name:   "Accumulator",
		OpsLen: 0,
		Format: func(ops []byte) string { return "A" },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			op(cpu, RegOperand{Reg: &cpu.Reg.A})
			cpu.Reg.PC++
			return
		},
	}

	Immediate = AddressingMode{
		Name:   "Immediate",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("#$%02x", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			op(cpu, ConstOperand{D: ops[0]})
			cpu.Reg.PC += 2
			return
		},
	}

	ZeroPage = AddressingMode{
		Name:   "ZeroPage",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: int(ops[0])})
			cpu.Reg.PC += 2
			return
		},
	}

	ZeroPageX = AddressingMode{
		Name:   "ZeroPageX",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x, X", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: int(ops[0] + cpu.Reg.X)})
			cpu.Reg.PC += 2
			return
		},
	}

	ZeroPageY = AddressingMode{
		Name:   "ZeroPageY",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x, Y", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: int(ops[0] + cpu.Reg.Y)})
			cpu.Reg.PC += 2
			return
		},
	}

	Relative = AddressingMode{
		Name:   "Relative",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			extraCycles = op(cpu, ConstOperand{D: ops[0]})
			cpu.Reg.PC += 2
			return
		},
	}

	Absolute = AddressingMode{
		Name:   "Absolute",
		OpsLen: 2,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x%02x", ops[1], ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			// We inc this beforehand so that JSR wont be incremented after
			// execution
			cpu.Reg.PC += 3

			addr := int(ops[0]) | int(ops[1])<<8
			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: addr})
			return
		},
	}

	AbsoluteX = AddressingMode{
		Name:   "AbsoluteX",
		OpsLen: 2,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x%02x, X", ops[1], ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			addr := int(ops[0]) | int(ops[1])<<8
			xAddr := addr + int(cpu.Reg.X)

			if addr/256 != xAddr/256 && pageBoundryCheck {
				extraCycles++
			}

			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: xAddr})
			cpu.Reg.PC += 3
			return
		},
	}

	AbsoluteY = AddressingMode{
		Name:   "AbsoluteY",
		OpsLen: 2,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x%02x, Y", ops[1], ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			addr := int(ops[0]) | int(ops[1])<<8
			yAddr := addr + int(cpu.Reg.Y)

			if addr/256 != yAddr/256 && pageBoundryCheck {
				extraCycles++
			}

			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: yAddr})
			cpu.Reg.PC += 3
			return
		},
	}

	Indirect = AddressingMode{
		Name:   "Indirect",
		OpsLen: 2,
		Format: func(ops []byte) string { return fmt.Sprintf("($%02x%02x)", ops[1], ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			addr := int(ops[0]) + int(ops[1])<<8

			adl := cpu.RAM.Read(addr)
			adh := cpu.RAM.Read(addr + 1)

			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: int(adl) | int(adh)<<8})
			return
		},
	}

	IndirectX = AddressingMode{
		Name:   "IndirectX",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("($%02x, X)", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			addr := int(ops[0] + cpu.Reg.X)

			adl := cpu.RAM.Read(addr)
			adh := cpu.RAM.Read((addr + 1) % 0x100)

			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: int(adl) | int(adh)<<8})
			cpu.Reg.PC += 2
			return
		},
	}

	IndirectY = AddressingMode{
		Name:   "IndirectY",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("($%02x), Y", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			addr := int(ops[0])

			adl := cpu.RAM.Read(addr)
			adh := cpu.RAM.Read((addr + 1) % 0x100)

			fetched := int(adl) | int(adh)<<8
			if fetched/256 != (fetched+int(cpu.Reg.Y))/256 && pageBoundryCheck {
				extraCycles++
			}
			fetched += int(cpu.Reg.Y)

			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: fetched})
			cpu.Reg.PC += 2
			return
		},
	}
)
