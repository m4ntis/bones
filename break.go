package bones

import (
	"github.com/m4ntis/bones/asm"
)

// Break describes the state of the NES and the running programme when the
// cpu breaks.
type Break struct {
	Code  asm.Code
	PCIdx int

	Err error
}
