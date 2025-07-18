// Package endianness provides access to the Endianness enumeration.
package endianness

import "unsafe"

// Endianness indicates the logical -byte- ordering of sequential bytes.  All binary data has a most significant side,
// where the binary placeholder has the highest relative value, as well as a least significant side.  The individual bits
// of a byte are colloquially manipulated in most→to→least significant order, but multiple -bytes- worth of information may
// be stored in least←to←most significant order while retaining the individual bit order of each byte.  There are two
// types of endianness -
//
// BigEndian, where the most significant bytes come first - or "raw" binary:
//
//	| Most Sig. Byte  |   Middle Byte   | Least Sig. Byte |
//	| 0 1 0 0 1 1 0 1 | 0 0 1 0 1 1 0 0 | 0 0 0 1 0 1 1 0 | (5,057,558)
//	|        4D       |        2C       |        16       |
//
// LittleEndian, where the least significant bytes come first - used by x86, AMD64, ARM, and the general world over:
//
//	| Least Sig. Byte |   Middle Byte   |  Most Sig. Byte |
//	| 0 0 0 1 0 1 1 0 | 0 0 1 0 1 1 0 0 | 0 1 0 0 1 1 0 1 | (5,057,558)
//	|        16       |        2C       |        4D       |
//	         ⬑  The byte's internal bits remain in most→to→least order
//
// NOTE: While some hardware may physically store bits in least←to←most order internally, Go's shift operators (<< and >>)
// are guaranteed by the language specification to always operate in most→to→least significant order. This, in turn, means
// that bit operations in tiny will -also- work with bits in most→to→least significant order regardless of the underlying
// architecture's physical bit storage order. When reading raw memory, only byte ordering needs to be handled explicitly.
//
// NOTE: Some protocols, like UART, traditionally transmit in least←to←most order, so you may also need to reverse bits
// within bytes when interfacing with such protocols - which we fully support =)
//
// See LittleEndian and BigEndian.
type Endianness byte

const (
	// LittleEndian indicates that bytes are handled in least←to←most significant order and is used by x86, AMD64, ARM, and the general
	// world over.
	//
	// See Endianness.
	LittleEndian Endianness = iota

	// BigEndian indicates that bytes are handled in most→to→least significant order and is often used in network protocols.
	//
	// See Endianness.
	BigEndian
)

// String prints an uppercase one-word representation of the Endianness.
func (e Endianness) String() string {
	switch e {
	case LittleEndian:
		return "LittleEndian"
	case BigEndian:
		return "BigEndian"
	default:
		return "Unknown"
	}
}

// StringFull prints an uppercase full two-word representation of the Endianness.
//
// You may optionally pass true for a lowercase representation.
func (e Endianness) StringFull(lowercase ...bool) string {
	lower := len(lowercase) > 0 && lowercase[0]
	switch e {
	case LittleEndian:
		if lower {
			return "little endian"
		}
		return "Little Endian"
	case BigEndian:
		if lower {
			return "big endian"
		}
		return "Big Endian"
	default:
		if lower {
			return "unknown"
		}
		return "Unknown"
	}
}

// GetEndianness returns the Endianness of the currently executing hardware.
func GetEndianness() Endianness {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf[0] {
	case 0xCD:
		return LittleEndian
	case 0xAB:
		return BigEndian
	default:
		panic("could not determine native endianness")
	}
}
