package tiny

import "fmt"

// Zero is an implicit Bit{0}.
const Zero Bit = 0

// One is an implicit Bit{1}.
const One Bit = 1

// ZeroZero is an implicit Crumb{00}.
const ZeroZero Crumb = 0

// ZeroOne is an implicit Crumb{01}.
const ZeroOne Crumb = 1

// OneZero is an implicit Crumb{10}.
const OneZero Crumb = 2

// OneOne is an implicit Crumb{11}.
const OneOne Crumb = 3

// MaxCrumb is the maximum value a Crumb can hold.
const MaxCrumb = 3

// MaxNote is the maximum value a Note can hold.
const MaxNote = 7

// MaxNibble is the maximum value a Nibble can hold.
const MaxNibble = 15

// MaxFlake is the maximum value a Flake can hold.
const MaxFlake = 31

// MaxMorsel is the maximum value a Morsel can hold.
const MaxMorsel = 63

// MaxShred is the maximum value a Shred can hold.
const MaxShred = 127

// MaxByte is the maximum value a byte can hold.
const MaxByte = (1 << 8) - 1

// MaxScale is the maximum value a Scale can hold.
const MaxScale = (1 << 12) - 1

// MaxMotif is the maximum value a Motif can hold.
const MaxMotif = (1 << 16) - 1

// MaxRiff is the maximum value a Riff can hold.
const MaxRiff = (1 << 24) - 1

// MaxCadence is the maximum value a Cadence can hold.
const MaxCadence = (1 << 32) - 1

// MaxHook is the maximum value a Hook can hold.
const MaxHook = (1 << 48) - 1

// WidthBit is the number of binary positions a Bit represents.
const WidthBit = 1

// WidthCrumb is the number of binary positions a Crumb represents.
const WidthCrumb = 2

// WidthNote is the number of binary positions a Note represents.
const WidthNote = 3

// WidthNibble is the number of binary positions a Nibble represents.
const WidthNibble = 4

// WidthFlake is the number of binary positions a Flake represents.
const WidthFlake = 5

// WidthMorsel is the number of binary positions a Morsel represents.
const WidthMorsel = 6

// WidthShred is the number of binary positions a Shred represents.
const WidthShred = 7

// WidthByte is the number of binary positions a Byte represents.
const WidthByte = 8

// WidthRun is the number of binary positions a Run represents
const WidthRun = 10

// WidthScale is the number of binary positions a Scale represents.
const WidthScale = 12

// WidthMotif is the number of binary positions a Motif represents.
const WidthMotif = 16

// WidthRiff is the number of binary positions a Riff represents.
const WidthRiff = 24

// WidthCadence is the number of binary positions a Cadence represents.
const WidthCadence = 32

// WidthHook is the number of binary positions a Hook represents.
const WidthHook = 48

/*
*
Error Messages
*/

const errorMeasurementLimit = "measurements are limited to the bit-width of your system's architecture"

const errorPassageLimit = "passages are limited to 256 bits in length"

const ErrorMsgEndOfBits = "no more bits to read"

var ErrorEndOfBits = fmt.Errorf(ErrorMsgEndOfBits)

/**
Passages
*/

// MaxPassage is the maximum length a passage can contain, which is specifically set to 2⁸.
//
// This value was chosen to keep the synthesis process extremely performant through concurrency.
const MaxPassage = 1 << 8

// NOTE: Update errorPassageLimit if you change this!!!

/**
Movements
*/

// MovementStart identifies the start region of a DNA file.
const MovementStart = "start"

// MovementPathway identifies the pathway region of a DNA file.
const MovementPathway = "pathway"

// MovementSeed identifies the seed region of a DNA file.
const MovementSeed = "seed"

/**
Relativity
*/

// Relativity represents the abstract logical relationship of two entities, 𝑎 and 𝑏.
//
// Rather than imbuing 'size', 'value', or 'position', this aims to describe that '𝑎' has
// a logical relationship with '𝑏' that's understood contextually by the caller.  Whether
// in an ordered list, comparing the physical dimensionality, or general timing - this provides
// a common language for describing the relationship between both entities.
//
// See Before, Same, After
type Relativity int

const (
	// Before indicates that 𝑎 logically comes before 𝑏.
	Before Relativity = -1
	// Same indicates that 𝑎 and 𝑏 are logically equal.
	Same = 0
	// After indicates that 𝑎 logically comes after 𝑏.
	After Relativity = 1
)

// NewRelativeSize creates a new Relativity structure.
//
// If the value is 0, Same is returned.
// If the value is positive, After is returned.
// If the value is negative, Before is returned.
func NewRelativeSize(value int) Relativity {
	switch {
	case value > 0:
		return After
	case value < 0:
		return Before
	default:
		return Same
	}
}
