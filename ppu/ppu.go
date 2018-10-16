// Package ppu implements the NES's Ricoh 2A03 ppu
package ppu

import (
	"image"
	"image/color"

	"github.com/m4ntis/bones/ines"
	"github.com/m4ntis/bones/ines/mapper"
)

const (
	FrontPriority = byte(0)
)

type sprite struct {
	dataLo byte
	dataHi byte

	attr byte

	x byte

	shifted int
}

// Displayer describes a place that the PPU outputs its frames to.
type Displayer interface {
	Display(image.Image)
}

// PPU implements the Ricoh 2A03.
//
// Before starting to run the PPU, it should first be initialized with a parsed
// NES ROM, containing graphic data for the PPU (CHR ROM).
//
// The PPU exports its VRAM which can be read and written to.
//
// The PPU also contains methods for reading and writing to it's registers, as
// they are interfaces via memory mapped i/o on the CPU RAM.
type PPU struct {
	VRAM *VRAM
	OAM  *OAM

	Regs *Registers

	NMI chan bool

	sprites [8]sprite
	sOAM    *secondaryOAM

	evaluatedSprNum   int
	foundSprCount     int
	spriteZeroPresent bool

	scanline int
	x        int

	mirror int

	frame *frame
	disp  Displayer
}

// New initializes a PPU instance and returns it.
//
// NMI is the channel on which the PPU publishes NMIs.
//
// disp is the where the PPU outputs its frames.
func New(mirror int, romMapper mapper.Mapper, disp Displayer) *PPU {
	vram := &VRAM{
		Mapper: romMapper,
	}

	oam := &OAM{}
	soam := &secondaryOAM{}

	nmi := make(chan bool)

	return &PPU{
		VRAM: vram,
		Regs: newRegisters(nmi, oam, vram),

		OAM:     oam,
		sOAM:    soam,
		sprites: [8]sprite{},

		mirror: mirror,

		NMI: nmi,

		frame: newFrame(),
		disp:  disp,
	}
}

//TODO: Take note of oamaddr

// DMA is copies 256 bytes of OAM data to the PPU's OAM
func (ppu *PPU) DMA(oamData [256]byte) {
	oam := OAM(oamData)
	ppu.OAM = &oam
}

// Cycle executes a single PPU cycle.
//
// Cycle may cause the ppu to generate an NMI or output a frame to the display.
func (ppu *PPU) Cycle() {
	if ppu.scanline >= 0 && ppu.scanline < 240 {
		ppu.visibleScanlineCycle()
	} else if ppu.scanline == 241 && ppu.x == 1 {
		ppu.vblankBegin()
	} else if ppu.scanline == 261 && ppu.x == 1 {
		ppu.vblankEnd()
	}

	ppu.incCoords()

	// TODO: What about sprite evaluation on last scanline for scanline 0?
}

// visibleScanlineCycle executes the ppu's logic for scanlines between 0 and
// 239, the visible scanlines.
func (ppu *PPU) visibleScanlineCycle() {
	ppu.evaluateSprites()

	if ppu.x < 256 {
		ppu.frame.push(pixel{
			x: ppu.x,
			y: ppu.scanline,

			color: ppu.calculatePixelValue(),
		})
	}
}

// vblankBegin sets vblank flags, publishes an NMI if nmi is enabled in PPUCTRL
// and pushes a frame to display.
func (ppu *PPU) vblankBegin() {
	// Set vblank flag internally
	ppu.Regs.vblank = true

	// Set bit 7 of PPUSTATUS - vblank flag
	ppu.Regs.ppuStatus |= 1 << 7

	// Check nmi enable flag in PPUCTRL and publish NMI
	if ppu.Regs.ppuCtrl>>7 == 1 {
		ppu.NMI <- true
	}

	// Push frame to display
	ppu.disp.Display(ppu.frame.create())
}

// vblankEnd clears the internal vblank flag and PPUSTATUS.
func (ppu *PPU) vblankEnd() {
	ppu.Regs.vblank = false
	ppu.Regs.ppuStatus = 0
}

// incCoords increments ppu's coordinate parameters for next cycle.
func (ppu *PPU) incCoords() {
	// TODO: skip cycle (0, 0) on odd frames

	ppu.x++
	if ppu.x > 340 {
		ppu.x = 0

		ppu.scanline++
		if ppu.scanline > 261 {
			ppu.scanline = 0
		}
	}
}

// evaluateSprites fetches sprite date for next scanline's sprites during
// visible scanlines.
//
// During each ppu cycle, a small bit of the evaluation happens, depending on
// the cycle's x coordinate.
func (ppu *PPU) evaluateSprites() {
	if ppu.x >= 1 && ppu.x <= 256 {
		if ppu.x == 1 {
			ppu.resetEvaluatedSprites()
		}

		// Evaluate a sprite every 4 cycles (total 64)
		if ppu.x%4 == 0 {
			ppu.evaluateSprite()
			ppu.evaluatedSprNum++
		}

		ppu.shiftSprites()
	} else if ppu.x >= 257 && ppu.x <= 320 {
		// Once every 8 cycles, colour and data is fetched for one out of the 8
		// sprites to be displayed next frame.
		if ppu.x%8 == 1 {
			ppu.renderSprite()
		}
	}
}

// resetEvaluatedSprites resets secondary OAM and evaluated sprite counters.
func (ppu *PPU) resetEvaluatedSprites() {
	ppu.evaluatedSprNum = 0
	ppu.foundSprCount = 0
	ppu.spriteZeroPresent = false

	ppu.sOAM = &secondaryOAM{}
}

// evaluateSprite checks if the currently evaluated sprite matches the next
// frame, and copies it to secondary OAM if so.
func (ppu *PPU) evaluateSprite() {
	sprData := ppu.OAM[ppu.evaluatedSprNum*sprDataSize : (ppu.evaluatedSprNum+1)*sprDataSize]
	sprY := int(sprData[0])

	// Check sprite in range of next scanline
	if ppu.scanline >= sprY && ppu.scanline < sprY+8 {
		if ppu.foundSprCount < 8 {
			// Copy sprite data from OAM to secondary OAM
			copy(ppu.sOAM[ppu.foundSprCount*sprDataSize:(ppu.foundSprCount+1)*sprDataSize],
				sprData)
			ppu.spriteZeroPresent = true
		} else {
			// Set overflow flag
			ppu.Regs.ppuStatus &= 1 << 5
		}

		ppu.foundSprCount++
	}
}

// renderSprite fetches the colours and data of a sprite to be displayed on
// the next frame.
func (ppu *PPU) renderSprite() {
	// Determine which sprite is being rendered (zero indexed)
	renderedSprNum := (ppu.x - 256) / 8
	sprData := ppu.sOAM[renderedSprNum*sprDataSize : (renderedSprNum+1)*sprDataSize]

	// Check whether the sprite is in range of the displayed sprites. If not,
	// all it's data will be 0xff.
	if ppu.foundSprCount >= renderedSprNum+1 {
		// Determine pattern table number and address to fetch sprite data from
		pt := int(ppu.Regs.ppuCtrl >> 3 & 1)
		ptAddr := pt*PTSize + int(sprData[1])*16

		// Calculate line of the sprite to be displayed on the scanline
		sprLine := ppu.scanline - int(sprData[0])

		// Fetch sprite data
		sprDataLo := ppu.VRAM.Read(ptAddr + sprLine)
		sprDataHi := ppu.VRAM.Read(ptAddr + sprLine + 8)

		// Fetch sprite attribute byte
		attr := sprData[2]

		// Invert sprite if horizontal invert bit is off
		if attr>>6&1 == 0 {
			sprDataLo = flip_byte(sprDataLo)
			sprDataHi = flip_byte(sprDataHi)
		}

		ppu.sprites[renderedSprNum] = sprite{
			attr: attr,

			dataHi: sprDataHi,
			dataLo: sprDataLo,

			x: sprData[3],
		}
	} else {
		ppu.sprites[renderedSprNum] = sprite{
			attr: 0xff,

			dataHi: 0xff,
			dataLo: 0xff,

			x: 0xff,
		}
	}
}

func (ppu *PPU) calculatePixelValue() color.RGBA {
	scrolledX := ppu.x + ppu.Regs.xScroll

	// Select Pattern Table, Name Table and Attribute Table
	pt := ppu.getPTAddr()
	nt := (ppu.scanline/8)*32 + scrolledX%256/8
	at := (ppu.scanline/32)*8 + scrolledX%256/32

	ntBase := getNTAddr(int(ppu.Regs.ppuCtrl)&3+scrolledX/256, ppu.mirror)
	ntByte := ppu.VRAM.Read(ntBase + nt)

	patternAddr := pt + int(ntByte)*16

	ptx := scrolledX % 8
	pty := ppu.scanline % 8
	ptLowByte := ppu.VRAM.Read(patternAddr + pty)
	ptHighByte := ppu.VRAM.Read(patternAddr + pty + 8)
	ptLowBit := ptLowByte >> uint(7-ptx) & 1
	ptHighBit := ptHighByte >> uint(7-ptx) & 1

	bgLo := ptLowBit + ptHighBit<<1

	atQuarter := scrolledX%32/16 + ppu.scanline%32/16<<1

	atBase := getATAddr(ntBase)
	atByte := ppu.VRAM.Read(atBase + at)

	bgHi := atByte >> uint(2*atQuarter) & 3

	sprLo, sprHi, sprPriority, sprN := ppu.calcSprForPixel()
	pAddr := ppu.calcPaletteAddr(bgLo, bgHi, sprLo, sprHi, sprPriority)

	if ppu.spriteZeroPresent && sprN == 0 && int(sprLo+sprHi<<2) != 0 && int(bgLo+bgHi<<2) != 0 {
		// Set sprite 0 hit flag
		ppu.Regs.ppuStatus |= 1 << 6
	}

	return Palette[pAddr]
}

func (ppu *PPU) getPTAddr() int {
	if int(ppu.Regs.ppuCtrl>>4&1) == 0 {
		return PT0Addr
	} else {
		return PT1Addr
	}
}

func getNTAddr(nt int, mirroring int) int {
	// Assuming either horizontal or vertical mirroring
	if mirroring == ines.HorizontalMirroring {
		if nt == 0 || nt == 1 {
			return NT0Addr
		} else {
			return NT2Addr
		}
	} else {
		if nt == 0 || nt == 2 {
			return NT0Addr
		} else {
			return NT1Addr
		}
	}
}

func getATAddr(nt int) int {
	return nt + 0x3c0
}

func (ppu *PPU) calcPaletteAddr(bgLo, bgHi, sprLo, sprHi, sprPriority byte) (pAddr byte) {
	// sprPriority == 255 when a sprite wasn't found for the scanline
	// bg 0 or sprite not opaque and with front priority
	if sprPriority != 255 && (bgLo == 0 || (sprLo != 0 && sprPriority == FrontPriority)) {
		return ppu.VRAM.Read(SprPaletteAddr + int(sprLo+sprHi<<2))
	}

	return ppu.VRAM.Read(BgrPaletteAddr + int(bgLo) + int(bgHi<<2))
}

// shiftSprites is implemented in a duff machine fashion for optimization
// purposes.
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

// calcSprForPixel is implemented in a duff machine fashion for optimization
// purposes.
func (ppu *PPU) calcSprForPixel() (sprLo, sprHi, priority byte, sprN int) {
	if ppu.sprites[0].x == 0 && ppu.sprites[0].shifted < 8 {
		sprLo := ppu.sprites[0].dataLo&1 + ppu.sprites[0].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[0].attr & 3, ppu.sprites[0].attr >> 5 & 1, 0
		}
	}
	if ppu.sprites[1].x == 0 && ppu.sprites[1].shifted < 8 {
		sprLo := ppu.sprites[1].dataLo&1 + ppu.sprites[1].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[1].attr & 3, ppu.sprites[1].attr >> 5 & 1, 1
		}
	}
	if ppu.sprites[2].x == 0 && ppu.sprites[2].shifted < 8 {
		sprLo := ppu.sprites[2].dataLo&1 + ppu.sprites[2].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[2].attr & 3, ppu.sprites[2].attr >> 5 & 1, 2
		}
	}
	if ppu.sprites[3].x == 0 && ppu.sprites[3].shifted < 8 {
		sprLo := ppu.sprites[3].dataLo&1 + ppu.sprites[3].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[3].attr & 3, ppu.sprites[3].attr >> 5 & 1, 3
		}
	}
	if ppu.sprites[4].x == 0 && ppu.sprites[4].shifted < 8 {
		sprLo := ppu.sprites[4].dataLo&1 + ppu.sprites[4].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[4].attr & 3, ppu.sprites[4].attr >> 5 & 1, 4
		}
	}
	if ppu.sprites[5].x == 0 && ppu.sprites[5].shifted < 8 {
		sprLo := ppu.sprites[5].dataLo&1 + ppu.sprites[5].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[5].attr & 3, ppu.sprites[5].attr >> 5 & 1, 5
		}
	}
	if ppu.sprites[6].x == 0 && ppu.sprites[6].shifted < 8 {
		sprLo := ppu.sprites[6].dataLo&1 + ppu.sprites[6].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[6].attr & 3, ppu.sprites[6].attr >> 5 & 1, 6
		}
	}
	if ppu.sprites[7].x == 0 && ppu.sprites[7].shifted < 8 {
		sprLo := ppu.sprites[7].dataLo&1 + ppu.sprites[7].dataHi&1<<1
		if sprLo != 0 {
			return sprLo, ppu.sprites[7].attr & 3, ppu.sprites[7].attr >> 5 & 1, 7
		}
	}
	return 0, 0, 255, -1
}

func flip_byte(d byte) byte {
	d = ((d >> 1) & 0x55) | ((d & 0x55) << 1)
	d = ((d >> 2) & 0x33) | ((d & 0x33) << 2)
	d = ((d >> 4) & 0x0F) | ((d & 0x0F) << 4)
	return d
}
