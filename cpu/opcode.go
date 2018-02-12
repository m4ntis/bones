package cpu

// OpCode defines an opcode of the 2a03.
//
// It has it's textual representation, it's addressing mode and the operation
// itself
type OpCode struct {
	name   string
	argLen int
	cycles int

	mode AddressingMode
	oper Operation
}

// Exec runs the opcode with the given arguments.
//
// It runs it's addressing mode, which in turn fetches the arguments and calls
// the operation
func (op OpCode) Exec(cpu *CPU, args []byte) {
	op.mode(cpu, op.oper, args)
}
