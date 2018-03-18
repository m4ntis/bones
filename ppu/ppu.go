package ppu

import "github.com/m4ntis/bones/models"

type PPU struct {
	RAM *RAM
}

func NewPPU() *PPU {
	var ram RAM
	return &PPU{
		RAM: &ram,
	}
}

func (ppu *PPU) LoadROM(rom *models.ROM) {
	// Load first 2 pages of PrgROM (not supporting mappers as of yet)
	copy(ppu.RAM.data[0x0:0x1000], rom.ChrROM[0][:])
	copy(ppu.RAM.data[0x1000:0x2000], rom.ChrROM[1][:])
}
