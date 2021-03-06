// Package asm provides disassembly fucntionality for NES roms
package asm

import (
	"fmt"

	"github.com/m4ntis/bones/cpu"
	"github.com/m4ntis/bones/ines"
)

// Instruction represents a single logical instruction in an NES rom.
//
type Instruction struct {
	// Addr holds the instuction's location in the rom
	Addr int
	// Code contains the coniguous bytes in rom representing the instruction
	// (opcode + operands)
	Code []byte
	// Text is the formatted textual representation of the instruction
	Text string
}

// Code is a slice of instructions, representing a program written for the mos
// 6502.
type Code []Instruction

// addrTable maps an address in RAM to it's index in the code
type addrTable map[int]int

// Disassembly contains a code object and logic for translation between an
// address and it's location in the code. This is the type returned when
// disassembling a rom.
type Disassembly struct {
	Code Code

	addrTable addrTable
}

// IndexOf is used to get the index of an address in the code of a disassembly.
// This is because you can't infer the location by using the address alone, as
// the operand length of an opcode varies.
//
// IndexOf returns -1 if the address isn't a beginning of an instruction
func (d Disassembly) IndexOf(addr int) int {
	i, ok := d.addrTable[addr]

	if ok {
		return i
	}
	return -1
}

// DisassembleROM takes a ROM and returnes it's disassembled code
func DisassembleROM(rom *ines.ROM) Disassembly {
	prgROM := rom.Mapper.GetPRGRom()

	asm := genContiguousAsm(prgROM)
	code := disassemble(asm)
	addrTable := genAddrTable(code)

	// Create an addressing table for the code
	return Disassembly{
		Code:      code,
		addrTable: addrTable,
	}
}

// DisassembleRAM disassembles a few instructions, starting from a point in RAM
//
// The disassembly starts from addr and disassembles count instructions
func DisassembleRAM(ram *cpu.RAM, addr, count int) Code {
	code := Code(make([]Instruction, count))

	for i := 0; i < count; i++ {
		code[i] = disassOne(ram, addr)
		addr += len(code[i].Code)
	}

	return code
}

func genContiguousAsm(prgROM []ines.PrgROMPage) []byte {
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
	if len(asm) == ines.PrgROMPageSize {
		loadAddr = 0xc000
	}

	code := Code(make([]Instruction, 0))
	for i := 0; i < len(asm); i++ {
		op := cpu.OpCodes[asm[i]]

		var inst Instruction
		if op.Name == "" {
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

func disassOne(ram *cpu.RAM, addr int) Instruction {
	op := cpu.OpCodes[ram.MustRead(addr)]

	if op.Name == "" {
		code := readSliceFromRAM(ram, addr, 1)
		return Instruction{
			Addr: addr,
			Code: code,
			Text: fmt.Sprintf(".byte %02x", code[0]),
		}
	}

	code := readSliceFromRAM(ram, addr, op.Mode.OpsLen+1)
	return Instruction{
		Addr: addr,
		Code: code,
		Text: fmt.Sprintf("%s %s", op.Name,
			op.Mode.Format(code[1:])),
	}
}

func readSliceFromRAM(ram *cpu.RAM, addr, n int) []byte {
	d := make([]byte, n)
	for i := 0; i < n; i++ {
		d[i] = ram.MustRead(addr + i)
	}

	return d
}

func genAddrTable(code Code) addrTable {
	addrTable := addrTable{}

	for i, inst := range code {
		addrTable[inst.Addr] = i
	}
	return addrTable
}
