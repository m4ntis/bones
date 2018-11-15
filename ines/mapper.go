package ines

import (
	"github.com/pkg/errors"
)

const (
	PrgROMPageSize = 16384 // 16k, 0x4000
	ChrROMPageSize = 8192  // 8K, 0x2000

	PrgRAMSize = 8192 // 8k, 0x2000
	ChrRAMSize = 8192 // 8K, 0x2000
)

type PrgROMPage [PrgROMPageSize]byte
type ChrROMPage [ChrROMPageSize]byte

type Mapper interface {
	Read(addr int) (byte, error)
	Write(addr int, d byte) error
	Observe(addr int) (byte, error)

	Populate([]PrgROMPage, []ChrROMPage)
	GetPRGRom() []PrgROMPage

	SetSram(bool)
}

func NewMapper(num int) (Mapper, error) {
	m, ok := mappers[num]
	if !ok {
		return nil, errors.Errorf("iNes Mapper %d not yet implemented", num)
	}

	return m, nil
}

var mappers = map[int]Mapper{
	0: &Mapper000{},
	1: &Mapper001{},
}
