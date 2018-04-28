package cpu

// Operand defines the object passed to an operation by the addressing mode.
//
// Read returns the byte value of the operand.
//
// Write writes a given byte value to the operand, returning the CPU cycles the
// operation took. This is only applicable when the write is to the OAMDMA PPU
// register, which starts a DMA and blocks the cpu for 513/514 cycles.
type Operand interface {
	Read() byte
	Write(byte) int
}

// RAMOperand is an Operand that reads and writes to RAM at a fixed location.
//
// Addr is the address that will be accessed in RAM when reading/writing to this
// operand. Addr should be populated by the calling addressing mode.
//
// RAMOperand is the most common, differing between the addressing modes only in
// the way that Addr is calculated.
type RAMOperand struct {
	RAM  *RAM
	Addr int
}

func (op RAMOperand) Read() byte {
	return op.RAM.Read(op.Addr)
}

func (op RAMOperand) Write(d byte) (cycles int) {
	return op.RAM.Write(op.Addr, d)
}

// RegOperand is an Operand containing a reference to a single CPU register,
// reading and writing to it.
type RegOperand struct {
	Reg *byte
}

func (op RegOperand) Read() byte {
	return *op.Reg
}

func (op RegOperand) Write(d byte) (cycles int) {
	*op.Reg = d
	return
}

// ConstOperand is an operand that represents a literal value passed to the
// opcode. Reading this operand will return its const byte value and writing
// doesn't affect it.
type ConstOperand struct {
	D byte
}

func (op ConstOperand) Read() byte {
	return op.D
}

func (op ConstOperand) Write(d byte) (cycles int) {
	// Writing to a const operand has no logical meaning
	return
}

// NilOperand is an operand used for implied addressing mode, where no operand
// is passed to the operation. Reading or writing to this operand does nothing.
type NilOperand struct{}

func (op NilOperand) Read() byte {
	// Reading a nil operand has no logical meaning
	return 0
}

func (op NilOperand) Write(d byte) (cycles int) {
	// Writing to a nil operand has no logical meaning
	return
}
