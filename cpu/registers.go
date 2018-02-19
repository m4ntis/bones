package cpu

const (
	SET   = 1
	CLEAR = 0
)

type Registers struct {
	pc int
	sp byte

	a byte
	x byte
	y byte

	c byte
	z byte
	i byte
	d byte
	b byte
	v byte
	n byte
}

func (reg *Registers) getP() byte {
	return reg.c | reg.z<<1 | reg.i<<2 | reg.d<<3 | reg.b<<4 | reg.v<<6 |
		reg.n<<7
}

func (reg *Registers) setP(p byte) {
	reg.c = p & 1
	reg.z = p & 2
	reg.i = p & 4
	reg.d = p & 8
	reg.b = p & 16
	reg.v = p & 64
	reg.n = p & 128
}
