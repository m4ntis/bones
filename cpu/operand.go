package cpu

type Operand interface {
	Read() byte
	Write(byte)
}

type RAMOperand struct {
	RAM  *RAM
	Addr int
}

func (op RAMOperand) Read() byte {
	return op.RAM.Read(op.Addr)
}

func (op RAMOperand) Write(d byte) {
	op.RAM.Write(op.Addr, d)
}

type RegOperand struct {
	Reg *byte
}

func (op RegOperand) Read() byte {
	return *op.Reg
}

func (op RegOperand) Write(d byte) {
	*op.Reg = d
}

type ConstOperand struct {
	D byte
}

func (op ConstOperand) Read() byte {
	return op.D
}

func (op ConstOperand) Write(d byte) {
	panic("Can't write to a const operand")
}
