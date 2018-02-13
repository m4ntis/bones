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
		cpu.reg.v = CLEAR
	}

	// Carry
	if int(res) != int(arg1)+int(arg2)+int(arg3) {
		cpu.reg.c = SET
	} else {
		cpu.reg.c = CLEAR
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

func BCC(cpu *CPU, args []byte) {
	if cpu.reg.c == CLEAR {
		cpu.reg.pc += int(int8(args[0]))
	}
}

func BCS(cpu *CPU, args []byte) {
	if cpu.reg.c == SET {
		cpu.reg.pc += int(int8(args[0]))
	}
}

func BEQ(cpu *CPU, args []byte) {
	if cpu.reg.z == SET {
		cpu.reg.pc += int(int8(args[0]))
	}
}

func BIT(cpu *CPU, args []byte) {
	setN(cpu.reg, args[0])
	cpu.reg.v = args[0] & 64

	res := cpu.reg.a & args[0]
	setZ(cpu.reg, res)
}

func BMI(cpu *CPU, args []byte) {
	if cpu.reg.n == SET {
		cpu.reg.pc += int(int8(args[0]))
	}
}

func BNE(cpu *CPU, args []byte) {
	if cpu.reg.z == CLEAR {
		cpu.reg.pc += int(int8(args[0]))
	}
}

func BPL(cpu *CPU, args []byte) {
	if cpu.reg.n == CLEAR {
		cpu.reg.pc += int(int8(args[0]))
	}
}

//TODO:BRK

func BVC(cpu *CPU, args []byte) {
	if cpu.reg.v == CLEAR {
		cpu.reg.pc += int(int8(args[0]))
	}
}

func BVS(cpu *CPU, args []byte) {
	if cpu.reg.v == SET {
		cpu.reg.pc += int(int8(args[0]))
	}
}

func CLC(cpu *CPU, args []byte) {
	cpu.reg.c = CLEAR
}

func CLD(cpu *CPU, args []byte) {
	cpu.reg.d = CLEAR
}

func CLI(cpu *CPU, args []byte) {
	cpu.reg.i = CLEAR
}

func CLV(cpu *CPU, args []byte) {
	cpu.reg.v = CLEAR
}

func CMP(cpu *CPU, args []byte) {
	res := cpu.reg.a - args[0]

	setN(cpu.reg, res)
	setZ(cpu.reg, res)
	if args[0] > cpu.reg.a {
		cpu.reg.c = SET
	} else {
		cpu.reg.c = CLEAR
	}
}

func CPX(cpu *CPU, args []byte) {
	res := cpu.reg.x - args[0]

	setN(cpu.reg, res)
	setZ(cpu.reg, res)
	if args[0] > cpu.reg.a {
		cpu.reg.c = SET
	} else {
		cpu.reg.c = CLEAR
	}
}

func CPY(cpu *CPU, args []byte) {
	res := cpu.reg.y - args[0]

	setN(cpu.reg, res)
	setZ(cpu.reg, res)
	if args[0] > cpu.reg.a {
		cpu.reg.c = SET
	} else {
		cpu.reg.c = CLEAR
	}
}

func DEC(cpu *CPU, args []byte) {
	args[0]--
	setN(cpu.reg, args[0])
	setZ(cpu.reg, args[0])
}

func DEX(cpu *CPU, args []byte) {
	cpu.reg.x--
	setN(cpu.reg, args[0])
	setZ(cpu.reg, args[0])
}

func DEY(cpu *CPU, args []byte) {
	cpu.reg.y--
	setN(cpu.reg, args[0])
	setZ(cpu.reg, args[0])
}

func EOR(cpu *CPU, args []byte) {
	cpu.reg.a ^= args[0]
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.z)
}

func INC(cpu *CPU, args []byte) {
	args[0]++
	setN(cpu.reg, args[0])
	setZ(cpu.reg, args[0])
}

func INX(cpu *CPU, args []byte) {
	cpu.reg.x++
	setN(cpu.reg, args[0])
	setZ(cpu.reg, args[0])
}

func INY(cpu *CPU, args []byte) {
	cpu.reg.y++
	setN(cpu.reg, args[0])
	setZ(cpu.reg, args[0])
}

//TODO: Think about how JMP fits within our design

func setZ(reg *Registers, val byte) {
	if val == 0x0 {
		reg.z = SET
		return
	}
	reg.z = CLEAR
}

func setN(reg *Registers, val byte) {
	reg.n = val & 128
}
