package ppu

import (
	"image/color"

	"github.com/m4ntis/bones/models"
)

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

func (ppu *PPU) Cycle(scanline int, x int) color.RGBA {
	if scanline >= 0 && scanline < 240 {
		pt := int(*ppu.Reg.PPUCTRL >> 4 & 1)
		nt := (scanline/8)*32 + x/8
		at := (scanline/32)*8 + x/32

		// For now we assume nametable 0
		ntByte := *ppu.RAM.Fetch(NT0_IDX + nt)

		patternAddr := 0x1000*pt + int(ntByte)*16

		ptx := x % 8
		pty := scanline % 8

		ptLowByte := *ppu.RAM.Fetch(patternAddr + pty)
		ptLowBit := ptLowByte >> uint(ptx) & 1
		ptHighByte := *ppu.RAM.Fetch(patternAddr + pty + 8)
		ptHighBit := ptHighByte >> uint(ptx) & 1

		peAddrLow := ptLowBit + ptHighBit<<1

		atQuarter := x%32/16 + scanline%32/16<<1

		// Assuming nametable 0, as mentioned above
		atByte := *ppu.RAM.Fetch(AT0_IDX + at)

		peAddrHigh := atByte >> uint(2*atQuarter) & 3

		peAddr := peAddrLow + peAddrHigh<<2

		pIdx := *ppu.RAM.Fetch(BGR_PALETTE_IDX + int(peAddr))

		return Palette[pIdx]
	}
	return color.RGBA{}
}
