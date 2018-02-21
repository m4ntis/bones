package cpu

// OpCode defines an opcode of the 2a03.
//
// Contains it's textual representation, addressing mode and the operation
// itself.
//
// The opcode also contains some informaition on the amount of cycles it takes
// to execute. cycles is the base cycle count, and checkPageBoundry tells the
// addressing mode whether a page boundry cross affects it's cycle count
type OpCode struct {
	name string

	cycles           int
	checkPageBoundry bool

	mode AddressingMode
	oper Operation
}

// Exec runs the opcode with the given arguments.
//
// It runs it's addressing mode, which in turn fetches the arguments and calls
// the operation
func (op OpCode) Exec(cpu *CPU, args ...*byte) (extraCycles int) {
	return op.mode(cpu, op.oper, args...)
}
