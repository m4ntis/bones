package mapper

import "github.com/pkg/errors"

type Mapper001 struct {
	sram bool

	prgROM []PrgROMPage
	chrROM []ChrROMPage

	prgRAM [PrgRAMSize]byte
	chrRAM [ChrRAMSize]byte

	sr         byte
	writeCount int

	ctrl byte
	chr0 byte
	chr1 byte
	prg  byte

	booted bool
}

func (m *Mapper001) SetSram(b bool) {
	m.sram = b
}

func (m *Mapper001) Read(addr int) (d byte, err error) {
	if addr >= 0 && addr < 0x2000 {
		if m.sram {
			return m.chrRAM[addr], nil
		}

		page, index := m.decodeChrROMAddr(addr)
		return m.chrROM[page][index], nil
	} else if addr >= 0x6000 && addr < 0x10000 {
		return m.readPrgROM(addr), nil
	}

	return 0, errors.Errorf("invalid mapper reading addr %04x", addr)
}

// Write handles writing to an address mapped by the mapper.
//
// Writing to a mapper address is used for writing to RAM areas,
// as well as writing to registers controlling the mapper.
//
// TODO: Ignore consecutive writes to the mapper, a cycle immediately the
// the previous.
func (m *Mapper001) Write(addr int, d byte) error {
	// Write to prg ram
	if addr >= 0x6000 && addr < 0x8000 {
		addr -= 0x6000
		m.prgRAM[addr] = d
	}

	if addr < 0x2000 {
		m.chrRAM[addr] = d
	}

	if addr >= 0x8000 && addr < 0x10000 {
		// Bit 7 of data is set, reset shift register
		if d&128 == 128 {
			m.sr = 0
			m.writeCount = 0
		} else {
			// First 5 writes to 0x8000 ~ 0xffff with bit 7 of data clear are
			// shifted onto the internal shift register
			m.sr = m.sr << 1
			m.sr |= d & 1
			m.writeCount++

			// 5th write flushes the data to an internal register
			//
			// Register number (identifier of the register to write to) is
			// determined by addr bit 13 and 14:
			// 0 -> write to ctrl reg
			// 1 -> write to chr0 reg
			// 2 -> write to chr1 reg
			// 3 -> write to prg reg
			if m.writeCount == 5 {
				regN := addr >> 13 & 3

				switch regN {
				case 0:
					m.booted = true
					m.ctrl = m.sr
				case 1:
					m.chr0 = m.sr
				case 2:
					m.chr1 = m.sr
				case 3:
					m.prg = m.sr
				}

				// Reset shift reg
				m.sr = 0
				m.writeCount = 0
			}
		}
	}

	return nil
}

func (m *Mapper001) Observe(addr int) (d byte, err error) {
	// There is no side effect to reading from mapper 001
	return m.Read(addr)
}

func (m *Mapper001) Populate(prgROM []PrgROMPage, chrROM []ChrROMPage) {
	// If only one page of prg rom, duplicate it
	if len(prgROM) == 1 {
		var pageCopy PrgROMPage
		copy(pageCopy[:], prgROM[0][:])

		prgROM = append(prgROM, pageCopy)
	}

	m.prgROM = prgROM
	m.chrROM = chrROM
}

func (m *Mapper001) GetPRGRom() []PrgROMPage {
	return m.prgROM
}

func (m *Mapper001) readPrgROM(addr int) byte {
	if addr >= 0x6000 && addr < 0x8000 {
		addr -= 0x6000
		return m.prgRAM[addr]
	}

	addr -= 0x8000
	page, index := m.decodePrgROMAddr(addr)
	return m.prgROM[page][index]
}

func (m *Mapper001) decodePrgROMAddr(addr int) (page, index int) {
	// Default to upper bank fixed
	if !m.booted && addr > PrgROMPageSize {
		return len(m.prgROM) - 1, addr % PrgROMPageSize
	}

	index = addr % PrgROMPageSize

	// Test 32kb prg mode
	if m.ctrl&2 == 0 {
		page = (int(flip_reg(m.prg)) & 7) / 2
		page += addr / PrgROMPageSize
	} else {
		// Test if the accessed bank is swappable
		if (addr < PrgROMPageSize && m.ctrl&4 == 4) || (addr >= PrgROMPageSize && m.ctrl&4 == 0) {
			page = int(flip_reg(m.prg)) & 15
		} else if addr >= PrgROMPageSize {
			// If upper bank is fixed, fix it to last page
			page = len(m.prgROM) - 1
		}
	}

	return page, index
}

func (m *Mapper001) decodeChrROMAddr(addr int) (page, index int) {
	// Test 8kb chr mode
	if m.ctrl&1 == 0 {
		// Ignoring MSB of page num (bit 5) in 8kb mode
		page = int(flip_reg(m.chr0)) & 0xf
		index = addr % ChrROMPageSize
	} else {
		index = addr % (ChrROMPageSize / 2)

		// Test for first or second chr 4k bank access
		if addr < 0x1000 {
			page = int(flip_reg(m.chr0))
		} else {
			page = int(flip_reg(m.chr1))
		}

		index += (page % 2) * (ChrROMPageSize / 2)
		page /= 2
	}

	return page, index
}

func flip_reg(d byte) byte {
	d = ((d >> 1) & 0x55) | ((d & 0x55) << 1)
	d = ((d >> 2) & 0x33) | ((d & 0x33) << 2)
	d = ((d >> 4) & 0x0F) | ((d & 0x0F) << 4)
	return d >> 3
}
