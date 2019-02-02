package ines

import (
	"github.com/pkg/errors"
)

type Mapper000 struct {
	prgROM []PrgROMPage
	sRAM   [SRAMSize]byte

	chrROM    []ChrROMPage
	chrRAM    [ChrRAMSize]byte
	useChrRAM bool
}

func (m *Mapper000) Read(addr int) (d byte, err error) {
	switch {
	case addr < 0x2000:
		if m.useChrRAM {
			return m.chrRAM[addr], nil
		}
		return m.readChrROM(addr), nil

	case addr >= 0x8000:
		return m.readPrgROM(addr - 0x8000), nil

	case addr >= 0x6000:
		return m.sRAM[addr-0x6000], nil
	}

	return 0, errors.Errorf("Invalid mapper reading addr %04x", addr)
}

func (m *Mapper000) Write(addr int, d byte) error {
	if m.useChrRAM && addr < 0x2000 {
		m.chrRAM[addr] = d
	}

	if addr >= 0x6000 && addr < 0x8000 {
		m.sRAM[addr-0x6000] = d
	}

	return nil
}

func (m *Mapper000) Observe(addr int) (d byte, err error) {
	// There is no side effect to reading from mapper 000
	return m.Read(addr)
}

func (m *Mapper000) Populate(prgROM []PrgROMPage, chrROM []ChrROMPage) {
	// If only one page of prg rom, duplicate it
	if len(prgROM) == 1 {
		var pageCopy PrgROMPage
		copy(pageCopy[:], prgROM[0][:])

		prgROM = append(prgROM, pageCopy)
	}

	if len(chrROM) == 0 {
		m.useChrRAM = true
	}

	m.prgROM = prgROM
	m.chrROM = chrROM
}

func (m *Mapper000) GetPRGRom() []PrgROMPage {
	return m.prgROM
}

func (m *Mapper000) readPrgROM(addr int) byte {
	return m.prgROM[addr/PrgROMPageSize][addr%PrgROMPageSize]
}

func (m *Mapper000) readChrROM(addr int) byte {
	return m.chrROM[addr/ChrROMPageSize][addr%ChrROMPageSize]
}
