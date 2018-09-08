// Package ines provides an api for iNES format parsing
package ines

import (
	"bytes"
	"encoding/hex"
	"io"

	"github.com/m4ntis/bones/ines/mapper"
	"github.com/pkg/errors"
)

const (
	InesHeaderSize = 16

	HorizontalMirroring = 0
	VerticalMirroring   = 1
)

type INESHeader struct {
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
func parseHeader(headerBuff []byte) (header INESHeader, err error) {
	if !bytes.Equal(headerBuff[:4], []byte{0x4e, 0x45, 0x53, 0x1a}) {
		return INESHeader{}, errors.Errorf("Incorrect iNes header prefix: %s",
			hex.Dump(headerBuff[:4]))
	}

	prgROMSize := headerBuff[4]
	if prgROMSize == 0 {
		return INESHeader{}, errors.New("PRG ROM size can't be 0")
	}

	chrROMSize := headerBuff[5]

	/*
		Flag 6
		76543210
		||||||||
		|||||||+- Mirroring: 0: horizontal (vertical arrangement) (CIRAM A10 = PPU A11)
		|||||||              1: vertical (horizontal arrangement) (CIRAM A10 = PPU A10)
		||||||+-- 1: Cartridge contains battery-backed PRG RAM ($6000-7FFF) or other persistent memory
		|||||+--- 1: 512-byte trainer at $7000-$71FF (stored before PRG data)
		||||+---- 1: Ignore mirroring control or above mirroring bit; instead provide four-screen VRAM
		++++----- Lower nybble of mapper number
	*/
	mirroring := headerBuff[6] & 1
	persistentMemory := headerBuff[6] & 2 >> 1
	trainer := headerBuff[6] & 4 >> 2
	ignoreMirror := headerBuff[6] & 8 >> 3
	mapperNumber := headerBuff[6] & 240 >> 4

	/*
		Flag 7
		76543210
		||||||||
		|||||||+- VS Unisystem
		||||||+-- PlayChoice-10 (8KB of Hint Screen data stored after CHR data)
		||||++--- If equal to 2, flags 8-15 are in NES 2.0 format
		++++----- Upper nybble of mapper number
	*/
	//version := headerBuff[7] & 12 >> 2
	mapperNumber += headerBuff[7] & 240 >> 4

	tvSystem := headerBuff[9] & 1
	if tvSystem == 1 {
		return INESHeader{}, errors.Errorf("This is a PAL ROM. BoNES only supports NTSC games")
	}

	return INESHeader{
		PrgROMSize: int(prgROMSize),
		ChrROMSize: int(chrROMSize),

		Mirroring:        int(mirroring),
		PersistentMemory: int(persistentMemory),
		Trainer:          int(trainer),
		IgnoreMirror:     int(ignoreMirror),
		MapperNumber:     int(mapperNumber),
	}, nil
}

// Parse reads an ines rom from r and populates a ROM struct with its data or
// returns an error.
func Parse(r io.Reader) (rom *ROM, err error) {
	// Read and parse header
	headerBuff, err := readHeader(r)
	if err != nil {
		return nil, errors.Wrap(err, "Error while reading iNes header")
	}
	header, err := parseHeader(headerBuff)
	if err != nil {
		return nil, errors.Wrap(err, "Error while parsing iNes header")
	}

	romMapper, err := mapper.New(header.MapperNumber)
	if err != nil {
		return nil, errors.Wrap(err, "Error while parsing iNes rom")
	}

	// Calculate ROM size and read it
	trainerSize := header.Trainer * TrainerSize
	prgROMSize := header.PrgROMSize * PrgROMPageSize
	chrROMSize := header.ChrROMSize * ChrROMPageSize
	romSize := trainerSize + prgROMSize + chrROMSize

	romBuff := make([]byte, romSize)
	n, err := r.Read(romBuff)
	if n > 0 && n < romSize {
		return nil, errors.Errorf("Not enough data in ROM, %d/%d", n, romSize)
	} else if err != nil {
		return nil, errors.Wrap(err, "Error while reading ROM")
	}

	// Populate ROM fields
	var trainer Trainer
	copy(trainer[:], romBuff[:trainerSize])

	prgROM := make([]PrgROMPage, header.PrgROMSize)
	for i := range prgROM {
		startIndex := trainerSize + i*PrgROMPageSize
		copy(prgROM[i][:],
			romBuff[startIndex:startIndex+PrgROMPageSize])
	}

	chrROM := make([]ChrROMPage, header.ChrROMSize)
	for i := range chrROM {
		startIndex := trainerSize + prgROMSize + i*ChrROMPageSize
		copy(chrROM[i][:],
			romBuff[startIndex:startIndex+ChrROMPageSize])
	}

	romMapper.Populate(prgROM, chrROM)

	return &ROM{Header: header, Trainer: trainer, Mapper: romMapper}, nil
}
