package disass

type Instruction struct {
	Data []byte
	Text string
}

type Disassebly []Instruction

// AddressTable maps an address in ram to it's index in the disassembly
type AddressTable map[int]int
