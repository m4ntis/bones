package disass

import (
	"fmt"

	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/models"
)

type Instruction struct {
	Addr int
	Code []byte
	Text string
}

type Code []Instruction

// AddrTable maps an address in RAM to it's index in the code
type AddrTable map[int]int

type Disassembly struct {
	Code Code

	addrTable AddrTable
}

func (d *Disassembly) IndexOf(addr int) int {
	return d.addrTable[addr]
}

func Disassemble(prgROM []models.PrgROMPage) Disassembly {
	// Create assembly - a byte slice of the whole program rom
	assembly := make([]byte, 0)
	for _, page := range prgROM {
		assembly = append(assembly, page[:]...)
	}

	// Extract code - a list of parsed instructions from the assembly
	code := Code(make([]Instruction, 0))
	for i := 0; i < len(assembly); i++ {
		op, ok := cpu.OpCodes[assembly[i]]

		var inst Instruction
		if !ok {
			inst = Instruction{
				Addr: i,
				Code: assembly[i : i+1],
				Text: fmt.Sprintf(".byte %02x", assembly[i]),
			}
		} else {

			inst = Instruction{
				Addr: i,
				Code: assembly[i : i+1+op.Mode.ArgsLen],
				Text: fmt.Sprintf("%s %s", op.Name,
					op.Mode.Format(assembly[i+1:i+1+op.Mode.ArgsLen])),
			}
		}

		code = append(code, inst)
		i += op.Mode.ArgsLen
	}

	// Create an addressing table for the code
	addrTable := AddrTable{}

	for i, inst := range code {
		addrTable[inst.Addr] = i
	}

	return Disassembly{
		Code:      code,
		addrTable: addrTable,
	}
}
