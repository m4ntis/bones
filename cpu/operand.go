package cpu

import "fmt"

const (
	ARegisterOperand = iota
)

const (
	CPURAMOperand = iota
	CPURegisterOperand
	ConstOperand
)

type OperandType int

type Operand struct {
	cpu *CPU

	Type  OperandType
	Ident int
}

func NewOperand(cpu *CPU, ot OperandType, identifier int) Operand {
	return Operand{
		cpu: cpu,

		Type:  ot,
		Ident: identifier,
	}
}

func (op Operand) Read() byte {
	switch op.Type {
	case CPURAMOperand:
		return op.cpu.RAM.Read(op.Ident)
	case CPURegisterOperand:
		if op.Ident == ARegisterOperand {
			return cpu.Regs.A
		} else {
			panic(fmt.Sprintf("Invalid cpu register identifier %d",
				op.Ident))
		}
	case ConstOperand:
		return byte(op.Ident)
	default:
		panic(fmt.Sprintf("Invalid operand type %d", op.ot))
	}
}

func (op Operand) Write(d byte) {
	switch op.Type {
	case CPURAMOperand:
		return op.cpu.RAM.Write(op.Ident, d)
	case CPURegisterOperand:
		if op.Ident == a {
			cpu.Regs.A = d
		} else {
			panic(fmt.Sprintf("Invalid cpu register identifier %d",
				op.Ident))
		}
	case ConstOperand:
		panic("Can't write to a const operand")
	default:
		panic(fmt.Sprintf("Invalid operand type %d", op.ot))
	}
}
