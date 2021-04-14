package ines

import (
	"encoding/binary"
	"errors"
	"io"
)

type Type byte

const (
	INES Type = iota
	INES2
	Archaic
)

const Magic = 0x1a53454e

var ErrBadMagicValue = errors.New("not an ines file")

type ControlBits [2]byte

type Mapper uint8

type Header struct {
	PRGRomSize  byte // in 16KB units
	CHRRomSize  byte // in 8kb units
	ControlBits ControlBits
	PRGRamSize  byte // in 8kb units

	Flags9  Flags9
	Flags10 Flags10
	Padding [5]byte // padding
}

type Flags9 byte

// https://wiki.nesdev.com/w/index.php/INES#Flags_9 bit 0
// Otherwise NTSC, rarely seen in wild
func (f Flags9) TVSystem() TVSystem {
	if f&0b1 == 0b1 {
		return PAL
	} else {
		return NTSC
	}
}

type Flags10 byte

type TVSystem byte

const (
	NTSC          TVSystem = 0
	DUAL_COMPAT_1          = 1
	PAL                    = 2
	DUAL_COMPAT_3          = 3
)

// https://wiki.nesdev.com/w/index.php/INES#Flags_10 bits 0-1
func (f Flags10) TVSystem() TVSystem {
	return TVSystem(f & 0b11)
}

// https://wiki.nesdev.com/w/index.php/INES#Flags_10 bit 2
func (f Flags10) PRGRamPresent() bool {
	return f&0b100 == 0b100
}

// https://wiki.nesdev.com/w/index.php/INES#Flags_10 bit 3
func (f Flags10) HasBusConflicts() bool {
	return f&0b1000 == 0b1000
}

// https://wiki.nesdev.com/w/index.php/INES#Flags_6 bit 0
func (c ControlBits) Mirroring() bool {
	return c[0]&0b1 == 0b1
}

// https://wiki.nesdev.com/w/index.php/INES#Flags_6 bit 1
func (c ControlBits) HasPRGRam() bool {
	return c[0]&0b10 == 0b10
}

// https://wiki.nesdev.com/w/index.php/INES#Flags_6 bit 2
func (c ControlBits) HasTrainer() bool {
	return c[0]&0b100 == 0b100
}

// https://wiki.nesdev.com/w/index.php/INES#Flags_6 bit 3
func (c ControlBits) IgnoreMirroringProvideFourscreenVram() bool {
	return c[0]&0b1000 == 0b1000
}

func (c ControlBits) Mapper() Mapper {
	return Mapper(c[0]>>4 | (c[1] & 0xF0))
}

// https://wiki.nesdev.com/w/index.php/INES#Flags_7 bit 0
func (c ControlBits) VSUnisystem() bool {
	return c[1]&0b1 == 0b1
}

// https://wiki.nesdev.com/w/index.php/INES#Flags_7 bit 1
func (c ControlBits) Playchoice10() bool {
	return c[1]&0b10 == 0b10
}

// https://wiki.nesdev.com/w/index.php/INES#Flags_7 bit 2-3
func (c ControlBits) Nes2FormatFlag() bool {
	return c[1]&0b1100 == 0b1000
}

// https://wiki.nesdev.com/w/index.php/INES#Flags_7 bit 2-3
func (c ControlBits) Nes1FormatFlag() bool {
	return c[1]&0b1100 == 0b0000
}

// Checking padding bits zero to differentiate archaic vs INES 1
func (h Header) LastFourZero() bool {
	for i := 1; i < len(h.Padding); i++ {
		if h.Padding[i] != 0 {
			return false
		}
	}
	return true
}

func (h Header) Type() Type {
	switch {
	case h.ControlBits.Nes2FormatFlag():
		return INES2
	case h.ControlBits.Nes1FormatFlag() && h.LastFourZero():
		return INES
	default:
		return Archaic
	}
}

func (h *Header) Read(r io.Reader) error {
	var magic uint32
	if err := binary.Read(r, binary.LittleEndian, &magic); err != nil {
		return err
	}
	if magic != Magic {
		return ErrBadMagicValue
	}
	if err := binary.Read(r, binary.LittleEndian, h); err != nil {
		return err
	}
	return nil
}
