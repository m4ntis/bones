package ppu

const (
	OamSize          = 0x100
	SecondaryOamSize = 0x20
)

type OAM [OamSize]byte
type SecondaryOAM [SecondaryOamSize]byte
