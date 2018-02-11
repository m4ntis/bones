package models

const (
	TRAINER_SIZE      = 512
	PRG_ROM_PAGE_SIZE = 16384
	CHR_ROM_PAGE_SIZE = 8192
)

type Trainer [TRAINER_SIZE]byte
type PrgROMPage [PRG_ROM_PAGE_SIZE]byte
type ChrROMPage [CHR_ROM_PAGE_SIZE]byte

type ROM struct {
	Trainer Trainer
	PrgROM  []PrgROMPage
	ChrROM  []ChrROMPage
}

func NewROM(trainer Trainer, prgROM []PrgROMPage, chrROM []ChrROMPage) *ROM {
	return &ROM{Trainer: trainer, PrgROM: prgROM, ChrROM: chrROM}
}
