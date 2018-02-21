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
type AddressingMode func(*CPU, Operation, bool, ...*byte) int

func ZeroPage(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	op(cpu, cpu.ram.Fetch(int(*args[0])))
	cpu.reg.pc += 2
	return
}

func ZeroPageX(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	op(cpu, cpu.ram.Fetch(int(*args[0]+cpu.reg.x)))
	cpu.reg.pc += 2
	return
}

func ZeroPageY(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	op(cpu, cpu.ram.Fetch(int(*args[0]+cpu.reg.y)))
	cpu.reg.pc += 2
	return
}

func Absolute(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	// ADL stored at *args[0], ADH at *args[1]
	op(cpu, cpu.ram.Fetch(int(*args[0])|int(*args[1])<<8))
	cpu.reg.pc += 3
	return
}

func AbsoluteX(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	addr := int(*args[0]) | int(*args[1])<<8
	if addr/256 != (addr+int(cpu.reg.x))/256 && pageBoundryCheck {
		extraCycles++
	}

	op(cpu, cpu.ram.Fetch(addr+int(cpu.reg.x)))
	cpu.reg.pc += 3
	return
}

func AbsoluteY(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	addr := int(*args[0]) | int(*args[1])<<8
	if addr/256 != (addr+int(cpu.reg.y))/256 && pageBoundryCheck {
		extraCycles++
	}

	op(cpu, cpu.ram.Fetch(addr+int(cpu.reg.y)))
	cpu.reg.pc += 3
	return
}

// AbsoluteJMP is here because the absolute JMP operation take the immediate
// arguments instead of their value at the new location like Absolute
func AbsoluteJMP(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	op(cpu, args[0], args[1])
	return
}

func Indirect(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	adl := cpu.ram.Fetch(int(*args[0]) | int(*args[1])<<8)
	adh := cpu.ram.Fetch(int(*args[0]) | int(*args[1])<<8 + 1)

	op(cpu, adl, adh)
	return
}

func IndirectX(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	adl := cpu.ram.Fetch(int(*args[0] + cpu.reg.x))
	adh := cpu.ram.Fetch(int(*args[0] + cpu.reg.x + 1))

	op(cpu, cpu.ram.Fetch(int(*adl)|int(*adh)<<8))
	cpu.reg.pc += 2
	return
}

func IndirectY(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	adl := cpu.ram.Fetch(int(*args[0]))
	adh := cpu.ram.Fetch(int(*args[0] + 1))

	addr := int(*adl) | int(*adh)<<8
	if addr/256 != (addr+int(cpu.reg.y))/256 && pageBoundryCheck {
		extraCycles++
	}

	op(cpu, cpu.ram.Fetch(addr+int(cpu.reg.y)))
	cpu.reg.pc += 2
	return
}

func Implied(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	op(cpu)
	cpu.reg.pc++
	return
}

func Accumulator(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	op(cpu, &cpu.reg.a)
	cpu.reg.pc++
	return
}

func Immediate(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	op(cpu, args[0])
	cpu.reg.pc += 2
	return
}

func Relative(cpu *CPU, op Operation, pageBoundryCheck bool, args ...*byte) (extraCycles int) {
	extraCycles = op(cpu, args[0])
	cpu.reg.pc += 2
	return
}
