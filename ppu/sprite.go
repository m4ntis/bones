package ppu

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
