package bones_test

import (
	"os"

	"github.com/m4ntis/bones"
	"github.com/m4ntis/bones/ines"
	"github.com/m4ntis/bones/io"
)

const (
	displayFPS = false
	scale      = 4.0
	filename   = "tetris.nes"
)

func Example_runROM() {
	// Open and parse ROM file
	f, err := os.Open(filename)
	panicOnErr(err)

	rom, err := ines.Parse(f)
	panicOnErr(err)

	// Init I/O components
	ctrl := new(io.Controller)
	disp := io.NewDisplay(ctrl, displayFPS, scale)

	// Init NES
	n := bones.New(disp, ctrl, bones.ModeRun)
	n.Load(rom)

	// Run ROM and display
	go n.Start()
	disp.Run()
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
