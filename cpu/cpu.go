package cpu

type CPU struct {
	ram *RAM

	pc int
	sp int

	a byte
	x byte
	y byte

	c int
	z int
	i int
	d int
	b int
	v int
	n int
}
