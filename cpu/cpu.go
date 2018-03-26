package cpu

import (
	"sync"

	"github.com/m4ntis/bones/models"
)

type CPU struct {
	RAM *RAM
	Reg *Registers

	cycles int

	irq   bool
	nmi   bool
	reset bool

	interruptMux *sync.Mutex
}

func New(ram *RAM) *CPU {
	return &CPU{
		RAM: ram,
		Reg: &Registers{
			PC: 0x8000,
		},

		irq:   false,
		nmi:   false,
		reset: false,

		interruptMux: &sync.Mutex{},
	}
}

func (cpu *CPU) LoadROM(rom *models.ROM) {
	if len(rom.PrgROM) > 1 {
		// Load first 2 pages of PrgROM (not supporting mappers as of yet)
		copy(cpu.RAM.data[0x8000:0x8000+models.PrgROMPageSize], rom.PrgROM[0][:])
		copy(cpu.RAM.data[0x8000+models.PrgROMPageSize:0x8000+2*models.PrgROMPageSize],
			rom.PrgROM[1][:])
	} else {
		// If there is only one page of prg rom, load it to $c000 ~ $ffff
		copy(cpu.RAM.data[0x8000+models.PrgROMPageSize:0x8000+2*models.PrgROMPageSize],
			rom.PrgROM[0][:])
	}

	// Init pc to the reset handler addr
	cpu.Reg.PC = int(cpu.RAM.Read(0xfffc)) | int(cpu.RAM.Read(0xfffd))<<8
}

func (cpu *CPU) ExecNext() (cycles int) {
	op := OpCodes[cpu.RAM.Read(cpu.Reg.PC)]

	cycles = op.cycles

	// We are doing this manually cus there are only 3 posibilities and writing
	// logic to describe this would be ugly IMO
	if op.Mode.OpsLen == 1 {
		cycles += op.Exec(cpu, cpu.RAM.Read(cpu.Reg.PC+1))
	} else if op.Mode.OpsLen == 2 {
		cycles += op.Exec(cpu, cpu.RAM.Read(cpu.Reg.PC+1), cpu.RAM.Read(cpu.Reg.PC+2))
	} else {
		cycles += op.Exec(cpu)
	}

	// TODO: decrement cycles after 1786830 cycles
	cpu.cycles += cycles
	return cycles
}

// Interrupt Handling
func (cpu *CPU) HandleInterupts() {
	cpu.interruptMux.Lock()
	defer cpu.interruptMux.Unlock()
	if cpu.reset {
		cpu.interrupt(0xfffc)
		cpu.reset = false
	} else if cpu.nmi {
		cpu.interrupt(0xfffa)
		cpu.nmi = false
	} else if cpu.irq {
		cpu.interrupt(0xfffe)
		cpu.irq = false
	}
}

func (cpu *CPU) IRQ() {
	if cpu.Reg.I == Clear {
		cpu.interruptMux.Lock()
		cpu.irq = true
		cpu.interruptMux.Unlock()
	}
}

func (cpu *CPU) NMI() {
	cpu.interruptMux.Lock()
	cpu.nmi = true
	cpu.interruptMux.Unlock()
}

func (cpu *CPU) Reset() {
	cpu.interruptMux.Lock()
	cpu.reset = true
	cpu.interruptMux.Unlock()
}

func (cpu *CPU) interrupt(handlerAddr int) {
	// push PCH
	cpu.push(byte(cpu.Reg.PC >> 8))
	// push PCL
	cpu.push(byte(cpu.Reg.PC & 0xff))
	// push P
	cpu.push(cpu.Reg.GetP())

	cpu.Reg.I = 1

	// fetch PCL from $fffe and PCH from $ffff
	cpu.Reg.PC = int(cpu.RAM.Read(handlerAddr)) |
		int(cpu.RAM.Read(handlerAddr+1))<<8
	return
}

// Stack operations
func (cpu *CPU) push(d byte) {
	cpu.RAM.Write(cpu.getStackAddr(), d)
	cpu.Reg.SP--
}

func (cpu *CPU) pull() byte {
	cpu.Reg.SP++
	return cpu.RAM.Read(cpu.getStackAddr())
}

func (cpu *CPU) getStackAddr() int {
	return int(cpu.Reg.SP) | (1 << 8)
}
