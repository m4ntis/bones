package ppu

import (
	"image/color"
)

var (
	Palette = map[byte]color.RGBA{
		0:  color.RGBA{124, 124, 124, 255},
		1:  color.RGBA{0, 0, 252, 255},
		2:  color.RGBA{0, 0, 188, 255},
		3:  color.RGBA{68, 40, 188, 255},
		4:  color.RGBA{148, 0, 132, 255},
		5:  color.RGBA{168, 0, 32, 255},
		6:  color.RGBA{168, 16, 0, 255},
		7:  color.RGBA{136, 20, 0, 255},
		8:  color.RGBA{80, 48, 0, 255},
		9:  color.RGBA{0, 120, 0, 255},
		10: color.RGBA{0, 104, 0, 255},
		11: color.RGBA{0, 88, 0, 255},
		12: color.RGBA{0, 64, 88, 255},
		13: color.RGBA{0, 0, 0, 255},
		14: color.RGBA{0, 0, 0, 255},
		15: color.RGBA{0, 0, 0, 255},
		16: color.RGBA{188, 188, 188, 255},
		17: color.RGBA{0, 120, 248, 255},
		18: color.RGBA{0, 88, 248, 255},
		19: color.RGBA{104, 68, 252, 255},
		20: color.RGBA{216, 0, 204, 255},
		21: color.RGBA{228, 0, 88, 255},
		22: color.RGBA{248, 56, 0, 255},
		23: color.RGBA{228, 92, 16, 255},
		24: color.RGBA{172, 124, 0, 255},
		25: color.RGBA{0, 184, 0, 255},
		26: color.RGBA{0, 168, 0, 255},
		27: color.RGBA{0, 168, 68, 255},
		28: color.RGBA{0, 136, 136, 255},
		29: color.RGBA{0, 0, 0, 255},
		30: color.RGBA{0, 0, 0, 255},
		31: color.RGBA{0, 0, 0, 255},
		32: color.RGBA{248, 248, 248, 255},
		33: color.RGBA{60, 188, 252, 255},
		34: color.RGBA{104, 136, 252, 255},
		35: color.RGBA{152, 120, 248, 255},
		36: color.RGBA{248, 120, 248, 255},
		37: color.RGBA{248, 88, 152, 255},
		38: color.RGBA{248, 120, 88, 255},
		39: color.RGBA{252, 160, 68, 255},
		40: color.RGBA{248, 184, 0, 255},
		41: color.RGBA{184, 248, 24, 255},
		42: color.RGBA{88, 216, 84, 255},
		43: color.RGBA{88, 248, 152, 255},
		44: color.RGBA{0, 232, 216, 255},
		45: color.RGBA{120, 120, 120, 255},
		46: color.RGBA{0, 0, 0, 255},
		47: color.RGBA{0, 0, 0, 255},
		48: color.RGBA{252, 252, 252, 255},
		49: color.RGBA{164, 228, 252, 255},
		50: color.RGBA{184, 184, 248, 255},
		51: color.RGBA{216, 184, 248, 255},
		52: color.RGBA{248, 184, 248, 255},
		53: color.RGBA{248, 164, 192, 255},
		54: color.RGBA{240, 208, 176, 255},
		55: color.RGBA{252, 224, 168, 255},
		56: color.RGBA{248, 216, 120, 255},
		57: color.RGBA{216, 248, 120, 255},
		58: color.RGBA{184, 248, 184, 255},
		59: color.RGBA{184, 248, 216, 255},
		60: color.RGBA{0, 252, 252, 255},
		61: color.RGBA{248, 216, 248, 255},
		62: color.RGBA{0, 0, 0, 255},
		63: color.RGBA{0, 0, 0, 255},
	}
)