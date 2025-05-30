package tiny

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

// MaxMeasurementBitLength is the maximum number of bits a Measurement can hold.
const MaxMeasurementBitLength = 32

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
const MaxByte = 255

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

/*
*
Error Messages
*/
const errorMeasurementLimit = "measurements are limited to a maximum of 32 bits wide"

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
Primitives
*/

// Primitive represents the primitive patterns, encoded as a note.
//
// See:
//
// 000 | 0000 0000  0000 0000 <- Pattern_Zero
//
// 001 | 1000 0000  0000 0000 <- Pattern_Light
//
// 010 | 1000 0000  1000 0000 <- Pattern_SemiLight
//
// 011 | 1001 0010  0100 1001 <- Pattern_100
//
// 100 | 1010 1010  1010 1010 <- Pattern_10
//
// 101 | 1101 1011  0110 1101 <- Pattern_110
//
// 110 | 1111 1111  0111 1111 <- Pattern_SemiDark
//
// 111 | 1111 1111  1111 1111 <- Pattern_Dark
type Primitive Note

const (
	// Pattern_Zero is all zeros -> 0000 0000  0000 0000
	Pattern_Zero Primitive = iota

	// Pattern_Light is all zeros except the first position -> 1000 0000  0000 0000
	Pattern_Light

	// Pattern_SemiLight is all zeros except the first index of each half -> 1000 0000  1000 0000
	Pattern_SemiLight

	// Pattern_100 is a repeated pattern of 100 -> 1001 0010  0100 1001
	Pattern_100

	// Pattern_100 is a repeated pattern of 10 -> 1010 1010  1010 1010
	Pattern_10

	// Pattern_110 is a repeated pattern of 110 -> 1101 1011  0110 1101
	Pattern_110

	// Pattern_SemiDark is all ones except the first index of the second half -> 1111 1111  0111 1111
	Pattern_SemiDark

	// Pattern_Dark is all ones -> 1111 1111  1111 1111
	Pattern_Dark
)
