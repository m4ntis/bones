package ppu

import "github.com/m4ntis/bones/models"

type PPU struct {
	RAM *RAM
	OAM *OAM
	Reg Regs
}

func New(ppuctrl *byte, ppumask *byte, ppustatus *byte, oamaddr *byte,
	oamdata *byte, ppuscroll *byte, ppuaddr *byte, ppudata *byte,
	oamdma *byte) *PPU {
	var ram RAM
	var oam OAM

	return &PPU{
		RAM: &ram,
		OAM: &oam,
		Reg: Regs{
			PPUCTRL:   ppuctrl,
			PPUMASK:   ppumask,
			PPUSTATUS: ppustatus,
			OAMADDR:   oamaddr,
			OAMDATA:   oamdata,
			PPUSCROLL: ppuscroll,
			PPUADDR:   ppuaddr,
			PPUDATA:   ppudata,
			OAMDMA:    oamdma,
		},
	}
}

func (ppu *PPU) LoadROM(rom *models.ROM) {
	// Load first 2 pages of PrgROM (not supporting mappers as of yet)
	copy(ppu.RAM.data[0x0:models.CHR_ROM_PAGE_SIZE], rom.ChrROM[0][:])
}

func (ppu *PPU) DMA(oamData [256]byte) {
	ppu.OAM = &OAM{oamData}
}
