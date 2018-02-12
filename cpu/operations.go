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
type Operation func(*CPU, ...byte)

func ADC(cpu *CPU, args ...byte) {
	// Calculate result and store in a
	arg1 := cpu.reg.a
	arg2 := args[0]
	arg3 := cpu.reg.c

	res := arg1 + arg2 + arg3
	cpu.reg.a = res

	// Set flags
	if res == 0x0 {
		cpu.reg.z = SET
	}
	cpu.reg.n = res & 128

	signed_arg1 := int8(arg1)
	signed_arg2 := int8(arg2)
	signed_arg3 := int8(arg3)
	if int(signed_arg1+signed_arg2+signed_arg3) != int(signed_arg1)+
		int(signed_arg2)+int(signed_arg3) {
		cpu.reg.v = SET
	} else {
		cpu.reg.v = RESET
	}

	if int(res) != int(arg1)+int(arg2)+int(arg3) {
		cpu.reg.c = SET
	} else {
		cpu.reg.c = RESET
	}
}
