package cpu

type CPU struct {
	RAM *RAM

	Reg *Registers
}

func (cpu *CPU) push(b byte) {
	*cpu.RAM.Fetch(cpu.getStackAddr()) = b
	cpu.Reg.SP--
}

func (cpu *CPU) pull() byte {
	cpu.Reg.SP++
	return *cpu.RAM.Fetch(cpu.getStackAddr())
}

func (cpu *CPU) getStackAddr() int {
	return int(cpu.Reg.SP) | (1 << 8)
}
