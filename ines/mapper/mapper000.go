package mapper

import "fmt"

type Mapper000 struct {
	prgROM []PrgROMPage
	chrROM []ChrROMPage
}

func (m *Mapper000) SetSram(b bool) {
	// No sram in ines mapper 000
}

func (m *Mapper000) Read(addr int) byte {
	if addr >= 0 && addr < 0x2000 {
		return m.readChrROM(addr)
	} else if addr >= 0x6000 && addr < 0x10000 {
		return m.readPrgROM(addr)
	}

	// TODO: don't panic
	panic(fmt.Sprintf("invalid mapper accessing addr %04x", addr))
}

func (m *Mapper000) Write(addr int, d byte) int {
	// ROM, no writing
	return 0
}

func (m *Mapper000) Observe(addr int) byte {
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

	m.prgROM = prgROM
	m.chrROM = chrROM
}

func (m *Mapper000) GetPRGRom() []PrgROMPage {
	return m.prgROM
}

func (m *Mapper000) readPrgROM(addr int) byte {
	if addr >= 0x6000 && addr < 0x8000 {
		return 0
	}

	addr -= 0x8000
	return m.prgROM[addr/PrgROMPageSize][addr%PrgROMPageSize]
}

func (m *Mapper000) readChrROM(addr int) byte {
	return m.chrROM[addr/ChrROMPageSize][addr%ChrROMPageSize]
}
