package disass

import (
	"fmt"
	"os"

	"github.com/m4ntis/bones/ines"
)

func Example() {
	f, _ := os.Open("path_to_rom.nes")
	rom, _ := ines.Parse(f)

	d := Disassemble(rom.PrgROM)
	for _, inst := range d.Code {
		fmt.Printf("%04x: %s\n", inst.Addr, inst.Text)
	}
	// Output:
	// 8000: SEI
	// 8001: CLD
	// 8002: LDA #$10
	// ...
}
