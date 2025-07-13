package tiny

import (
	"fmt"
	"strconv"
)

/**
Constants
*/

// Bit represents one binary place. [0 - 1]
type Bit byte

// Zero is an implicit Bit{0}.
const Zero Bit = 0

// One is an implicit Bit{1}.
const One Bit = 1

// WordWidth is the bit width of a standard int, which for all reasonable intents and purposes matches the architecture's word width.
const WordWidth = strconv.IntSize // NOTE: While this could mismatch the architecture's word in some cases, the performance implications are minimal.

// Endianness indicates the logical -byte- ordering of sequential bytes.  All binary data has a most significant side,
// where the binary placeholder has the highest relative value, as well as a least significant side.  The individual bits
// of a byte are always stored in most→to→least significant order, but multiple bytes worth of information may be stored
// in least←to←most significant order while retaining the individual bit order of each byte.  There are two types of
// endianness -
//
// BigEndian, where the most significant bytes come first - often used in network protocols:
//
//	| Most Significant                  Least Significant |
//	| 0 1 0 0 1 1 0 1 - 0 0 1 0 1 1 0 0 - 0 0 0 1 0 1 1 0 | (5,057,558)
//	|        77       |        44       |        22       |
//
// LittleEndian, where the least significant bytes come first - used by x86, AMD64, ARM, and many more:
//
//	| Least Significant                  Most Significant |
//	| 0 0 0 1 0 1 1 0 - 0 0 1 0 1 1 0 0 - 0 1 0 0 1 1 0 1 | (5,057,558)
//	|        22       |        44       |        77       |
type Endianness int

const (
	LittleEndian Endianness = iota
	BigEndian
)

/**
Errors
*/

const errorMsgMeasurementLimit = "measurements are limited to the bit width of an int"

// ErrorMeasurementLimit is emitted whenever a Measurement attempts to hold more than a WordWidth of binary information.
var ErrorMeasurementLimit = fmt.Errorf(errorMsgMeasurementLimit)

const errorMsgEndOfBits = "no more bits to read"

// ErrorEndOfBits is emitted whenever a read operation has requested to read beyond the available binary data's width.
// This error can absolutely be ignored, but also allows one to implicitly read until the end of bits has been reached.
var ErrorEndOfBits = fmt.Errorf(errorMsgEndOfBits)

const errorMsgNotABit = "bits must be 0 or 1"

// ErrorNotABit is emitted whenever a method expecting a Bit is provided with any other byte value than 1 or 0, as Bit is a byte underneath.
var ErrorNotABit = fmt.Errorf(errorMsgNotABit)
