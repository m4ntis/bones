package mapper

import (
	"github.com/pkg/errors"
)

type Mapper interface {
	Read(addr int) byte
	Populate([]PrgROMPage, []ChrROMPage)
}

func New(num int) (mapper Mapper, err error) {
	mapper, ok := mappers[num]
	if !ok {
		return nil, errors.Errorf("iNes Mapper %d not implemented yet", num)
	}
}

var mappers = map[int]Mapper{}
