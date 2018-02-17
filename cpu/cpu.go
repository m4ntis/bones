package cpu

type CPU struct {
	ram *RAM

	reg *Registers
}

func (cpu *CPU) push(b byte) {
	cpu.ram.Write(getStackAddr())
	cpu.reg.sp--
}

func (cpu *CPU) pull() byte {
	cpu.reg.sp++
	return cpu.ram.Read(getStackAddr())
}

func (cpu *CPU) getStackAddr() int {
	return int(cpu.reg.sp) | (1 << 8)
}
