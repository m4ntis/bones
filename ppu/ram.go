package ppu

const (
	PT0_IDX            = 0x0
	PT1_IDX            = 0x1000
	NT0_IDX            = 0x2000
	AT0_IDX            = 0x23c0
	NT1_IDX            = 0x2400
	AT1_IDX            = 0x27c0
	NT2_IDX            = 0x2800
	AT2_IDX            = 0x2bc0
	NT3_IDX            = 0x2c00
	AT3_IDX            = 0x2fc0
	TABLES_MIRROR_IDX  = 0x3000
	IMG_PALETTE_IDX    = 0x3f00
	SPRITE_PALETTE_IDX = 0x3f10
	PALETTE_MIRROR_IDX = 0x3f20
	RAM_MIRROR_IDX     = 0x4000
	RAM_SIZE           = 0x10000
)

type RAM struct {
	data [RAM_SIZE]byte
}

func getIndex(index int) int {
	if index < 0 || index > RAM_SIZE {
		panic("RAM accessing index out of range")
	}

	if index >= TABLES_MIRROR_IDX && index < IMG_PALETTE_IDX {
		return index - 0x1000
	}
	if index >= PALETTE_MIRROR_IDX && index < RAM_MIRROR_IDX {
		return (index-IMG_PALETTE_IDX)%0x20 + IMG_PALETTE_IDX
	}
	if index >= RAM_MIRROR_IDX {
		return index % 0x4000
	}
	return index
}

func (r *RAM) Fetch(index int) *byte {
	return &r.data[getIndex(index)]
}
