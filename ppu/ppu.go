package ppu

import "github.com/m4ntis/bones/models"

type PPU struct {
	RAM *RAM
}

func New() *PPU {
	var ram RAM
	return &PPU{
		RAM: &ram,
	}
}

func (ppu *PPU) LoadROM(rom *models.ROM) {
	// Load first 2 pages of PrgROM (not supporting mappers as of yet)
	copy(ppu.RAM.data[0x0:models.CHR_ROM_PAGE_SIZE], rom.ChrROM[0][:])
}
