package tiny

import (
	"fmt"
	"strconv"
	"unsafe"
)

/**
Constants
*/

// Bit represents one binary place. [0 - 1]
type Bit byte

// String converts the provided Bit to a string 1 or 0, or panics if the found value is anything else
func (b Bit) String() string {
	if b > 1 {
		// PLEASE inform us if there is an issue whenever possible!!!
		panic(ErrorNotABit)
	}
	if b == 0 {
		return "0"
	}
	return "1"
}

// Zero is an implicit Bit{0}.
const Zero Bit = 0

// One is an implicit Bit{1}.
const One Bit = 1

// Encoding represents the encoding scheme of a Phrase of Measurement points.
type Encoding int

const (
	// Raw indicates this is simply a phrase of arbitrarily long binary information.
	Raw Encoding = iota

	// Logical indicates this phrase entirely consists of byte measurements.
	Logical

	// Signed indicates the first measurement is a sign, followed by the value.
	Signed

	// Float indicates the first measurement is a sign, followed by an exponent, and lastly a mantissa.
	Float

	// Index indicates the phrase entirely consists of logical binary data bound to a fixed width.
	Index
)

// WordWidth is the bit width of a standard int, which for all reasonable intents and purposes matches the architecture's word width.
const WordWidth = strconv.IntSize // NOTE: While this could mismatch the architecture's word in some cases, the performance implications are minimal.

// Endianness indicates the logical -byte- ordering of sequential bytes.  All binary data has a most significant side,
// where the binary placeholder has the highest relative value, as well as a least significant side.  The individual bits
// of a byte are colloquially manipulated in most→to→least significant order, but multiple -bytes- worth of information may
// be stored in least←to←most significant order while retaining the individual bit order of each byte.  There are two
// types of endianness -
//
// BigEndian, where the most significant bytes come first - often used in network protocols:
//
//		| Most Sig. Byte  |   Middle Byte   | Least Sig. Byte |
//		| 0 1 0 0 1 1 0 1 | 0 0 1 0 1 1 0 0 | 0 0 0 1 0 1 1 0 | (5,057,558)
//		|        77       |        44       |        22       |
//	          ⬑ I'm -only- showing the values here as an
//	            identifier, not as part of the order
//
// LittleEndian, where the least significant bytes come first - used by x86, AMD64, ARM, and many more:
//
//	| Least Sig. Byte |   Middle Byte   |  Most Sig. Byte |
//	| 0 0 0 1 0 1 1 0 | 0 0 1 0 1 1 0 0 | 0 1 0 0 1 1 0 1 | (5,057,558)
//	|        22       |        44       |        77       |
//
// NOTE: While some hardware may physically store bits in least←to←most order internally, Go's shift operators (>> and <<)
// are guaranteed by the language specification to always operate in most→to→least significant order. This, in turn, means
// that bit operations in tiny will -also- work with bits in most→to→least significant order regardless of the underlying
// architecture's physical bit storage order. When reading raw memory, only byte ordering needs to be handled explicitly.
//
// NOTE: Some protocols like UART traditionally transmit the least significant bit first, so you may need to reverse bits
// within bytes when interfacing with such protocols.
type Endianness int

const (
	LittleEndian Endianness = iota
	BigEndian
)

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

/**
Errors
*/

const errorMsgEndOfBits = "no more bits to read"

// ErrorEndOfBits is emitted whenever a read operation has requested to read beyond the available binary data's width.
// This error can absolutely be ignored, but also allows one to implicitly read until the end of bits has been reached.
var ErrorEndOfBits = fmt.Errorf(errorMsgEndOfBits)

const errorMsgNotABit = "bits must be 0 or 1"

// ErrorNotABit is emitted whenever a method expecting a Bit is provided with any other byte value than 1 or 0, as Bit is a byte underneath.
var ErrorNotABit = fmt.Errorf(errorMsgNotABit)
