package tiny

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
	default:
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
