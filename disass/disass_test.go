package disass

import (
	"fmt"

	"github.com/m4ntis/bones/disass"
)

func Example() {
	d := disass.Disassemble(rom.PrgROM)
	for _, inst := range d.Code {
		fmt.Printf("%04x: %s\n", inst.Addr, inst.Text)
	}
	// Output:
	// 8000: SEI
	// 8001: CLD
	// 8002: LDA #$10
	// ...
}
