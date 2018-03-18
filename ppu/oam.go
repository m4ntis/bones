package ppu

const (
	OAM_SIZE = 0x100
)

type OAM struct {
	data [OAM_SIZE]byte
}
