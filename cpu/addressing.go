package cpu

import (
	"fmt"

	"github.com/pkg/errors"
)

// TODO: Consider optimal method of handling PC register overflow as a result of
// bad NES code / faulty ROM parsing.
//
// For example, AbsoluteX/Y addressing like 0xFFFF + X/Y (where X/Y > 0), would
// result in an error raised in RAM accessing, ONLY when and IF the operation
// calls Read / Write of the operand.
//
// This holds both for RAM Operand instanciation and for RAM reading of an
// address (IE indirect $ffff, adh would be in $0000).
//
// Optimally BoNES should panic if the user is using `bones run` on the ROM, but
// perhaps just return an error msg when doing something like debugging.

// AddressingMode defines one of the mos 6502's ways of addressing operands.
//
// Each addressing mode is responsible of fetching the operands in it's way,
// and calling the operation with them.
//
// The bool tells whether there needs to be a page boundry check.
//
// The addressing mode returns the amount of extra cycles the opration took
// because of page boundry or branching.
type AddressingMode struct {
	Name   string
	OpsLen int
	Format func([]byte) string

	Address func(*CPU, Operation, bool, ...byte) (int, error)
}

var (
	Implied = AddressingMode{
		Name:   "Implied",
		OpsLen: 0,
		Format: func(ops []byte) string { return "" },

		Address: func(cpu *CPU, op Operation, pageBoundryCheck bool,
			ops ...byte) (extraCycles int, err error) {
			cpu.Reg.PC++
			op(cpu, NilOperand{})
			return extraCycles, nil
		},
	}

	Accumulator = AddressingMode{
		Name:   "Accumulator",
		OpsLen: 0,
		Format: func(ops []byte) string { return "A" },

		Address: func(cpu *CPU, op Operation, pageBoundryCheck bool,
			ops ...byte) (extraCycles int, err error) {
			op(cpu, RegOperand{Reg: &cpu.Reg.A})
			cpu.Reg.PC++
			return extraCycles, nil
		},
	}

	Immediate = AddressingMode{
		Name:   "Immediate",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("#$%02x", ops[0]) },

		Address: func(cpu *CPU, op Operation, pageBoundryCheck bool,
			ops ...byte) (extraCycles int, err error) {
			op(cpu, ConstOperand{D: ops[0]})
			cpu.Reg.PC += 2
			return extraCycles, nil
		},
	}

	ZeroPage = AddressingMode{
		Name:   "ZeroPage",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x", ops[0]) },

		Address: func(cpu *CPU, op Operation, pageBoundryCheck bool,
			ops ...byte) (extraCycles int, err error) {
			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: int(ops[0])})
			cpu.Reg.PC += 2
			return extraCycles, nil
		},
	}

	ZeroPageX = AddressingMode{
		Name:   "ZeroPageX",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x, X", ops[0]) },

		Address: func(cpu *CPU, op Operation, pageBoundryCheck bool,
			ops ...byte) (extraCycles int, err error) {
			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: int(ops[0] + cpu.Reg.X)})
			cpu.Reg.PC += 2
			return extraCycles, nil
		},
	}

	ZeroPageY = AddressingMode{
		Name:   "ZeroPageY",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x, Y", ops[0]) },

		Address: func(cpu *CPU, op Operation, pageBoundryCheck bool,
			ops ...byte) (extraCycles int, err error) {
			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: int(ops[0] + cpu.Reg.Y)})
			cpu.Reg.PC += 2
			return extraCycles, nil
		},
	}

	Relative = AddressingMode{
		Name:   "Relative",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x", ops[0]) },

		Address: func(cpu *CPU, op Operation, pageBoundryCheck bool,
			ops ...byte) (extraCycles int, err error) {
			extraCycles = op(cpu, ConstOperand{D: ops[0]})
			cpu.Reg.PC += 2
			return extraCycles, nil
		},
	}

	Absolute = AddressingMode{
		Name:   "Absolute",
		OpsLen: 2,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x%02x", ops[1], ops[0]) },

		Address: func(cpu *CPU, op Operation, pageBoundryCheck bool,
			ops ...byte) (extraCycles int, err error) {
			// We inc this beforehand so that JSR wont be incremented after
			// execution
			cpu.Reg.PC += 3

			addr := int(ops[0]) | int(ops[1])<<8
			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: addr})
			return extraCycles, nil
		},
	}

	AbsoluteX = AddressingMode{
		Name:   "AbsoluteX",
		OpsLen: 2,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x%02x, X", ops[1], ops[0]) },

		Address: func(cpu *CPU, op Operation, pageBoundryCheck bool,
			ops ...byte) (extraCycles int, err error) {
			addr := int(ops[0]) | int(ops[1])<<8
			xAddr := addr + int(cpu.Reg.X)

			if addr/256 != xAddr/256 && pageBoundryCheck {
				extraCycles++
			}

			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: xAddr})
			cpu.Reg.PC += 3
			return extraCycles, nil
		},
	}

	AbsoluteY = AddressingMode{
		Name:   "AbsoluteY",
		OpsLen: 2,
		Format: func(ops []byte) string { return fmt.Sprintf("$%02x%02x, Y", ops[1], ops[0]) },

		Address: func(cpu *CPU, op Operation, pageBoundryCheck bool,
			ops ...byte) (extraCycles int, err error) {
			addr := int(ops[0]) | int(ops[1])<<8
			yAddr := addr + int(cpu.Reg.Y)

			if addr/256 != yAddr/256 && pageBoundryCheck {
				extraCycles++
			}

			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: yAddr})
			cpu.Reg.PC += 3
			return extraCycles, nil
		},
	}

	Indirect = AddressingMode{
		Name:   "Indirect",
		OpsLen: 2,
		Format: func(ops []byte) string { return fmt.Sprintf("($%02x%02x)", ops[1], ops[0]) },

		Address: func(cpu *CPU, op Operation, pageBoundryCheck bool,
			ops ...byte) (extraCycles int, err error) {
			addr := int(ops[0]) + int(ops[1])<<8

			adl, err := cpu.RAM.Read(addr)
			if err != nil {
				return 0, errors.Wrap(err, "Couldn't read indirected address")
			}

			adh, err := cpu.RAM.Read(addr + 1)
			if err != nil {
				return 0, errors.Wrap(err, "Couldn't read indirected address")
			}

			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: int(adl) | int(adh)<<8})
			return extraCycles, nil
		},
	}

	IndirectX = AddressingMode{
		Name:   "IndirectX",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("($%02x, X)", ops[0]) },

		Address: func(cpu *CPU, op Operation, pageBoundryCheck bool,
			ops ...byte) (extraCycles int, err error) {
			addr := int(ops[0] + cpu.Reg.X)

			adl, err := cpu.RAM.Read(addr)
			if err != nil {
				return 0, errors.Wrap(err, "Couldn't read indirected address")
			}

			adh, err := cpu.RAM.Read((addr + 1) % 0x100)
			if err != nil {
				return 0, errors.Wrap(err, "Couldn't read indirected address")
			}

			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: int(adl) | int(adh)<<8})
			cpu.Reg.PC += 2
			return extraCycles, nil
		},
	}

	IndirectY = AddressingMode{
		Name:   "IndirectY",
		OpsLen: 1,
		Format: func(ops []byte) string { return fmt.Sprintf("($%02x), Y", ops[0]) },

		Address: func(cpu *CPU, op Operation, pageBoundryCheck bool,
			ops ...byte) (extraCycles int, err error) {
			addr := int(ops[0])

			adl, err := cpu.RAM.Read(addr)
			if err != nil {
				return 0, errors.Wrap(err, "Couldn't read indirected address")
			}

			adh, err := cpu.RAM.Read((addr + 1) % 0x100)
			if err != nil {
				return 0, errors.Wrap(err, "Couldn't read indirected address")
			}

			fetched := int(adl) | int(adh)<<8
			if fetched/256 != (fetched+int(cpu.Reg.Y))/256 && pageBoundryCheck {
				extraCycles++
			}
			fetched += int(cpu.Reg.Y)

			op(cpu, RAMOperand{RAM: cpu.RAM, Addr: fetched})
			cpu.Reg.PC += 2
			return extraCycles, nil
		},
	}
)
