package cpu

// Operation defines an operation that the CPU executes in one or more of its
// opcodes.
//
// The operation receives a reference to the cpu in order to read or set its ram
// or registers' state.
//
// It also receives a single operand, given to it by the addressing mode that
// calls it. This abstraction enables to define the opcode's logic once,
// separating it from the way the operand is fetched, and leaving that logic to
// the addressing mode.
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
		cpu.Reg.V = set
	} else {
		cpu.Reg.V = clear
	}

	// Carry
	if int(arg1)+int(arg2)+int(arg3) > 255 {
		cpu.Reg.C = set
	} else {
		cpu.Reg.C = clear
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

	if cpu.Reg.C == clear {
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

	if cpu.Reg.C == set {
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

	if cpu.Reg.Z == set {
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
		cpu.Reg.Z = set
		return
	}
	cpu.Reg.Z = clear
	return
}

func BMI(cpu *CPU, op Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.N == set {
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

	if cpu.Reg.Z == clear {
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

	if cpu.Reg.N == clear {
		extraCycles++
		cpu.Reg.PC += int(int8(op.Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func BRK(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.PC++
	cpu.IRQ()
	return
}

func BVC(cpu *CPU, op Operand) (extraCycles int) {
	initPC := cpu.Reg.PC

	if cpu.Reg.V == clear {
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

	if cpu.Reg.V == set {
		extraCycles++
		cpu.Reg.PC += int(int8(op.Read()))

		if initPC/256 != cpu.Reg.PC/256 {
			extraCycles++
		}
	}
	return
}

func CLC(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.C = clear
	return
}

func CLD(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.D = clear
	return
}

func CLI(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.I = clear
	return
}

func CLV(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.V = clear
	return
}

func CMP(cpu *CPU, op Operand) (extraCycles int) {
	d := op.Read()

	res := cpu.Reg.A - d

	setNZ(cpu, res)
	if d <= cpu.Reg.A {
		cpu.Reg.C = set
	} else {
		cpu.Reg.C = clear
	}
	return
}

func CPX(cpu *CPU, op Operand) (extraCycles int) {
	d := op.Read()

	res := cpu.Reg.X - d

	setNZ(cpu, res)
	if d <= cpu.Reg.X {
		cpu.Reg.C = set
	} else {
		cpu.Reg.C = clear
	}
	return
}

func CPY(cpu *CPU, op Operand) (extraCycles int) {
	d := op.Read()

	res := cpu.Reg.Y - d

	setNZ(cpu, res)
	if d <= cpu.Reg.Y {
		cpu.Reg.C = set
	} else {
		cpu.Reg.C = clear
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
	cpu.Reg.I = clear
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
	arg1 := cpu.Reg.A
	arg2 := op.Read()
	arg3 := byte(0)
	if cpu.Reg.C == set {
		arg3 = 1
	}

	cpu.Reg.A = arg1 - arg2 - (1 - arg3)

	// Set flags
	setNZ(cpu, cpu.Reg.A)

	if int(arg1)-int(arg2)-int(1-arg3) >= 0 {
		cpu.Reg.C = set
	} else {
		cpu.Reg.C = clear
	}

	// Overflow
	if (arg1^arg2)&0x80 != 0 && (arg1^cpu.Reg.A)&0x80 != 0 {
		cpu.Reg.V = set
	} else {
		cpu.Reg.V = clear
	}
	return
}

func SEC(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.C = set
	return
}

func SED(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.D = set
	return
}

func SEI(cpu *CPU, op Operand) (extraCycles int) {
	cpu.Reg.I = set
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
	return
}

func setNZ(cpu *CPU, d byte) {
	cpu.Reg.N = d >> 7

	if d == 0x0 {
		cpu.Reg.Z = set
		return
	}
	cpu.Reg.Z = clear
}
