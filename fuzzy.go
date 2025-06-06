package tiny

import (
	"math/big"
)

// _fuzzy is a factory for creating or referencing fuzzy projection functions.
type _fuzzy struct {
	SixtyFour      _sixtyFour
	Five           _five
	FiveCumulative _five
	Power          _power
	ZLE            _zle
}

type _sixtyFour struct{}
type _five struct{}
type _fiveCumulative struct{}
type _power struct{}
type _zle struct{}

// Read uses the below map to parse a value from the next bits in the provided phrase:
//
// The Fuzzy.SixtyFour Map
//
//	    Key | Projection Bit Range
//	      1 | 4
//	    0 1 | 8
//	  0 0 1 | 16
//	0 0 0 1 | 32
//	0 0 0 0 | 64
func (_ _sixtyFour) Read(data Phrase) (value int, remainder Phrase) {
	zeros, remainder := data.ReadUntilOne(4)
	var projectionRange int

	switch zeros {
	case 0:
		projectionRange = 4
	case 1:
		projectionRange = 8
	case 2:
		projectionRange = 16
	case 3:
		projectionRange = 32
	case 4:
		projectionRange = 64
	}

	projection, remainder := remainder.Read(projectionRange)
	return To.Number(projectionRange, projection.Bits()...), remainder
}

// Encode uses the below map to encode a ZLE key and projection from the provided value.
//
// The Fuzzy.SixtyFour Map
//
//	    Key | Projection Bit Range
//	      1 | 4
//	    0 1 | 8
//	  0 0 1 | 16
//	0 0 0 1 | 32
//	0 0 0 0 | 64
func (_ _sixtyFour) Encode(value int) (key Phrase, projection Phrase) {
	var bitLength int
	switch {
	case value < 1<<4:
		bitLength = 4
		key = NewPhraseFromBits(1)
	case value < 1<<8:
		bitLength = 8
		key = NewPhraseFromBits(0, 1)
	case value < 1<<16:
		bitLength = 16
		key = NewPhraseFromBits(0, 0, 1)
	case value < 1<<32:
		bitLength = 32
		key = NewPhraseFromBits(0, 0, 0, 1)
	case value < 1<<64:
		bitLength = 64
		key = NewPhraseFromBits(0, 0, 0, 0)
	}

	return key, NewPhraseFromBits(From.Number(value, bitLength)...)
}

// Read uses the below map to parse a value from the next bits in the provided phrase:
//
// The Fuzzy.Five Map
//
//	    Key | Projection | Range
//	      1 |     1      | 0 - 1
//	    0 1 |     2      | 0 - 3
//	  0 0 1 |     3      | 0 - 7
//	0 0 0 1 |     4      | 0 - 15
//	0 0 0 0 |     5      | 0 - 31
func (_ _five) Read(data Phrase) (value int, remainder Phrase) {
	zeros, remainder := data.ReadUntilOne(4)
	var projectionRange int

	switch zeros {
	case 0:
		projectionRange = 1
	case 1:
		projectionRange = 2
	case 2:
		projectionRange = 3
	case 3:
		projectionRange = 4
	case 4:
		projectionRange = 5
	}

	projection, remainder := remainder.Read(projectionRange)
	return To.Number(projectionRange, projection.Bits()...), remainder
}

// Encode uses the below map to encode a ZLE key and projection from the provided value.
//
// The Fuzzy.Five Map
//
//	    Key | Projection | Range
//	      1 |     1      | 0 - 1
//	    0 1 |     2      | 0 - 3
//	  0 0 1 |     3      | 0 - 7
//	0 0 0 1 |     4      | 0 - 15
//	0 0 0 0 |     5      | 0 - 31
func (_ _five) Encode(value int) (key Phrase, projection Phrase) {
	var bitLength int
	switch {
	case value < 1<<1:
		bitLength = 1
		key = NewPhraseFromBits(1)
	case value < 1<<2:
		bitLength = 2
		key = NewPhraseFromBits(0, 1)
	case value < 1<<3:
		bitLength = 3
		key = NewPhraseFromBits(0, 0, 1)
	case value < 1<<4:
		bitLength = 4
		key = NewPhraseFromBits(0, 0, 0, 1)
	case value < 1<<5:
		bitLength = 5
		key = NewPhraseFromBits(0, 0, 0, 0)
	default:
		panic("input value too large for a five-bit map")
	}

	return key, NewPhraseFromBits(From.Number(value, bitLength)...)
}

// Read uses the below map to parse a value from the next bits in the provided phrase:
//
// The Fuzzy.Five Map - Cumulative
//
//	    Key | Projection | Range  | Cumulative Interpretation
//	      1 |     1      | 0 - 1  |  0 - 1
//	    0 1 |     2      | 0 - 3  |  2 - 5
//	  0 0 1 |     3      | 0 - 7  |  6 - 13
//	0 0 0 1 |     4      | 0 - 15 | 14 - 29
//	0 0 0 0 |     5      | 0 - 31 | 30 - 61
func (_ _fiveCumulative) Read(data Phrase) (value int, remainder Phrase) {
	zeros, remainder := data.ReadUntilOne(4)
	var projectionRange int
	var shim int

	switch zeros {
	case 0:
		projectionRange = 1
	case 1:
		projectionRange = 2
		shim += 2
	case 2:
		projectionRange = 3
		shim += 6
	case 3:
		projectionRange = 4
		shim += 14
	case 4:
		projectionRange = 5
		shim += 30
	}

	projection, remainder := remainder.Read(projectionRange)
	return To.Number(projectionRange, projection.Bits()...) + shim, remainder
}

// Encode uses the below map to encode a ZLE key and projection from the provided value.
//
// The Fuzzy.Five Map - Cumulative
//
//	    Key | Projection | Range  | Cumulative Interpretation
//	      1 |     1      | 0 - 1  |  0 - 1
//	    0 1 |     2      | 0 - 3  |  2 - 5
//	  0 0 1 |     3      | 0 - 7  |  6 - 13
//	0 0 0 1 |     4      | 0 - 15 | 14 - 29
//	0 0 0 0 |     5      | 0 - 31 | 30 - 61
func (_ _fiveCumulative) Encode(value int) (key Phrase, projection Phrase) {
	var bitLength int
	switch {
	case value < 2:
		bitLength = 1
		key = NewPhraseFromBits(1)
	case value < 6:
		bitLength = 2
		value -= 2
		key = NewPhraseFromBits(0, 1)
	case value < 14:
		bitLength = 3
		value -= 6
		key = NewPhraseFromBits(0, 0, 1)
	case value < 30:
		bitLength = 4
		value -= 30
		key = NewPhraseFromBits(0, 0, 0, 1)
	case value < 62:
		bitLength = 5
		key = NewPhraseFromBits(0, 0, 0, 0)
	default:
		panic("input value too large for a cumulative five-bit map")
	}

	return key, NewPhraseFromBits(From.Number(value, bitLength)...)
}

// Read uses the below map to parse a value from the next bits in the provided phrase:
//
// The Fuzzy.Power Map
//
//	    Key | Projection | Value Range | Power Interpretation
//	      1 |      2     |   1 - 4     |      2ⁿ - 1
//	    0 1 |      3     |   1 - 8     |      2ⁿ - 1
//	  0 0 1 |      4     |   1 - 16    |      2ⁿ - 1
//	0 0 0 1 |      5     |   1 - 32    |      2ⁿ - 1
//	0 0 0 0 |      6     |   1 - 64    |      2ⁿ - 1
func (_ _power) Read(data Phrase) (value int, remainder Phrase) {
	zeros, remainder := data.ReadUntilOne(4)
	var projectionRange int

	switch zeros {
	case 0:
		projectionRange = 2
	case 1:
		projectionRange = 3
	case 2:
		projectionRange = 4
	case 3:
		projectionRange = 5
	case 4:
		projectionRange = 6
	}

	projection, remainder := remainder.Read(projectionRange)
	power := To.Number(projectionRange, projection.Bits()...)
	power += 1
	return 1<<power - 1, remainder
}

// Encode uses the below map to encode a ZLE key and projection from the provided value.
//
// NOTE: When encoding this value, you provide the exponent as the value.
//
// The Fuzzy.Power Map
//
//	    Key | Projection | Value Range | Power Interpretation
//	      1 |      2     |   1 - 4     |      2ⁿ - 1
//	    0 1 |      3     |   1 - 8     |      2ⁿ - 1
//	  0 0 1 |      4     |   1 - 16    |      2ⁿ - 1
//	0 0 0 1 |      5     |   1 - 32    |      2ⁿ - 1
//	0 0 0 0 |      6     |   1 - 64    |      2ⁿ - 1
func (_ _power) Encode(power int) (key Phrase, projection Phrase) {
	var bitLength int
	power -= 1

	switch {
	case power < 1<<2:
		bitLength = 1
		key = NewPhraseFromBits(1)
	case power < 1<<3:
		bitLength = 2
		key = NewPhraseFromBits(0, 1)
	case power < 1<<4:
		bitLength = 3
		key = NewPhraseFromBits(0, 0, 1)
	case power < 1<<5:
		bitLength = 4
		key = NewPhraseFromBits(0, 0, 0, 1)
	case power < 1<<6:
		bitLength = 5
		key = NewPhraseFromBits(0, 0, 0, 0)
	default:
		panic("input value too large for a five-bit map")
	}

	return key, NewPhraseFromBits(From.Number(power, bitLength)...)
}

// Read uses the below map to parse a value from the next bits in the provided phrase:
//
// The Fuzzy.ZLE Map
//
//	NOTE: This will overflow if you let it read too far =)
//
//	        Key | Projection
//	          1 | 1 [2⁰]
//	        0 1 | 2 [2¹]
//	      0 0 1 | 4 [2²]
//	    0 0 0 1 | 8 [2³]
//	           ...
//	      𝑛   1 | 2ⁿ
func (_ _zle) Read(data Phrase) (value int, remainder Phrase) {
	zeros, remainder := data.ReadUntilOne()
	projectionRange := 1 << zeros
	projection, remainder := remainder.Read(projectionRange)
	return To.Number(projectionRange, projection.Bits()...), remainder
}

// Encode uses the below map to encode a ZLE key and projection from the provided value.
//
// NOTE: This always returns an empty projection.
//
// The Fuzzy.ZLE Map
//
//	NOTE: This will overflow if you let it read too far =)
//
//	        Key | Projection
//	          1 | 1 [2⁰]
//	        0 1 | 2 [2¹]
//	      0 0 1 | 4 [2²]
//	    0 0 0 1 | 8 [2³]
//	           ...
//	      𝑛   1 | 2ⁿ
func (_ _zle) Encode(power int) (key Phrase, projection Phrase) {
	return Synthesize.Zeros(power).AppendBits(1), Phrase{}
}

// Approximation represents an synthetically generated approximate value.
//
// Indices - Provides the four index points to synthesize a known bit range with.
//
//	 Index 0 represents the first ⅛th of the target bits and is approximated at 4x the resolution
//		Index 1 represents the second ⅛th and is approximated at 4x the resolution
//		Index 2 represents the second ¼ and is approximated at 2x the resolution
//		Index 3 represents the final ½ and is approximated at standard resolution
//
// Value - Provides the value of the synthesized binary data.
//
// Target - Gives the target value this approximation attempted to fuzzily replicate.
//
// Delta - Gives the absolute value of the difference between the Value and Target.
//
// Relativity - Dictates if the approximated value is relativistically smaller or larger than the target.
type Approximation struct {
	Indices    Passage
	Value      *big.Int
	Target     *big.Int
	Delta      *big.Int
	Relativity RelativeSize
}

// Approximation creates a synthetic approximation of the target's bits at four different scales and returns
// the approximation indices, approximation, delta, and whether the approximation is larger or smaller than the target.
//
// The approximation itself is a phrase of four indices, each representing the subdivision index of that
// particular index's region.
//
// NOTE: The standard minimum resolution bit width is a Note (3 bits) but you may provide your own bit
// width, if desired.
// The resolution bit width defines the maximum value of the minimum resolution to subdivide at.
//
//	Index 0 represents the first ⅛th of the target bits and is approximated at 4x the resolution
//	Index 1 represents the second ⅛th and is approximated at 4x the resolution
//	Index 2 represents the second ¼ and is approximated at 2x the resolution
//	Index 3 represents the final ½ and is approximated at the minimum bit width's resolution
//
// This yields the following breakdown for a 64-bit melody:
//
//	|                             64 Bit Melody                             |
//	 10110100 10101101 00100110 10010101 00101110 10100101 10100100 00111011
//	|Index 0 | Index 1|    Index 2      |             Index 3               |
//
// NOTE: The indices bit-widths are subdivided using flooring, meaning the last index always holds the excess bits.
//
// For example, with a 67 bit input:
//
//	|                             64 Bit Melody                             |   |
//	 10110100 10101101 00100110 10010101 00101110 10100101 10100100 00111011 110
//	|Index 0 | Index 1|    Index 2      |               Index 3                 |
//
// Above, 67/8 = 8.375 so the ⅛ indices are 8 bits while 67/4 = 16.75 so the ¼ index is 16 bits.
// Finally, the ½ index picks up whatever remaining bits are leftover.
//
// Whereas, with a 68 bit input:
//
//	|                              64 Bit Melody                             |    |
//	 10110100 10101101 00100110 10010101 0 0101110 10100101 10100100 00111011 1101
//	|Index 0 | Index 1|     Index 2       |               Index 3                 |
//
// Above, 68/8 = 8.5 so the ⅛ indices are still 8 bits while 68/4 = 17 so the ¼ index grows to 17 bits.
// Finally, the ½ index picks up whatever remaining bits are leftover.
func (_ _fuzzy) Approximation(target *big.Int, minResolution ...int) Approximation {
	var approx Approximation
	approx.Target = target

	bitWidth := 3
	if len(minResolution) > 0 {
		bitWidth = minResolution[0]
	}

	bitWidth2x := bitWidth + 1
	bitWidth4x := bitWidth + 2

	resolutionMax := To.Number(bitWidth, Synthesize.Ones(bitWidth).Bits()...)
	resolutionMax2x := To.Number(bitWidth2x, Synthesize.Ones(bitWidth2x).Bits()...)
	resolutionMax4x := To.Number(bitWidth4x, Synthesize.Ones(bitWidth4x).Bits()...)
	bitLength := target.BitLen()

	eighth := bitLength / 8
	quarter := bitLength / 4
	phrase := NewPhraseFromBigInt(target)

	region0, phrase := phrase.Read(eighth)
	region1, phrase := phrase.Read(eighth)
	region2, phrase := phrase.Read(quarter)
	region3 := phrase

	fuzzy0, index0 := Synthesize.Approximation(region0.AsBigInt(), resolutionMax4x, eighth)
	fuzzy1, index1 := Synthesize.Approximation(region1.AsBigInt(), resolutionMax4x, eighth)
	fuzzy2, index2 := Synthesize.Approximation(region2.AsBigInt(), resolutionMax2x, quarter)
	fuzzy3, index3 := Synthesize.Approximation(region3.AsBigInt(), resolutionMax, bitLength-eighth-eighth-quarter)

	indexBits0 := From.Number(index0, bitWidth4x)
	indexBits1 := From.Number(index1, bitWidth4x)
	indexBits2 := From.Number(index2, bitWidth2x)
	indexBits3 := From.Number(index3, bitWidth)

	approx.Value = NewPhraseFromBits(fuzzy0...).AppendBits(fuzzy1...).AppendBits(fuzzy2...).AppendBits(fuzzy3...).AsBigInt()
	approx.Indices = NewPassage(NewPhraseFromBits(indexBits0...), NewPhraseFromBits(indexBits1...), NewPhraseFromBits(indexBits2...), NewPhraseFromBits(indexBits3...))

	approx.Relativity = NewRelativeSize(approx.Value.Cmp(target))
	if approx.Relativity == Equal {
		approx.Delta = new(big.Int)
	} else if approx.Relativity == Smaller {
		approx.Delta = new(big.Int).Sub(target, approx.Value)
	} else {
		approx.Delta = new(big.Int).Sub(approx.Value, target)
	}

	return approx
}
