package cpu

// OpCode defines an opcode of the 2a03.
//
// Contains it's textual representation, addressing mode and the operation
// itself.
//
// The opcode also contains some informaition on the amount of cycles it takes
// to execute. cycles is the base cycle count, and pageBoundryCheck tells the
// addressing mode whether a page boundry cross affects it's cycle count
type OpCode struct {
	Name string

	cycles           int
	pageBoundryCheck bool

	Mode AddressingMode
	Oper Operation
}

// Exec runs the opcode with the given arguments.
//
// It runs it's addressing mode, which in turn fetches the arguments and calls
// the operation
func (op OpCode) Exec(cpu *CPU, ops ...byte) (extraCycles int) {
	return op.Mode.address(cpu, op.Oper, op.pageBoundryCheck, ops...)
}
