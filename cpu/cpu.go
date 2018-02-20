package cpu

type CPU struct {
	ram *RAM

	reg *Registers
}

func (cpu *CPU) push(b byte) {
	*cpu.ram.Fetch(cpu.getStackAddr()) = b
	cpu.reg.sp--
}

func (cpu *CPU) pull() byte {
	cpu.reg.sp++
	return *cpu.ram.Fetch(cpu.getStackAddr())
}

func (cpu *CPU) getStackAddr() int {
	return int(cpu.reg.sp) | (1 << 8)
}
