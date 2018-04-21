// Package ines provides an api for ines format parsing
package ines

import (
	"bytes"
	"encoding/hex"
	"io"

	"github.com/m4ntis/bones/models"
	"github.com/pkg/errors"
)

const (
	InesHeaderSize = 16

	HorizontalMirroring = 0
	VerticalMirroring   = 1
)

type inesHeader struct {
	PrgROMSize int
	ChrROMSize int

	Mirroring        int
	PersistentMemory int
	Trainer          int
	IgnoreMirror     int
	MapperNumber     int
}

func readHeader(r io.Reader) (header []byte, err error) {
	header = make([]byte, InesHeaderSize)
	n, err := r.Read(header)
	if n < InesHeaderSize && err == nil {
		return nil, errors.Errorf("Couldn't read enough bytes, read %d/%d", n,
			InesHeaderSize)
	} else if err != nil {
		return nil, errors.WithStack(err)
	}
	return header, nil
}

// parseHeader parses the slice it gets into an inesHeader struct.
//
// This method expects the slice to be of size 16, and panics if shorter and
// disregards the trailing data if longer, so the caller must ensure that the
// sent buffer is of len() >= 16.
func parseHeader(headerBuff []byte) (header *inesHeader, err error) {
	if !bytes.Equal(headerBuff[:4], []byte{0x4e, 0x45, 0x53, 0x1a}) {
		return nil, errors.Errorf("Incorrect iNes header prefix: %s",
			hex.Dump(headerBuff[:4]))
	}

	prgROMSize := headerBuff[4]
	if prgROMSize == 0 {
		return nil, errors.New("PRG ROM size can't be 0")
	}

	chrROMSize := headerBuff[5]
	if chrROMSize == 0 {
		chrROMSize = 1
	}

	mirroring := headerBuff[6] & 1
	persistentMemory := headerBuff[6] & 2
	trainer := headerBuff[6] & 4
	ignoreMirror := headerBuff[6] & 8
	mapperNumber := headerBuff[6] & 240

	return &inesHeader{
		PrgROMSize: int(prgROMSize),
		ChrROMSize: int(chrROMSize),

		Mirroring:        int(mirroring),
		PersistentMemory: int(persistentMemory),
		Trainer:          int(trainer),
		IgnoreMirror:     int(ignoreMirror),
		MapperNumber:     int(mapperNumber),
	}, nil
}

// Parse reads an ines file from r and populates a ROM struct or returns an
// error.
func Parse(r io.Reader) (rom *models.ROM, err error) {
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
	trainerSize := header.Trainer * models.TrainerSize
	prgROMSize := header.PrgROMSize * models.PrgROMPageSize
	chrROMSize := header.ChrROMSize * models.ChrROMPageSize
	romSize := trainerSize + prgROMSize + chrROMSize

	romBuff := make([]byte, romSize)
	n, err := r.Read(romBuff)
	if n > 0 && n < romSize {
		return nil, errors.Errorf("Not enough data in ROM, %d/%d", n, romSize)
	} else if err != nil {
		return nil, errors.Wrap(err, "Error while reading ROM")
	}

	// Populate ROM fields
	var trainer models.Trainer
	copy(trainer[:], romBuff[:trainerSize])

	prgROM := make([]models.PrgROMPage, header.PrgROMSize)
	for i := range prgROM {
		startIndex := trainerSize + i*models.PrgROMPageSize
		copy(prgROM[i][:],
			romBuff[startIndex:startIndex+models.PrgROMPageSize])
	}

	chrROM := make([]models.ChrROMPage, header.ChrROMSize)
	for i := range chrROM {
		startIndex := trainerSize + prgROMSize + i*models.ChrROMPageSize
		copy(chrROM[i][:],
			romBuff[startIndex:startIndex+models.ChrROMPageSize])
	}

	return models.NewROM(trainer, prgROM, chrROM), nil
}
