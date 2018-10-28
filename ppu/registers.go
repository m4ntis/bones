package ppu

type Registers struct {
	ppuCtrl   byte
	ppuMask   byte
	ppuStatus byte

	oamAddr byte

	scrollFirstWrite bool
	xScroll          int
	yScroll          int

	addrFirstWrite bool
	ppuAddr        int

	ppuData    byte
	ppuDataBuf byte

	vblank bool
	nmi    chan bool

	oam  *OAM
	vram *VRAM
}

func newRegisters(nmi chan bool, oam *OAM, vram *VRAM) *Registers {
	return &Registers{
		scrollFirstWrite: true,
		addrFirstWrite:   true,

		vblank: false,
		nmi:    nmi,

		oam:  oam,
		vram: vram,
	}
}

func (r *Registers) PPUCtrlWrite(data byte) {
	// If V Flag set (and changed) while in vblank
	if r.ppuCtrl>>7 == 0 && data>>7 == 1 && r.vblank {
		r.nmi <- true
	}

	r.ppuCtrl = data
}

func (r *Registers) PPUMaskWrite(data byte) {
	r.ppuMask = data
}

func (r *Registers) PPUStatusRead() byte {
	defer func() {
		// Clear bit 7
		r.ppuStatus &= 0x7f

		r.scrollFirstWrite = true
		r.xScroll = 0
		r.yScroll = 0

		r.ppuAddr = 0
	}()

	return r.ppuStatus
}

func (r *Registers) OAMAddrWrite(data byte) {
	r.oamAddr = data
}

func (r *Registers) OAMDataRead() byte {
	return (*r.oam)[r.oamAddr]
}

func (r *Registers) OAMDataWrite(data byte) {
	r.oamAddr++

	(*r.oam)[r.oamAddr] = data
}

// TODO: Changes made to the vertical scroll during rendering will only take
// effect on the next frame

func (r *Registers) PPUScrollWrite(data byte) {
	defer func() { r.scrollFirstWrite = !r.scrollFirstWrite }()

	if r.scrollFirstWrite {
		r.xScroll = int(data)
		return
	}

	r.yScroll = int(data)
}

func (r *Registers) PPUAddrWrite(data byte) {
	defer func() { r.addrFirstWrite = !r.addrFirstWrite }()

	if r.addrFirstWrite {
		r.ppuAddr = int(data) << 8
		return
	}
	r.ppuAddr |= int(data)
}

func (r *Registers) PPUDataRead() byte {
	defer r.incAddr()

	// If the read is from palette data, it is immediatelly put on the data bus
	if stripMirror(r.ppuAddr) >= bgrPaletteAddr {
		// TODO: Reading the palettes still updates the internal buffer though,
		// but the data placed in it is the mirrored nametable data that would
		// appear "underneath" the palette.
		return r.vram.Read(r.ppuAddr)
	}

	defer func() { r.ppuDataBuf = r.vram.Read(r.ppuAddr) }()
	return r.ppuDataBuf
}

func (r *Registers) PPUDataWrite(d byte) {
	defer r.incAddr()

	r.vram.Write(r.ppuAddr, d)
}

func (r *Registers) incAddr() {
	r.ppuAddr += int(1 + (r.ppuCtrl>>2&1)*31)
}
