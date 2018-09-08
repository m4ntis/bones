package mapper

import (
	"github.com/pkg/errors"
)

const (
	PrgROMPageSize = 16384 // 16k, 0x4000
	ChrROMPageSize = 8192  // 8K, 0x2000
)

type PrgROMPage [PrgROMPageSize]byte
type ChrROMPage [ChrROMPageSize]byte

type Mapper interface {
	Read(addr int) byte
	Write(addr int, d byte) int
	Observe(addr int) byte

	Populate([]PrgROMPage, []ChrROMPage)
}

func New(num int) (mapper Mapper, err error) {
	mapper, ok := mappers[num]
	if !ok {
		return nil, errors.Errorf("iNes Mapper %d not implemented yet", num)
	}

	return mapper, nil
}

var mappers = map[int]Mapper{}
