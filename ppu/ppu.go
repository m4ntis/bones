package ppu

import (
	"image/color"

	"github.com/m4ntis/bones/models"
)

type PPU struct {
	VRAM *VRAM
	OAM  *OAM

	scanline int
	x        int

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
	var vram VRAM
	var oam OAM

	return &PPU{
		VRAM: &vram,
		OAM:  &oam,

		scrollFirstWrite: true,
		addrFirstWrite:   true,

		vblank: false,
		nmi:    nmi,
	}
}

func (ppu *PPU) LoadROM(rom *models.ROM) {
	// Load first 2 pages of ChrROM (not supporting mappers as of yet)
	copy(ppu.VRAM.data[0x0:models.ChrROMPageSize], rom.ChrROM[0][:])
}

func (ppu *PPU) Cycle() models.Pixel {
	defer ppu.incCoords()

	if ppu.scanline >= 0 && ppu.scanline < 240 {
		return models.Pixel{
			X: ppu.x,
			Y: ppu.scanline,

			Color: ppu.visibleFrameCycle(),
		}
	} else if ppu.scanline == 241 && ppu.x == 1 {
		ppu.vblank = true
		ppu.ppuStatus |= 1 << 7
		if ppu.ppuCtrl>>7 == 1 {
			ppu.nmi <- true
		}
	} else if ppu.scanline == 261 && ppu.x == 1 {
		ppu.ppuStatus = 0
		ppu.vblank = false
	}

	return models.Pixel{
		X:     ppu.x,
		Y:     ppu.scanline,
		Color: color.RGBA{},
	}
}

func (ppu *PPU) PPUCtrlWrite(data byte) {
	// If V Flag set (and changed) while in vblank
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
		ppu.ppuAddr = int(data) << 8
		return
	}
	ppu.ppuAddr |= int(data)
}

func (ppu *PPU) PPUDataRead() byte {
	defer ppu.incAddr()

	// If the read is from palette data, it is immediatelly put on the data bus
	if getAddr(ppu.ppuAddr) >= BgrPaletteIdx {
		// TODO: Reading the palettes still updates the internal buffer though,
		// but the data placed in it is the mirrored nametable data that would
		// appear "underneath" the palette.
		return ppu.VRAM.Read(ppu.ppuAddr)
	}

	defer func() { ppu.ppuDataBuf = ppu.VRAM.Read(ppu.ppuAddr) }()
	return ppu.ppuDataBuf
}

func (ppu *PPU) PPUDataWrite(d byte) {
	defer ppu.incAddr()

	ppu.VRAM.Write(ppu.ppuAddr, d)
}

//TODO: Take note of oamaddr
func (ppu *PPU) DMA(oamData [256]byte) {
	oam := OAM(oamData)
	ppu.OAM = &oam
}

func (ppu *PPU) incAddr() {
	ppu.ppuAddr += int(1 + (ppu.ppuCtrl>>2&1)*31)
}

// TODO: skip cycle (0, 0) on odd frames
func (ppu *PPU) incCoords() {
	ppu.x++
	if ppu.x > 340 {
		ppu.x = 0

		ppu.scanline++
		if ppu.scanline > 261 {
			ppu.scanline = 0
		}
	}
}

func (ppu *PPU) visibleFrameCycle() color.RGBA {
	pt := int(ppu.ppuCtrl >> 4 & 1)
	nt := (ppu.scanline/8)*32 + ppu.x/8
	at := (ppu.scanline/32)*8 + ppu.x/32

	// For now we assume nametable 0
	ntByte := ppu.VRAM.Read(NT0Idx + nt)

	patternAddr := 0x1000*pt + int(ntByte)*16

	ptx := ppu.x % 8
	pty := ppu.scanline % 8
	ptLowByte := ppu.VRAM.Read(patternAddr + pty)
	ptLowBit := ptLowByte >> uint(ptx) & 1
	ptHighByte := ppu.VRAM.Read(patternAddr + pty + 8)
	ptHighBit := ptHighByte >> uint(ptx) & 1

	peAddrLow := ptLowBit + ptHighBit<<1

	atQuarter := ppu.x%32/16 + ppu.scanline%32/16<<1

	// Assuming nametable 0, as mentioned above
	atByte := ppu.VRAM.Read(AT0Idx + at)

	peAddrHigh := atByte >> uint(2*atQuarter) & 3

	peAddr := peAddrLow + peAddrHigh<<2

	pIdx := ppu.VRAM.Read(BgrPaletteIdx + int(peAddr))

	return Palette[pIdx]
}
