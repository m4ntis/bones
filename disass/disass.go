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

// addrTable maps an address in RAM to it's index in the code
type addrTable map[int]int

type Disassembly struct {
	Code Code

	addrTable addrTable
}

func (d Disassembly) IndexOf(addr int) int {
	i, ok := d.addrTable[addr]

	if ok {
		return i
	}
	return -1
}

func Disassemble(prgROM []models.PrgROMPage) Disassembly {
	asm := genContiguousAsm(prgROM)
	code := disassemble(asm)
	addrTable := genAddrTable(code)

	// Create an addressing table for the code
	return Disassembly{
		Code:      code,
		addrTable: addrTable,
	}
}

func genContiguousAsm(prgROM []models.PrgROMPage) []byte {
	asm := make([]byte, 0)

	for _, page := range prgROM {
		asm = append(asm, page[:]...)
	}
	return asm
}

func disassemble(asm []byte) Code {
	loadAddr := 0x8000

	// If only single page of prg rom, it is loaded to $c000 instead of the
	// usual $8000
	if len(asm) == models.PrgROMPageSize {
		loadAddr = 0xc000
	}

	code := Code(make([]Instruction, 0))
	for i := 0; i < len(asm); i++ {
		op, ok := cpu.OpCodes[asm[i]]

		var inst Instruction
		if !ok {
			inst = Instruction{
				Addr: i + loadAddr,
				Code: asm[i : i+1],
				Text: fmt.Sprintf(".byte %02x", asm[i]),
			}
		} else {
			inst = Instruction{
				Addr: i + loadAddr,
				Code: asm[i : i+1+op.Mode.OpsLen],
				Text: fmt.Sprintf("%s %s", op.Name,
					op.Mode.Format(asm[i+1:i+1+op.Mode.OpsLen])),
			}
		}

		code = append(code, inst)
		i += op.Mode.OpsLen
	}

	return code
}

func genAddrTable(code Code) addrTable {
	addrTable := addrTable{}

	for i, inst := range code {
		addrTable[inst.Addr] = i
	}
	return addrTable
}
