package cpu

import (
	"sync"

	"github.com/m4ntis/bones/models"
)

type CPU struct {
	RAM *RAM
	Reg *Registers

	irq   bool
	nmi   bool
	reset bool

	interruptMux *sync.Mutex
}

func NewCPU() *CPU {
	var ram RAM
	return &CPU{
		RAM: &ram,
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
	// Load first 2 pages of PrgROM (not supporting mappers as of yet)
	copy(cpu.RAM.data[0x8000:0x8000+models.PRG_ROM_PAGE_SIZE], rom.PrgROM[0][:])
	copy(cpu.RAM.data[0x8000+models.PRG_ROM_PAGE_SIZE:0x8000+2*models.PRG_ROM_PAGE_SIZE],
		rom.PrgROM[1][:])
}

func (cpu *CPU) ExecNext() (cycles int) {
	op := OpCodes[*cpu.RAM.Fetch(cpu.Reg.PC)]

	cycles = op.cycles

	// We are doing this manually cus there are only 3 posibilities and writing
	// logic to describe this would be ugly IMO
	if op.Mode.ArgsLen == 1 {
		cycles += op.Exec(cpu, cpu.RAM.Fetch(cpu.Reg.PC+1))
	} else if op.Mode.ArgsLen == 2 {
		cycles += op.Exec(cpu, cpu.RAM.Fetch(cpu.Reg.PC+1), cpu.RAM.Fetch(cpu.Reg.PC+2))
	} else {
		cycles += op.Exec(cpu)
	}

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
	if cpu.Reg.I == CLEAR {
		cpu.interruptMux.Lock()
		cpu.irq = true
		cpu.interruptMux.Unlock()
	}
}

func (cpu *CPU) NMI() {
	if int(*cpu.RAM.Fetch(0x2000)&1<<7) == CLEAR {
		cpu.interruptMux.Lock()
		cpu.nmi = true
		cpu.interruptMux.Unlock()
	}
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
	cpu.Reg.PC = int(*cpu.RAM.Fetch(handlerAddr)) |
		int(*cpu.RAM.Fetch(handlerAddr + 1))<<8
	return
}

// Stack operations
func (cpu *CPU) push(b byte) {
	*cpu.RAM.Fetch(cpu.getStackAddr()) = b
	cpu.Reg.SP--
}

func (cpu *CPU) pull() byte {
	cpu.Reg.SP++
	return *cpu.RAM.Fetch(cpu.getStackAddr())
}

func (cpu *CPU) getStackAddr() int {
	return int(cpu.Reg.SP) | (1 << 8)
}
