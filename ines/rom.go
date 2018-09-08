package ines

const (
	TrainerSize    = 512
	PrgROMPageSize = 16384 // 16k, 0x4000
	ChrROMPageSize = 8192  // 8K, 0x2000
)

type Trainer [TrainerSize]byte
type PrgROMPage [PrgROMPageSize]byte
type ChrROMPage [ChrROMPageSize]byte

// ROM represents a whole NES rom, containing the program rom, chr rom and the
// optional trainer.
type ROM struct {
	Header INESHeader

	Trainer Trainer
	Mapper  Mapper
}
