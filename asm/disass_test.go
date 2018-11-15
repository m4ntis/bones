package asm

import (
	"fmt"
	"os"

	"github.com/m4ntis/bones/ines"
)

func Example() {
	f, err := os.Open("path_to_rom.nes")
	if err != nil {
		return
	}
	rom, err := ines.Parse(f)
	if err != nil {
		return
	}

	d := DisassembleROM(rom)
	for _, inst := range d.Code {
		fmt.Printf("%04x: %s\n", inst.Addr, inst.Text)
	}
	// 8000: SEI
	// 8001: CLD
	// 8002: LDA #$10
	// ...
}
