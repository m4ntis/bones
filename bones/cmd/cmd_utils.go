package cmd

import (
	"fmt"
	"os"

	"github.com/m4ntis/bones/ines"
)

func openRom(cmdName string, args []string) *ines.ROM {
	if len(args) != 1 {
		fmt.Printf("Usage:\n  bones %s <romname>.nes\n", cmdName)
		os.Exit(1)
	}

	filename := args[0]
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening file %s:\n%s\n", filename, err)
		os.Exit(1)
	}

	rom, err := ines.Parse(f)
	if err != nil {
		fmt.Printf("Error parsing iNES file %s:\n%s\n", filename, err)
		os.Exit(1)
	}

	return rom
}
