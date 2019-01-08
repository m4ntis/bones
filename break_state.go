package bones

import (
	"github.com/m4ntis/bones/asm"
	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/ppu"
)

// BreakState describes the state of the NES and the running programme when the
// cpu breaks.
type BreakState struct {
	Reg  *cpu.Registers
	RAM  *cpu.RAM
	VRAM *ppu.VRAM

	Code  asm.Code
	PCIdx int

	Err error
}
