package cpu

type Operand interface {
	Read() byte
	Write(byte) int
}

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

type ConstOperand struct {
	D byte
}

func (op ConstOperand) Read() byte {
	return op.D
}

func (op ConstOperand) Write(d byte) (cycles int) {
	panic("Can't write const operand")
}

type NilOperand struct{}

func (op NilOperand) Read() byte {
	panic("Can't read nil operand")
}

func (op NilOperand) Write(d byte) (cycles int) {
	panic("Can't write nil operand")
}
