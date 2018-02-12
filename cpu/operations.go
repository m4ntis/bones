package cpu

// Operation defines an operation that the CPU executes in one or more of it's
// opcodes.
//
// The byte values it received are it's arguments. Arguments can be of any
// length, depending on the operation. There isn't a gurantee that the
// operation will check for the correct number of arguments, so make sure that
// you pass in the correct amount.
//
// The operation also gets a reference to the cpu so it can test and change the
// registers and RAM.
type Operation func(*CPU, []byte)

func ADC(cpu *CPU, args []byte) {
	// Calculate result and store in a
	arg1 := cpu.reg.a
	arg2 := args[0]
	arg3 := cpu.reg.c

	res := arg1 + arg2 + arg3
	cpu.reg.a = res

	// Set flags
	setZ(cpu.reg, res)
	setN(cpu.reg, res)

	// Overflow
	signed_arg1 := int8(arg1)
	signed_arg2 := int8(arg2)
	signed_arg3 := int8(arg3)
	if int(signed_arg1+signed_arg2+signed_arg3) != int(signed_arg1)+
		int(signed_arg2)+int(signed_arg3) {
		cpu.reg.v = SET
	} else {
		cpu.reg.v = RESET
	}

	// Carry
	if int(res) != int(arg1)+int(arg2)+int(arg3) {
		cpu.reg.c = SET
	} else {
		cpu.reg.c = RESET
	}
}

func AND(cpu *CPU, args []byte) {
	cpu.reg.a &= args[0]
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.z)
}

func ASL(cpu *CPU, args []byte) {
	cpu.reg.c = args[0] & 128
	args[0] <<= 1
	setN(cpu.reg, args[0])
	setZ(cpu.reg, args[0])
}

func setZ(reg *Registers, val byte) {
	if val == 0x0 {
		reg.z = SET
		return
	}
	reg.z = RESET
}

func setN(reg *Registers, val byte) {
	reg.n = val & 128
}
