package cpu

// AddressingMode defines one of the 2a03's ways of addressing the operands of
// the different opcodes.
//
// Each addressing mode is responsible of fetching the operands in it's way,
// and calling the operation with them
type AddressingMode func(*CPU, Operation, ...*byte)

func ZeroPage(cpu *CPU, op Operation, args ...*byte) {
	op(cpu, cpu.ram.Fetch(int(*args[0])))
	cpu.reg.pc += 2
}

func ZeroPageX(cpu *CPU, op Operation, args ...*byte) {
	op(cpu, cpu.ram.Fetch(int(*args[0]+cpu.reg.x)))
	cpu.reg.pc += 2
}

func ZeroPageY(cpu *CPU, op Operation, args ...*byte) {
	op(cpu, cpu.ram.Fetch(int(*args[0]+cpu.reg.y)))
	cpu.reg.pc += 2
}

func Absolute(cpu *CPU, op Operation, args ...*byte) {
	// ADL stored at *args[0], ADH at *args[1]
	op(cpu, cpu.ram.Fetch(int(*args[0])|int(*args[1])<<8))
	cpu.reg.pc += 3
}

func AbsoluteX(cpu *CPU, op Operation, args ...*byte) {
	// fetch ADL | ADH << 8 + X
	op(cpu, cpu.ram.Fetch(int(*args[0])|int(*args[1])<<8+int(cpu.reg.x)))
	cpu.reg.pc += 3
}

func AbsoluteY(cpu *CPU, op Operation, args ...*byte) {
	// fetch ADL | ADH << 8 + Y
	op(cpu, cpu.ram.Fetch(int(*args[0])|int(*args[1])<<8+int(cpu.reg.y)))
	cpu.reg.pc += 3
}

// AbsoluteJMP is here because the absolute JMP operation take the immediate
// arguments instead of their value at the new location like Absolute
func AbsoluteJMP(cpu *CPU, op Operation, args ...*byte) {
	op(cpu, args[0], args[1])
}

func Indirect(cpu *CPU, op Operation, args ...*byte) {
	adl := cpu.ram.Fetch(int(*args[0]) | int(*args[1])<<8)
	adh := cpu.ram.Fetch(int(*args[0]) | int(*args[1])<<8 + 1)

	op(cpu, adl, adh)
}

func IndirectX(cpu *CPU, op Operation, args ...*byte) {
	adl := cpu.ram.Fetch(int(*args[0] + cpu.reg.x))
	adh := cpu.ram.Fetch(int(*args[0] + cpu.reg.x + 1))

	op(cpu, cpu.ram.Fetch(int(*adl)|int(*adh)<<8))
	cpu.reg.pc += 2
}

func IndirectY(cpu *CPU, op Operation, args ...*byte) {
	adl := cpu.ram.Fetch(int(*args[0]))
	adh := cpu.ram.Fetch(int(*args[0] + 1))

	op(cpu, cpu.ram.Fetch(int(*adl)|int(*adh)<<8+int(cpu.reg.y)))
	cpu.reg.pc += 2
}

func Implied(cpu *CPU, op Operation, args ...*byte) {
	op(cpu)
	cpu.reg.pc++
}

func Accumulator(cpu *CPU, op Operation, args ...*byte) {
	op(cpu, &cpu.reg.a)
	cpu.reg.pc++
}

func Immediate(cpu *CPU, op Operation, args ...*byte) {
	op(cpu, args[0])
	cpu.reg.pc += 2
}

func Relative(cpu *CPU, op Operation, args ...*byte) {
	op(cpu, args[0])
	cpu.reg.pc += 2
}
