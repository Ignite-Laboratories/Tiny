package tiny

import (
	"crypto/rand"
	"math/big"
)

type _synthesize struct{}

// ForEach calls the provided function the desired number of times and then builds
// a Phrase from the collected results of all invocations.
// For example, to synthesize a string of 5 ones:
// Synthesize.WalkBits(5, func(i int) Bit { return 1 })
func (_ _synthesize) ForEach(count int, f func(int) Bit) Phrase {
	var bytes []byte
	var bits []Bit
	subI := 0
	for i := 0; i < count; i++ {
		bits = append(bits, f(i))
		subI++
		if subI == 8 {
			bytes = append(bytes, To.Byte(bits...))

			bits = make([]Bit, 0)
			subI = 0
		}
	}
	return NewPhraseFromBytesAndBits(bytes, bits...)
}

// Ones creates a slice of '1's of the requested length.
func (s _synthesize) Ones(count int) Phrase {
	return s.ForEach(count, func(i int) Bit { return One })
}

// Zeros creates a slice of '0's of the requested length.
func (s _synthesize) Zeros(count int) Phrase {
	return s.ForEach(count, func(i int) Bit { return Zero })
}

// Repeating repeats the provided pattern the desired number of times.
// Use Repeating when you want the entire pattern emitted a fixed number of times.
// Use Pattern when you want the pattern to fit within a specified length.
func (s _synthesize) Repeating(count int, pattern ...Bit) Phrase {
	patternI := 0
	return s.ForEach(count*len(pattern), func(_ int) Bit {
		bit := pattern[patternI]
		patternI++
		if patternI == len(pattern) {
			patternI = 0
		}
		return bit
	})
}

// Pattern repeats the provided pattern up to the desired length.
// Use Pattern when you want the pattern to fit within a specified length.
// Use Repeating when you want the entire pattern emitted a fixed number of times.
func (s _synthesize) Pattern(length int, pattern ...Bit) Phrase {
	patternI := 0
	return s.ForEach(length, func(_ int) Bit {
		bit := pattern[patternI]
		patternI++
		if patternI == len(pattern) {
			patternI = 0
		}
		return bit
	})
}

// Random creates a basic random sequence of 1s and 0s as a phrase.
//
// If you would prefer to implement your own bit-for-bit randomness, you may optionally provide
// a function that dynamically generates each Bit.
//
// NOTE: This will panic if given a length greater than your architecture's bit width.
func (s _synthesize) Random(length int, generator ...func(int) Bit) Phrase {
	g := func(_ int) Bit {
		var b [1]byte
		_, _ = rand.Read(b[:])
		return Bit(b[0] % 2)
	}
	if len(generator) > 0 && generator[0] != nil {
		g = generator[0]
	}

	if length == 0 {
		return NewPhrase()
	}
	if length > GetArchitectureBitWidth() {
		panic(errorMeasurementLimit)
	}
	for {
		result := s.ForEach(length, g)
		bits := result.Bits()
		if len(bits) > 2 {
			ones := Analyze.Repetition(bits, 1)
			zeros := Analyze.Repetition(bits, 0)
			oneZeros := Analyze.Repetition(bits, 1, 0)
			zeroOnes := Analyze.Repetition(bits, 0, 1)

			if !ones && !zeros && !oneZeros && !zeroOnes {
				return result
			}
		} else {
			// Two digit (or less) requests are always "random"
			return result
		}
	}
}

// RandomPhrase creates a random phrase containing the provided number of measurements, each initialized
// with 8 random bits.
//
// If you would prefer different sized measurements, you may optionally provide the measurement width.
//
// NOTE: This will panic if you provide a measurement width of 0 or less, or greater
// than your architecture's bit width.
func (s _synthesize) RandomPhrase(length int, measurementWidth ...int) Phrase {
	var phrase Phrase
	if length == 0 {
		return phrase
	}
	width := 8
	if len(measurementWidth) > 0 {
		width = measurementWidth[0]
		if width > GetArchitectureBitWidth() {
			panic(errorMeasurementLimit)
		}
		if width <= 0 {
			panic("cannot synthesize measurements of 0 or negative bit lengths")
		}
	}
	for i := 0; i < length; i++ {
		phrase = append(phrase, s.Random(width)...)
	}
	return phrase
}

// Subdivided returns back a synthetic set of binary digits of the provided bit width.
//
// The numeric range is then subdivided at the provided resolution and a value is synthesized
// at the provided index.
//
// NOTE: If you request an index outsize of the subdivision range, this "clamps" it into that range.
func (s _synthesize) Subdivided(width int, index int, resolution int) []Bit {
	if index < 0 {
		index = 0
	}
	if index > resolution {
		index = resolution
	}

	upper, _ := new(big.Int).SetString(Synthesize.Ones(width).StringBinary(), 2)
	step := new(big.Float).Quo(new(big.Float).SetInt(upper), big.NewFloat(float64(resolution)))
	value := new(big.Float).Mul(step, big.NewFloat(float64(index)))
	truncated, _ := value.Int(nil)
	return From.BigInt(truncated, width)
}

// Approximate subdivides the target's bit-width range and then finds the closest index to the provided target.
//
// If no approximation with is provided, the bit length of the target is used.
//
// The result is a slice of bits and the index of the closest match.
func (s _synthesize) Approximate(target *big.Int, resolution int, width ...int) ([]Bit, int) {
	w := target.BitLen()
	if len(width) > 0 {
		w = width[0]
	}
	targetFloat := new(big.Float).SetInt(target)
	upper, _ := new(big.Int).SetString(Synthesize.Ones(w).StringBinary(), 2)
	step := new(big.Float).Quo(new(big.Float).SetInt(upper), big.NewFloat(float64(resolution)))
	index, _ := new(big.Float).Quo(targetFloat, step).Int(nil)
	indexInt := int(index.Int64())
	return Synthesize.Subdivided(w, indexInt, resolution), indexInt
}

// FuzzyApproximation calls Approximate on the target's bits at four different scales and returns
// the approximation phrase, delta, and whether the approximation is larger or smaller than the target.
//
// The approximation itself is a Scale (12 bits) of four indices, each representing the subdivision index
// of that particular index's region.
//
// NOTE: Standard resolution is a Note (3 bits) but you may provide an optional resolution, if desired.
// The resolution is how many bits each index occupies.
//
//	Index 0 represents the first ⅛th of the target bits
//	Index 1 represents the second ⅛th
//	Index 2 represents the second ¼
//	Index 3 represents the final ½
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
func (s _synthesize) FuzzyApproximation(target *big.Int, resolution ...int) (indices Phrase, approximation *big.Int, delta *big.Int, comparison RelativeSize) {
	r := 3
	if len(resolution) > 0 {
		r = resolution[0]
	}

	eighth := target.BitLen() / 8
	quarter := target.BitLen() / 4
	phrase := NewPhraseFromBigInt(target)

	region0, phrase := phrase.Read(eighth)
	region1, phrase := phrase.Read(eighth)
	region2, phrase := phrase.Read(quarter)
	region3 := phrase

	fuzzy0, index0 := Synthesize.Approximate(region0.AsBigInt(), r)
	fuzzy1, index1 := Synthesize.Approximate(region1.AsBigInt(), r)
	fuzzy2, index2 := Synthesize.Approximate(region2.AsBigInt(), r)
	fuzzy3, index3 := Synthesize.Approximate(region3.AsBigInt(), r)

	indexBits0 := From.Number(index0, r)
	indexBits1 := From.Number(index1, r)
	indexBits2 := From.Number(index2, r)
	indexBits3 := From.Number(index3, r)

	approximation = NewPhraseFromBits(fuzzy0...).AppendBits(fuzzy1...).AppendBits(fuzzy2...).AppendBits(fuzzy3...).AsBigInt()
	indices = NewPhraseFromBits(indexBits0...).AppendBits(indexBits1...).AppendBits(indexBits2...).AppendBits(indexBits3...)
	comparison = NewRelativeSize(approximation.Cmp(target))

	if comparison == 0 {
		delta = new(big.Int)
	} else if comparison < 0 {
		delta = new(big.Int).Sub(target, approximation)
	} else {
		delta = new(big.Int).Sub(approximation, target)
	}

	return indices, approximation, delta, comparison
}

// Approximation synthesizes the input approximation phrase at the provided bit width resolution.
func (s _synthesize) Approximation(indices Phrase, resolution int) Phrase {
	return nil
}
