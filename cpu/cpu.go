package cpu

type CPU struct {
	ram *RAM

	reg *Registers
}

func (cpu *CPU) push(b byte) {
	*cpu.ram.Fetch(cpu.getStackAddr()) = b
	cpu.reg.SP--
}

func (cpu *CPU) pull() byte {
	cpu.reg.SP++
	return *cpu.ram.Fetch(cpu.getStackAddr())
}

func (cpu *CPU) getStackAddr() int {
	return int(cpu.reg.SP) | (1 << 8)
}
