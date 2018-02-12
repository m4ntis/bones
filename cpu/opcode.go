package cpu

// Operation defines an operation that the CPU executes in one or more of it's
// opcodes.
//
// The byte values it received are it's arguments. Both arguments and return
// value might be nil, depends on the operation. There isn't a gurantee that
// the operation will check for the correct number of arguments, so make sure
// you pass in the correct amount.
type Operation func(...byte) interface{}

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
