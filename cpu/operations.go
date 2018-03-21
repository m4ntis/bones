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
type Operation func(*CPU, ...Operand) int

func ADC(cpu *CPU, ops ...Operand) (extraCycles int) {
	// Calculate result and store in a
	arg1 := cpu.Reg.A
	arg2 := ops[0].Read()
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

func AND(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.A &= ops[0].Read()
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.Z)
	return
}

func ASL(cpu *CPU, ops ...Operand) (extraCycles int) {
	op := ops[0].Read()

	cpu.Reg.C = op >> 7
	op <<= 1

	setN(cpu.Reg, op)
	setZ(cpu.Reg, op)

	ops[0].Write(op)
	return
}

func BCC(cpu *CPU, ops ...Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.C == Clear {
		extraCycles++
		cpu.Reg.PC += int(int8(ops[0].Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BCS(cpu *CPU, ops ...Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.C == Set {
		extraCycles++
		cpu.Reg.PC += int(int8(ops[0].Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BEQ(cpu *CPU, ops ...Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.Z == Set {
		extraCycles++
		cpu.Reg.PC += int(int8(ops[0].Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BIT(cpu *CPU, ops ...Operand) (extraCycles int) {
	op := ops[0].Read()

	setN(cpu.Reg, op)
	cpu.Reg.V = (op >> 6) & 1

	res := cpu.Reg.A & op
	setZ(cpu.Reg, res)
	return
}

func BMI(cpu *CPU, ops ...Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.N == Set {
		extraCycles++
		cpu.Reg.PC += int(int8(ops[0].Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BNE(cpu *CPU, ops ...Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.Z == Clear {
		extraCycles++
		cpu.Reg.PC += int(int8(ops[0].Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BPL(cpu *CPU, ops ...Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.N == Clear {
		extraCycles++
		cpu.Reg.PC += int(int8(ops[0].Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BRK(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.IRQ()
	return
}

func BVC(cpu *CPU, ops ...Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.V == Clear {
		extraCycles++
		cpu.Reg.PC += int(int8(ops[0].Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BVS(cpu *CPU, ops ...Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.V == Set {
		extraCycles++
		cpu.Reg.PC += int(int8(ops[0].Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func CLC(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.C = Clear
	return
}

func CLD(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.D = Clear
	return
}

func CLI(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.I = Clear
	return
}

func CLV(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.V = Clear
	return
}

func CMP(cpu *CPU, ops ...Operand) (extraCycles int) {
	op := ops[0].Read()

	res := cpu.Reg.A - op

	setN(cpu.Reg, res)
	setZ(cpu.Reg, res)
	if op > cpu.Reg.A {
		cpu.Reg.C = Set
	} else {
		cpu.Reg.C = Clear
	}
	return
}

func CPX(cpu *CPU, ops ...Operand) (extraCycles int) {
	op := ops[0].Read()

	res := cpu.Reg.X - op

	setN(cpu.Reg, res)
	setZ(cpu.Reg, res)
	if op > cpu.Reg.A {
		cpu.Reg.C = Set
	} else {
		cpu.Reg.C = Clear
	}
	return
}

func CPY(cpu *CPU, ops ...Operand) (extraCycles int) {
	op := ops[0].Read()

	res := cpu.Reg.Y - op

	setN(cpu.Reg, res)
	setZ(cpu.Reg, res)
	if op > cpu.Reg.A {
		cpu.Reg.C = Set
	} else {
		cpu.Reg.C = Clear
	}
	return
}

func DEC(cpu *CPU, ops ...Operand) (extraCycles int) {
	op = ops[0].Read() - 1

	setN(cpu.Reg, op)
	setZ(cpu.Reg, op)

	ops[0].Write(op)
	return
}

func DEX(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.X--
	setN(cpu.Reg, cpu.Reg.X)
	setZ(cpu.Reg, cpu.Reg.X)
	return
}

func DEY(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.Y--
	setN(cpu.Reg, cpu.Reg.Y)
	setZ(cpu.Reg, cpu.Reg.Y)
	return
}

func EOR(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.A ^= ops[0].Read()
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.A)
	return
}

func INC(cpu *CPU, ops ...Operand) (extraCycles int) {
	op = ops[0].Read() + 1

	setN(cpu.Reg, op)
	setZ(cpu.Reg, op)

	ops[0].Write(op)
	return
}

func INX(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.X++
	setN(cpu.Reg, cpu.Reg.X)
	setZ(cpu.Reg, cpu.Reg.X)
	return
}

func INY(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.Y++
	setN(cpu.Reg, cpu.Reg.Y)
	setZ(cpu.Reg, cpu.Reg.Y)
	return
}

func JMP(cpu *CPU, ops ...Operand) (extraCycles int) {
	jmpPC := int(ops[0].Read()) | int(ops[1].Read())<<8
	cpu.Reg.PC = jmpPC
	return
}

func JSR(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.PC += 2
	// push PCH
	cpu.push(byte(cpu.Reg.PC >> 8))
	// push PCL
	cpu.push(byte(cpu.Reg.PC & 0xff))

	jmpPC := int(ops[0]).Read() | int(ops[1].Read())<<8
	cpu.Reg.PC = jmpPC
	return
}

func LDA(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.A = ops[0].Read()
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.A)
	return
}

func LDX(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.X = ops[0].Read()
	setN(cpu.Reg, cpu.Reg.X)
	setZ(cpu.Reg, cpu.Reg.X)
	return
}

func LDY(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.Y = ops[0].Read()
	setN(cpu.Reg, cpu.Reg.Y)
	setZ(cpu.Reg, cpu.Reg.Y)
	return
}

func LSR(cpu *CPU, ops ...Operand) (extraCycles int) {
	op := ops[0].Read()

	cpu.Reg.C = op & 1
	op >>= 1

	cpu.Reg.N = Clear
	setZ(cpu.Reg, op)

	ops[0].Write(op)
	return
}

func NOP(cpu *CPU, ops ...Operand) (extraCycles int) {
	return
}

func ORA(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.A |= ops[0].Read()
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.A)
	return
}

func PHA(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.push(cpu.Reg.A)
	return
}

func PHP(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.push(cpu.Reg.GetP())
	return
}

func PLA(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.A = cpu.pull()
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.A)
	return
}

func PLP(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.SetP(cpu.pull())
	return
}

func ROL(cpu *CPU, ops ...Operand) (extraCycles int) {
	op := ops[0].Read()

	carry := cpu.Reg.C
	cpu.Reg.C = op >> 7

	op <<= 1
	op |= carry

	setN(cpu.Reg, op)
	setZ(cpu.Reg, op)

	ops[0].Read()
	return
}

func ROR(cpu *CPU, ops ...Operand) (extraCycles int) {
	op := ops[0].Read()

	carry := cpu.Reg.C
	cpu.Reg.C = op & 1

	op >>= 1
	op |= (carry << 7)

	setN(cpu.Reg, op)
	setZ(cpu.Reg, op)

	ops[0].Read()
	return
}

func RTI(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.SetP(cpu.pull())
	// pull PCL and then PHC
	cpu.Reg.PC = int(cpu.pull()) | int(cpu.pull())<<8
	return
}

func RTS(cpu *CPU, ops ...Operand) (extraCycles int) {
	// pull PCL and then PHC
	cpu.Reg.PC = int(cpu.pull()) | int(cpu.pull())<<8
	cpu.Reg.PC++
	return
}

func SBC(cpu *CPU, ops ...Operand) (extraCycles int) {
	// Calculate result and store in a
	arg1 := int8(cpu.Reg.A)
	arg2 := int8(ops[0].Read())

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

func SEC(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.C = Set
	return
}

func SED(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.D = Set
	return
}

func SEI(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.I = Set
	return
}

func STA(cpu *CPU, ops ...Operand) (extraCycles int) {
	ops[0].Write(cpu.Reg.A)
	return
}

func STX(cpu *CPU, ops ...Operand) (extraCycles int) {
	ops[0].Write(cpu.Reg.X)
	return
}

func STY(cpu *CPU, ops ...Operand) (extraCycles int) {
	ops[0].Write(cpu.Reg.Y)
	return
}

func TAX(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.X = cpu.Reg.A
	setN(cpu.Reg, cpu.Reg.X)
	setZ(cpu.Reg, cpu.Reg.X)
	return
}

func TAY(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.Y = cpu.Reg.A
	setN(cpu.Reg, cpu.Reg.Y)
	setZ(cpu.Reg, cpu.Reg.Y)
	return
}

func TSX(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.X = cpu.Reg.SP
	setN(cpu.Reg, cpu.Reg.X)
	setZ(cpu.Reg, cpu.Reg.X)
	return
}

func TXA(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.A = cpu.Reg.X
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.A)
	return
}

func TYA(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.A = cpu.Reg.Y
	setN(cpu.Reg, cpu.Reg.A)
	setZ(cpu.Reg, cpu.Reg.A)
	return
}

func TXS(cpu *CPU, ops ...Operand) (extraCycles int) {
	cpu.Reg.SP = cpu.Reg.X
	setN(cpu.Reg, cpu.Reg.SP)
	setZ(cpu.Reg, cpu.Reg.SP)
	return
}

func setZ(Reg *Registers, d byte) {
	if d == 0x0 {
		Reg.Z = Set
		return
	}
	Reg.Z = Clear
}

func setN(Reg *Registers, d byte) {
	Reg.N = d >> 7
}
