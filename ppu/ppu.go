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

var (
	nilSprite = sprite{
		low:     0xff,
		high:    0xff,
		palette: 0xff,

		x:       0xff,
		shifted: 0xff,

		priority:   false,
		spriteZero: false,
	}
)

// sprite is data structure internal to the ppu, holding data about a loaded
// sprite.
type sprite struct {
	low     byte
	high    byte
	palette byte

	x       byte
	shifted int

	priority   bool
	spriteZero bool
}

func (spr sprite) getColor() byte {
	return spr.low&1 + (spr.high&1)<<1
}

func (spr sprite) getData() int {
	return int(spr.getColor() + spr.palette<<2)
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

			// On x == 4 the ppu evaluates sprite zero. If so, set the flag
			if ppu.x == 4 {
				ppu.spriteZeroPresent = true
			}
		} else {
			// Set overflow flag
			ppu.Regs.ppuStatus |= (1 << 5)
		}

		ppu.foundSprCount++
	}
}

// renderSprite fetches the colours and data of a sprite to be displayed on
// the next frame.
func (ppu *PPU) renderSprite() {
	// Determine which sprite is being rendered (zero indexed, 0 ~ 7)
	renderedSprNum := (ppu.x - 256) / 8
	sprData := ppu.sOAM[renderedSprNum*sprDataSize : (renderedSprNum+1)*sprDataSize]

	// Check whether this sprite slot is used for next frame
	if ppu.foundSprCount >= renderedSprNum+1 {
		// Determine pattern table number and address to fetch sprite data from
		pt := int((ppu.Regs.ppuCtrl >> 3) & 1)
		ptAddr := pt*PTSize + int(sprData[1])*16

		// Calculate line of the sprite to be displayed on the scanline
		sprLine := ppu.scanline - int(sprData[0])

		// Fetch sprite data
		dataLow := ppu.VRAM.Read(ptAddr + sprLine)
		dataHigh := ppu.VRAM.Read(ptAddr + sprLine + 8)

		// Invert sprite if horizontal invert bit of attribute byte is off
		if (sprData[2]>>6)&1 == 0 {
			dataLow = flip_byte(dataLow)
			dataHigh = flip_byte(dataHigh)
		}

		// Fill sprite slot with sprite data
		ppu.sprites[renderedSprNum] = sprite{
			low:     dataLow,
			high:    dataHigh,
			palette: sprData[2] & 3,

			x:       sprData[3],
			shifted: 0,

			priority:   (sprData[2]>>5)&1 == FrontPriority,
			spriteZero: renderedSprNum == 0 && ppu.spriteZeroPresent,
		}
	} else {
		ppu.sprites[renderedSprNum] = nilSprite
	}
}

func (ppu *PPU) calculatePixelValue() color.RGBA {
	bgr := ppu.calcBackground()
	spr := ppu.matchSprite()

	paletteAddr := ppu.muxPixel(bgr, spr)

	if spr.spriteZero && spr.getColor() != 0 && bgr&3 != 0 {
		// Set sprite 0 hit flag
		ppu.Regs.ppuStatus |= (1 << 6)
	}

	return Palette[paletteAddr]
}

func (ppu *PPU) calcBackground() (bgr int) {
	// Add x fine scroll to ppu.x
	scrolledX := ppu.x + ppu.Regs.xScroll

	// Fetch byte from NT
	baseNTAddr := getNTAddr(int(ppu.Regs.ppuCtrl)&3+scrolledX/256, ppu.mirror)
	ntIdx := (ppu.scanline/8)*32 + scrolledX%0x100/8
	byteFromNT := ppu.VRAM.Read(baseNTAddr + ntIdx)

	// Fetch pattern table address
	basePTAddr := ppu.getPTAddr()
	patternAddr := basePTAddr + int(byteFromNT)*16

	// Fetch pattern line from PT (Y coordinate)
	pty := ppu.scanline % 8
	ptLowByte := ppu.VRAM.Read(patternAddr + pty)
	ptHighByte := ppu.VRAM.Read(patternAddr + pty + 8)

	// Fetch pixel data from pattern line (X coordinate)
	ptx := scrolledX % 8
	ptLowBit := (ptLowByte >> uint(7-ptx)) & 1
	ptHighBit := (ptHighByte >> uint(7-ptx)) & 1

	// BGR low 2 bits are added together (the pattern bits)
	bgrLow := ptLowBit + ptHighBit<<1

	// Fetch byte from AT
	baseATAddr := getATAddr(baseNTAddr)
	atIdx := (ppu.scanline/32)*8 + scrolledX%0x100/32
	byteFromAT := ppu.VRAM.Read(baseATAddr + atIdx)

	// Calculate which 2 bits of the byte from AT are relevant
	atQuarter := scrolledX%32/16 + ppu.scanline%32/16<<1

	// Fetch BGR high 2 bits from the byte from AT
	bgrHigh := (byteFromAT >> uint(2*atQuarter)) & 3

	return int(bgrLow + (bgrHigh&3)<<2)
}

func getNTAddr(ntNum int, mirroring int) int {
	// Assuming either horizontal or vertical mirroring
	if mirroring == ines.HorizontalMirroring {
		if ntNum == 0 || ntNum == 1 {
			return NT0Addr
		} else {
			return NT2Addr
		}
	} else {
		if ntNum == 0 || ntNum == 2 {
			return NT0Addr
		} else {
			return NT1Addr
		}
	}
}

func (ppu *PPU) getPTAddr() int {
	if int((ppu.Regs.ppuCtrl>>4)&1) == 0 {
		return PT0Addr
	} else {
		return PT1Addr
	}
}

func getATAddr(ntAddr int) int {
	return ntAddr + 0x3c0
}

func (ppu *PPU) muxPixel(bgr int, spr sprite) (paletteAddr byte) {
	// No sprite was found for current pixel, display background
	if spr == nilSprite {
		return ppu.VRAM.Read(BgrPaletteAddr + bgr)
	}

	// Background clear, display sprite
	if bgr&3 == 0 {
		return ppu.VRAM.Read(SprPaletteAddr + spr.getData())
	}

	// Sprite isn't clear and with priority, display sprite
	if spr.getColor() != 0 && spr.priority {
		return ppu.VRAM.Read(SprPaletteAddr + spr.getData())
	}

	// Display background
	return ppu.VRAM.Read(BgrPaletteAddr + bgr)
}

// shiftSprites is implemented in a duff machine fashion for optimization
// purposes.
func (ppu *PPU) shiftSprites() {
	if ppu.sprites[0].shifted < 8 {
		if ppu.sprites[0].x > 0 {
			ppu.sprites[0].x--
		} else {
			ppu.sprites[0].high >>= 1
			ppu.sprites[0].low >>= 1
			ppu.sprites[0].shifted++
		}
	}
	if ppu.sprites[1].shifted < 8 {
		if ppu.sprites[1].x > 0 {
			ppu.sprites[1].x--
		} else {
			ppu.sprites[1].high >>= 1
			ppu.sprites[1].low >>= 1
			ppu.sprites[1].shifted++
		}
	}
	if ppu.sprites[2].shifted < 8 {
		if ppu.sprites[2].x > 0 {
			ppu.sprites[2].x--
		} else {
			ppu.sprites[2].high >>= 1
			ppu.sprites[2].low >>= 1
			ppu.sprites[2].shifted++
		}
	}
	if ppu.sprites[3].shifted < 8 {
		if ppu.sprites[3].x > 0 {
			ppu.sprites[3].x--
		} else {
			ppu.sprites[3].high >>= 1
			ppu.sprites[3].low >>= 1
			ppu.sprites[3].shifted++
		}
	}
	if ppu.sprites[4].shifted < 8 {
		if ppu.sprites[4].x > 0 {
			ppu.sprites[4].x--
		} else {
			ppu.sprites[4].high >>= 1
			ppu.sprites[4].low >>= 1
			ppu.sprites[4].shifted++
		}
	}
	if ppu.sprites[5].shifted < 8 {
		if ppu.sprites[5].x > 0 {
			ppu.sprites[5].x--
		} else {
			ppu.sprites[5].high >>= 1
			ppu.sprites[5].low >>= 1
			ppu.sprites[5].shifted++
		}
	}
	if ppu.sprites[6].shifted < 8 {
		if ppu.sprites[6].x > 0 {
			ppu.sprites[6].x--
		} else {
			ppu.sprites[6].high >>= 1
			ppu.sprites[6].low >>= 1
			ppu.sprites[6].shifted++
		}
	}
	if ppu.sprites[7].shifted < 8 {
		if ppu.sprites[7].x > 0 {
			ppu.sprites[7].x--
		} else {
			ppu.sprites[7].high >>= 1
			ppu.sprites[7].low >>= 1
			ppu.sprites[7].shifted++
		}
	}
}

// matchSprite goes over the ppu spries to be loaded this frame and returns the
// the first non-clear sprite in x coordinate range.
//
// matchSprite is implemented in a duff machine fashion for optimization
// purposes.
func (ppu *PPU) matchSprite() sprite {
	if ppu.sprites[0].x == 0 && ppu.sprites[0].shifted < 8 {
		if ppu.sprites[0].getColor() != 0 {
			return ppu.sprites[0]
		}
	}
	if ppu.sprites[1].x == 0 && ppu.sprites[1].shifted < 8 {
		if ppu.sprites[1].getColor() != 0 {
			return ppu.sprites[1]
		}
	}
	if ppu.sprites[2].x == 0 && ppu.sprites[2].shifted < 8 {
		if ppu.sprites[2].getColor() != 0 {
			return ppu.sprites[2]
		}
	}
	if ppu.sprites[3].x == 0 && ppu.sprites[3].shifted < 8 {
		if ppu.sprites[3].getColor() != 0 {
			return ppu.sprites[3]
		}
	}
	if ppu.sprites[4].x == 0 && ppu.sprites[4].shifted < 8 {
		if ppu.sprites[4].getColor() != 0 {
			return ppu.sprites[4]
		}
	}
	if ppu.sprites[5].x == 0 && ppu.sprites[5].shifted < 8 {
		if ppu.sprites[5].getColor() != 0 {
			return ppu.sprites[5]
		}
	}
	if ppu.sprites[6].x == 0 && ppu.sprites[6].shifted < 8 {
		if ppu.sprites[6].getColor() != 0 {
			return ppu.sprites[6]
		}
	}
	if ppu.sprites[7].x == 0 && ppu.sprites[7].shifted < 8 {
		if ppu.sprites[7].getColor() != 0 {
			return ppu.sprites[7]
		}
	}

	return nilSprite
}

func flip_byte(d byte) byte {
	d = ((d >> 1) & 0x55) | ((d & 0x55) << 1)
	d = ((d >> 2) & 0x33) | ((d & 0x33) << 2)
	d = ((d >> 4) & 0x0F) | ((d & 0x0F) << 4)
	return d
}
