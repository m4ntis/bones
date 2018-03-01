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
	Name    string
	ArgsLen int
	Format  func([]byte) string

	address func(*CPU, Operation, bool, ...*byte) int
}

var (
	ZeroPage = AddressingMode{
		Name:    "ZeroPage",
		ArgsLen: 1,
		Format:  func(args []byte) string { return fmt.Sprintf("$%02x", args[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu, cpu.RAM.Fetch(int(*args[0])))
			cpu.Reg.PC += 2
			return
		},
	}

	ZeroPageX = AddressingMode{
		Name:    "ZeroPageX",
		ArgsLen: 1,
		Format:  func(args []byte) string { return fmt.Sprintf("$%02x, X", args[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu, cpu.RAM.Fetch(int(*args[0]+cpu.Reg.X)))
			cpu.Reg.PC += 2
			return
		},
	}

	ZeroPageY = AddressingMode{
		Name:    "ZeroPageY",
		ArgsLen: 1,
		Format:  func(args []byte) string { return fmt.Sprintf("$%02x, Y", args[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu, cpu.RAM.Fetch(int(*args[0]+cpu.Reg.Y)))
			cpu.Reg.PC += 2
			return
		},
	}

	Absolute = AddressingMode{
		Name:    "Absolute",
		ArgsLen: 2,
		Format:  func(args []byte) string { return fmt.Sprintf("$%02x%02x", args[1], args[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			// ADL stored at *args[0], ADH at *args[1]
			op(cpu, cpu.RAM.Fetch(int(*args[0])|int(*args[1])<<8))
			cpu.Reg.PC += 3
			return
		},
	}

	AbsoluteX = AddressingMode{
		Name:    "AbsoluteX",
		ArgsLen: 2,
		Format:  func(args []byte) string { return fmt.Sprintf("$%02x%02x, X", args[1], args[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			addr := int(*args[0]) | int(*args[1])<<8
			if addr/256 != (addr+int(cpu.Reg.X))/256 && pageBoundryCheck {
				extraCycles++
			}

			op(cpu, cpu.RAM.Fetch(addr+int(cpu.Reg.X)))
			cpu.Reg.PC += 3
			return
		},
	}

	AbsoluteY = AddressingMode{
		Name:    "AbsoluteY",
		ArgsLen: 2,
		Format:  func(args []byte) string { return fmt.Sprintf("$%02x%02x, Y", args[1], args[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			addr := int(*args[0]) | int(*args[1])<<8
			if addr/256 != (addr+int(cpu.Reg.Y))/256 && pageBoundryCheck {
				extraCycles++
			}

			op(cpu, cpu.RAM.Fetch(addr+int(cpu.Reg.Y)))
			cpu.Reg.PC += 3
			return
		},
	}

	// AbsoluteJMP is here because the absolute JMP operation take the immediate
	// arguments instead of their value at the new location like Absolute
	AbsoluteJMP = AddressingMode{
		Name:    "AbsoluteJMP",
		ArgsLen: 2,
		Format:  func(args []byte) string { return fmt.Sprintf("$%02x%02x", args[1], args[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu, args[0], args[1])
			return
		},
	}

	Indirect = AddressingMode{
		Name:    "Indirect",
		ArgsLen: 2,
		Format:  func(args []byte) string { return fmt.Sprintf("($%02x%02x)", args[1], args[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			adl := cpu.RAM.Fetch(int(*args[0]) | int(*args[1])<<8)
			adh := cpu.RAM.Fetch(int(*args[0]) | int(*args[1])<<8 + 1)

			op(cpu, adl, adh)
			return
		},
	}

	IndirectX = AddressingMode{
		Name:    "IndirectX",
		ArgsLen: 1,
		Format:  func(args []byte) string { return fmt.Sprintf("($%02x, X)", args[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			adl := cpu.RAM.Fetch(int(*args[0] + cpu.Reg.X))
			adh := cpu.RAM.Fetch(int(*args[0] + cpu.Reg.X + 1))

			op(cpu, cpu.RAM.Fetch(int(*adl)|int(*adh)<<8))
			cpu.Reg.PC += 2
			return
		},
	}

	IndirectY = AddressingMode{
		Name:    "IndirectY",
		ArgsLen: 1,
		Format:  func(args []byte) string { return fmt.Sprintf("($%02x), Y", args[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			adl := cpu.RAM.Fetch(int(*args[0]))
			adh := cpu.RAM.Fetch(int(*args[0] + 1))

			addr := int(*adl) | int(*adh)<<8
			if addr/256 != (addr+int(cpu.Reg.Y))/256 && pageBoundryCheck {
				extraCycles++
			}

			op(cpu, cpu.RAM.Fetch(addr+int(cpu.Reg.Y)))
			cpu.Reg.PC += 2
			return
		},
	}

	Implied = AddressingMode{
		Name:    "Implied",
		ArgsLen: 0,
		Format:  func(args []byte) string { return "" },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu)
			cpu.Reg.PC++
			return
		},
	}

	Accumulator = AddressingMode{
		Name:    "Accumulator",
		ArgsLen: 0,
		Format:  func(args []byte) string { return "" },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu, &cpu.Reg.A)
			cpu.Reg.PC++
			return
		},
	}

	Immediate = AddressingMode{
		Name:    "Immediate",
		ArgsLen: 1,
		Format:  func(args []byte) string { return fmt.Sprintf("#$%02x", args[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu, args[0])
			cpu.Reg.PC += 2
			return
		},
	}

	Relative = AddressingMode{
		Name:    "Relative",
		ArgsLen: 1,
		Format:  func(args []byte) string { return fmt.Sprintf("$%02x", args[0]) },

		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			extraCycles = op(cpu, args[0])
			cpu.Reg.PC += 2
			return
		},
	}
)
