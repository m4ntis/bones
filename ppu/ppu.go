package ppu

import (
	"fmt"
	"image/color"

	"github.com/m4ntis/bones/models"
)

type Sprite struct {
	dataLo byte
	dataHi byte

	attr byte

	x byte

	shifted int
}

type PPU struct {
	VRAM *VRAM

	OAM     *OAM
	sOAM    *SecondaryOAM
	sprites [8]Sprite

	evalSprN      int
	foundSprCount int

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

	pixelc chan models.Pixel
}

func New(nmi chan bool, pixelc chan models.Pixel) *PPU {
	var vram VRAM
	var oam OAM
	var soam SecondaryOAM

	return &PPU{
		VRAM: &vram,

		OAM:     &oam,
		sOAM:    &soam,
		sprites: [8]Sprite{},

		scrollFirstWrite: true,
		addrFirstWrite:   true,

		vblank: false,
		nmi:    nmi,

		pixelc: pixelc,
	}
}

func (ppu *PPU) LoadROM(rom *models.ROM) {
	// Load first 2 pages of ChrROM (not supporting mappers as of yet)
	copy(ppu.VRAM.data[0x0:models.ChrROMPageSize], rom.ChrROM[0][:])
}

func (ppu *PPU) Cycle() {
	defer ppu.incCoords()

	if ppu.scanline >= 0 && ppu.scanline < 240 {
		ppu.pixelc <- models.Pixel{
			X: ppu.x,
			Y: ppu.scanline,

			Color: ppu.visibleFrameCycle(),
		}

		ppu.spriteEval()
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
	ptHighByte := ppu.VRAM.Read(patternAddr + pty + 8)
	ptLowBit := ptLowByte >> uint(7-ptx) & 1
	ptHighBit := ptHighByte >> uint(7-ptx) & 1

	bgLo := ptLowBit + ptHighBit<<1

	atQuarter := ppu.x%32/16 + ppu.scanline%32/16<<1

	// Assuming nametable 0, as mentioned above
	atByte := ppu.VRAM.Read(AT0Idx + at)

	bgHi := atByte >> uint(2*atQuarter) & 3

	pIdx := ppu.calcPaletteIdx(int(bgLo), int(bgHi), ppu.calcSprForPixel())

	return Palette[pIdx]
}

func (ppu *PPU) calcPaletteIdx(bgLo, bgHi, sprIdx int) (pIdx byte) {
	if sprIdx == -1 {
		return ppu.VRAM.Read(BgrPaletteIdx + bgLo + bgHi<<2)
	}

	spr := ppu.sprites[sprIdx]
	sprData := spr.dataLo&1 + spr.dataHi&1<<1

	// bg 0 or sprite not opaque and with front priority
	if bgLo == 0 && sprData != 0 && spr.attr>>5&1 == 0 {
		return ppu.VRAM.Read(SprPaletteIdx + int(sprData+spr.attr&3<<2))
	}

	return ppu.VRAM.Read(BgrPaletteIdx + bgLo + bgHi<<2)
}

func (ppu *PPU) spriteEval() {
	if ppu.x >= 1 && ppu.x <= 256 {
		if ppu.x == 1 {
			// Reset ppu sprite eval flags
			ppu.evalSprN = 0
			ppu.foundSprCount = 0
			ppu.sOAM = &SecondaryOAM{}
		}
		if ppu.evalSprN < 64 {
			// Sprite in scanline range
			if ppu.scanline >= int(ppu.OAM[ppu.evalSprN*4]) &&
				ppu.scanline < int(ppu.OAM[ppu.evalSprN*4])+8 {
				if ppu.foundSprCount >= 8 {
					// Set overflow flag
					ppu.ppuStatus &= 1 << 5
				} else {
					// Copy 4 bytes of sprite data from OAM to secondary OAM
					copy(ppu.sOAM[ppu.foundSprCount*4:ppu.foundSprCount*4+4], ppu.OAM[ppu.evalSprN*4:ppu.evalSprN*4+4])
					ppu.foundSprCount++
				}
			}

			ppu.evalSprN++
		}

		ppu.shiftSprites()
	} else if ppu.x >= 257 && ppu.x <= 320 {
		// Gotta do this once every 8 cycles
		if ppu.x%8 == 1 {
			sprN := (ppu.x - 256) / 8

			if ppu.foundSprCount >= sprN+1 {
				// Determine pattern table
				pt := int(ppu.ppuCtrl >> 3 & 1)
				patternAddr := 0x1000*pt + int(ppu.sOAM[sprN*4+1])*16

				// Line of the sprite to be displayed on the scanline
				sprLine := ppu.scanline - int(ppu.sOAM[sprN*4])

				// TODO: Might gotta invert this
				sprDataLo := ppu.VRAM.Read(patternAddr + sprLine)
				sprDataHi := ppu.VRAM.Read(patternAddr + sprLine + 8)

				ppu.sprites[sprN] = Sprite{
					attr: ppu.sOAM[sprN*4+2],

					dataHi: sprDataHi,
					dataLo: sprDataLo,

					x: ppu.sOAM[sprN*4+3],
				}
				fmt.Println(ppu.sprites[sprN])
			} else {
				ppu.sprites[sprN] = Sprite{
					attr: 0xff,

					dataHi: 0xff,
					dataLo: 0xff,

					x: 0xff,
				}
			}
		}
	}
}

func (ppu *PPU) shiftSprites() {
	for i := range ppu.sprites {
		if ppu.sprites[i].shifted < 8 {
			if ppu.sprites[i].x > 0 {
				ppu.sprites[i].x--
			} else {
				ppu.sprites[i].dataHi >>= 1
				ppu.sprites[i].dataLo >>= 1
				ppu.sprites[i].shifted++
			}
		}
	}
}

func (ppu *PPU) calcSprForPixel() (idx int) {
	for i, spr := range ppu.sprites {
		if spr.x == 0 && spr.shifted < 8 && spr.attr>>5&1 == 0 {
			return i
		}
	}
	return -1
}
