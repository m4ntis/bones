package ppu

const (
	OamSize          = 0x100
	secondaryOamSize = 0x20
)

type OAM [OamSize]byte
type secondaryOAM [secondaryOamSize]byte
