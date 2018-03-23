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
	return op.ram.Read(op.addr)
}

func (op RAMOperand) Write(d byte) {
	op.ram.Write(d)
}

type RegOperand struct {
	Reg *byte
}

func (op RegOperand) Read() byte {
	return *op.reg
}

func (op RegOperand) Write(d byte) {
	*op.reg = d
}

type ConstOperand struct {
	D byte
}

func (op ConstOperand) Read() byte {
	return op.d
}

func (op ConstOperand) Write(d byte) {
	panic("Can't write to a const operand")
}
