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
type Operation func(*Registers, Operand) int

func ADC(reg *Registers, op Operand) (extraCycles int) {
	// Calculate result and store in a
	arg1 := reg.A
	arg2 := op.Read()
	arg3 := reg.C

	res := arg1 + arg2 + arg3
	reg.A = res

	// Set flags
	setNZ(reg, res)

	// Overflow
	signed_arg1 := int8(arg1)
	signed_arg2 := int8(arg2)
	signed_arg3 := int8(arg3)
	if int(signed_arg1+signed_arg2+signed_arg3) != int(signed_arg1)+
		int(signed_arg2)+int(signed_arg3) {
		reg.V = Set
	} else {
		reg.V = Clear
	}

	// Carry
	if int(arg1)+int(arg2)+int(arg3) > 255 {
		reg.C = Set
	} else {
		reg.C = Clear
	}
	return
}

func AND(reg *Registers, op Operand) (extraCycles int) {
	reg.A &= op.Read()
	setNZ(reg, reg.A)
	return
}

func ASL(reg *Registers, op Operand) (extraCycles int) {
	d := op.Read()

	reg.C = d >> 7
	d <<= 1

	setNZ(reg, d)

	op.Write(d)
	return
}

func BCC(reg *Registers, op Operand) (extraCycles int) {
	initPC := reg.PC

	if reg.C == Clear {
		extraCycles++
		reg.PC += int(int8(op.Read()))

		if initPC/256 != reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BCS(reg *Registers, op Operand) (extraCycles int) {
	initPC := reg.PC

	if reg.C == Set {
		extraCycles++
		reg.PC += int(int8(op.Read()))

		if initPC/256 != reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BEQ(reg *Registers, op Operand) (extraCycles int) {
	initPC := reg.PC

	if reg.Z == Set {
		extraCycles++
		reg.PC += int(int8(op.Read()))

		if initPC/256 != reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BIT(reg *Registers, op Operand) (extraCycles int) {
	d := op.Read()

	reg.N = d >> 7
	reg.V = (d >> 6) & 1

	res := reg.A & d
	if res == 0x0 {
		reg.Z = Set
		return
	}
	reg.Z = Clear
	return
}

func BMI(reg *Registers, op Operand) (extraCycles int) {
	initPC := reg.PC

	if reg.N == Set {
		extraCycles++
		reg.PC += int(int8(op.Read()))

		if initPC/256 != reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BNE(reg *Registers, op Operand) (extraCycles int) {
	initPC := reg.PC

	if reg.Z == Clear {
		extraCycles++
		reg.PC += int(int8(op.Read()))

		if initPC/256 != reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BPL(reg *Registers, op Operand) (extraCycles int) {
	initPC := reg.PC

	if reg.N == Clear {
		extraCycles++
		reg.PC += int(int8(op.Read()))

		if initPC/256 != reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BRK(reg *Registers, op Operand) (extraCycles int) {
	cpu.IRQ()
	return
}

func BVC(reg *Registers, op Operand) (extraCycles int) {
	initPC := reg.PC

	if reg.V == Clear {
		extraCycles++
		reg.PC += int(int8(op.Read()))

		if initPC/256 != reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BVS(reg *Registers, op Operand) (extraCycles int) {
	initPC := reg.PC

	if reg.V == Set {
		extraCycles++
		reg.PC += int(int8(op.Read()))

		if initPC/256 != reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func CLC(reg *Registers, op Operand) (extraCycles int) {
	reg.C = Clear
	return
}

func CLD(reg *Registers, op Operand) (extraCycles int) {
	reg.D = Clear
	return
}

func CLI(reg *Registers, op Operand) (extraCycles int) {
	reg.I = Clear
	return
}

func CLV(reg *Registers, op Operand) (extraCycles int) {
	reg.V = Clear
	return
}

func CMP(reg *Registers, op Operand) (extraCycles int) {
	d := op.Read()

	res := reg.A - d

	setNZ(reg, res)
	if d > reg.A {
		reg.C = Set
	} else {
		reg.C = Clear
	}
	return
}

func CPX(reg *Registers, op Operand) (extraCycles int) {
	d := op.Read()

	res := reg.X - d

	setNZ(reg, res)
	if d > reg.A {
		reg.C = Set
	} else {
		reg.C = Clear
	}
	return
}

func CPY(reg *Registers, op Operand) (extraCycles int) {
	d := op.Read()

	res := reg.Y - d

	setNZ(reg, res)
	if d > reg.A {
		reg.C = Set
	} else {
		reg.C = Clear
	}
	return
}

func DEC(reg *Registers, op Operand) (extraCycles int) {
	d := op.Read() - 1

	setNZ(reg, d)

	op.Write(d)
	return
}

func DEX(reg *Registers, op Operand) (extraCycles int) {
	reg.X--
	setNZ(reg, reg.X)
	return
}

func DEY(reg *Registers, op Operand) (extraCycles int) {
	reg.Y--
	setNZ(reg, reg.Y)
	return
}

func EOR(reg *Registers, op Operand) (extraCycles int) {
	reg.A ^= op.Read()
	setNZ(reg, reg.A)
	return
}

func INC(reg *Registers, op Operand) (extraCycles int) {
	d := op.Read() + 1
	setNZ(reg, d)
	op.Write(d)
	return
}

func INX(reg *Registers, op Operand) (extraCycles int) {
	reg.X++
	setNZ(reg, reg.X)
	return
}

func INY(reg *Registers, op Operand) (extraCycles int) {
	reg.Y++
	setNZ(reg, reg.Y)
	return
}

func JMP(reg *Registers, op Operand) (extraCycles int) {
	jmpPC := op.Ident
	reg.PC = jmpPC
	return
}

func JSR(reg *Registers, op Operand) (extraCycles int) {
	reg.PC += 2
	// push PCH
	cpu.push(byte(reg.PC >> 8))
	// push PCL
	cpu.push(byte(reg.PC & 0xff))

	jmpPC := op.Ident
	reg.PC = jmpPC
	return
}

func LDA(reg *Registers, op Operand) (extraCycles int) {
	reg.A = op.Read()
	setNZ(reg, reg.A)
	return
}

func LDX(reg *Registers, op Operand) (extraCycles int) {
	reg.X = op.Read()
	setNZ(reg, reg.X)
	return
}

func LDY(reg *Registers, op Operand) (extraCycles int) {
	reg.Y = op.Read()
	setNZ(reg, reg.Y)
	return
}

func LSR(reg *Registers, op Operand) (extraCycles int) {
	d := op.Read()

	reg.C = d & 1
	d >>= 1

	setNZ(reg, d)

	op.Write(d)
	return
}

func NOP(reg *Registers, op Operand) (extraCycles int) {
	return
}

func ORA(reg *Registers, op Operand) (extraCycles int) {
	reg.A |= op.Read()
	setNZ(reg, reg.A)
	return
}

func PHA(reg *Registers, op Operand) (extraCycles int) {
	cpu.push(reg.A)
	return
}

func PHP(reg *Registers, op Operand) (extraCycles int) {
	cpu.push(reg.GetP())
	return
}

func PLA(reg *Registers, op Operand) (extraCycles int) {
	reg.A = cpu.pull()
	setNZ(reg, reg.A)
	return
}

func PLP(reg *Registers, op Operand) (extraCycles int) {
	reg.SetP(cpu.pull())
	return
}

func ROL(reg *Registers, op Operand) (extraCycles int) {
	d := op.Read()

	carry := reg.C
	reg.C = d >> 7

	d <<= 1
	d |= carry

	setNZ(reg, d)

	op.Write(d)
	return
}

func ROR(reg *Registers, op Operand) (extraCycles int) {
	d := op.Read()

	carry := reg.C
	reg.C = d & 1

	d >>= 1
	d |= (carry << 7)

	setNZ(reg, d)

	op.Write(d)
	return
}

func RTI(reg *Registers, op Operand) (extraCycles int) {
	reg.SetP(cpu.pull())
	// pull PCL and then PHC
	reg.PC = int(cpu.pull()) | int(cpu.pull())<<8
	return
}

func RTS(reg *Registers, op Operand) (extraCycles int) {
	// pull PCL and then PHC
	reg.PC = int(cpu.pull()) | int(cpu.pull())<<8
	reg.PC++
	return
}

func SBC(reg *Registers, op Operand) (extraCycles int) {
	// Calculate result and store in a
	arg1 := int8(reg.A)
	arg2 := int8(op.Read())

	res := byte(arg1 - arg2)
	reg.A = res

	// Set flags
	setNZ(reg, res)

	reg.C = (res >> 7) ^ 1

	// Overflow
	if math.Abs(float64(arg1)-float64(arg2)) > 127 {
		reg.V = Set
	} else {
		reg.V = Clear
	}
	return
}

func SEC(reg *Registers, op Operand) (extraCycles int) {
	reg.C = Set
	return
}

func SED(reg *Registers, op Operand) (extraCycles int) {
	reg.D = Set
	return
}

func SEI(reg *Registers, op Operand) (extraCycles int) {
	reg.I = Set
	return
}

func STA(reg *Registers, op Operand) (extraCycles int) {
	op.Write(reg.A)
	return
}

func STX(reg *Registers, op Operand) (extraCycles int) {
	op.Write(reg.X)
	return
}

func STY(reg *Registers, op Operand) (extraCycles int) {
	op.Write(reg.Y)
	return
}

func TAX(reg *Registers, op Operand) (extraCycles int) {
	reg.X = reg.A
	setNZ(reg, reg.X)
	return
}

func TAY(reg *Registers, op Operand) (extraCycles int) {
	reg.Y = reg.A
	setNZ(reg, reg.Y)
	return
}

func TSX(reg *Registers, op Operand) (extraCycles int) {
	reg.X = reg.SP
	setNZ(reg, reg.X)
	return
}

func TXA(reg *Registers, op Operand) (extraCycles int) {
	reg.A = reg.X
	setNZ(reg, reg.A)
	return
}

func TYA(reg *Registers, op Operand) (extraCycles int) {
	reg.A = reg.Y
	setNZ(reg, reg.A)
	return
}

func TXS(reg *Registers, op Operand) (extraCycles int) {
	reg.SP = reg.X
	setNZ(reg, reg.SP)
	return
}

func setNZ(reg *Registers, d byte) {
	if d == 0x0 {
		reg.Z = Set
		return
	}
	reg.Z = Clear
	reg.N = d >> 7
}
