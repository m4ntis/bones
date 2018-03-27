package models

const (
	TrainerSize    = 512
	PrgROMPageSize = 16384
	ChrROMPageSize = 8192
)

type Trainer [TrainerSize]byte
type PrgROMPage [PrgROMPageSize]byte
type ChrROMPage [ChrROMPageSize]byte

type ROM struct {
	Trainer Trainer
	PrgROM  []PrgROMPage
	ChrROM  []ChrROMPage
}

func NewROM(trainer Trainer, prgROM []PrgROMPage, chrROM []ChrROMPage) *ROM {
	return &ROM{Trainer: trainer, PrgROM: prgROM, ChrROM: chrROM}
}
