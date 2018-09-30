package cpu

// OpCode defines an opcode of the mos 6502.
//
// Name is the opcode's textual name in assembly language.
//
// Mode holds the opcode's addressing mode, or the way to address the operands.
//
// Oper is the operation itself, the actual logic being executed.
type OpCode struct {
	Name string

	// cycles contains the base cycle count for the opcode.
	cycles int

	// pageBoundryCheck tells the addressing mode whether a page boundry cross
	// affects it's cycle count.
	pageBoundryCheck bool

	Mode AddressingMode
	Oper Operation
}

// Exec runs the opcode with the given arguments.
//
// It runs it's addressing mode, which in turn fetches operands if necessary and
// calls the operation.
func (op OpCode) Exec(cpu *CPU, ops ...byte) (cycles int, err error) {
	cycles, err = op.Mode.Address(cpu, op.Oper, op.pageBoundryCheck, ops...)
	return cycles + op.cycles, err
}
