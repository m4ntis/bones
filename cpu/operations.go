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
type Operation func(*CPU, Operand) int

func ADC(cpu *CPU, op Operand) (extraCycles int) {
	// Calculate result and store in a
	arg1 := cpu.Reg.A
	arg2 := op.Read()
	arg3 := cpu.Reg.C

	res := arg1 + arg2 + arg3
	cpu.Reg.A = res

	// Set flags
	setNZ(cpu, res)

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

func AND(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.A &= op.Read()
	setNZ(cpu, cpu.Reg.A)
	return
}

func ASL(cpu *CPU, op Operand) (extraCycles int) {
	d := op.Read()

	cpu.Reg.C = d >> 7
	d <<= 1

	setNZ(cpu, d)

	extraCycles += op.Write(d)
	return
}

func BCC(cpu *CPU, op Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.C == Clear {
		extraCycles++
		cpu.Reg.PC += int(int8(op.Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BCS(cpu *CPU, op Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.C == Set {
		extraCycles++
		cpu.Reg.PC += int(int8(op.Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BEQ(cpu *CPU, op Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.Z == Set {
		extraCycles++
		cpu.Reg.PC += int(int8(op.Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BIT(cpu *CPU, op Operand) (extraCycles int) {
	d := op.Read()

	cpu.Reg.N = d >> 7
	cpu.Reg.V = (d >> 6) & 1

	res := cpu.Reg.A & d
	if res == 0x0 {
		cpu.Reg.Z = Set
		return
	}
	cpu.Reg.Z = Clear
	return
}

func BMI(cpu *CPU, op Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.N == Set {
		extraCycles++
		cpu.Reg.PC += int(int8(op.Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BNE(cpu *CPU, op Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.Z == Clear {
		extraCycles++
		cpu.Reg.PC += int(int8(op.Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BPL(cpu *CPU, op Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.N == Clear {
		extraCycles++
		cpu.Reg.PC += int(int8(op.Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BRK(cpu *CPU, op Operand) (extraCycles int) {
	cpu.IRQ()
	return
}

func BVC(cpu *CPU, op Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.V == Clear {
		extraCycles++
		cpu.Reg.PC += int(int8(op.Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BVS(cpu *CPU, op Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.V == Set {
		extraCycles++
		cpu.Reg.PC += int(int8(op.Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func CLC(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.C = Clear
	return
}

func CLD(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.D = Clear
	return
}

func CLI(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.I = Clear
	return
}

func CLV(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.V = Clear
	return
}

func CMP(cpu *CPU, op Operand) (extraCycles int) {
	d := op.Read()

	res := cpu.Reg.A - d

	setNZ(cpu, res)
	if d > cpu.Reg.A {
		cpu.Reg.C = Set
	} else {
		cpu.Reg.C = Clear
	}
	return
}

func CPX(cpu *CPU, op Operand) (extraCycles int) {
	d := op.Read()

	res := cpu.Reg.X - d

	setNZ(cpu, res)
	if d > cpu.Reg.A {
		cpu.Reg.C = Set
	} else {
		cpu.Reg.C = Clear
	}
	return
}

func CPY(cpu *CPU, op Operand) (extraCycles int) {
	d := op.Read()

	res := cpu.Reg.Y - d

	setNZ(cpu, res)
	if d > cpu.Reg.A {
		cpu.Reg.C = Set
	} else {
		cpu.Reg.C = Clear
	}
	return
}

func DEC(cpu *CPU, op Operand) (extraCycles int) {
	d := op.Read() - 1

	setNZ(cpu, d)

	extraCycles += op.Write(d)
	return
}

func DEX(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.X--
	setNZ(cpu, cpu.Reg.X)
	return
}

func DEY(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.Y--
	setNZ(cpu, cpu.Reg.Y)
	return
}

func EOR(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.A ^= op.Read()
	setNZ(cpu, cpu.Reg.A)
	return
}

func INC(cpu *CPU, op Operand) (extraCycles int) {
	d := op.Read() + 1
	setNZ(cpu, d)
	extraCycles += op.Write(d)
	return
}

func INX(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.X++
	setNZ(cpu, cpu.Reg.X)
	return
}

func INY(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.Y++
	setNZ(cpu, cpu.Reg.Y)
	return
}

func JMP(cpu *CPU, op Operand) (extraCycles int) {
	jmpPC := op.(RAMOperand).Addr
	cpu.Reg.PC = jmpPC
	return
}

func JSR(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.PC -= 1
	// push PCH
	cpu.push(byte(cpu.Reg.PC >> 8))
	// push PCL
	cpu.push(byte(cpu.Reg.PC & 0xff))

	jmpPC := op.(RAMOperand).Addr
	cpu.Reg.PC = jmpPC
	return
}

func LDA(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.A = op.Read()
	setNZ(cpu, cpu.Reg.A)
	return
}

func LDX(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.X = op.Read()
	setNZ(cpu, cpu.Reg.X)
	return
}

func LDY(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.Y = op.Read()
	setNZ(cpu, cpu.Reg.Y)
	return
}

func LSR(cpu *CPU, op Operand) (extraCycles int) {
	d := op.Read()

	cpu.Reg.C = d & 1
	d >>= 1

	setNZ(cpu, d)

	extraCycles += op.Write(d)
	return
}

func NOP(cpu *CPU, op Operand) (extraCycles int) {
	return
}

func ORA(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.A |= op.Read()
	setNZ(cpu, cpu.Reg.A)
	return
}

func PHA(cpu *CPU, op Operand) (extraCycles int) {
	cpu.push(cpu.Reg.A)
	return
}

func PHP(cpu *CPU, op Operand) (extraCycles int) {
	cpu.push(cpu.Reg.GetP())
	return
}

func PLA(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.A = cpu.pull()
	setNZ(cpu, cpu.Reg.A)
	return
}

func PLP(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.SetP(cpu.pull())
	return
}

func ROL(cpu *CPU, op Operand) (extraCycles int) {
	d := op.Read()

	carry := cpu.Reg.C
	cpu.Reg.C = d >> 7

	d <<= 1
	d |= carry

	setNZ(cpu, d)

	extraCycles += op.Write(d)
	return
}

func ROR(cpu *CPU, op Operand) (extraCycles int) {
	d := op.Read()

	carry := cpu.Reg.C
	cpu.Reg.C = d & 1

	d >>= 1
	d |= (carry << 7)

	setNZ(cpu, d)

	extraCycles += op.Write(d)
	return
}

func RTI(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.SetP(cpu.pull())
	// pull PCL and then PHC
	cpu.Reg.PC = int(cpu.pull()) | int(cpu.pull())<<8
	return
}

func RTS(cpu *CPU, op Operand) (extraCycles int) {
	// pull PCL and then PHC
	cpu.Reg.PC = int(cpu.pull()) | int(cpu.pull())<<8
	cpu.Reg.PC++
	return
}

func SBC(cpu *CPU, op Operand) (extraCycles int) {
	// Calculate result and store in a
	arg1 := int8(cpu.Reg.A)
	arg2 := int8(op.Read())

	res := byte(arg1 - arg2)
	cpu.Reg.A = res

	// Set flags
	setNZ(cpu, res)

	cpu.Reg.C = (res >> 7) ^ 1

	// Overflow
	if math.Abs(float64(arg1)-float64(arg2)) > 127 {
		cpu.Reg.V = Set
	} else {
		cpu.Reg.V = Clear
	}
	return
}

func SEC(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.C = Set
	return
}

func SED(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.D = Set
	return
}

func SEI(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.I = Set
	return
}

func STA(cpu *CPU, op Operand) (extraCycles int) {
	extraCycles += op.Write(cpu.Reg.A)
	return
}

func STX(cpu *CPU, op Operand) (extraCycles int) {
	extraCycles += op.Write(cpu.Reg.X)
	return
}

func STY(cpu *CPU, op Operand) (extraCycles int) {
	extraCycles += op.Write(cpu.Reg.Y)
	return
}

func TAX(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.X = cpu.Reg.A
	setNZ(cpu, cpu.Reg.X)
	return
}

func TAY(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.Y = cpu.Reg.A
	setNZ(cpu, cpu.Reg.Y)
	return
}

func TSX(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.X = cpu.Reg.SP
	setNZ(cpu, cpu.Reg.X)
	return
}

func TXA(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.A = cpu.Reg.X
	setNZ(cpu, cpu.Reg.A)
	return
}

func TYA(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.A = cpu.Reg.Y
	setNZ(cpu, cpu.Reg.A)
	return
}

func TXS(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.SP = cpu.Reg.X
	setNZ(cpu, cpu.Reg.SP)
	return
}

func setNZ(cpu *CPU, d byte) {
	cpu.Reg.N = d >> 7

	if d == 0x0 {
		cpu.Reg.Z = Set
		return
	}
	cpu.Reg.Z = Clear
}
