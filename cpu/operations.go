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
// arguments too, changing the underlying RAM or Register.
//
// The operation also gets a reference to the cpu so it can test and change the
// Registers and RAM.
//
// Similar to addressing modes, opcodes too return whether the operation's
// execution took extra cycles. This happens on the operation level only in
// branching operations.
type Operation func(*CPU, ...*byte) int

func ADC(cpu *CPU, args ...*byte) (extraCycles int) {
	// Calculate result and store in a
	arg1 := cpu.Reg.A
	arg2 := *args[0]
	arg3 := cpu.Reg.C

	res := arg1 + arg2 + arg3
	cpu.Reg.A = res

	// Set flags
	setZ(cpu.Reg, res)
	setN(cpu.Reg, res)

	// Overflow
	signed_arg1 := int8(arg1)
	signed_arg2 := int8(arg2)
	signed_arg3 := int8(arg3)
	if int(signed_arg1+signed_arg2+signed_arg3) != int(signed_arg1)+
		int(signed_arg2)+int(signed_arg3) {
		cpu.Reg.V = Set
	} else {
		cpu.Reg.V = Clear
	}

	// Carry
	if int(arg1)+int(arg2)+int(arg3) > 255 {
		cpu.Reg.C = Set
	} else {
		cpu.Reg.C = Clear
	}
	return
}

func AND(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.A &= *args[0]
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.Z)
	return
}

func ASL(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.C = *args[0] >> 7
	*args[0] <<= 1
	setN(cpu.Reg, *args[0])
	setZ(cpu.Reg, *args[0])
	return
}

func BCC(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.C == Clear {
		extraCycles++
		cpu.Reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BCS(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.C == Set {
		extraCycles++
		cpu.Reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BEQ(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.Z == Set {
		extraCycles++
		cpu.Reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BIT(cpu *CPU, args ...*byte) (extraCycles int) {
	setN(cpu.Reg, *args[0])
	cpu.Reg.V = (*args[0] >> 6) & 1

	res := cpu.Reg.A & *args[0]
	setZ(cpu.Reg, res)
	return
}

func BMI(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.N == Set {
		extraCycles++
		cpu.Reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BNE(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.Z == Clear {
		extraCycles++
		cpu.Reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BPL(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.N == Clear {
		extraCycles++
		cpu.Reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BRK(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.IRQ()
	return
}

func BVC(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.V == Clear {
		extraCycles++
		cpu.Reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BVS(cpu *CPU, args ...*byte) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.V == Set {
		extraCycles++
		cpu.Reg.PC += int(int8(*args[0]))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func CLC(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.C = Clear
	return
}

func CLD(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.D = Clear
	return
}

func CLI(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.I = Clear
	return
}

func CLV(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.V = Clear
	return
}

func CMP(cpu *CPU, args ...*byte) (extraCycles int) {
	res := cpu.Reg.A - *args[0]

	setN(cpu.Reg, res)
	setZ(cpu.Reg, res)
	if *args[0] > cpu.Reg.A {
		cpu.Reg.C = Set
	} else {
		cpu.Reg.C = Clear
	}
	return
}

func CPX(cpu *CPU, args ...*byte) (extraCycles int) {
	res := cpu.Reg.X - *args[0]

	setN(cpu.Reg, res)
	setZ(cpu.Reg, res)
	if *args[0] > cpu.Reg.A {
		cpu.Reg.C = Set
	} else {
		cpu.Reg.C = Clear
	}
	return
}

func CPY(cpu *CPU, args ...*byte) (extraCycles int) {
	res := cpu.Reg.Y - *args[0]

	setN(cpu.Reg, res)
	setZ(cpu.Reg, res)
	if *args[0] > cpu.Reg.A {
		cpu.Reg.C = Set
	} else {
		cpu.Reg.C = Clear
	}
	return
}

func DEC(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0]--
	setN(cpu.Reg, *args[0])
	setZ(cpu.Reg, *args[0])
	return
}

func DEX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.X--
	setN(cpu.Reg, cpu.Reg.X)
	setZ(cpu.Reg, cpu.Reg.X)
	return
}

func DEY(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.Y--
	setN(cpu.Reg, cpu.Reg.Y)
	setZ(cpu.Reg, cpu.Reg.Y)
	return
}

func EOR(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.A ^= *args[0]
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.A)
	return
}

func INC(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0]++
	setN(cpu.Reg, *args[0])
	setZ(cpu.Reg, *args[0])
	return
}

func INX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.X++
	setN(cpu.Reg, cpu.Reg.X)
	setZ(cpu.Reg, cpu.Reg.X)
	return
}

func INY(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.Y++
	setN(cpu.Reg, cpu.Reg.Y)
	setZ(cpu.Reg, cpu.Reg.Y)
	return
}

func JMP(cpu *CPU, args ...*byte) (extraCycles int) {
	jmpPC := int(*args[0]) | int(*args[1])<<8
	cpu.Reg.PC = jmpPC
	return
}

func JSR(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.PC += 2
	// push PCH
	cpu.push(byte(cpu.Reg.PC >> 8))
	// push PCL
	cpu.push(byte(cpu.Reg.PC & 0xff))

	jmpPC := int(*args[0]) | int(*args[1])<<8
	cpu.Reg.PC = jmpPC
	return
}

func LDA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.A = *args[0]
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.A)
	return
}

func LDX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.X = *args[0]
	setN(cpu.Reg, cpu.Reg.X)
	setZ(cpu.Reg, cpu.Reg.X)
	return
}

func LDY(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.Y = *args[0]
	setN(cpu.Reg, cpu.Reg.Y)
	setZ(cpu.Reg, cpu.Reg.Y)
	return
}

func LSR(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.C = *args[0] & 1
	*args[0] >>= 1
	cpu.Reg.N = Clear
	setZ(cpu.Reg, *args[0])
	return
}

func NOP(cpu *CPU, args ...*byte) (extraCycles int) {
	return
}

func ORA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.A |= *args[0]
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.A)
	return
}

func PHA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.push(cpu.Reg.A)
	return
}

func PHP(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.push(cpu.Reg.GetP())
	return
}

func PLA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.A = cpu.pull()
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.A)
	return
}

func PLP(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.SetP(cpu.pull())
	return
}

func ROL(cpu *CPU, args ...*byte) (extraCycles int) {
	carry := cpu.Reg.C
	cpu.Reg.C = *args[0] >> 7

	*args[0] <<= 1
	*args[0] |= carry

	setN(cpu.Reg, *args[0])
	setZ(cpu.Reg, *args[0])
	return
}

func ROR(cpu *CPU, args ...*byte) (extraCycles int) {
	carry := cpu.Reg.C
	cpu.Reg.C = *args[0] & 1

	*args[0] >>= 1
	*args[0] |= (carry << 7)

	setN(cpu.Reg, *args[0])
	setZ(cpu.Reg, *args[0])
	return
}

func RTI(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.SetP(cpu.pull())
	// pull PCL and then PHC
	cpu.Reg.PC = int(cpu.pull()) | int(cpu.pull())<<8
	return
}

func RTS(cpu *CPU, args ...*byte) (extraCycles int) {
	// pull PCL and then PHC
	cpu.Reg.PC = int(cpu.pull()) | int(cpu.pull())<<8
	cpu.Reg.PC++
	return
}

func SBC(cpu *CPU, args ...*byte) (extraCycles int) {
	// Calculate result and store in a
	arg1 := int8(cpu.Reg.A)
	arg2 := int8(*args[0])

	res := byte(arg1 - arg2)
	cpu.Reg.A = res

	// Set flags
	setZ(cpu.Reg, res)
	setN(cpu.Reg, res)

	cpu.Reg.C = (res >> 7) ^ 1

	// Overflow
	if math.Abs(float64(arg1)-float64(arg2)) > 127 {
		cpu.Reg.V = Set
	} else {
		cpu.Reg.V = Clear
	}
	return
}

func SEC(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.C = Set
	return
}

func SED(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.D = Set
	return
}

func SEI(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.I = Set
	return
}

func STA(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0] = cpu.Reg.A
	return
}

func STX(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0] = cpu.Reg.X
	return
}

func STY(cpu *CPU, args ...*byte) (extraCycles int) {
	*args[0] = cpu.Reg.Y
	return
}

func TAX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.X = cpu.Reg.A
	setN(cpu.Reg, cpu.Reg.X)
	setZ(cpu.Reg, cpu.Reg.X)
	return
}

func TAY(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.Y = cpu.Reg.A
	setN(cpu.Reg, cpu.Reg.Y)
	setZ(cpu.Reg, cpu.Reg.Y)
	return
}

func TSX(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.X = cpu.Reg.SP
	setN(cpu.Reg, cpu.Reg.X)
	setZ(cpu.Reg, cpu.Reg.X)
	return
}

func TXA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.A = cpu.Reg.X
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.A)
	return
}

func TYA(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.A = cpu.Reg.Y
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.A)
	return
}

func TXS(cpu *CPU, args ...*byte) (extraCycles int) {
	cpu.Reg.SP = cpu.Reg.X
	setN(cpu.Reg, cpu.Reg.SP)
	setZ(cpu.Reg, cpu.Reg.SP)
	return
}

func setZ(Reg *Registers, val byte) {
	if val == 0x0 {
		Reg.Z = Set
		return
	}
	Reg.Z = Clear
}

func setN(Reg *Registers, val byte) {
	Reg.N = val >> 7
}
