package cpu

var OpCodes = map[byte]OpCode{
	0x69: OpCode{
		name: "ADC",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Immediate,
		oper: ADC,
	},
	0x65: OpCode{
		name: "ADC",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: ADC,
	},
	0x75: OpCode{
		name: "ADC",

		cycles:           4,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: ADC,
	},
	0x6d: OpCode{
		name: "ADC",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: ADC,
	},
	0x7d: OpCode{
		name: "ADC",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteX,
		oper: ADC,
	},
	0x79: OpCode{
		name: "ADC",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteY,
		oper: ADC,
	},
	0x61: OpCode{
		name: "ADC",

		cycles:           6,
		pageBoundryCheck: false,

		mode: IndirectX,
		oper: ADC,
	},
	0x71: OpCode{
		name: "ADC",

		cycles:           5,
		pageBoundryCheck: true,

		mode: IndirectY,
		oper: ADC,
	},
	0x29: OpCode{
		name: "AND",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Immediate,
		oper: AND,
	},
	0x25: OpCode{
		name: "AND",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: AND,
	},
	0x35: OpCode{
		name: "AND",

		cycles:           4,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: AND,
	},
	0x2d: OpCode{
		name: "AND",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: AND,
	},
	0x3d: OpCode{
		name: "AND",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteX,
		oper: AND,
	},
	0x39: OpCode{
		name: "AND",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteY,
		oper: AND,
	},
	0x21: OpCode{
		name: "AND",

		cycles:           6,
		pageBoundryCheck: false,

		mode: IndirectX,
		oper: AND,
	},
	0x31: OpCode{
		name: "AND",

		cycles:           5,
		pageBoundryCheck: true,

		mode: IndirectY,
		oper: AND,
	},
	0x0a: OpCode{
		name: "ASL",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Accumulator,
		oper: ASL,
	},
	0x06: OpCode{
		name: "ASL",

		cycles:           5,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: ASL,
	},
	0x16: OpCode{
		name: "ASL",

		cycles:           6,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: ASL,
	},
	0x0e: OpCode{
		name: "ASL",

		cycles:           6,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: ASL,
	},
	0x1e: OpCode{
		name: "ASL",

		cycles:           7,
		pageBoundryCheck: false,

		mode: AbsoluteX,
		oper: ASL,
	},
	0x90: OpCode{
		name: "BCC",

		cycles:           2,
		pageBoundryCheck: true,

		mode: Relative,
		oper: BCC,
	},
	0xb0: OpCode{
		name: "BCS",

		cycles:           2,
		pageBoundryCheck: true,

		mode: Relative,
		oper: BCS,
	},
	0xf0: OpCode{
		name: "BEQ",

		cycles:           2,
		pageBoundryCheck: true,

		mode: Relative,
		oper: BEQ,
	},
	0x24: OpCode{
		name: "BIT",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: BIT,
	},
	0x2c: OpCode{
		name: "BIT",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: BIT,
	},
	0x30: OpCode{
		name: "BMI",

		cycles:           2,
		pageBoundryCheck: true,

		mode: Relative,
		oper: BMI,
	},
	0xd0: OpCode{
		name: "BNE",

		cycles:           2,
		pageBoundryCheck: true,

		mode: Relative,
		oper: BNE,
	},
	0x10: OpCode{
		name: "BPL",

		cycles:           2,
		pageBoundryCheck: true,

		mode: Relative,
		oper: BPL,
	},
	0x00: OpCode{
		name: "BRK",

		cycles:           7,
		pageBoundryCheck: false,

		mode: Implied,
		oper: BRK,
	},
	0x50: OpCode{
		name: "BVC",

		cycles:           2,
		pageBoundryCheck: true,

		mode: Relative,
		oper: BVC,
	},
	0x70: OpCode{
		name: "BVS",

		cycles:           2,
		pageBoundryCheck: true,

		mode: Relative,
		oper: BVS,
	},
	0x18: OpCode{
		name: "CLC",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: CLC,
	},
	0xd8: OpCode{
		name: "CLD",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: CLD,
	},
	0x58: OpCode{
		name: "CLI",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: CLI,
	},
	0xb8: OpCode{
		name: "CLV",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: CLV,
	},
	0xc9: OpCode{
		name: "CMP",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Immediate,
		oper: CMP,
	},
	0xc5: OpCode{
		name: "CMP",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: CMP,
	},
	0xd5: OpCode{
		name: "CMP",

		cycles:           4,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: CMP,
	},
	0xcd: OpCode{
		name: "CMP",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: CMP,
	},
	0xdd: OpCode{
		name: "CMP",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteX,
		oper: CMP,
	},
	0xd9: OpCode{
		name: "CMP",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteY,
		oper: CMP,
	},
	0xc1: OpCode{
		name: "CMP",

		cycles:           6,
		pageBoundryCheck: false,

		mode: IndirectX,
		oper: CMP,
	},
	0xd1: OpCode{
		name: "CMP",

		cycles:           5,
		pageBoundryCheck: true,

		mode: IndirectY,
		oper: CMP,
	},
	0xe0: OpCode{
		name: "CPX",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Immediate,
		oper: CPX,
	},
	0xe4: OpCode{
		name: "CPX",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: CPX,
	},
	0xec: OpCode{
		name: "CPX",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: CPX,
	},
	0xc0: OpCode{
		name: "CPY",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Immediate,
		oper: CPY,
	},
	0xc4: OpCode{
		name: "CPY",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: CPY,
	},
	0xcc: OpCode{
		name: "CPY",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: CPY,
	},
	0xc6: OpCode{
		name: "DEC",

		cycles:           5,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: DEC,
	},
	0xd6: OpCode{
		name: "DEC",

		cycles:           6,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: DEC,
	},
	0xce: OpCode{
		name: "DEC",

		cycles:           6,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: DEC,
	},
	0xde: OpCode{
		name: "DEC",

		cycles:           7,
		pageBoundryCheck: false,

		mode: AbsoluteX,
		oper: DEC,
	},
	0xca: OpCode{
		name: "DEX",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: DEX,
	},
	0x88: OpCode{
		name: "DEY",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: DEY,
	},
	0x49: OpCode{
		name: "EOR",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Immediate,
		oper: EOR,
	},
	0x45: OpCode{name: "EOR",
		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: EOR,
	},
	0x55: OpCode{
		name: "EOR",

		cycles:           4,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: EOR,
	},
	0x4d: OpCode{
		name: "EOR",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: EOR,
	},
	0x5d: OpCode{
		name: "EOR",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteX,
		oper: EOR,
	},
	0x59: OpCode{
		name: "EOR",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteY,
		oper: EOR,
	},
	0x41: OpCode{
		name: "EOR",

		cycles:           6,
		pageBoundryCheck: false,

		mode: IndirectX,
		oper: EOR,
	},
	0x51: OpCode{
		name: "EOR",

		cycles:           5,
		pageBoundryCheck: true,

		mode: IndirectY,
		oper: EOR,
	},
	0xe6: OpCode{
		name: "INC",

		cycles:           5,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: INC,
	},
	0xf6: OpCode{
		name: "INC",

		cycles:           6,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: INC,
	},
	0xee: OpCode{
		name: "INC",

		cycles:           6,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: INC,
	},
	0xfe: OpCode{
		name: "INC",

		cycles:           7,
		pageBoundryCheck: false,

		mode: AbsoluteX,
		oper: INC,
	},
	0xe8: OpCode{
		name: "INX",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: INX,
	},
	0xc8: OpCode{
		name: "INY",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: INY,
	},
	0x4c: OpCode{
		name: "JMP",

		cycles:           3,
		pageBoundryCheck: false,

		mode: AbsoluteJMP,
		oper: JMP,
	},
	0x6c: OpCode{
		name: "JMP",

		cycles:           5,
		pageBoundryCheck: false,

		mode: Indirect,
		oper: JMP,
	},
	0x20: OpCode{
		name: "JSR",

		cycles:           6,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: JSR,
	},
	0xa9: OpCode{
		name: "LDA",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Immediate,
		oper: LDA,
	},
	0xa5: OpCode{
		name: "LDA",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: LDA,
	},
	0xb5: OpCode{
		name: "LDA",

		cycles:           4,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: LDA,
	},
	0xad: OpCode{
		name: "LDA",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: LDA,
	},
	0xbd: OpCode{
		name: "LDA",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteX,
		oper: LDA,
	},
	0xb9: OpCode{
		name: "LDA",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteY,
		oper: LDA,
	},
	0xa1: OpCode{
		name: "LDA",

		cycles:           6,
		pageBoundryCheck: false,

		mode: IndirectX,
		oper: LDA,
	},
	0xb1: OpCode{
		name: "LDA",

		cycles:           5,
		pageBoundryCheck: true,

		mode: IndirectY,
		oper: LDA,
	},
	0xa2: OpCode{
		name: "LDX",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Immediate,
		oper: LDX,
	},
	0xa6: OpCode{
		name: "LDX",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: LDX,
	},
	0xb6: OpCode{
		name: "LDX",

		cycles:           4,
		pageBoundryCheck: false,

		mode: ZeroPageY,
		oper: LDX,
	},
	0xae: OpCode{
		name: "LDX",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: LDX,
	},
	0xbe: OpCode{
		name: "LDX",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteY,
		oper: LDX,
	},
	0xa0: OpCode{
		name: "LDY",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Immediate,
		oper: LDY,
	},
	0xa4: OpCode{
		name: "LDY",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: LDY,
	},
	0xb4: OpCode{
		name: "LDY",

		cycles:           4,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: LDY,
	},
	0xac: OpCode{
		name: "LDY",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: LDY,
	},
	0xbc: OpCode{
		name: "LDY",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteX,
		oper: LDY,
	},
	0x4a: OpCode{
		name: "LSR",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Accumulator,
		oper: LSR,
	},
	0x46: OpCode{
		name: "LSR",

		cycles:           5,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: LSR,
	},
	0x56: OpCode{
		name: "LSR",

		cycles:           6,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: LSR,
	},
	0x4e: OpCode{
		name: "LSR",

		cycles:           6,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: LSR,
	},
	0x5e: OpCode{
		name: "LSR",

		cycles:           7,
		pageBoundryCheck: false,

		mode: AbsoluteX,
		oper: LSR,
	},
	0xea: OpCode{
		name: "NOP",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: NOP,
	},
	0x09: OpCode{
		name: "ORA",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Immediate,
		oper: ORA,
	},
	0x05: OpCode{
		name: "ORA",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: ORA,
	},
	0x15: OpCode{
		name: "ORA",

		cycles:           4,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: ORA,
	},
	0x0d: OpCode{
		name: "ORA",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: ORA,
	},
	0x1d: OpCode{
		name: "ORA",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteX,
		oper: ORA,
	},
	0x19: OpCode{
		name: "ORA",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteY,
		oper: ORA,
	},
	0x01: OpCode{
		name: "ORA",

		cycles:           6,
		pageBoundryCheck: false,

		mode: IndirectX,
		oper: ORA,
	},
	0x11: OpCode{
		name: "ORA",

		cycles:           5,
		pageBoundryCheck: true,

		mode: IndirectY,
		oper: ORA,
	},
	0x48: OpCode{
		name: "PHA",

		cycles:           3,
		pageBoundryCheck: false,

		mode: Implied,
		oper: PHA,
	},
	0x08: OpCode{
		name: "PHP",

		cycles:           3,
		pageBoundryCheck: false,

		mode: Implied,
		oper: PHP,
	},
	0x68: OpCode{
		name: "PLA",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Implied,
		oper: PLA,
	},
	0x28: OpCode{
		name: "PLP",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Implied,
		oper: PLP,
	},
	0x2a: OpCode{
		name: "ROL",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Accumulator,
		oper: ROL,
	},
	0x26: OpCode{
		name: "ROL",

		cycles:           5,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: ROL,
	},
	0x36: OpCode{
		name: "ROL",

		cycles:           6,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: ROL,
	},
	0x2e: OpCode{
		name: "ROL",

		cycles:           6,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: ROL,
	},
	0x3e: OpCode{
		name: "ROL",

		cycles:           7,
		pageBoundryCheck: false,

		mode: AbsoluteX,
		oper: ROL,
	},
	0x6a: OpCode{
		name: "ROR",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Accumulator,
		oper: ROR,
	},
	0x66: OpCode{
		name: "ROR",

		cycles:           5,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: ROR,
	},
	0x76: OpCode{
		name: "ROR",

		cycles:           6,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: ROR,
	},
	0x6e: OpCode{
		name: "ROR",

		cycles:           6,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: ROR,
	},
	0x7e: OpCode{
		name: "ROR",

		cycles:           7,
		pageBoundryCheck: false,

		mode: AbsoluteX,
		oper: ROR,
	},
	0x40: OpCode{
		name: "RTI",

		cycles:           6,
		pageBoundryCheck: false,

		mode: Implied,
		oper: RTI,
	},
	0x60: OpCode{
		name: "RTS",

		cycles:           6,
		pageBoundryCheck: false,

		mode: Implied,
		oper: RTS,
	},
	0xe9: OpCode{
		name: "SBC",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Immediate,
		oper: SBC,
	},
	0xe5: OpCode{
		name: "SBC",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: SBC,
	},
	0xf5: OpCode{
		name: "SBC",

		cycles:           4,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: SBC,
	},
	0xed: OpCode{
		name: "SBC",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: SBC,
	},
	0xfd: OpCode{
		name: "SBC",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteX,
		oper: SBC,
	},
	0xf9: OpCode{
		name: "SBC",

		cycles:           4,
		pageBoundryCheck: true,

		mode: AbsoluteY,
		oper: SBC,
	},
	0xe1: OpCode{
		name: "SBC",

		cycles:           6,
		pageBoundryCheck: false,

		mode: IndirectX,
		oper: SBC,
	},
	0xf1: OpCode{
		name: "SBC",

		cycles:           5,
		pageBoundryCheck: true,

		mode: IndirectY,
		oper: SBC,
	},
	0x38: OpCode{
		name: "SEC",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: SEC,
	},
	0xf8: OpCode{
		name: "SED",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: SED,
	},
	0x78: OpCode{
		name: "SEI",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: SEI,
	},
	0x85: OpCode{
		name: "STA",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: STA,
	},
	0x95: OpCode{
		name: "STA",

		cycles:           4,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: STA,
	},
	0x8d: OpCode{
		name: "STA",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: STA,
	},
	0x9d: OpCode{
		name: "STA",

		cycles:           5,
		pageBoundryCheck: true,

		mode: AbsoluteX,
		oper: STA,
	},
	0x99: OpCode{
		name: "STA",

		cycles:           5,
		pageBoundryCheck: false,

		mode: AbsoluteY,
		oper: STA,
	},
	0x81: OpCode{
		name: "STA",

		cycles:           6,
		pageBoundryCheck: false,

		mode: IndirectX,
		oper: STA,
	},
	0x91: OpCode{
		name: "STA",

		cycles:           6,
		pageBoundryCheck: false,

		mode: IndirectY,
		oper: STA,
	},
	0x86: OpCode{
		name: "STX",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: STX,
	},
	0x96: OpCode{
		name: "STX",

		cycles:           4,
		pageBoundryCheck: false,

		mode: ZeroPageY,
		oper: STX,
	},
	0x8e: OpCode{
		name: "STX",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: STX,
	},
	0x84: OpCode{
		name: "STY",

		cycles:           3,
		pageBoundryCheck: false,

		mode: ZeroPage,
		oper: STY,
	},
	0x94: OpCode{
		name: "STY",

		cycles:           4,
		pageBoundryCheck: false,

		mode: ZeroPageX,
		oper: STY,
	},
	0x8c: OpCode{
		name: "STY",

		cycles:           4,
		pageBoundryCheck: false,

		mode: Absolute,
		oper: STY,
	},
	0xaa: OpCode{
		name: "TAX",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: TAX,
	},
	0xa8: OpCode{
		name: "TAY",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: TAY,
	},
	0xba: OpCode{
		name: "TSX",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: TSX,
	},
	0x8a: OpCode{
		name: "TXA",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: TXA,
	},
	0x9a: OpCode{
		name: "TXS",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: TXS,
	},
	0x98: OpCode{
		name: "TYA",

		cycles:           2,
		pageBoundryCheck: false,

		mode: Implied,
		oper: TYA,
	},
}
