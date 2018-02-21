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
	arg1 := cpu.reg.A
	arg2 := *args[0]
	arg3 := cpu.reg.C

	res := arg1 + arg2 + arg3
	cpu.reg.A = res

	// Set flags
	setZ(cpu.reg, res)
	setN(cpu.reg, res)

	// Overflow
	signed_arg1 := int8(arg1)
	signed_arg2 := int8(arg2)
	signed_arg3 := int8(arg3)
	if int(signed_arg1+signed_arg2+signed_arg3) != int(signed_arg1)+
		int(signed_arg2)+int(signed_arg3) {
		cpu.reg.V = SET
	} else {
		cpu.reg.V = CLEAR
	}

	// Carry
	if int(arg1)+int(arg2)+int(arg3) > 255 {
		cpu.reg.C = SET
	} else {
		cpu.reg.C = CLEAR
	}
	return
}

func AND(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.A &= *args[0]
	setN(cpu.reg, cpu.reg.A)
	setZ(cpu.reg, cpu.reg.Z)
	return
}

func ASL(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.C = *args[0] >> 7
	*args[0] <<= 1
	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
	return
}

func BCC(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.PC

	if cpu.reg.C == CLEAR {
		extraCycles++
		cpu.reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BCS(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.PC

	if cpu.reg.C == SET {
		extraCycles++
		cpu.reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BEQ(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.PC

	if cpu.reg.Z == SET {
		extraCycles++
		cpu.reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BIT(cpu *CPU, args ...*byte) (extraCycles int) {
	setN(cpu.reg, *args[0])
	cpu.reg.V = (*args[0] >> 6) & 1

	res := cpu.reg.A & *args[0]
	setZ(cpu.reg, res)
	return
}

func BMI(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.PC

	if cpu.reg.N == SET {
		extraCycles++
		cpu.reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BNE(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.PC

	if cpu.reg.Z == CLEAR {
		extraCycles++
		cpu.reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BPL(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.PC

	if cpu.reg.N == CLEAR {
		extraCycles++
		cpu.reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BRK(cpu *CPU, args ...*byte) (extraCycles int) {
	// push PCH
	cpu.push(byte(cpu.reg.PC >> 8))
	// push PCL
	cpu.push(byte(cpu.reg.PC & 0xff))
	// push P
	cpu.push(cpu.reg.GetP())

	// fetch PCL from $fffe and PCH from $ffff
	cpu.reg.PC = int(*cpu.ram.Fetch(0xfffe)) | int(*cpu.ram.Fetch(0xffff))<<8
	return
}

func BVC(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.PC

	if cpu.reg.V == CLEAR {
		extraCycles++
		cpu.reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BVS(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.reg.PC

	if cpu.reg.V == SET {
		extraCycles++
		cpu.reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func CLC(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.C = CLEAR
	return
}

func CLD(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.D = CLEAR
	return
}

func CLI(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.I = CLEAR
	return
}

func CLV(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.V = CLEAR
	return
}

func CMP(cpu *CPU, args ...*byte) (extraCycles int) {
	res := cpu.reg.A - *args[0]

	setN(cpu.reg, res)
	setZ(cpu.reg, res)
	if *args[0] > cpu.reg.A {
		cpu.reg.C = SET
	} else {
		cpu.reg.C = CLEAR
	}
	return
}

func CPX(cpu *CPU, args ...*byte) (extraCycles int) {
	res := cpu.reg.X - *args[0]

	setN(cpu.reg, res)
	setZ(cpu.reg, res)
	if *args[0] > cpu.reg.A {
		cpu.reg.C = SET
	} else {
		cpu.reg.C = CLEAR
	}
	return
}

func CPY(cpu *CPU, args ...*byte) (extraCycles int) {
	res := cpu.reg.Y - *args[0]

	setN(cpu.reg, res)
	setZ(cpu.reg, res)
	if *args[0] > cpu.reg.A {
		cpu.reg.C = SET
	} else {
		cpu.reg.C = CLEAR
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
	cpu.reg.X--
	setN(cpu.reg, cpu.reg.X)
	setZ(cpu.reg, cpu.reg.X)
	return
}

func DEY(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.Y--
	setN(cpu.reg, cpu.reg.Y)
	setZ(cpu.reg, cpu.reg.Y)
	return
}

func EOR(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.A ^= *args[0]
	setN(cpu.reg, cpu.reg.A)
	setZ(cpu.reg, cpu.reg.A)
	return
}

func INC(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0]++
	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
	return
}

func INX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.X++
	setN(cpu.reg, cpu.reg.X)
	setZ(cpu.reg, cpu.reg.X)
	return
}

func INY(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.Y++
	setN(cpu.reg, cpu.reg.Y)
	setZ(cpu.reg, cpu.reg.Y)
	return
}

func JMP(cpu *CPU, args ...*byte) (extraCycles int) {
	jmpPC := int(*args[0]) | int(*args[1])<<8
	cpu.reg.PC = jmpPC
	return
}

func JSR(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.PC += 2
	// push PCH
	cpu.push(byte(cpu.reg.PC >> 8))
	// push PCL
	cpu.push(byte(cpu.reg.PC & 0xff))

	jmpPC := int(*args[0]) | int(*args[1])<<8
	cpu.reg.PC = jmpPC
	return
}

func LDA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.A = *args[0]
	setN(cpu.reg, cpu.reg.A)
	setZ(cpu.reg, cpu.reg.A)
	return
}

func LDX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.X = *args[0]
	setN(cpu.reg, cpu.reg.X)
	setZ(cpu.reg, cpu.reg.X)
	return
}

func LDY(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.Y = *args[0]
	setN(cpu.reg, cpu.reg.Y)
	setZ(cpu.reg, cpu.reg.Y)
	return
}

func LSR(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.C = *args[0] & 1
	*args[0] >>= 1
	cpu.reg.N = CLEAR
	setZ(cpu.reg, *args[0])
	return
}

func NOP(cpu *CPU, args ...*byte) (extraCycles int) {
	return
}

func ORA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.A |= *args[0]
	setN(cpu.reg, cpu.reg.A)
	setZ(cpu.reg, cpu.reg.A)
	return
}

func PHA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.push(cpu.reg.A)
	return
}

func PHP(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.push(cpu.reg.GetP())
	return
}

func PLA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.A = cpu.pull()
	setN(cpu.reg, cpu.reg.A)
	setZ(cpu.reg, cpu.reg.A)
	return
}

func PLP(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.SetP(cpu.pull())
	return
}

func ROL(cpu *CPU, args ...*byte) (extraCycles int) {
	carry := cpu.reg.C
	cpu.reg.C = *args[0] >> 7

	*args[0] <<= 1
	*args[0] |= carry

	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
	return
}

func ROR(cpu *CPU, args ...*byte) (extraCycles int) {
	carry := cpu.reg.C
	cpu.reg.C = *args[0] & 1

	*args[0] >>= 1
	*args[0] |= (carry << 7)

	setN(cpu.reg, *args[0])
	setZ(cpu.reg, *args[0])
	return
}

func RTI(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.SetP(cpu.pull())
	// pull PCL and then PHC
	cpu.reg.PC = int(cpu.pull()) | int(cpu.pull())<<8
	return
}

func RTS(cpu *CPU, args ...*byte) (extraCycles int) {
	// pull PCL and then PHC
	cpu.reg.PC = int(cpu.pull()) | int(cpu.pull())<<8
	cpu.reg.PC++
	return
}

func SBC(cpu *CPU, args ...*byte) (extraCycles int) {
	// Calculate result and store in a
	arg1 := int8(cpu.reg.A)
	arg2 := int8(*args[0])

	res := byte(arg1 - arg2)
	cpu.reg.A = res

	// Set flags
	setZ(cpu.reg, res)
	setN(cpu.reg, res)

	cpu.reg.C = (res >> 7) ^ 1

	// Overflow
	if math.Abs(float64(arg1)-float64(arg2)) > 127 {
		cpu.reg.V = SET
	} else {
		cpu.reg.V = CLEAR
	}
	return
}

func SEC(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.C = SET
	return
}

func SED(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.D = SET
	return
}

func SEI(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.I = SET
	return
}

func STA(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0] = cpu.reg.A
	return
}

func STX(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0] = cpu.reg.X
	return
}

func STY(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0] = cpu.reg.Y
	return
}

func TAX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.X = cpu.reg.A
	setN(cpu.reg, cpu.reg.X)
	setZ(cpu.reg, cpu.reg.X)
	return
}

func TAY(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.Y = cpu.reg.A
	setN(cpu.reg, cpu.reg.Y)
	setZ(cpu.reg, cpu.reg.Y)
	return
}

func TSX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.X = cpu.reg.SP
	setN(cpu.reg, cpu.reg.X)
	setZ(cpu.reg, cpu.reg.X)
	return
}

func TXA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.A = cpu.reg.X
	setN(cpu.reg, cpu.reg.A)
	setZ(cpu.reg, cpu.reg.A)
	return
}

func TYA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.A = cpu.reg.Y
	setN(cpu.reg, cpu.reg.A)
	setZ(cpu.reg, cpu.reg.A)
	return
}

func TXS(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.reg.SP = cpu.reg.X
	setN(cpu.reg, cpu.reg.SP)
	setZ(cpu.reg, cpu.reg.SP)
	return
}

func setZ(reg *Registers, val byte) {
	if val == 0x0 {
		reg.Z = SET
		return
	}
	reg.Z = CLEAR
}

func setN(reg *Registers, val byte) {
	reg.N = val >> 7
}
