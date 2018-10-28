package cpu

const (
	set   = 1
	clear = 0
)

// Registers is a simple struct containing all the CPU's registers.
//
// The P register is separated into it's different bits for ease of accessing.
// Reading from or writing to the P register at whole should instead be done
// using the Get and Set opreations.
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

// GetP returns the value of the P register, calculated from all the status
// bit registers.
func (reg *Registers) GetP() byte {
	return reg.C | reg.Z<<1 | reg.I<<2 | reg.D<<3 | reg.B<<4 | reg.V<<6 |
		reg.N<<7
}

// SetP sets a value to the P register, translated to all the separate status
// bits it contains.
func (reg *Registers) SetP(p byte) {
	reg.C = p & 1
	reg.Z = (p >> 1) & 1
	reg.I = (p >> 2) & 1
	reg.D = (p >> 3) & 1
	reg.B = (p >> 4) & 1
	reg.V = (p >> 6) & 1
	reg.N = (p >> 7) & 1
}
