package tiny

import (
	"fmt"
	"strconv"
	"unsafe"
)

// binary represents the types that tiny supports conversion to bits from.
type binary interface {
	Measurement | Phrase | byte | Bit
}

/**
Global Constants
*/

// Unlimited represents a constantly referencable integer value which can be considered a reasonably "unlimited" width.
var Unlimited = ^uint(0)

// Zero is an implicit Bit{0}.
const Zero Bit = 0

// One is an implicit Bit{1}.
const One Bit = 1

// SingleZero is an implicit []Bit{0}.
var SingleZero = []Bit{Zero}

// DoubleZero is an implicit []Bit{0, 0}.
var DoubleZero = []Bit{Zero, Zero}

// SingleOne is an implicit []Bit{1}.
var SingleOne = []Bit{One}

// DoubleOne is an implicit []Bit{1, 1}.
var DoubleOne = []Bit{One, One}

// Nil is an implicit Bit{219} - this allows bits to intentionally be left out of slices and still stand out visibly amongst
// the other bits, as our Bit type is technically a byte in memory.  For example -
//
//	[ 0 0 0 0 0 0 0 0 ]   (0) ← A zero bit
//	[ 0 0 0 0 0 0 0 1 ]   (1) ← A one bit
//	[ 1 1 0 1 1 0 1 1 ] (219) ← A nil bit
//	    ⬑ Darkness is instantly recognizable =)
//
// This also makes logical sense!  If you accidentally overflow or underflow your bit's value by ±1 or ±2, the system won't
// consider it to be in a logically acceptable "nil" state - instead, it -should- panic immediately from a sanity check.
//
// NOTE: Nil is not used in low-level calculations, only in higher level abstractions.
const Nil Bit = 219

// True is a constantly referenceable true.
var True bool = true

// False is a constantly referenceable false.
var False bool = false

// WordWidth is the bit width of a standard int, which for all reasonable intents and purposes matches the architecture's word width.
const WordWidth = strconv.IntSize // NOTE: While this could mismatch on exotic hardware, this is just a convenience value.

/**
Bit
*/

// Bit represents one binary place. [0 - 1]
//
// NOTE: This has a memory footprint of 8 bits.
type Bit byte

// String converts the provided Bit to a string "1", "0", or "nil" - or panics if the found value is anything else.
func (b Bit) String() string {
	switch b {
	case Zero:
		return "0"
	case One:
		return "1"
	case Nil:
		return "nil"
	default:
		panic(ErrorNotABit)
	}
}

/**
Directionality
*/

// Direction represents general directionality and includes both cardinal and abstract reference points in time and space.
//
// Cardinal references consider that the target (the result of calculation) is always relatively "down," or "towards the enemies gate,"
// no matter YOUR orientation in space.  Mentally this may be the direction of "gravity" while standing up and writing calculations on
// a whiteboard, but I think Ender described it best.
//
// Abstract references consider the target's relationality towards YOU as you float through the void of time and spatial calculation.
//
// Concrete references each have an implicit contextual purpose -
//
//	  South - Calculation
//	   West - Scale
//	  North - Accumulation
//	   East - Reduction
//	 Future - Anticipation
//	Present - Experience
//	   Past - Reflection
//
// See South, West, North, East, Future, Present, Past, Up, Up, Down, Down, Left, Right, Left, Right, B, A, Start
type Direction int

const (
	// South represents the cardinal Direction "down" - which is the -target- of all calculation.
	//
	// See West, North, East.
	South Direction = iota

	// West represents the cardinal Direction "left" - which is the direction of scale.
	//
	// See South, North, East.
	West

	// North represents the cardinal Direction "up" - which is the direction of accumulation.
	//
	// See South, West, East.
	North

	// East represents the cardinal Direction "right" - which is the direction of reduction.
	//
	// See South, West, North.
	East

	// Future represents the abstract temporal Direction of "eminently" - which is the direction of anticipation.
	//
	// See Present, Past.
	Future

	// Present represents the abstract temporal Direction of "currently" - which is the direction of experience.
	//
	// See Future, Past.
	Present

	// Past represents the abstract temporal Direction of "historically" - which is the direction of reflection.
	//
	// See Future, Present.
	Past

	// Up represents the abstract Direction of presently relative "up."
	//
	// See Down, Left, Right, Forward, Backward.
	Up // Up Down Down Left Right A B Start

	// Down represents the abstract Direction of presently relative "down."
	//
	// See Up, Left, Right, Forward, Backward.
	Down // Down Left Right A B Start

	// Left represents the abstract Direction of presently relative "left."
	//
	// See Up, Down, Right, Forward, Backward.
	Left // Right A B Start

	// Right represents the abstract Direction of presently relative "right."
	//
	// See Up, Down, Left, Forward, Backward.
	Right // A B Start

	// Forward represents the abstract Direction of presently relative "forward."
	//
	// See Up, Down, Left, Right, Backward.
	Forward

	// Backward represents the abstract Direction of presently relative "backward."
	//
	// See Up, Down, Left, Right, Forward.
	Backward
)

// String prints a single-character representation of the Direction -
//
//	   South: S
//	    West: W
//	   North: N
//	    East: E
//
//	  Future: ⏭
//	 Present: ⏸
//	    Past: ⏮
//
//	      Up: ↑
//	    Down: ↓
//	    Left: ←
//	   Right: →
//	 Forward: ↷
//	Backward: ↶
func (d Direction) String() string {
	switch d {
	case South:
		return "S"
	case West:
		return "W"
	case North:
		return "N"
	case East:
		return "E"
	case Future:
		return "⏭"
	case Present:
		return "⏸"
	case Past:
		return "⏮"
	case Up:
		return "↑"
	case Down:
		return "↓"
	case Left:
		return "←"
	case Right:
		return "→"
	case Forward:
		return "↷"
	case Backward:
		return "↶"
	default:
		return "Unknown"
	}
}

/**
Encoding
*/

// Encoding represents the encoding scheme of a Phrase of Measurement points.
//
// Raw indicates this is simply a phrase of arbitrarily long binary information.
//
// Logical indicates this phrase entirely consists of data measurements.
//
// Signed indicates the first measurement is a sign, followed by a value.
//
// Float indicates the first measurement is a sign, followed by an exponent, and lastly a mantissa.
//
// Index indicates the phrase entirely consists of logical binary data bound to a fixed width.
type Encoding int

const (
	// Raw indicates this is simply a phrase of arbitrarily long binary information.
	Raw Encoding = iota

	// Logical indicates this phrase entirely consists of data measurements.
	Logical

	// Signed indicates the first measurement is a sign, followed by a value.
	Signed

	// Float indicates the first measurement is a sign, followed by an exponent, and lastly a mantissa.
	Float

	// Index indicates the phrase entirely consists of logical binary data bound to a fixed width.
	Index
)

// String prints the full word representation of the encoding scheme.
func (e Encoding) String() string {
	switch e {
	case Raw:
		return "Raw"
	case Logical:
		return "Logical"
	case Signed:
		return "Signed"
	case Float:
		return "Float"
	case Index:
		return "Index"
	default:
		return "Unknown"
	}
}

/**
Endianness
*/

// Endianness indicates the logical -byte- ordering of sequential bytes.  All binary data has a most significant side,
// where the binary placeholder has the highest relative value, as well as a least significant side.  The individual bits
// of a byte are colloquially manipulated in most→to→least significant order, but multiple -bytes- worth of information may
// be stored in least←to←most significant order while retaining the individual bit order of each byte.  There are two
// types of endianness -
//
// BigEndian, where the most significant bytes come first - often used in network protocols:
//
//	| Most Sig. Byte  |   Middle Byte   | Least Sig. Byte |
//	| 0 1 0 0 1 1 0 1 | 0 0 1 0 1 1 0 0 | 0 0 0 1 0 1 1 0 | (5,057,558)
//	|        4D       |        2C       |        16       |
//
// LittleEndian, where the least significant bytes come first - used by x86, AMD64, ARM, and the general world over:
//
//		| Least Sig. Byte |   Middle Byte   |  Most Sig. Byte |
//		| 0 0 0 1 0 1 1 0 | 0 0 1 0 1 1 0 0 | 0 1 0 0 1 1 0 1 | (5,057,558)
//		|        16       |        2C       |        4D       |
//	          ⬑  The byte's internal bits remain in most→to→least order
//
// NOTE: While some hardware may physically store bits in least←to←most order internally, Go's shift operators (<< and >>)
// are guaranteed by the language specification to always operate in most→to→least significant order. This, in turn, means
// that bit operations in tiny will -also- work with bits in most→to→least significant order regardless of the underlying
// architecture's physical bit storage order. When reading raw memory, only byte ordering needs to be handled explicitly.
//
// NOTE: Some protocols, like UART, traditionally transmit in least←to←most order, so you may also need to reverse bits
// within bytes when interfacing with such protocols - which we fully support =)
type Endianness int

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

const errorMsgNotABit = "bits must be 0, 1, or 219 (nil)"

// ErrorNotABit is emitted whenever a method expecting a Bit is provided with any other byte value than 1, 0, or 219 (nil) - as Bit is a byte underneath.
var ErrorNotABit = fmt.Errorf(errorMsgNotABit)
