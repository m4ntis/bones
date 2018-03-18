package ppu

type Regs struct {
	PPUCTRL   *byte
	PPUMASK   *byte
	PPUSTATUS *byte
	OAMADDR   *byte
	OAMDATA   *byte
	PPUSCROLL *byte
	PPUADDR   *byte
	PPUDATA   *byte
	OAMDMA    *byte
}
