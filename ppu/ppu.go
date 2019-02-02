// Package ppu implements the NES's Ricoh 2A03 ppu
package ppu

import (
	"image"
	"image/color"

	"github.com/m4ntis/bones/ines"
)

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

	// Sprites data for current and next frame
	sprites [8]sprite
	sOAM    *secondaryOAM

	// Sprite evaluation counters
	evaluatedSprNum   int
	foundSprCount     int
	spriteZeroPresent bool

	// Current pixel coordinates
	scanline int
	x        int
	oddCycle bool

	// Mirroring type
	mirror int

	// Output
	frame *frame
	disp  Displayer
}

// New initializes a PPU instance and returns it.
//
// disp is the where the PPU outputs its frames.
func New(disp Displayer) *PPU {
	vram := &VRAM{}
	oam := &OAM{}
	soam := &secondaryOAM{}

	nmi := make(chan bool)

	return &PPU{
		VRAM: vram,
		Regs: newRegisters(nmi, oam, vram),

		OAM:     oam,
		sOAM:    soam,
		sprites: [8]sprite{},

		NMI: nmi,

		frame: newFrame(),
		disp:  disp,
	}
}

// Load connects a CPU's RAM to a ROM mapper and inits the CPU's PC to the
// ROM's reset vector.
func (ppu *PPU) Load(rom *ines.ROM) {
	ppu.VRAM.Mapper = rom.Mapper
	ppu.mirror = rom.Header.Mirroring
}

//TODO: Take note of oamaddr when performing DMA

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

			color: ppu.calcPixelValue(),
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
//
// incCoords will skip from (339,261) to (0,0) on odd ppu frames.
func (ppu *PPU) incCoords() {
	ppu.x++
	if ppu.x > 340 || (ppu.scanline == 261 && ppu.x == 339 || ppu.oddCycle) {
		ppu.x = 0

		ppu.scanline++
		if ppu.scanline > 261 {
			ppu.scanline = 0
			ppu.oddCycle = !ppu.oddCycle
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
	// TODO: Use OAMAddr as starting point of sprite 0 on x==65. See:
	// https://wiki.nesdev.com/w/index.php/PPU_registers#Values_during_rendering

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
		ptAddr := pt*ptSize + int(sprData[1])*16

		// Calculate line of the sprite to be displayed on the scanline
		sprLine := ppu.scanline - int(sprData[0])

		// Fetch sprite data
		dataLow := ppu.VRAM.Read(ptAddr + sprLine)
		dataHigh := ppu.VRAM.Read(ptAddr + sprLine + 8)

		// Invert sprite if horizontal invert bit of attribute byte is off
		if (sprData[2]>>6)&1 == 0 {
			dataLow = flipByte(dataLow)
			dataHigh = flipByte(dataHigh)
		}

		// Fill sprite slot with sprite data
		ppu.sprites[renderedSprNum] = sprite{
			low:     dataLow,
			high:    dataHigh,
			palette: sprData[2] & 3,

			x:       sprData[3],
			shifted: 0,

			priority:   (sprData[2]>>5)&1 == frontPriority,
			spriteZero: renderedSprNum == 0 && ppu.spriteZeroPresent,
		}
	} else {
		ppu.sprites[renderedSprNum] = nilSprite
	}
}

// calcPixelValue is called once per visible cycle (0 <= scanline < 240 &&
// 0 <= x < 256) and calculates the pixel value.
func (ppu *PPU) calcPixelValue() color.RGBA {
	bgr := ppu.calcBgrValue()
	spr := ppu.matchSprite()

	paletteAddr := ppu.muxPixel(bgr, spr)

	if spr.spriteZero && spr.getColor() != 0 && bgr&3 != 0 {
		// Set sprite 0 hit flag
		ppu.Regs.ppuStatus |= (1 << 6)
	}

	return Palette[paletteAddr]
}

// calcBgrValue calculates bgr value for current pixel.
//
// calcBgrValue return a nibble bgr value. The lower 2 bits (pattern) are
// fetched from a pattern table, and the upper 2 (colour) are fetched from an
// attribute table.
func (ppu *PPU) calcBgrValue() (bgr int) {
	// Return 0 if background rendering is disabled
	if ppu.Regs.ppuMask&(1<<3) == 0 {
		return 0
	}

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

// getATAddr returns the base address for nametable 1~4 based on a nametable
// number and a number representing mirroring mode.
func getNTAddr(ntNum int, mirroring int) int {
	// Assuming either horizontal or vertical mirroring
	if mirroring == ines.HorizontalMirroring {
		if ntNum == 0 || ntNum == 1 {
			return nt0Addr
		} else {
			return nt2Addr
		}
	} else {
		if ntNum == 0 || ntNum == 2 {
			return nt0Addr
		} else {
			return nt1Addr
		}
	}
}

// getPTAddr returns the base address for a pattern table according to bit 4 of
// PPUCTRL.
func (ppu *PPU) getPTAddr() int {
	if int((ppu.Regs.ppuCtrl>>4)&1) == 0 {
		return pt0Addr
	} else {
		return pt1Addr
	}
}

// getATAddr returns the base address of an attribute table of a nametable,
// based on the latter's base address.
func getATAddr(ntAddr int) int {
	return ntAddr + 0x3c0
}

// muxPixel multiplexes a background value and a sprite struct and returns the
// palette address for the pixel.
func (ppu *PPU) muxPixel(bgr int, spr sprite) (paletteAddr byte) {
	// No sprite was found for current pixel, display background
	if spr == nilSprite {
		if bgr&3 == 0 {
			// BGR palette addresses 0x0, 0x4, 0x8, 0xc are mirrored down to 0x0
			// when rendering
			bgr = 0
		}
		return ppu.VRAM.Read(bgrPaletteAddr + bgr)
	}

	// Background clear, display sprite
	if bgr&3 == 0 {
		return ppu.VRAM.Read(sprPaletteAddr + spr.getData())
	}

	// Sprite isn't clear and with priority, display sprite
	if spr.getColor() != 0 && spr.priority {
		return ppu.VRAM.Read(sprPaletteAddr + spr.getData())
	}

	// Display background
	if bgr&3 == 0 {
		// BGR palette addresses 0x0, 0x4, 0x8, 0xc are mirrored down to 0x0
		// when rendering
		bgr = 0
	}
	return ppu.VRAM.Read(bgrPaletteAddr + bgr)
}

// shiftSprites is called once per visible cycle and shifts each sprite's x
// coordinate by one.
//
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
	// Return nil sprite if sprite rendering is disabled
	if ppu.Regs.ppuMask&(1<<4) == 0 {
		return nilSprite
	}

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

// flipByte flips a byte. duh?
func flipByte(d byte) byte {
	d = ((d >> 1) & 0x55) | ((d & 0x55) << 1)
	d = ((d >> 2) & 0x33) | ((d & 0x33) << 2)
	d = ((d >> 4) & 0x0F) | ((d & 0x0F) << 4)
	return d
}
