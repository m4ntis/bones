package cpu

const (
	Set   = 1
	Clear = 0
)

type Registers struct {
	PC int
	SP byte

	A byte
	X byte
	Y byte

	C byte
	Z byte
	I byte
	D byte
	B byte
	V byte
	N byte
}

func (reg *Registers) GetP() byte {
	return reg.C | reg.Z<<1 | reg.I<<2 | reg.D<<3 | reg.B<<4 | reg.V<<6 |
		reg.N<<7
}

func (reg *Registers) SetP(p byte) {
	reg.C = p & 1
	reg.Z = p & 2
	reg.I = p & 4
	reg.D = p & 8
	reg.B = p & 16
	reg.V = p & 64
	reg.N = p & 128
}
