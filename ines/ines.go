package ines

import (
	"bytes"
	"encoding/hex"
	"io"

	"github.com/m4ntis/bones/models"
	"github.com/pkg/errors"
)

const (
	INES_HEADER_SIZE = 16

	HORIZONTAL_MIRRORING = 0
	VERTICAL_MIRRORING   = 1
)

type inesHeader struct {
	PrgROMSize int
	ChrROMSize int

	Mirroring        int
	PersistentMemory bool
	Trainer          bool
	IgnoreMirror     bool
	MapperNumber     int
}

func readHeader(r *io.Reader) (header []byte, err error) {
	header := make([]byte, INES_HEADER_SIZE)
	n, err := r.Read(header)
	if n < INES_HEADER_SIZE && err == nil {
		return nil, errors.Errorf("Couldn't read enough bytes, read %d/%d", n,
			INES_HEADER_SIZE)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}
	return b, nil
}

// parseHeader parses the slice it gets into an inesHeader struct.
//
// This method expects the slice to be of size 16, and panics if shorter and
// disregards the trailing data if longer, so the caller must ensure that the
// sent buffer is of len() >= 16.
func parseHeader(header []byte) (header inesHeader, err error) {
	if !bytes.Equal(header[:3], []byte{0x4e, 0x45, 0x53, 0x1a}) {
		return nil, errors.Errorf("Incorrect iNes header prefix: %s",
			hex.Dump(header[:3]))
	}

	prgROMSize := header[4]
	if prgROMSize == 0 {
		return nil, errors.New("PRG ROM size can't be 0")
	}

	chrROMSize := header[5]
	if chrROMSize == 0 {
		chrROMSize = 1
	}

	mirroring := header[6] & 1
	persistentMemory := header[6] & 2
	trainer := header[6] & 4
	ignoreMirror := header[6] & 8
	mapperNumber := header[6] & 240

	return inesHeader{
		PrgROMSize: prgROMSize,
		ChrROMSize: chrROMSize,

		Mirroring:        mirroring,
		PersistentMemory: persistentMemory,
		Trainer:          trainer,
		IgnoreMirror:     ignoreMirror,
		MapperNumber:     mapperNumber,
	}
}

// Parse reads from the reader, parses the data and populates a ROM
// struct accordingly.
//
// The errors it returns may relate to reading errors or ines format errors
func Parse(r *io.Reader) (rom *models.ROM, err error) {
	// Read and parse header
	headerBuff, err := readHeader(r)
	if err != nil {
		return nil, errors.Wrap(err, "Error while reading iNes header")
	}
	header, err := parseHeader(headerBuff)
	if err != nil {
		return nil, errors.Wrap(err, "Error while parsing iNes header")
	}

	// Calculate ROM size and read it
	trainerSize := header.Trainer * models.TRAINER_SIZE
	prgROMSize := header.PrgROMSize * models.PRG_ROM_PAGE_SIZE
	chrROMSize := header.ChrROMSize * models.CHR_ROM_PAGE_SIZE
	romSize := trainerSize + prgROMSize + chrROMSize

	romBuff := make([]byte, romSize)
	n := copy(romBuff, r)
	if n < romSize {
		return nil, errors.Errorf("Not enough data in rom, %d/%d", n, romSize)
	}

	// Populate ROM fields
	var trainer models.Trainer
	copy(trainer[:], romBuff[:models.TRAINER_SIZE])

	prgROM := make([]PrgROMPage, header.PrgROMSize)
	for i := range prgROM {
		startIndex := trainerSize + i*models.PRG_ROM_PAGE_SIZE
		copy(prgROM[i][:], r[startIndex:startIndex+models.PRG_ROM_PAGE_SIZE])
	}

	chrROM := make([]ChrROMPage, header.ChrROMSize)
	for i := range chrROM {
		startIndex := trainerSize + prgROMSize + i*models.CHR_ROM_PAGE_SIZE
		copy(chrROM[i][:], r[startIndex:startIndex+models.CHR_ROM_PAGE_SIZE])
	}

	return NewRom(trainer, prgROM, chrROM), nil
}
