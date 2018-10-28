package ppu

const (
	oamSize          = 0x100
	secondaryOamSize = 0x20

	sprDataSize = 0x4
)

type OAM [oamSize]byte
type secondaryOAM [secondaryOamSize]byte
