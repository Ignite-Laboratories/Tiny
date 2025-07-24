package tiny

import (
	"fmt"
	"math/big"
	"strconv"
)

// Primitive represents the primitive types provided by Go, all of which are convertable between Measurement and primitive form.
//
// In addition, the big.Int and big.Float types are considered "primitive" as they are fully interoperable with the matrix engine.
type Primitive interface {
	*big.Int | *big.Float |
		int8 | int16 | int32 | int64 |
		uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		complex64 | complex128 |
		int | uint | uintptr |
		bool
}

// Binary represents the basic binary types that compose all Operable types.
//
// See Bit, Measurement, Phrase, Natural, Real, Complex, and Index
type Binary interface {
	Bit | byte | Measurement | Phrase
}

// Operable represents any type that can be implicitly aligned and named for performing logical operations.
//
// See Natural, Real, Complex, Index, and Binary
type Operable interface {
	Natural | Real | Complex | Index | Binary
	GetName() string
	Named(string)
}

/**
Global Constants
*/

// EmptyPhrase represents a raw Phrase with no data named "3MP7Y".
var EmptyPhrase = NewPhraseNamed("3MP7Y")

// Unlimited represents a constantly referencable integer value which can be considered a reasonably "unlimited" width.
var Unlimited = ^uint(0)

// Start is a constantly referencable uint{0}.
//
// For a slice, please use Initial.
//
// See Initial, Zero, One, Nil, ZeroOne, OneZero, SingleZero, DoubleZero, SingleOne, and DoubleOne.
var Start uint = 0

// Initial is a constantly referencable []uint{0}.
//
// For a non-slice, please use Start.
//
// See Start, Zero, One, Nil, ZeroOne, OneZero, SingleZero, DoubleZero, SingleOne, and DoubleOne.
var Initial []uint = []uint{0}

// Zero is an implicit Bit{0}.
//
// See Start, Initial, One, Nil, ZeroOne, OneZero, SingleZero, DoubleZero, SingleOne, and DoubleOne.
const Zero Bit = 0

// One is an implicit Bit{1}.
//
// See Start, Initial, Zero, Nil, ZeroOne, OneZero, SingleZero, DoubleZero, SingleOne, and DoubleOne.
const One Bit = 1

// OneZero is an implicit Bit{1, 0}.
//
// See Start, Initial, Zero, Nil, ZeroOne, OneZero, SingleZero, DoubleZero, SingleOne, and DoubleOne.
const OneZero Bit = 1

// ZeroOne is an implicit Bit{0, 1}.
//
// See Start, Initial, Zero, Nil, ZeroOne, OneZero, SingleZero, DoubleZero, SingleOne, and DoubleOne.
const ZeroOne Bit = 1

// SingleZero is an implicit []Bit{0}.
//
// See Start, Initial, Zero, One, Nil, ZeroOne, OneZero, DoubleZero, SingleOne, and DoubleOne.
var SingleZero = []Bit{Zero}

// DoubleZero is an implicit []Bit{0, 0}.
//
// See Start, Initial, Zero, One, Nil, ZeroOne, OneZero, SingleZero, SingleOne, and DoubleOne.
var DoubleZero = []Bit{Zero, Zero}

// SingleOne is an implicit []Bit{1}.
//
// See Start, Initial, Zero, One, Nil, ZeroOne, OneZero, SingleZero, DoubleZero, and DoubleOne.
var SingleOne = []Bit{One}

// DoubleOne is an implicit []Bit{1, 1}.
//
// See Start, Initial, Zero, One, Nil, ZeroOne, OneZero, SingleZero, DoubleZero, and SingleOne.
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
//
// See Start, Initial, Zero, One, ZeroOne, OneZero, SingleZero, DoubleZero, SingleOne, and DoubleOne.
const Nil Bit = 219

// True is a constantly referenceable true.
//
// See False
var True bool = true

// False is a constantly referenceable false.
//
// See True
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
Errors
*/

const errorMsgNotABit = "bits must be 0 or 1 in value"
const errorMsgNotABitWithNil = "bits must be 0, 1, or 219 (nil) in value"

// ErrorNotABit is emitted whenever a method expecting a Bit is provided with any other byte value than 1, 0 - as Bit is a byte underneath.
var ErrorNotABit = fmt.Errorf(errorMsgNotABit)

// ErrorNotABitWithNil is emitted whenever a method expecting a Bit is provided with any other byte value than 1, 0, or 219 (nil) - as Bit is a byte underneath.
var ErrorNotABitWithNil = fmt.Errorf(errorMsgNotABitWithNil)
