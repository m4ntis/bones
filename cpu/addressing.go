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
	opsLen int
	Format func([]byte) string

	address func(*CPU, Operation, bool, ...byte) int
}

var (
	Implied = AddressingMode{
		Name:   "Implied",
		opsLen: 0,
		Format: func(ops []byte) string { return "" },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			op(cpu)
			cpu.Reg.PC++
			return
		},
	}

	Accumulator = AddressingMode{
		Name:   "Accumulator",
		opsLen: 0,
		Format: func(ops []byte) string { return "A" },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			op(cpu, NewOperand(cpu, CPURegisterOperand, ARegisterOperand))
			cpu.Reg.PC++
			return
		},
	}

	Immediate = AddressingMode{
		Name:   "Immediate",
		opsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("#$%02x", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			op(cpu, NewOperand(cpu, ConstOperand, int(ops[0])))
			cpu.Reg.PC += 2
			return
		},
	}

	ZeroPage = AddressingMode{
		Name:   "ZeroPage",
		opsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			op(cpu, NewOperand(cpu, CPURAMOperand, int(ops[0])))
			cpu.Reg.PC += 2
			return
		},
	}

	ZeroPageX = AddressingMode{
		Name:   "ZeroPageX",
		opsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x, X", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			op(cpu, NewOperand(cpu, CPURAMOperand, int(ops[0]+cpu.Reg.X)))
			cpu.Reg.PC += 2
			return
		},
	}

	ZeroPageY = AddressingMode{
		Name:   "ZeroPageY",
		opsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x, Y", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			op(cpu, NewOperand(cpu, CPURAMOperand, int(ops[0]+cpu.Reg.Y)))
			cpu.Reg.PC += 2
			return
		},
	}

	Relative = AddressingMode{
		Name:   "Relative",
		opsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			extraCycles = op(cpu, NewOperand(cpu, CPUConstOperand, ops[0]))
			cpu.Reg.PC += 2
			return
		},
	}

	Absolute = AddressingMode{
		Name:   "Absolute",
		opsLen: 2,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x%02x", ops[1], ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			addr := int(ops[0]) | int(ops[1])<<8
			op(cpu, NewOperand(cpu, CPURAMOperand, addr))
			cpu.Reg.PC += 3
			return
		},
	}

	AbsoluteX = AddressingMode{
		Name:   "AbsoluteX",
		opsLen: 2,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x%02x, X", ops[1], ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			addr := int(ops[0]) | int(ops[1])<<8
			xAddr := addr + int(cpu.Reg.X)

			if addr/256 != xAddr/256 && pageBoundryCheck {
				extraCycles++
			}

			op(cpu, NewOperand(cpu, CPURAMOperand, xAddr))
			cpu.Reg.PC += 3
			return
		},
	}

	AbsoluteY = AddressingMode{
		Name:   "AbsoluteY",
		opsLen: 2,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x%02x, Y", ops[1], ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			addr := int(ops[0]) | int(ops[1])<<8
			yAddr := addr + int(cpu.Reg.Y)

			if addr/256 != yAddr/256 && pageBoundryCheck {
				extraCycles++
			}

			op(cpu, NewOperand(cpu, CPURAMOperand, yAddr))
			cpu.Reg.PC += 3
			return
		},
	}

	Indirect = AddressingMode{
		Name:   "Indirect",
		opsLen: 2,
		Format: func(ops []byte) string { return fmt.Sprintf("($%02x%02x)", ops[1], ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			addr := int(ops[0]) + int(ops[1])<<8

			adl := cpu.RAM.Read(addr)
			adh := cpu.RAM.Read(addr + 1)

			op(cpu, NewOperand(cpu, CPURAMOperand, int(adl)|int(adh)<<8))
			return
		},
	}

	IndirectX = AddressingMode{
		Name:   "IndirectX",
		opsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("($%02x, X)", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			addr := int(ops[0] + cpu.Reg.X)

			adl := cpu.RAM.Read(addr)
			adh := cpu.RAM.Read((addr + 1) % 0x100)

			op(cpu, NewOperand(cpu, CPURAMOperand, int(adl)|int(adh)<<8))
			cpu.Reg.PC += 2
			return
		},
	}

	IndirectY = AddressingMode{
		Name:   "IndirectY",
		opsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("($%02x), Y", ops[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, ops ...byte) (extraCycles int) {
			addr := int(ops[0])

			adl := cpu.RAM.Read(addr)
			adh := cpu.RAM.Read((addr + 1) % 0x100)

			fetched := int(*adl) | int(*adh)<<8
			if fetched/256 != (fetched+int(cpu.Reg.Y))/256 && pageBoundryCheck {
				extraCycles++
			}
			fetched += cpu.Reg.Y

			op(cpu, NewOperand(cpu, CPURAMOperand, fetched))
			cpu.Reg.PC += 2
			return
		},
	}
)
