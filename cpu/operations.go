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
type Operation func(*CPU, ...*byte)

func ADC(cpu *CPU, args ...*byte) {
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
}

func AND(cpu *CPU, args ...*byte) {
	cpu.reg.a &= *args[0]
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.z)
}

func ASL(cpu *CPU, args ...*byte) {
	cpu.reg.c = *args[0] >> 7
	*args[0] <<= 1
	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
}

func BCC(cpu *CPU, args ...*byte) {
	if cpu.reg.c == CLEAR {
		cpu.reg.pc += int(int8(*args[0]))
	}
}

func BCS(cpu *CPU, args ...*byte) {
	if cpu.reg.c == SET {
		cpu.reg.pc += int(int8(*args[0]))
	}
}

func BEQ(cpu *CPU, args ...*byte) {
	if cpu.reg.z == SET {
		cpu.reg.pc += int(int8(*args[0]))
	}
}

func BIT(cpu *CPU, args ...*byte) {
	setN(cpu.reg, *args[0])
	cpu.reg.v = (*args[0] >> 6) & 1

	res := cpu.reg.a & *args[0]
	setZ(cpu.reg, res)
}

func BMI(cpu *CPU, args ...*byte) {
	if cpu.reg.n == SET {
		cpu.reg.pc += int(int8(*args[0]))
	}
}

func BNE(cpu *CPU, args ...*byte) {
	if cpu.reg.z == CLEAR {
		cpu.reg.pc += int(int8(*args[0]))
	}
}

func BPL(cpu *CPU, args ...*byte) {
	if cpu.reg.n == CLEAR {
		cpu.reg.pc += int(int8(*args[0]))
	}
}

func BRK(cpu *CPU, args ...*byte) {
	// push PCH
	cpu.push(byte(cpu.reg.pc >> 8))
	// push PCL
	cpu.push(byte(cpu.reg.pc & 0xff))
	// push P
	cpu.push(cpu.reg.getP())

	// fetch PCL from $fffe and PCH from $ffff
	cpu.reg.pc = int(*cpu.ram.Fetch(0xfffe)) | int(*cpu.ram.Fetch(0xffff))<<8
}

func BVC(cpu *CPU, args ...*byte) {
	if cpu.reg.v == CLEAR {
		cpu.reg.pc += int(int8(*args[0]))
	}
}

func BVS(cpu *CPU, args ...*byte) {
	if cpu.reg.v == SET {
		cpu.reg.pc += int(int8(*args[0]))
	}
}

func CLC(cpu *CPU, args ...*byte) {
	cpu.reg.c = CLEAR
}

func CLD(cpu *CPU, args ...*byte) {
	cpu.reg.d = CLEAR
}

func CLI(cpu *CPU, args ...*byte) {
	cpu.reg.i = CLEAR
}

func CLV(cpu *CPU, args ...*byte) {
	cpu.reg.v = CLEAR
}

func CMP(cpu *CPU, args ...*byte) {
	res := cpu.reg.a - *args[0]

	setN(cpu.reg, res)
	setZ(cpu.reg, res)
	if *args[0] > cpu.reg.a {
		cpu.reg.c = SET
	} else {
		cpu.reg.c = CLEAR
	}
}

func CPX(cpu *CPU, args ...*byte) {
	res := cpu.reg.x - *args[0]

	setN(cpu.reg, res)
	setZ(cpu.reg, res)
	if *args[0] > cpu.reg.a {
		cpu.reg.c = SET
	} else {
		cpu.reg.c = CLEAR
	}
}

func CPY(cpu *CPU, args ...*byte) {
	res := cpu.reg.y - *args[0]

	setN(cpu.reg, res)
	setZ(cpu.reg, res)
	if *args[0] > cpu.reg.a {
		cpu.reg.c = SET
	} else {
		cpu.reg.c = CLEAR
	}
}

func DEC(cpu *CPU, args ...*byte) {
	*args[0]--
	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
}

func DEX(cpu *CPU, args ...*byte) {
	cpu.reg.x--
	setN(cpu.reg, cpu.reg.x)
	setZ(cpu.reg, cpu.reg.x)
}

func DEY(cpu *CPU, args ...*byte) {
	cpu.reg.y--
	setN(cpu.reg, cpu.reg.y)
	setZ(cpu.reg, cpu.reg.y)
}

func EOR(cpu *CPU, args ...*byte) {
	cpu.reg.a ^= *args[0]
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.a)
}

func INC(cpu *CPU, args ...*byte) {
	*args[0]++
	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
}

func INX(cpu *CPU, args ...*byte) {
	cpu.reg.x++
	setN(cpu.reg, cpu.reg.x)
	setZ(cpu.reg, cpu.reg.x)
}

func INY(cpu *CPU, args ...*byte) {
	cpu.reg.y++
	setN(cpu.reg, cpu.reg.y)
	setZ(cpu.reg, cpu.reg.y)
}

func JMP(cpu *CPU, args ...*byte) {
	jmpPC := int(*args[0]) | int(*args[1])<<8
	cpu.reg.pc = jmpPC
}

func JSR(cpu *CPU, args ...*byte) {
	cpu.reg.pc += 2
	// push PCH
	cpu.push(byte(cpu.reg.pc >> 8))
	// push PCL
	cpu.push(byte(cpu.reg.pc & 0xff))

	jmpPC := int(*args[0]) | int(*args[1])<<8
	cpu.reg.pc = jmpPC
}

func LDA(cpu *CPU, args ...*byte) {
	cpu.reg.a = *args[0]
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.a)
}

func LDX(cpu *CPU, args ...*byte) {
	cpu.reg.x = *args[0]
	setN(cpu.reg, cpu.reg.x)
	setZ(cpu.reg, cpu.reg.x)
}

func LDY(cpu *CPU, args ...*byte) {
	cpu.reg.y = *args[0]
	setN(cpu.reg, cpu.reg.y)
	setZ(cpu.reg, cpu.reg.y)
}

func LSR(cpu *CPU, args ...*byte) {
	cpu.reg.c = *args[0] & 1
	*args[0] >>= 1
	cpu.reg.n = CLEAR
	setZ(cpu.reg, *args[0])
}

func NOP(cpu *CPU, args ...*byte) {
}

func ORA(cpu *CPU, args ...*byte) {
	cpu.reg.a |= *args[0]
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.a)
}

func PHA(cpu *CPU, args ...*byte) {
	cpu.push(cpu.reg.a)
}

func PHP(cpu *CPU, args ...*byte) {
	cpu.push(cpu.reg.getP())
}

func PLA(cpu *CPU, args ...*byte) {
	cpu.reg.a = cpu.pull()
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.a)
}

func PLP(cpu *CPU, args ...*byte) {
	cpu.reg.setP(cpu.pull())
}

func ROL(cpu *CPU, args ...*byte) {
	carry := cpu.reg.c
	cpu.reg.c = *args[0] >> 7

	*args[0] <<= 1
	*args[0] |= carry

	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
}

func ROR(cpu *CPU, args ...*byte) {
	carry := cpu.reg.c
	cpu.reg.c = *args[0] & 1

	*args[0] >>= 1
	*args[0] |= (carry << 7)

	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
}

func RTI(cpu *CPU, args ...*byte) {
	cpu.reg.setP(cpu.pull())
	// pull PCL and then PHC
	cpu.reg.pc = int(cpu.pull()) | int(cpu.pull())<<8
}

func RTS(cpu *CPU, args ...*byte) {
	// pull PCL and then PHC
	cpu.reg.pc = int(cpu.pull()) | int(cpu.pull())<<8
	cpu.reg.pc++
}

func SBC(cpu *CPU, args ...*byte) {
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
}

func SEC(cpu *CPU, args ...*byte) {
	cpu.reg.c = SET
}

func SED(cpu *CPU, args ...*byte) {
	cpu.reg.d = SET
}

func SEI(cpu *CPU, args ...*byte) {
	cpu.reg.i = SET
}

func STA(cpu *CPU, args ...*byte) {
	*args[0] = cpu.reg.a
}

func STX(cpu *CPU, args ...*byte) {
	*args[0] = cpu.reg.x
}

func STY(cpu *CPU, args ...*byte) {
	*args[0] = cpu.reg.y
}

func TAX(cpu *CPU, args ...*byte) {
	cpu.reg.x = cpu.reg.a
	setN(cpu.reg, cpu.reg.x)
	setZ(cpu.reg, cpu.reg.x)
}

func TAY(cpu *CPU, args ...*byte) {
	cpu.reg.y = cpu.reg.a
	setN(cpu.reg, cpu.reg.y)
	setZ(cpu.reg, cpu.reg.y)
}

func TSX(cpu *CPU, args ...*byte) {
	cpu.reg.x = cpu.reg.sp
	setN(cpu.reg, cpu.reg.x)
	setZ(cpu.reg, cpu.reg.x)
}

func TXA(cpu *CPU, args ...*byte) {
	cpu.reg.a = cpu.reg.x
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.a)
}

func TYA(cpu *CPU, args ...*byte) {
	cpu.reg.a = cpu.reg.y
	setN(cpu.reg, cpu.reg.a)
	setZ(cpu.reg, cpu.reg.a)
}

func TXS(cpu *CPU, args ...*byte) {
	cpu.reg.sp = cpu.reg.x
	setN(cpu.reg, cpu.reg.sp)
	setZ(cpu.reg, cpu.reg.sp)
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
