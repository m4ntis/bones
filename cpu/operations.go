package cpu

import (
	"math"
)

// Operation defines an operation that the CPU executes in one or more of it's
// opcodes.
//
// The byte values it received are it's arguments. Arguments can be of any
// length, depending on the operation. There isn't a gurantee that the
// operation will check for the correct number of arguments, so make sure that
// you pass in the correct amount. Note that the args are in the form of
// pointers to bytes. This is for the operation to be able to write to the
// arguments too, changing the underlying ram or register.
//
// The operation also gets a reference to the cpu so it can test and change the
// registers and RAM.
//
// Similar to addressing modes, opcodes too return whether the operation's
// execution took extra cycles. This happens on the operation level only in
// branching operations.
type Operation func(*CPU, ...*byte) int

func ADC(cpu *CPU, args ...*byte) (extraCycles int) {
	// Calculate result and store in a
	arg1 := cpu.reg.a
	arg2 := *args[0]
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
	if int(arg1)+int(arg2)+int(arg3) > 255 {
		cpu.reg.c = SET
	} else {
		cpu.reg.c = CLEAR
	}
	return
}

func AND(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.a &= *args[0]
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.z)
	return
}

func ASL(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.c = *args[0] >> 7
	*args[0] <<= 1
	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
	return
}

func BCC(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.pc

	if cpu.reg.c == CLEAR {
		extraCycles++
		cpu.reg.pc += int(int8(*args[0]))

		if initPC/256 != cpu.reg.pc/256 {
			extraCycles++
		}
	}
	return
}

func BCS(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.pc

	if cpu.reg.c == SET {
		extraCycles++
		cpu.reg.pc += int(int8(*args[0]))

		if initPC/256 != cpu.reg.pc/256 {
			extraCycles++
		}
	}
	return
}

func BEQ(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.pc

	if cpu.reg.z == SET {
		extraCycles++
		cpu.reg.pc += int(int8(*args[0]))

		if initPC/256 != cpu.reg.pc/256 {
			extraCycles++
		}
	}
	return
}

func BIT(cpu *CPU, args ...*byte) (extraCycles int) {
	setN(cpu.reg, *args[0])
	cpu.reg.v = (*args[0] >> 6) & 1

	res := cpu.reg.a & *args[0]
	setZ(cpu.reg, res)
	return
}

func BMI(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.pc

	if cpu.reg.n == SET {
		extraCycles++
		cpu.reg.pc += int(int8(*args[0]))

		if initPC/256 != cpu.reg.pc/256 {
			extraCycles++
		}
	}
	return
}

func BNE(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.pc

	if cpu.reg.z == CLEAR {
		extraCycles++
		cpu.reg.pc += int(int8(*args[0]))

		if initPC/256 != cpu.reg.pc/256 {
			extraCycles++
		}
	}
	return
}

func BPL(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.pc

	if cpu.reg.n == CLEAR {
		extraCycles++
		cpu.reg.pc += int(int8(*args[0]))

		if initPC/256 != cpu.reg.pc/256 {
			extraCycles++
		}
	}
	return
}

func BRK(cpu *CPU, args ...*byte) (extraCycles int) {
	// push PCH
	cpu.push(byte(cpu.reg.pc >> 8))
	// push PCL
	cpu.push(byte(cpu.reg.pc & 0xff))
	// push P
	cpu.push(cpu.reg.getP())

	// fetch PCL from $fffe and PCH from $ffff
	cpu.reg.pc = int(*cpu.ram.Fetch(0xfffe)) | int(*cpu.ram.Fetch(0xffff))<<8
	return
}

func BVC(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.pc

	if cpu.reg.v == CLEAR {
		extraCycles++
		cpu.reg.pc += int(int8(*args[0]))

		if initPC/256 != cpu.reg.pc/256 {
			extraCycles++
		}
	}
	return
}

func BVS(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.pc

	if cpu.reg.v == SET {
		extraCycles++
		cpu.reg.pc += int(int8(*args[0]))

		if initPC/256 != cpu.reg.pc/256 {
			extraCycles++
		}
	}
	return
}

func CLC(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.c = CLEAR
	return
}

func CLD(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.d = CLEAR
	return
}

func CLI(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.i = CLEAR
	return
}

func CLV(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.v = CLEAR
	return
}

func CMP(cpu *CPU, args ...*byte) (extraCycles int) {
	res := cpu.reg.a - *args[0]

	setN(cpu.reg, res)
	setZ(cpu.reg, res)
	if *args[0] > cpu.reg.a {
		cpu.reg.c = SET
	} else {
		cpu.reg.c = CLEAR
	}
	return
}

func CPX(cpu *CPU, args ...*byte) (extraCycles int) {
	res := cpu.reg.x - *args[0]

	setN(cpu.reg, res)
	setZ(cpu.reg, res)
	if *args[0] > cpu.reg.a {
		cpu.reg.c = SET
	} else {
		cpu.reg.c = CLEAR
	}
	return
}

func CPY(cpu *CPU, args ...*byte) (extraCycles int) {
	res := cpu.reg.y - *args[0]

	setN(cpu.reg, res)
	setZ(cpu.reg, res)
	if *args[0] > cpu.reg.a {
		cpu.reg.c = SET
	} else {
		cpu.reg.c = CLEAR
	}
	return
}

func DEC(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0]--
	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
	return
}

func DEX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.x--
	setN(cpu.reg, cpu.reg.x)
	setZ(cpu.reg, cpu.reg.x)
	return
}

func DEY(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.y--
	setN(cpu.reg, cpu.reg.y)
	setZ(cpu.reg, cpu.reg.y)
	return
}

func EOR(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.a ^= *args[0]
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.a)
	return
}

func INC(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0]++
	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
	return
}

func INX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.x++
	setN(cpu.reg, cpu.reg.x)
	setZ(cpu.reg, cpu.reg.x)
	return
}

func INY(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.y++
	setN(cpu.reg, cpu.reg.y)
	setZ(cpu.reg, cpu.reg.y)
	return
}

func JMP(cpu *CPU, args ...*byte) (extraCycles int) {
	jmpPC := int(*args[0]) | int(*args[1])<<8
	cpu.reg.pc = jmpPC
	return
}

func JSR(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.pc += 2
	// push PCH
	cpu.push(byte(cpu.reg.pc >> 8))
	// push PCL
	cpu.push(byte(cpu.reg.pc & 0xff))

	jmpPC := int(*args[0]) | int(*args[1])<<8
	cpu.reg.pc = jmpPC
	return
}

func LDA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.a = *args[0]
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.a)
	return
}

func LDX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.x = *args[0]
	setN(cpu.reg, cpu.reg.x)
	setZ(cpu.reg, cpu.reg.x)
	return
}

func LDY(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.y = *args[0]
	setN(cpu.reg, cpu.reg.y)
	setZ(cpu.reg, cpu.reg.y)
	return
}

func LSR(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.c = *args[0] & 1
	*args[0] >>= 1
	cpu.reg.n = CLEAR
	setZ(cpu.reg, *args[0])
	return
}

func NOP(cpu *CPU, args ...*byte) (extraCycles int) {
	return
}

func ORA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.a |= *args[0]
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.a)
	return
}

func PHA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.push(cpu.reg.a)
	return
}

func PHP(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.push(cpu.reg.getP())
	return
}

func PLA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.a = cpu.pull()
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.a)
	return
}

func PLP(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.setP(cpu.pull())
	return
}

func ROL(cpu *CPU, args ...*byte) (extraCycles int) {
	carry := cpu.reg.c
	cpu.reg.c = *args[0] >> 7

	*args[0] <<= 1
	*args[0] |= carry

	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
	return
}

func ROR(cpu *CPU, args ...*byte) (extraCycles int) {
	carry := cpu.reg.c
	cpu.reg.c = *args[0] & 1

	*args[0] >>= 1
	*args[0] |= (carry << 7)

	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
	return
}

func RTI(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.setP(cpu.pull())
	// pull PCL and then PHC
	cpu.reg.pc = int(cpu.pull()) | int(cpu.pull())<<8
	return
}

func RTS(cpu *CPU, args ...*byte) (extraCycles int) {
	// pull PCL and then PHC
	cpu.reg.pc = int(cpu.pull()) | int(cpu.pull())<<8
	cpu.reg.pc++
	return
}

func SBC(cpu *CPU, args ...*byte) (extraCycles int) {
	// Calculate result and store in a
	arg1 := int8(cpu.reg.a)
	arg2 := int8(*args[0])

	res := byte(arg1 - arg2)
	cpu.reg.a = res

	// Set flags
	setZ(cpu.reg, res)
	setN(cpu.reg, res)

	cpu.reg.c = (res >> 7) ^ 1

	// Overflow
	if math.Abs(float64(arg1)-float64(arg2)) > 127 {
		cpu.reg.v = SET
	} else {
		cpu.reg.v = CLEAR
	}
	return
}

func SEC(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.c = SET
	return
}

func SED(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.d = SET
	return
}

func SEI(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.i = SET
	return
}

func STA(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0] = cpu.reg.a
	return
}

func STX(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0] = cpu.reg.x
	return
}

func STY(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0] = cpu.reg.y
	return
}

func TAX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.x = cpu.reg.a
	setN(cpu.reg, cpu.reg.x)
	setZ(cpu.reg, cpu.reg.x)
	return
}

func TAY(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.y = cpu.reg.a
	setN(cpu.reg, cpu.reg.y)
	setZ(cpu.reg, cpu.reg.y)
	return
}

func TSX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.x = cpu.reg.sp
	setN(cpu.reg, cpu.reg.x)
	setZ(cpu.reg, cpu.reg.x)
	return
}

func TXA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.a = cpu.reg.x
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.a)
	return
}

func TYA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.a = cpu.reg.y
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.a)
	return
}

func TXS(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.sp = cpu.reg.x
	setN(cpu.reg, cpu.reg.sp)
	setZ(cpu.reg, cpu.reg.sp)
	return
}

func setZ(reg *Registers, val byte) {
	if val == 0x0 {
		reg.z = SET
		return
	}
	reg.z = CLEAR
}

func setN(reg *Registers, val byte) {
	reg.n = val >> 7
}
