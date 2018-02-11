package cpu

// Operation defines an operation that the CPU executes in one or more of it's
// opcodes.
//
// The byte values it received are it's operands
type Operation func(...byte)

// OpCode defines an opcode of the 2a03.
//
// It has it's textual representation, it's addressing mode and the operation
// itself
type OpCode struct {
	name string

	mode AddressingMode
	op   Operation
}

// Exec runs the opcode with the given arguments.
//
// It runs it's addressing mode, which in turn fetches the operands and calls
// the operation
func (op OpCode) Exec(args []byte) {
	mode(op, args)
}
