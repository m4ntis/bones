package ppu

import "github.com/m4ntis/bones/models"

type PPU struct {
	RAM *RAM
	OAM *OAM
}

func New() *PPU {
	var ram RAM
	var oam OAM

	return &PPU{
		RAM: &ram,
		OAM: &oam,
	}
}

func (ppu *PPU) LoadROM(rom *models.ROM) {
	// Load first 2 pages of PrgROM (not supporting mappers as of yet)
	copy(ppu.RAM.data[0x0:models.CHR_ROM_PAGE_SIZE], rom.ChrROM[0][:])
}

func (ppu *PPU) DMA(oamData [256]byte) {
	ppu.OAM = &OAM{oamData}
}
