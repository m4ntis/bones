package ppu

import (
	"image"
	"image/color"

	"github.com/m4ntis/bones/ines"
)

const (
	FrontPriority = byte(0)
)

type Sprite struct {
	dataLo byte
	dataHi byte

	attr byte

	x byte

	shifted int
}

type Displayer interface {
	Display(image.Image)
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

	frame *frame
	disp  Displayer
}

func New(nmi chan bool, disp Displayer) *PPU {
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

		frame: newFrame(),
		disp:  disp,
	}
}

func (ppu *PPU) LoadROM(rom *ines.ROM) {
	// Load first 2 pages of ChrROM (not supporting mappers as of yet)
	copy(ppu.VRAM.data[0x0:ines.ChrROMPageSize], rom.ChrROM[0][:])
}

func (ppu *PPU) Cycle() {
	if ppu.scanline >= 0 && ppu.scanline < 240 {
		if ppu.x < 256 {
			ppu.frame.push(pixel{
				x: ppu.x,
				y: ppu.scanline,

				color: ppu.visibleFrameCycle(),
			})
		}

		ppu.spriteEval()
	} else if ppu.scanline == 241 && ppu.x == 1 {
		ppu.vblank = true
		ppu.ppuStatus |= 1 << 7
		if ppu.ppuCtrl>>7 == 1 {
			ppu.nmi <- true
		}

		ppu.disp.Display(ppu.frame.create())
	} else if ppu.scanline == 261 && ppu.x == 1 {
		ppu.ppuStatus = 0
		ppu.vblank = false
	}

	ppu.incCoords()
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

	sprLo, sprHi, sprPriority := ppu.calcSprForPixel()
	pIdx := ppu.calcPaletteIdx(bgLo, bgHi, sprLo, sprHi, sprPriority)

	return Palette[pIdx]
}

func (ppu *PPU) calcPaletteIdx(bgLo, bgHi, sprLo, sprHi, sprPriority byte) (pIdx byte) {
	// bg 0 or sprite not opaque and with front priority
	if bgLo == 0 || (sprLo != 0 && sprPriority == FrontPriority) {
		return ppu.VRAM.Read(SprPaletteIdx + int(sprLo+sprHi<<2))
	}

	return ppu.VRAM.Read(BgrPaletteIdx + int(bgLo) + int(bgHi<<2))
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

				attr := ppu.sOAM[sprN*4+2]
				// Check horizontal invert bit off
				if attr>>6&1 == 0 {
					sprDataLo = flip_byte(sprDataLo)
					sprDataHi = flip_byte(sprDataHi)
				}

				ppu.sprites[sprN] = Sprite{
					attr: ppu.sOAM[sprN*4+2],

					dataHi: sprDataHi,
					dataLo: sprDataLo,

					x: ppu.sOAM[sprN*4+3],
				}
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

// shiftSprites is implemented in a duff machine fashion for optimization
// purposes
func (ppu *PPU) shiftSprites() {
	if ppu.sprites[0].shifted < 8 {
		if ppu.sprites[0].x > 0 {
			ppu.sprites[0].x--
		} else {
			ppu.sprites[0].dataHi >>= 1
			ppu.sprites[0].dataLo >>= 1
			ppu.sprites[0].shifted++
		}
	}
	if ppu.sprites[1].shifted < 8 {
		if ppu.sprites[1].x > 0 {
			ppu.sprites[1].x--
		} else {
			ppu.sprites[1].dataHi >>= 1
			ppu.sprites[1].dataLo >>= 1
			ppu.sprites[1].shifted++
		}
	}
	if ppu.sprites[2].shifted < 8 {
		if ppu.sprites[2].x > 0 {
			ppu.sprites[2].x--
		} else {
			ppu.sprites[2].dataHi >>= 1
			ppu.sprites[2].dataLo >>= 1
			ppu.sprites[2].shifted++
		}
	}
	if ppu.sprites[3].shifted < 8 {
		if ppu.sprites[3].x > 0 {
			ppu.sprites[3].x--
		} else {
			ppu.sprites[3].dataHi >>= 1
			ppu.sprites[3].dataLo >>= 1
			ppu.sprites[3].shifted++
		}
	}
	if ppu.sprites[4].shifted < 8 {
		if ppu.sprites[4].x > 0 {
			ppu.sprites[4].x--
		} else {
			ppu.sprites[4].dataHi >>= 1
			ppu.sprites[4].dataLo >>= 1
			ppu.sprites[4].shifted++
		}
	}
	if ppu.sprites[5].shifted < 8 {
		if ppu.sprites[5].x > 0 {
			ppu.sprites[5].x--
		} else {
			ppu.sprites[5].dataHi >>= 1
			ppu.sprites[5].dataLo >>= 1
			ppu.sprites[5].shifted++
		}
	}
	if ppu.sprites[6].shifted < 8 {
		if ppu.sprites[6].x > 0 {
			ppu.sprites[6].x--
		} else {
			ppu.sprites[6].dataHi >>= 1
			ppu.sprites[6].dataLo >>= 1
			ppu.sprites[6].shifted++
		}
	}
	if ppu.sprites[7].shifted < 8 {
		if ppu.sprites[7].x > 0 {
			ppu.sprites[7].x--
		} else {
			ppu.sprites[7].dataHi >>= 1
			ppu.sprites[7].dataLo >>= 1
			ppu.sprites[7].shifted++
		}
	}
}

func (ppu *PPU) calcSprForPixel() (sprLo, sprHi, priority byte) {
	if ppu.sprites[0].x == 0 && ppu.sprites[0].shifted < 8 {
		sprLo := ppu.sprites[0].dataLo&1 + ppu.sprites[0].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[0].attr & 3, ppu.sprites[0].attr >> 5 & 1
		}
	}
	if ppu.sprites[1].x == 0 && ppu.sprites[1].shifted < 8 {
		sprLo := ppu.sprites[1].dataLo&1 + ppu.sprites[1].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[1].attr & 3, ppu.sprites[1].attr >> 5 & 1
		}
	}
	if ppu.sprites[2].x == 0 && ppu.sprites[2].shifted < 8 {
		sprLo := ppu.sprites[2].dataLo&1 + ppu.sprites[2].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[2].attr & 3, ppu.sprites[2].attr >> 5 & 1
		}
	}
	if ppu.sprites[3].x == 0 && ppu.sprites[3].shifted < 8 {
		sprLo := ppu.sprites[3].dataLo&1 + ppu.sprites[3].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[3].attr & 3, ppu.sprites[3].attr >> 5 & 1
		}
	}
	if ppu.sprites[4].x == 0 && ppu.sprites[4].shifted < 8 {
		sprLo := ppu.sprites[4].dataLo&1 + ppu.sprites[4].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[4].attr & 3, ppu.sprites[4].attr >> 5 & 1
		}
	}
	if ppu.sprites[5].x == 0 && ppu.sprites[5].shifted < 8 {
		sprLo := ppu.sprites[5].dataLo&1 + ppu.sprites[5].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[5].attr & 3, ppu.sprites[5].attr >> 5 & 1
		}
	}
	if ppu.sprites[6].x == 0 && ppu.sprites[6].shifted < 8 {
		sprLo := ppu.sprites[6].dataLo&1 + ppu.sprites[6].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[6].attr & 3, ppu.sprites[6].attr >> 5 & 1
		}
	}
	if ppu.sprites[7].x == 0 && ppu.sprites[7].shifted < 8 {
		sprLo := ppu.sprites[7].dataLo&1 + ppu.sprites[7].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[7].attr & 3, ppu.sprites[7].attr >> 5 & 1
		}
	}
	return 0, 0, 1
}

func flip_byte(d byte) byte {
	d = ((d >> 1) & 0x55) | ((d & 0x55) << 1)
	d = ((d >> 2) & 0x33) | ((d & 0x33) << 2)
	d = ((d >> 4) & 0x0F) | ((d & 0x0F) << 4)
	return d
}
