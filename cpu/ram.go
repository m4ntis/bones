package cpu

import "sync"

const (
	ZERO_PAGE_BEGIN_IDX           = 0x0
	STACK_BEGIN_IDX               = 0x100
	RAM_BEGIN_IDX                 = 0x200
	RAM_MIRROR_BEGIN_IDX          = 0x800
	LOWER_IO_REG_BEGIN_IDX        = 0x2000
	LOWER_IO_REG_MIRROR_BEGIN_IDX = 0x2008
	UPPER_IO_REG_BEGIN_IDX        = 0x4000
	EXPANSION_ROM_BEGIN_IDX       = 0x4020
	SRAM_BEGIN_IDX                = 0x4000
	PRG_ROM_LOWER_BEGIN_IDX       = 0x8000
	PRG_ROM_UPPER_BEGIN_IDX       = 0xc000
	RAM_SIZE                      = 0x10000
)

type RAM struct {
	data [RAM_SIZE]byte
	mux  sync.RWMutex
}

func getIndex(index int) int {
	if index < 0 || index > RAM_SIZE {
		panic("RAM accessing index out of range")
	}

	if index >= RAM_MIRROR_BEGIN_IDX && index < LOWER_IO_REG_BEGIN_IDX {
		return index % 0x800
	}
	if index >= LOWER_IO_REG_MIRROR_BEGIN_IDX &&
		index < UPPER_IO_REG_BEGIN_IDX {
		return (index-LOWER_IO_REG_BEGIN_IDX)%0x8 +
			LOWER_IO_REG_BEGIN_IDX
	}
	return index
}

func (r *RAM) Read(index int) byte {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return r.data[getIndex(index)]
}

func (r *RAM) Write(index int, val byte) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.data[getIndex(index)] = val
}
