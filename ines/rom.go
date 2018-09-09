package ines

import "github.com/m4ntis/bones/ines/mapper"

const (
	TrainerSize = 512
)

type Trainer [TrainerSize]byte

// ROM represents a whole NES rom, containing the program rom, chr rom and the
// optional trainer.
type ROM struct {
	Header INESHeader

	Trainer Trainer
	Mapper  mapper.Mapper
}
