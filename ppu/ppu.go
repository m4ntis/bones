package ppu

import (
	"image/color"

	"github.com/m4ntis/bones/models"
)

type PPU struct {
	RAM *RAM
	OAM *OAM

	ppuCtrl   byte
	ppuMask   byte
	ppuStatus byte

	oamAddr byte

	scrollFirstWrite bool
	ppuScroll        int

	addrFirstWrite bool
	ppuAddr        int

	ppuData    byte
	ppuDataBuf byte

	vblank bool
	nmi    chan bool
}

func New(nmi chan bool) *PPU {
	var ram RAM
	var oam OAM

	return &PPU{
		RAM: &ram,
		OAM: &oam,

		scrollFirstWrite: true,

		vblank: false,
		nmi:    nmi,
	}
}

func (ppu *PPU) LoadROM(rom *models.ROM) {
	// Load first 2 pages of PrgROM (not supporting mappers as of yet)
	copy(ppu.RAM.data[0x0:models.ChrROMPageSize], rom.ChrROM[0][:])
}

//TODO: Take note of oamaddr
func (ppu *PPU) DMA(oamData [256]byte) {
	oam := OAM(oamData)
	ppu.OAM = &oam
}

func (ppu *PPU) Cycle(scanline int, x int) color.RGBA {
	if scanline >= 0 && scanline < 240 {
		pt := int(ppu.ppuCtrl >> 4 & 1)
		nt := (scanline/8)*32 + x/8
		at := (scanline/32)*8 + x/32

		// For now we assume nametable 0
		ntByte := ppu.RAM.Read(NT0Idx + nt)

		patternAddr := 0x1000*pt + int(ntByte)*16

		ptx := x % 8
		pty := scanline % 8
		ptLowByte := ppu.RAM.Read(patternAddr + pty)
		ptLowBit := ptLowByte >> uint(ptx) & 1
		ptHighByte := ppu.RAM.Read(patternAddr + pty + 8)
		ptHighBit := ptHighByte >> uint(ptx) & 1

		peAddrLow := ptLowBit + ptHighBit<<1

		atQuarter := x%32/16 + scanline%32/16<<1

		// Assuming nametable 0, as mentioned above
		atByte := ppu.RAM.Read(AT0Idx + at)

		peAddrHigh := atByte >> uint(2*atQuarter) & 3

		peAddr := peAddrLow + peAddrHigh<<2

		pIdx := ppu.RAM.Read(BgrPaletteIdx + int(peAddr))

		return Palette[pIdx]
	}
	return color.RGBA{}
}

func (ppu *PPU) PPUCtrlWrite(data byte) {
	// If V Flag set while in vblank
	if ppu.ppuCtrl>>7 == 0 && data>>7 == 1 && ppu.vblank {
		ppu.nmi <- true
	}

	ppu.ppuCtrl = data
}

func (ppu *PPU) PPUMaskWrite(data byte) {
	ppu.ppuMask = data
}

func (ppu *PPU) PPUStatusRead() byte {
	defer func() {
		// Clear bit 7
		ppu.ppuStatus &= 0x7f

		ppu.scrollFirstWrite = true
		ppu.ppuScroll &= 0

		ppu.ppuAddr &= 0
	}()

	return ppu.ppuStatus
}

func (ppu *PPU) OAMAddrWrite(data byte) {
	ppu.oamAddr = data
}

func (ppu *PPU) OAMDataRead() byte {
	return (*ppu.OAM)[ppu.oamAddr]
}

func (ppu *PPU) OAMDataWrite(data byte) {
	ppu.oamAddr++

	(*ppu.OAM)[ppu.oamAddr] = data
}

// TODO: Changes made to the vertical scroll during rendering will only take
// effect on the next frame
func (ppu *PPU) PPUScrollWrite(data byte) {
	defer func() { ppu.scrollFirstWrite = !ppu.scrollFirstWrite }()

	if ppu.scrollFirstWrite {
		ppu.ppuScroll = int(data)
		return
	}

	ppu.ppuScroll |= int(data) << 8
}

func (ppu *PPU) PPUAddrWrite(data byte) {
	defer func() { ppu.addrFirstWrite = !ppu.addrFirstWrite }()

	if ppu.addrFirstWrite {
		ppu.ppuAddr = int(data)
		return
	}

	ppu.ppuAddr |= int(data) << 8
}

func (ppu *PPU) PPUDataRead() byte {
	defer ppu.incAddr()

	// If the read is from palette data, it is immediatelly put on the data bus
	if getAddr(ppu.ppuAddr) >= BgrPaletteIdx {
		// TODO: Reading the palettes still updates the internal buffer though,
		// but the data placed in it is the mirrored nametable data that would
		// appear "underneath" the palette.
		return ppu.RAM.Read(ppu.ppuAddr)
	}

	defer func() { ppu.ppuDataBuf = ppu.RAM.Read(ppu.ppuAddr) }()
	return ppu.ppuDataBuf
}

func (ppu *PPU) PPUDataWrite(d byte) {
	defer ppu.incAddr()

	ppu.RAM.Write(ppu.ppuAddr, d)
}

func (ppu *PPU) incAddr() {
	ppu.ppuAddr += int(1 + (ppu.ppuCtrl>>2&1)*31)
}
