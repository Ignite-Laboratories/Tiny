package tiny

import (
	"math"
)

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

// MaxMelody is the maximum value a Melody can hold.
const MaxMelody = (1 << 64) - 1

// MaxVerse is the maximum value a Verse can hold.
const MaxVerse = (1 << 128) - 1

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

// WidthMelody is the number of binary positions a Melody represents.
const WidthMelody = 64

// WidthVerse is the number of binary positions a Verse represents.
const WidthVerse = 128

/*
*
Error Messages
*/
const errorMeasurementLimit = "measurements are limited to the bit-width of your system's architecture"

/**
Passages
*/

// PassageMaxVerses defines the maximum number of Verses a Passage may hold.
//
// A Verse contains 16 bytes.
const PassageMaxVerses = 64

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
RelativeSize
*/

// RelativeSize represents the relationship of two values (A and B).
//
// See Smaller, Equal, Larger
type RelativeSize int

const (
	// Smaller indicates that A is smaller than B.
	Smaller RelativeSize = -1
	// Equal indicates that A and B are equal.
	Equal = 0
	// Larger indicates that A is larger than B.
	Larger RelativeSize = 1
)

// NewRelativeSize creates a new RelativeSize structure.
//
// If the value is 0, Equal is returned.
// If the value is positive, Larger is returned.
// If the value is negative, Smaller is returned.
func NewRelativeSize(value int) RelativeSize {
	switch {
	case value > 0:
		return Larger
	case value < 0:
		return Smaller
	default:
		return Equal
	}
}

/**
Correction Factor Focus
*/

var factors0 = [][]float64{
	GenerateCorrectionFactors(focus000Log, 7),
	GenerateCorrectionFactors(focus001Log, 15),
	GenerateCorrectionFactors(focus010Log, 63),
	GenerateCorrectionFactors(focus011Log, 255),
}
var factors1 = [][]float64{
	GenerateCorrectionFactors(focus100Log, 7),
	GenerateCorrectionFactors(focus101Log, 15),
	GenerateCorrectionFactors(focus110Log, 63),
	GenerateCorrectionFactors(focus111Log, 255),
}

var focus000Log = (math.Log(1.05) - math.Log(1.0)) / float64(7)
var focus001Log = (math.Log(1.05) - math.Log(1.0)) / float64(15)
var focus010Log = (math.Log(1.05) - math.Log(1.0)) / float64(63)
var focus011Log = (math.Log(1.05) - math.Log(1.0)) / float64(255)
var focus100Log = (math.Log(0.95) - math.Log(1.0)) / float64(7)
var focus101Log = (math.Log(0.95) - math.Log(1.0)) / float64(15)
var focus110Log = (math.Log(0.95) - math.Log(1.0)) / float64(63)
var focus111Log = (math.Log(0.95) - math.Log(1.0)) / float64(255)

func GenerateCorrectionFactors(log float64, max int) []float64 {
	out := make([]float64, max)
	for i := 0; i < max; i++ {
		out[i] = math.Exp(float64(i) * log)
	}
	return out
}
