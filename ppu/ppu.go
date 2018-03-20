package ppu

import (
	"image/color"

	"github.com/m4ntis/bones/models"
)

type PPU struct {
	RAM *RAM
	OAM *OAM

	ppuctrl   byte
	ppumask   byte
	ppustatus byte
	oamaddr   byte
	oamdata   byte
	ppuscroll byte
	ppuaddr   byte
	ppudata   byte
	oamdma    byte

	scrollSecondWrite bool
	xScroll           byte
	yScroll           byte

	vblank bool
	nmi    chan bool
}

func New(nmi chan bool) *PPU {
	var ram RAM
	var oam OAM

	return &PPU{
		RAM: &ram,
		OAM: &oam,

		vblank: false,
		nmi:    nmi,
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
		pt := int(*ppu.PPUCTRL >> 4 & 1)
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

func (ppu *PPU) PPUCtrlWrite(data byte) {
	// If V Flag set while in vblank
	if ppu.ppuctrl>>7 == 0 && data>>7 == 1 && ppu.vblank {
		ppu.nmi <- true
	}

	ppu.ppuctrl = data
}

func (ppu *PPU) PPUMaskWrite(data byte) {
	ppu.ppumask = data
}

func (ppu *PPU) PPUStatusRead() byte {
	defer func() {
		// Clear bit 7
		ppu.ppustatus &= 0x7f

		ppu.ppuscroll &= 0
		ppu.ppuaddr &= 0
	}()

	return ppu.ppustatus
}

func (ppu *PPU) OAMAddrWrite(data byte) {
	ppu.oamaddr = data
}

func (ppu *PPU) OAMDataRead(data byte) {
	ppu.oamaddr++
}

func (ppu *PPU) OAMDataWrite(data byte) {
	ppu.oamaddr++
}

func (ppu *PPU) PpuScrollWrite(data byte) {
	defer func() { ppu.scrollSecondWrite = !ppu.scrollSecondWrite }()

	if ppu.scrollSecondWrite {
		ppu.yScroll = data
		return
	}

	ppu.xScroll = data
}

func (ppu *PPU) PPUADDRWrite(data byte) {
}

func (ppu *PPU) PPUDATARead() {
}

func (ppu *PPU) PPUDATAWrite(data byte) {
}

func (ppu *PPU) OAMDMAWrite(data byte) {
}
