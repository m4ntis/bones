package cpu

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
	name string

	argLen  int
	address func(*CPU, Operation, bool, ...*byte) int
}

var (
	ZeroPage = AddressingMode{
		name: "ZeroPage",

		argLen: 1,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu, cpu.ram.Fetch(int(*args[0])))
			cpu.reg.pc += 2
			return
		},
	}

	ZeroPageX = AddressingMode{
		name: "ZeroPageX",

		argLen: 1,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu, cpu.ram.Fetch(int(*args[0]+cpu.reg.x)))
			cpu.reg.pc += 2
			return
		},
	}

	ZeroPageY = AddressingMode{
		name: "ZeroPageY",

		argLen: 1,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu, cpu.ram.Fetch(int(*args[0]+cpu.reg.y)))
			cpu.reg.pc += 2
			return
		},
	}

	Absolute = AddressingMode{
		name: "Absolute",

		argLen: 2,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			// ADL stored at *args[0], ADH at *args[1]
			op(cpu, cpu.ram.Fetch(int(*args[0])|int(*args[1])<<8))
			cpu.reg.pc += 3
			return
		},
	}

	AbsoluteX = AddressingMode{
		name: "AbsoluteX",

		argLen: 2,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			addr := int(*args[0]) | int(*args[1])<<8
			if addr/256 != (addr+int(cpu.reg.x))/256 && pageBoundryCheck {
				extraCycles++
			}

			op(cpu, cpu.ram.Fetch(addr+int(cpu.reg.x)))
			cpu.reg.pc += 3
			return
		},
	}

	AbsoluteY = AddressingMode{
		name: "AbsoluteY",

		argLen: 2,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			addr := int(*args[0]) | int(*args[1])<<8
			if addr/256 != (addr+int(cpu.reg.y))/256 && pageBoundryCheck {
				extraCycles++
			}

			op(cpu, cpu.ram.Fetch(addr+int(cpu.reg.y)))
			cpu.reg.pc += 3
			return
		},
	}

	// AbsoluteJMP is here because the absolute JMP operation take the immediate
	// arguments instead of their value at the new location like Absolute
	AbsoluteJMP = AddressingMode{
		name: "AbsoluteJMP",

		argLen: 2,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu, args[0], args[1])
			return
		},
	}

	Indirect = AddressingMode{
		name: "Indirect",

		argLen: 2,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			adl := cpu.ram.Fetch(int(*args[0]) | int(*args[1])<<8)
			adh := cpu.ram.Fetch(int(*args[0]) | int(*args[1])<<8 + 1)

			op(cpu, adl, adh)
			return
		},
	}

	IndirectX = AddressingMode{
		name: "IndirectX",

		argLen: 1,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			adl := cpu.ram.Fetch(int(*args[0] + cpu.reg.x))
			adh := cpu.ram.Fetch(int(*args[0] + cpu.reg.x + 1))

			op(cpu, cpu.ram.Fetch(int(*adl)|int(*adh)<<8))
			cpu.reg.pc += 2
			return
		},
	}

	IndirectY = AddressingMode{
		name: "IndirectY",

		argLen: 1,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			adl := cpu.ram.Fetch(int(*args[0]))
			adh := cpu.ram.Fetch(int(*args[0] + 1))

			addr := int(*adl) | int(*adh)<<8
			if addr/256 != (addr+int(cpu.reg.y))/256 && pageBoundryCheck {
				extraCycles++
			}

			op(cpu, cpu.ram.Fetch(addr+int(cpu.reg.y)))
			cpu.reg.pc += 2
			return
		},
	}

	Implied = AddressingMode{
		name: "Implied",

		argLen: 0,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu)
			cpu.reg.pc++
			return
		},
	}

	Accumulator = AddressingMode{
		name: "Accumulator",

		argLen: 0,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu, &cpu.reg.a)
			cpu.reg.pc++
			return
		},
	}

	Immediate = AddressingMode{
		name: "Immediate",

		argLen: 1,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			op(cpu, args[0])
			cpu.reg.pc += 2
			return
		},
	}

	Relative = AddressingMode{
		name: "Relative",

		argLen: 1,
		address: func(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
			extraCycles = op(cpu, args[0])
			cpu.reg.pc += 2
			return
		},
	}
)
