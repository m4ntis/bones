package cpu

var OpCodes = map[byte]OpCode{
	0x69: OpCode{
		Name: "ADC",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Immediate,
		Oper: ADC,
	},
	0x65: OpCode{
		Name: "ADC",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: ADC,
	},
	0x75: OpCode{
		Name: "ADC",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: ADC,
	},
	0x6d: OpCode{
		Name: "ADC",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: ADC,
	},
	0x7d: OpCode{
		Name: "ADC",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteX,
		Oper: ADC,
	},
	0x79: OpCode{
		Name: "ADC",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteY,
		Oper: ADC,
	},
	0x61: OpCode{
		Name: "ADC",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: IndirectX,
		Oper: ADC,
	},
	0x71: OpCode{
		Name: "ADC",

		cycles:           5,
		pageBoundryCheck: true,

		Mode: IndirectY,
		Oper: ADC,
	},
	0x29: OpCode{
		Name: "AND",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Immediate,
		Oper: AND,
	},
	0x25: OpCode{
		Name: "AND",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: AND,
	},
	0x35: OpCode{
		Name: "AND",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: AND,
	},
	0x2d: OpCode{
		Name: "AND",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: AND,
	},
	0x3d: OpCode{
		Name: "AND",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteX,
		Oper: AND,
	},
	0x39: OpCode{
		Name: "AND",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteY,
		Oper: AND,
	},
	0x21: OpCode{
		Name: "AND",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: IndirectX,
		Oper: AND,
	},
	0x31: OpCode{
		Name: "AND",

		cycles:           5,
		pageBoundryCheck: true,

		Mode: IndirectY,
		Oper: AND,
	},
	0x0a: OpCode{
		Name: "ASL",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Accumulator,
		Oper: ASL,
	},
	0x06: OpCode{
		Name: "ASL",

		cycles:           5,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: ASL,
	},
	0x16: OpCode{
		Name: "ASL",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: ASL,
	},
	0x0e: OpCode{
		Name: "ASL",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: ASL,
	},
	0x1e: OpCode{
		Name: "ASL",

		cycles:           7,
		pageBoundryCheck: false,

		Mode: AbsoluteX,
		Oper: ASL,
	},
	0x90: OpCode{
		Name: "BCC",

		cycles:           2,
		pageBoundryCheck: true,

		Mode: Relative,
		Oper: BCC,
	},
	0xb0: OpCode{
		Name: "BCS",

		cycles:           2,
		pageBoundryCheck: true,

		Mode: Relative,
		Oper: BCS,
	},
	0xf0: OpCode{
		Name: "BEQ",

		cycles:           2,
		pageBoundryCheck: true,

		Mode: Relative,
		Oper: BEQ,
	},
	0x24: OpCode{
		Name: "BIT",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: BIT,
	},
	0x2c: OpCode{
		Name: "BIT",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: BIT,
	},
	0x30: OpCode{
		Name: "BMI",

		cycles:           2,
		pageBoundryCheck: true,

		Mode: Relative,
		Oper: BMI,
	},
	0xd0: OpCode{
		Name: "BNE",

		cycles:           2,
		pageBoundryCheck: true,

		Mode: Relative,
		Oper: BNE,
	},
	0x10: OpCode{
		Name: "BPL",

		cycles:           2,
		pageBoundryCheck: true,

		Mode: Relative,
		Oper: BPL,
	},
	0x00: OpCode{
		Name: "BRK",

		cycles:           7,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: BRK,
	},
	0x50: OpCode{
		Name: "BVC",

		cycles:           2,
		pageBoundryCheck: true,

		Mode: Relative,
		Oper: BVC,
	},
	0x70: OpCode{
		Name: "BVS",

		cycles:           2,
		pageBoundryCheck: true,

		Mode: Relative,
		Oper: BVS,
	},
	0x18: OpCode{
		Name: "CLC",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: CLC,
	},
	0xd8: OpCode{
		Name: "CLD",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: CLD,
	},
	0x58: OpCode{
		Name: "CLI",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: CLI,
	},
	0xb8: OpCode{
		Name: "CLV",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: CLV,
	},
	0xc9: OpCode{
		Name: "CMP",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Immediate,
		Oper: CMP,
	},
	0xc5: OpCode{
		Name: "CMP",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: CMP,
	},
	0xd5: OpCode{
		Name: "CMP",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: CMP,
	},
	0xcd: OpCode{
		Name: "CMP",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: CMP,
	},
	0xdd: OpCode{
		Name: "CMP",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteX,
		Oper: CMP,
	},
	0xd9: OpCode{
		Name: "CMP",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteY,
		Oper: CMP,
	},
	0xc1: OpCode{
		Name: "CMP",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: IndirectX,
		Oper: CMP,
	},
	0xd1: OpCode{
		Name: "CMP",

		cycles:           5,
		pageBoundryCheck: true,

		Mode: IndirectY,
		Oper: CMP,
	},
	0xe0: OpCode{
		Name: "CPX",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Immediate,
		Oper: CPX,
	},
	0xe4: OpCode{
		Name: "CPX",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: CPX,
	},
	0xec: OpCode{
		Name: "CPX",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: CPX,
	},
	0xc0: OpCode{
		Name: "CPY",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Immediate,
		Oper: CPY,
	},
	0xc4: OpCode{
		Name: "CPY",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: CPY,
	},
	0xcc: OpCode{
		Name: "CPY",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: CPY,
	},
	0xc6: OpCode{
		Name: "DEC",

		cycles:           5,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: DEC,
	},
	0xd6: OpCode{
		Name: "DEC",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: DEC,
	},
	0xce: OpCode{
		Name: "DEC",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: DEC,
	},
	0xde: OpCode{
		Name: "DEC",

		cycles:           7,
		pageBoundryCheck: false,

		Mode: AbsoluteX,
		Oper: DEC,
	},
	0xca: OpCode{
		Name: "DEX",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: DEX,
	},
	0x88: OpCode{
		Name: "DEY",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: DEY,
	},
	0x49: OpCode{
		Name: "EOR",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Immediate,
		Oper: EOR,
	},
	0x45: OpCode{Name: "EOR",
		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: EOR,
	},
	0x55: OpCode{
		Name: "EOR",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: EOR,
	},
	0x4d: OpCode{
		Name: "EOR",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: EOR,
	},
	0x5d: OpCode{
		Name: "EOR",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteX,
		Oper: EOR,
	},
	0x59: OpCode{
		Name: "EOR",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteY,
		Oper: EOR,
	},
	0x41: OpCode{
		Name: "EOR",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: IndirectX,
		Oper: EOR,
	},
	0x51: OpCode{
		Name: "EOR",

		cycles:           5,
		pageBoundryCheck: true,

		Mode: IndirectY,
		Oper: EOR,
	},
	0xe6: OpCode{
		Name: "INC",

		cycles:           5,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: INC,
	},
	0xf6: OpCode{
		Name: "INC",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: INC,
	},
	0xee: OpCode{
		Name: "INC",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: INC,
	},
	0xfe: OpCode{
		Name: "INC",

		cycles:           7,
		pageBoundryCheck: false,

		Mode: AbsoluteX,
		Oper: INC,
	},
	0xe8: OpCode{
		Name: "INX",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: INX,
	},
	0xc8: OpCode{
		Name: "INY",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: INY,
	},
	0x4c: OpCode{
		Name: "JMP",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: JMP,
	},
	0x6c: OpCode{
		Name: "JMP",

		cycles:           5,
		pageBoundryCheck: false,

		Mode: Indirect,
		Oper: JMP,
	},
	0x20: OpCode{
		Name: "JSR",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: JSR,
	},
	0xa9: OpCode{
		Name: "LDA",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Immediate,
		Oper: LDA,
	},
	0xa5: OpCode{
		Name: "LDA",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: LDA,
	},
	0xb5: OpCode{
		Name: "LDA",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: LDA,
	},
	0xad: OpCode{
		Name: "LDA",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: LDA,
	},
	0xbd: OpCode{
		Name: "LDA",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteX,
		Oper: LDA,
	},
	0xb9: OpCode{
		Name: "LDA",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteY,
		Oper: LDA,
	},
	0xa1: OpCode{
		Name: "LDA",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: IndirectX,
		Oper: LDA,
	},
	0xb1: OpCode{
		Name: "LDA",

		cycles:           5,
		pageBoundryCheck: true,

		Mode: IndirectY,
		Oper: LDA,
	},
	0xa2: OpCode{
		Name: "LDX",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Immediate,
		Oper: LDX,
	},
	0xa6: OpCode{
		Name: "LDX",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: LDX,
	},
	0xb6: OpCode{
		Name: "LDX",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: ZeroPageY,
		Oper: LDX,
	},
	0xae: OpCode{
		Name: "LDX",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: LDX,
	},
	0xbe: OpCode{
		Name: "LDX",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteY,
		Oper: LDX,
	},
	0xa0: OpCode{
		Name: "LDY",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Immediate,
		Oper: LDY,
	},
	0xa4: OpCode{
		Name: "LDY",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: LDY,
	},
	0xb4: OpCode{
		Name: "LDY",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: LDY,
	},
	0xac: OpCode{
		Name: "LDY",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: LDY,
	},
	0xbc: OpCode{
		Name: "LDY",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteX,
		Oper: LDY,
	},
	0x4a: OpCode{
		Name: "LSR",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Accumulator,
		Oper: LSR,
	},
	0x46: OpCode{
		Name: "LSR",

		cycles:           5,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: LSR,
	},
	0x56: OpCode{
		Name: "LSR",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: LSR,
	},
	0x4e: OpCode{
		Name: "LSR",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: LSR,
	},
	0x5e: OpCode{
		Name: "LSR",

		cycles:           7,
		pageBoundryCheck: false,

		Mode: AbsoluteX,
		Oper: LSR,
	},
	0xea: OpCode{
		Name: "NOP",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: NOP,
	},
	0x09: OpCode{
		Name: "ORA",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Immediate,
		Oper: ORA,
	},
	0x05: OpCode{
		Name: "ORA",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: ORA,
	},
	0x15: OpCode{
		Name: "ORA",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: ORA,
	},
	0x0d: OpCode{
		Name: "ORA",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: ORA,
	},
	0x1d: OpCode{
		Name: "ORA",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteX,
		Oper: ORA,
	},
	0x19: OpCode{
		Name: "ORA",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteY,
		Oper: ORA,
	},
	0x01: OpCode{
		Name: "ORA",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: IndirectX,
		Oper: ORA,
	},
	0x11: OpCode{
		Name: "ORA",

		cycles:           5,
		pageBoundryCheck: true,

		Mode: IndirectY,
		Oper: ORA,
	},
	0x48: OpCode{
		Name: "PHA",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: PHA,
	},
	0x08: OpCode{
		Name: "PHP",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: PHP,
	},
	0x68: OpCode{
		Name: "PLA",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: PLA,
	},
	0x28: OpCode{
		Name: "PLP",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: PLP,
	},
	0x2a: OpCode{
		Name: "ROL",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Accumulator,
		Oper: ROL,
	},
	0x26: OpCode{
		Name: "ROL",

		cycles:           5,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: ROL,
	},
	0x36: OpCode{
		Name: "ROL",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: ROL,
	},
	0x2e: OpCode{
		Name: "ROL",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: ROL,
	},
	0x3e: OpCode{
		Name: "ROL",

		cycles:           7,
		pageBoundryCheck: false,

		Mode: AbsoluteX,
		Oper: ROL,
	},
	0x6a: OpCode{
		Name: "ROR",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Accumulator,
		Oper: ROR,
	},
	0x66: OpCode{
		Name: "ROR",

		cycles:           5,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: ROR,
	},
	0x76: OpCode{
		Name: "ROR",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: ROR,
	},
	0x6e: OpCode{
		Name: "ROR",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: ROR,
	},
	0x7e: OpCode{
		Name: "ROR",

		cycles:           7,
		pageBoundryCheck: false,

		Mode: AbsoluteX,
		Oper: ROR,
	},
	0x40: OpCode{
		Name: "RTI",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: RTI,
	},
	0x60: OpCode{
		Name: "RTS",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: RTS,
	},
	0xe9: OpCode{
		Name: "SBC",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Immediate,
		Oper: SBC,
	},
	0xe5: OpCode{
		Name: "SBC",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: SBC,
	},
	0xf5: OpCode{
		Name: "SBC",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: SBC,
	},
	0xed: OpCode{
		Name: "SBC",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: SBC,
	},
	0xfd: OpCode{
		Name: "SBC",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteX,
		Oper: SBC,
	},
	0xf9: OpCode{
		Name: "SBC",

		cycles:           4,
		pageBoundryCheck: true,

		Mode: AbsoluteY,
		Oper: SBC,
	},
	0xe1: OpCode{
		Name: "SBC",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: IndirectX,
		Oper: SBC,
	},
	0xf1: OpCode{
		Name: "SBC",

		cycles:           5,
		pageBoundryCheck: true,

		Mode: IndirectY,
		Oper: SBC,
	},
	0x38: OpCode{
		Name: "SEC",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: SEC,
	},
	0xf8: OpCode{
		Name: "SED",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: SED,
	},
	0x78: OpCode{
		Name: "SEI",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: SEI,
	},
	0x85: OpCode{
		Name: "STA",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: STA,
	},
	0x95: OpCode{
		Name: "STA",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: STA,
	},
	0x8d: OpCode{
		Name: "STA",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: STA,
	},
	0x9d: OpCode{
		Name: "STA",

		cycles:           5,
		pageBoundryCheck: true,

		Mode: AbsoluteX,
		Oper: STA,
	},
	0x99: OpCode{
		Name: "STA",

		cycles:           5,
		pageBoundryCheck: false,

		Mode: AbsoluteY,
		Oper: STA,
	},
	0x81: OpCode{
		Name: "STA",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: IndirectX,
		Oper: STA,
	},
	0x91: OpCode{
		Name: "STA",

		cycles:           6,
		pageBoundryCheck: false,

		Mode: IndirectY,
		Oper: STA,
	},
	0x86: OpCode{
		Name: "STX",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: STX,
	},
	0x96: OpCode{
		Name: "STX",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: ZeroPageY,
		Oper: STX,
	},
	0x8e: OpCode{
		Name: "STX",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: STX,
	},
	0x84: OpCode{
		Name: "STY",

		cycles:           3,
		pageBoundryCheck: false,

		Mode: ZeroPage,
		Oper: STY,
	},
	0x94: OpCode{
		Name: "STY",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: ZeroPageX,
		Oper: STY,
	},
	0x8c: OpCode{
		Name: "STY",

		cycles:           4,
		pageBoundryCheck: false,

		Mode: Absolute,
		Oper: STY,
	},
	0xaa: OpCode{
		Name: "TAX",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: TAX,
	},
	0xa8: OpCode{
		Name: "TAY",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: TAY,
	},
	0xba: OpCode{
		Name: "TSX",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: TSX,
	},
	0x8a: OpCode{
		Name: "TXA",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: TXA,
	},
	0x9a: OpCode{
		Name: "TXS",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: TXS,
	},
	0x98: OpCode{
		Name: "TYA",

		cycles:           2,
		pageBoundryCheck: false,

		Mode: Implied,
		Oper: TYA,
	},
}
