package cpu

const (
	SET   = 1
	CLEAR = 0
)

type Registers struct {
	pc int
	sp int

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
