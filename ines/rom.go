// Package ines contains some common types that are used throughout the
// several packages, and don't contain any logic.
package ines

const (
	TrainerSize    = 512
	PrgROMPageSize = 16384
	ChrROMPageSize = 8192
)

type Trainer [TrainerSize]byte
type PrgROMPage [PrgROMPageSize]byte
type ChrROMPage [ChrROMPageSize]byte

// ROM represents a whole NES rom, containing the program rom, chr rom and the
// optional trainer
type ROM struct {
	Trainer Trainer
	PrgROM  []PrgROMPage
	ChrROM  []ChrROMPage
}
