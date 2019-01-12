package dbg

import (
	"fmt"

	"github.com/m4ntis/bones"
)

func List(b bones.Break) {
	for i, inst := range b.Code {
		if i == b.PCIdx {
			fmt.Printf("=> %04x: %s\n", inst.Addr, inst.Text)
			continue
		}
		fmt.Printf("   %04x: %s\n", inst.Addr, inst.Text)
	}

	if b.Err != nil {
		fmt.Println(b.Err)
	}
}
