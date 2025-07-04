package backup

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
	"math"
	"math/big"
)

// _fuzzy is a factory for creating or referencing fuzzy projection functions.
type _fuzzy int

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
	Indices    tiny.Passage
	Value      *big.Int
	Target     *big.Int
	Delta      *big.Int
	Relativity tiny.RelativeSize
}

// CorrectionFactor represents a synthetic value with three parts:
//
// Threshold - A threshold of '1' indicates that the correction factor
// is above 1.0, while a threshold of '0' indicates it's below 1.0.
//
// Focus - This crumb indicates how this region should be "focused in" on.
// This defines both the resolution bit width for the value as well as the
// target factor range to subdivide a synthetic value from.
// The focus crumb key can be decoded using the below table:
//
//	Focus | Bit Width | Resolution
//	  00  |     3     | 8
//	  01  |     4     | 16
//	  10  |     6     | 64
//	  11  |     8     | 256
//
// The above key also interprets the factor range boundary, when combined with the threshold:
//
//	  Threshold   | Focus | Boundary
//	0 (below 1.0) |   00  |  1.050
//	0 (below 1.0) |   01  |  1.100
//	0 (below 1.0) |   10  |  1.250
//	0 (below 1.0) |   11  |  1.500
//	--------------------------------
//	1 (above 1.0) |   00  |  0.750
//	1 (above 1.0) |   01  |  0.875
//	1 (above 1.0) |   10  |  0.900
//	1 (above 1.0) |   11  |  0.950
//
// Value - This is a variable width region of bits that indicates the
// subdivision index to factor against the approximation.
//
// The value is interpreted as a logarithmic step between 1.0 and the defined Boundary, focused
// towards 1.0 as much as possible.
// Calculating this value is simple:
//
//		             𝑛 = bit width
//			resolution = 2ⁿ
//		             x = value ⋅ ( (ln(boundary) - ln(1.0)) / resolution )
//	         result = ℯˣ
//
// Because of this, the stride log values have already been precalculated and are available as constants.
// See Focus000Log through Focus111Log
//
// In order to go from a stride value to the desired factor value, use CalculateFocusValue
type CorrectionFactor struct {
	Threshold tiny.Bit
	Focus     tiny.Crumb
	Value     []tiny.Bit
}

// Count returns a function that will return true the requested number of times.
func (_ _fuzzy) Count(value int) func(tiny.Bit) bool {
	i := 0
	return func(b tiny.Bit) bool {
		i++
		return i < value
	}
}

// ZLEKey reads up to four bits or until a value of 1 is reached.
// This will yield a Zero Length Encoding key that can be parsed using tiny.Fuzzy.ParseZLE64
//
// If you would like to read a ZLE key longer than 4 bits, you may provide an upper limit.
//
// If you wish for no upper limit (just read until EOD or a 1) then provide <= 0 as the upper limit..
func (_ _fuzzy) ZLEKey(upperLimit ...int) func(tiny.Bit) bool {
	limit := 4
	if len(upperLimit) > 0 {
		limit = upperLimit[0]
	}

	i := 0
	return func(b tiny.Bit) bool {
		i++
		if limit <= 0 {
			return b == tiny.Zero
		}
		return b == tiny.Zero && i < limit
	}
}

// ParseZLEScaled uses the provided Zero Length Encoding key to calculate how many more bits to read.
//
// This particular flavor of ZLE will yield an addressable range up to 64 bits wide while prioritizing
// a minimal number of bits under that length.
//
// NOTE: This function merely gives the bit ranges for each key entry.
//
//	ZLE Key | Bit Range | Value Range
//	      1 |     2     |   0-3
//	    0 1 |     3     |   0-2³ + 4 (4-11)
//	  0 0 1 |     8     |   0-2⁸ + 12 (12-267)
//	0 0 0 0 |    16     |   0-2¹⁶
//	0 0 0 1 |    64     |   0-2⁶⁴
func (_ _fuzzy) ParseZLEScaled(key tiny.Measurement) int {
	switch bits := key.Bits; {
	case len(bits) == 1 && key.Value() == 1:
		return 2
	case len(bits) == 2 && key.Value() == 1:
		return 3
	case len(bits) == 3 && key.Value() == 1:
		return 8
	case len(bits) == 4 && key.Value() == 0:
		return 16
	case len(bits) == 4 && key.Value() == 1:
		return 64
	default:
		panic(fmt.Sprintf("invalid scaled ZLE key: %v", key.Bits))
	}
}

// InterpretZLEScaled returns the -interpreted- value of a scaled ZLE phrase.
//
//	ZLE Key | Bit Range | Value Range
//	      1 |     2     |   0-3
//	    0 1 |     3     |   0-2³ + 4 (4-11)
//	  0 0 1 |     8     |   0-2⁸ + 12 (12-267)
//	0 0 0 0 |    16     |   0-2¹⁶
//	0 0 0 1 |    64     |   0-2⁶⁴
func (_ _fuzzy) InterpretZLEScaled(passage tiny.Passage) int {
	key := passage[0][0]
	projection := passage[1]
	switch bits := key.Bits; {
	case len(bits) == 1 && key.Value() == 1:
		return tiny.To.Number(2, projection.Bits()...)
	case len(bits) == 2 && key.Value() == 1:

		return tiny.To.Number(3, projection.Bits()...) + 4
	case len(bits) == 3 && key.Value() == 1:
		return tiny.To.Number(8, projection.Bits()...) + 12
	case len(bits) == 4 && key.Value() == 0:
		return tiny.To.Number(16, projection.Bits()...)
	case len(bits) == 4 && key.Value() == 1:
		return int(projection.AsBigInt().Int64())
	default:
		panic(fmt.Sprintf("invalid scaled ZLE key: %v", key.Bits))
	}
}

// ParseZLE64 uses the provided Zero Length Encoding key to calculate how many more bits to read.
//
// This particular flavor of ZLE will yield an addressable range up to 64 bits wide.
//
//	ZLE Key | Bit Range
//	      1 | 4
//	    0 1 | 8
//	  0 0 1 | 16
//	0 0 0 0 | 32
//	0 0 0 1 | 64
func (_ _fuzzy) ParseZLE64(key tiny.Measurement) int {
	switch bits := key.Bits; {
	case len(bits) == 1 && key.Value() == 1:
		return 4
	case len(bits) == 2 && key.Value() == 1:
		return 8
	case len(bits) == 3 && key.Value() == 1:
		return 16
	case len(bits) == 4 && key.Value() == 0:
		return 32
	case len(bits) == 4 && key.Value() == 1:
		return 64
	default:
		panic(fmt.Sprintf("invalid ZLE key: %v", key.Bits))
	}
}

// ParseZLE5 uses the provided Zero Length Encoding key to calculate how many more bits to read.
//
// This particular flavor of ZLE will yield an addressable range up to 5 bits wide.
//
//	ZLE Key | Bit Range
//	      1 | 1
//	    0 1 | 2
//	  0 0 1 | 3
//	0 0 0 0 | 4
//	0 0 0 1 | 5
func (_ _fuzzy) ParseZLE5(key tiny.Measurement) int {
	switch bits := key.Bits; {
	case len(bits) == 1 && key.Value() == 1:
		return 1
	case len(bits) == 2 && key.Value() == 1:
		return 2
	case len(bits) == 3 && key.Value() == 1:
		return 3
	case len(bits) == 4 && key.Value() == 0:
		return 4
	case len(bits) == 4 && key.Value() == 1:
		return 5
	default:
		panic(fmt.Sprintf("invalid micro ZLE key: %v", key.Bits))
	}
}

// ParseZLE uses the provided Zero Length Encoding key to calculate how many more bits to read.
//
// This returns back 2ⁿ - where 𝑛 is the number of zeros found.
//
// NOTE: This will overflow if you let it read too far =)
//
//		ZLE Key | Bit Range
//		      1 | 0
//		    0 1 | 2
//		  0 0 1 | 4
//		0 0 0 1 | 8
//	           ...
//	      𝑛   1 | 2ⁿ
func (_ _fuzzy) ParseZLE(key tiny.Measurement) int {
	count := 0
	for _, b := range key.Bits {
		if b == tiny.Zero {
			count++
		} else {
			break
		}
	}
	if count == 0 {
		return 0
	}
	return int(math.Pow(2, float64(count)))
}

// SixtyFour uses the key Measurement value to calculate a bit range of up to six bits, yielding 64 unique values.
//
// NOTE: This will still return a bit length of 6 if provided a key value greater than 3.
//
//	Key | Bit Range
//	  0 | 0
//	  1 | 2
//	  2 | 4
//	  3 | 6
func (_ _fuzzy) SixtyFour(key tiny.Measurement) int {
	switch v := key.Value(); v {
	case 0:
		return 0
	case 1:
		return 2
	case 2:
		return 4
	case 3:
		return 6
	default:
		return 6
	}
}

// Window creates a fuzzy projection function using a window width.
// The window width is multiplied against the value of the key measurement (plus one) to build
// the continuation projection range.
//
// NOTE: This will panic if provided a window width <= 0.
//
// For example:
//
//		 FuzzyRead(2, tiny.Fuzzy.Window(3))
//
//	Value-> 1  |   1                                             ← Window Occurances
//	     | 0 0 | 0 0 1 | 1 0 1 0 0 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | ← Raw bits
//	     | Key | Cont  |            Remainder                  | ← Fuzzy read
//
//	Value-> 2  |   1       2                                       ← Window Occurances
//	     | 0 1 | 0 0 1 - 1 0 1 | 0 0 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | ← Raw bits
//	     | Key |  Continuation |         Remainder               | ← Fuzzy read
//
//	Value-> 3  |   1       2       3                                 ← Window Occurances
//	     | 1 0 | 0 0 1 - 1 0 1 - 0 0 0 | 1 0 1 1 0 0 0 1 0 0 0 0 1 | ← Raw bits
//	     | Key |     Continuation      |        Remainder          | ← Fuzzy read
//
//	Value-> 4  |   1       2       3       4                           ← Window Occurances
//	     | 1 1 | 0 0 1 - 1 0 1 - 0 0 0 - 1 0 1 | 1 0 0 0 1 0 0 0 0 1 | ← Raw bits
//	     | Key |          Continuation         |      Remainder      | ← Fuzzy read
func (_ _fuzzy) Window(windowWidth int) func(tiny.Measurement) int {
	if windowWidth <= 0 {
		panic("fuzzy.Window: window width must be greater than zero")
	}

	return func(key tiny.Measurement) int {
		return (key.Value() + 1) * windowWidth
	}
}

// PowerWindow creates a fuzzy projection function using a window width.
// The window width is multiplied against the value of the key measurement to build
// the continuation projection range.
// This differs from a standard Window operation in that the key value is considered to be a
// 'power of 2', providing further projection into the data (with 0 being a single occurrance.)
//
// NOTE: This will panic if provided a window width <= 0.
//
// For example:
//
//		 FuzzyRead(2, tiny.Fuzzy.PowerWindow(2))
//
//	Value-> 1  |  1                                              ← Window Occurances
//	     | 0 0 | 0 0 | 1 1 0 1 0 0 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | ← Raw bits
//	     | Key |  C  |             Remainder                   | ← Fuzzy read
//
//	Value-> 2  |  1     2                                          ← Window Occurances
//	     | 0 1 | 0 0 - 1 1 | 0 1 0 0 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | ← Raw bits
//	     | Key |   Cont    |           Remainder                 | ← Fuzzy read
//
//	Value-> 4  |  1     2     3     4                                  ← Window Occurances
//	     | 1 0 | 0 0 - 1 1 - 0 1 - 0 0 | 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | ← Raw bits
//	     | Key |      Continuation     |          Remainder          | ← Fuzzy read
//
//	Value-> 8  |  1     2     3     4     5     6     7     8                  ← Window Occurances
//	     | 1 1 | 0 0 - 1 1 - 0 1 - 0 0 - 0 1 - 0 1 - 1 0 - 0 0 | 1 0 0 0 0 1 | ← Raw bits
//	     | Key |                  Continuation                 |  Remainder  | ← Fuzzy read
//
// NOTE: The output continuation phrase will be aligned to its source form, but a call to Phrase.Align(windowWidth)
// will yield even measurements as demonstrated above.
func (_ _fuzzy) PowerWindow(windowWidth int) func(tiny.Measurement) int {
	if windowWidth <= 0 {
		panic("fuzzy.Window: window width must be greater than zero")
	}

	return func(key tiny.Measurement) int {
		v := key.Value()
		if v == 0 {
			return windowWidth
		}
		return int(math.Pow(2, float64(v))) * windowWidth
	}
}

// EncodeZLE64Value stores the provided integer value using the below ZLE scheme.
//
//	ZLE Key | Bit Range To Store
//	      1 | 4
//	    0 1 | 8
//	  0 0 1 | 16
//	0 0 0 0 | 32
//	0 0 0 1 | 64
func (_ _fuzzy) EncodeZLE64Value(x int) (key tiny.Phrase, projection tiny.Phrase) {
	input := tiny.NewPhraseFromBits(tiny.From.Number(x)...)
	bitLen := input.BitLength()

	if bitLen <= 4 {
		value := make([]tiny.Bit, 4-bitLen)
		value = append(value, input.Bits()...)
		return tiny.NewPhraseFromBits(1), tiny.NewPhraseFromBits(value...)
	}
	if bitLen <= 8 {
		value := make([]tiny.Bit, 8-bitLen)
		value = append(value, input.Bits()...)
		return tiny.NewPhraseFromBits(0, 1), tiny.NewPhraseFromBits(value...)
	}
	if bitLen <= 16 {
		value := make([]tiny.Bit, 16-bitLen)
		value = append(value, input.Bits()...)
		return tiny.NewPhraseFromBits(0, 0, 1), tiny.NewPhraseFromBits(value...)
	}
	if bitLen <= 32 {
		value := make([]tiny.Bit, 32-bitLen)
		value = append(value, input.Bits()...)
		return tiny.NewPhraseFromBits(0, 0, 0, 0), tiny.NewPhraseFromBits(value...)
	}
	if bitLen <= 64 {
		value := make([]tiny.Bit, 64-bitLen)
		value = append(value, input.Bits()...)
		return tiny.NewPhraseFromBits(0, 0, 0, 1), tiny.NewPhraseFromBits(value...)
	}
	panic(fmt.Sprintf("invalid 64-bit ZLE key: %v", input.Bits))
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

	resolutionMax := tiny.To.Number(bitWidth, tiny.Synthesize.Ones(bitWidth).Bits()...)
	resolutionMax2x := tiny.To.Number(bitWidth2x, tiny.Synthesize.Ones(bitWidth2x).Bits()...)
	resolutionMax4x := tiny.To.Number(bitWidth4x, tiny.Synthesize.Ones(bitWidth4x).Bits()...)
	bitLength := target.BitLen()

	eighth := bitLength / 8
	quarter := bitLength / 4
	phrase := tiny.NewPhraseFromBigInt(target)

	region0, phrase := phrase.Read(eighth)
	region1, phrase := phrase.Read(eighth)
	region2, phrase := phrase.Read(quarter)
	region3 := phrase

	fuzzy0, index0 := tiny.Synthesize.Approximation(region0.AsBigInt(), resolutionMax4x, eighth)
	fuzzy1, index1 := tiny.Synthesize.Approximation(region1.AsBigInt(), resolutionMax4x, eighth)
	fuzzy2, index2 := tiny.Synthesize.Approximation(region2.AsBigInt(), resolutionMax2x, quarter)
	fuzzy3, index3 := tiny.Synthesize.Approximation(region3.AsBigInt(), resolutionMax, bitLength-eighth-eighth-quarter)

	indexBits0 := tiny.From.Number(index0, bitWidth4x)
	indexBits1 := tiny.From.Number(index1, bitWidth4x)
	indexBits2 := tiny.From.Number(index2, bitWidth2x)
	indexBits3 := tiny.From.Number(index3, bitWidth)

	approx.Value = tiny.NewPhraseFromBits(fuzzy0...).AppendBits(fuzzy1...).AppendBits(fuzzy2...).AppendBits(fuzzy3...).AsBigInt()
	approx.Indices = tiny.NewPassage(tiny.NewPhraseFromBits(indexBits0...), tiny.NewPhraseFromBits(indexBits1...), tiny.NewPhraseFromBits(indexBits2...), tiny.NewPhraseFromBits(indexBits3...))

	approx.Relativity = tiny.NewRelativeSize(approx.Value.Cmp(target))
	if approx.Relativity == tiny.Equal {
		approx.Delta = new(big.Int)
	} else if approx.Relativity == tiny.Smaller {
		approx.Delta = new(big.Int).Sub(target, approx.Value)
	} else {
		approx.Delta = new(big.Int).Sub(approx.Value, target)
	}

	return approx
}

func (_ _fuzzy) Correct(approximation Approximation) (Approximation, CorrectionFactor) {
	var out CorrectionFactor
	valueF := new(big.Float).SetInt(approximation.Value)

	bitWidths := []int{3, 4, 6, 8}
	factors := make([][]float64, 4)

	if approximation.Relativity == tiny.Larger {
		out.Threshold = tiny.One
		factors = tiny.factors1
	} else {
		out.Threshold = tiny.Zero
		factors = tiny.factors0
	}

	for ii := 0; ii < 4; ii++ {
		for i := 0; i < len(factors[ii]); i++ {
			factor := new(big.Float).SetFloat64(factors[ii][i])
			a, _ := new(big.Float).Mul(valueF, factor).Int(nil)
			r := tiny.NewRelativeSize(a.Cmp(approximation.Target))
			var d *big.Int

			if r == tiny.Larger {
				d = new(big.Int).Sub(a, approximation.Target)
			} else {
				d = new(big.Int).Sub(approximation.Target, a)
			}

			if d.Cmp(approximation.Delta) < 0 {
				approximation.Value = a
				approximation.Delta = d
				approximation.Relativity = r
				out.Focus = tiny.Crumb(ii)
				out.Value = tiny.From.Number(i, bitWidths[ii])
			}
		}
	}

	return approximation, out
}
