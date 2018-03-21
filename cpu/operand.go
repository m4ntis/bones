package cpu

import "fmt"

const (
	a = iota
)

const (
	CPURAMOperand = iota
	CPURegisterOperand
)

type OperandType int

type Operand struct {
	cpu *CPU

	ot    OperandType
	ident int
}

func NewOperand(cpu *CPU, ot OperandType, identifier int) Operand {
	return Operand{
		cpu: cpu,

		ot:    ot,
		ident: identifier,
	}
}

func (op Operand) Read() byte {
	switch op.Type {
	case CPURAMOperand:
		return op.cpu.RAM.Read(op.Identifier)
	case CPURegisterOperand:
		if op.Identifier == a {
			return cpu.Regs.A
		} else {
			panic(fmt.Sprintf("Invalid cpu register identifier %d",
				op.Identifier))
		}
	default:
		panic(fmt.Sprintf("Invalid operand type %d", op.ot))
	}
}

func (op Operand) Write(d byte) {
	switch op.Type {
	case CPURAMOperand:
		return op.cpu.RAM.Write(op.Identifier, d)
	case CPURegisterOperand:
		if op.Identifier == a {
			cpu.Regs.A = d
		} else {
			panic(fmt.Sprintf("Invalid cpu register identifier %d",
				op.Identifier))
		}
	default:
		panic(fmt.Sprintf("Invalid operand type %d", op.ot))
	}
}
