package tiny

import (
	"crypto/rand"
	"math/big"
)

type _synthesize struct{}

// ForEach calls the provided function the desired number of times and then builds
// a Phrase from the collected results of all invocations.
//
// NOTE: This efficiently converts every 8 bits into a full byte automatically.
//
// For example, to synthesize a phrase of 5 ones:
//
//	Synthesize.ForEach(5, func(i int) Bit { return One })
//
// Or, to synthesize a phrase of zeros except every nth bit:
//
//	Synthesize.ForEach(1024, func(i int) Bit {
//	 if i % 𝑛 == 0 {
//	  return One
//	 }
//	 return Zero
//	})
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

// Passage encodes a traditional binary Phrase using the below scheme.
//
//   - For the full bit width of the target, synthesize the mid-point value
//   - Calculate the "Delta" from the midpoint to the target
//   - Save off the Delta's sign to the passage's signature
//   - Repeat the process for one bit width smaller
//   - Continue this operation until you reach a value that is the same bit width as "deltaWidth"
//   - Store the Delta on the passage
//
// You can reconstruct the original information by "performing" the passage.
//
// This may or may not get an overall reduction in bits - however, on average, you will gain 2 bits =)
func (s _synthesize) Passage(target Phrase, deltaWidth int) Passage {
	p := Passage{
		Signature:    NewPhrase(),
		Delta:        NewPhrase(),
		DeltaWidth:   deltaWidth,
		InitialWidth: target.BitLength(),
	}

	delta := target.AsBigInt()

	for i := p.InitialWidth; i >= deltaWidth; i-- {
		midpoint := Synthesize.Midpoint(i)

		delta = new(big.Int).Sub(delta, midpoint.AsBigInt())
		if delta.Sign() < 0 {
			p.Signature = p.Signature.AppendBits(1)
		} else {
			p.Signature = p.Signature.AppendBits(0)
		}
		p.Signature = p.Signature.Align()
		delta = new(big.Int).Abs(delta)
		p.Delta = NewPhraseFromBigInt(delta)
	}
	return p
}

// Movement gathers all of the movements into a single seed value which can be re-performed.
func (s _synthesize) Movement(target Phrase, bitWidth int) Movement {
	m := Movement{}
	for {
		p := s.Passage(target, bitWidth)
		target = p.AsPhrase()
		m.Cycles++

		if target.BitLength() < 16 {
			m.Signature = p.Signature
			m.Delta = p.Delta
			break
		}
	}
	return m
}

// Ones creates a slice of '1's of the requested length.
func (s _synthesize) Ones(count int) Phrase {
	return s.ForEach(count, func(i int) Bit { return One })
}

// Zeros creates a slice of '0's of the requested length.
func (s _synthesize) Zeros(count int) Phrase {
	return s.ForEach(count, func(i int) Bit { return Zero })
}

// TrailingZeros creates a slice of '1's of the requested length, except for the trailing bits.
//
// This has a mathematical purpose!  A dark index of data has a max addressable value of (2ⁿ)-1 - but
// introducing zeros to the right of a fully dark index is equivalent to 2ⁿ-2ᶻ, with 'z' being
// the number of zeros introduced.
//
// @formatter:off
//
// For example:
//
//	1 0 0 0 0 0 0 ← 64 (2⁶)
//	  1 1 1 1 1 1 ← 63 (2⁶-1)  [2⁰]
//	  1 1 1 1 1 0 ← 62 (2⁶-2)  [2¹]
//	  1 1 1 1 0 0 ← 60 (2⁶-4)  [2²]
//	  1 1 1 0 0 0 ← 56 (2⁶-8)  [2³]
//	  1 1 0 0 0 0 ← 48 (2⁶-16) [2⁴]
//	  1 0 0 0 0 0 ← 32 (2⁶-32) [2⁵]
//	  0 0 0 0 0 0 ← 0  (2⁶-64) [2⁶]
//
// This allows us to decay a dark value exponentially using a light value.
//
// @formatter:on
func (s _synthesize) TrailingZeros(count int, zeros int) Phrase {
	remainder := count
	return s.ForEach(count, func(i int) Bit {
		remainder--
		if remainder < zeros {
			return Zero
		}
		return One
	})
}

// Midpoint creates a slice with a '1' in the first position and zeros in all subsequent positions.
func (s _synthesize) Midpoint(width int) Phrase {
	return s.ForEach(width, func(i int) Bit {
		if i == 0 {
			return One
		}
		return Zero
	})
}

// Point synthesizes the target point in the provided index.  This is restricted to a "practical limit"
// of 64 bits, as the input type is unable to address anything larger - if you wish to directly synthesize
// a wider point, please consider an alternative means.
func (s _synthesize) Point(i int, index int) Phrase {
	return NewPhraseFromBits(From.Number(i, index)...)
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
// NOTE: This ensures it will never return all 1s, 0s, or repeating [ 1 0 ] [ 0 1 ].
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
func (s _synthesize) RandomPhrase(length int, measurementWidth ...int) (phrase Phrase) {
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

// Boundary synthesizes a binary boundary position.  These are positions where the most significant
// bits are defined, followed by a repeating bit until the end of the bit-width.
//
// The provided width is the overall bit-width of the resulting phrase.
//
// @formatter:off
//
// For example - a note can subdivide 8 boundary positions of a byte index:
//
//	 |← Overall Width ->|
//		|   ⬐ The MSBs      |
//		| 1 1 1 - 1 1 1 1 1 | ← Dark boundary
//		| 1 1 1 - 0 0 0 0 0 |  ⅞
//		| 1 1 0 - 0 0 0 0 0 |  ¾
//		| 1 0 1 - 0 0 0 0 0 |  ⅝
//		| 1 0 0 - 0 0 0 0 0 |  Mid-point
//		| 0 1 1 - 0 0 0 0 0 |  ⅜
//		| 0 1 0 - 0 0 0 0 0 |  ¼
//		| 0 0 1 - 0 0 0 0 0 |  ⅛
//		| 0 0 0 - 0 0 0 0 0 |  Light boundary
//		              ⬑ The repetend
//
// @formatter:on
func (s _synthesize) Boundary(msbs []Bit, repetend Bit, width int) Phrase {
	if width == 0 {
		return Phrase{}
	}

	x := 0
	return s.ForEach(width, func(i int) Bit {
		if x >= len(msbs) {
			return repetend
		}
		out := msbs[x]
		x++
		return out
	})
}

// AllBoundaries generates all of the boundaries for the provided depth at the specified bit width.
//
// The depth value defines the bit width of subdivision - for instance, a value of 3 will create
// a 3-bit wide range of boundary points.
//
// See Boundary for more information on that process.
//
// NOTE: This will panic if provided a negative depth or width.
func (s _synthesize) AllBoundaries(depth int, width int) (boundaries []Phrase) {
	if depth < 0 {
		panic("cannot synthesize boundaries with a negative depth")
	}
	if width <= 0 {
		if width == 0 {
			return []Phrase{}
		}
		panic("cannot synthesize boundaries with a negative width")
	}

	i := 0
	for {
		// Create the MSBs
		bits := From.Number(i, depth)

		// Check if they are all 1
		reachedEnd := true
		for ii := 0; ii < len(bits); ii++ {
			if bits[ii] == Zero {
				// This is the final iteration
				reachedEnd = false
			}
		}

		// Synthesize the boundary point
		boundaries = append(boundaries, s.Boundary(bits, Zero, width))

		i++
		if reachedEnd {
			// Synthesize the final 'dark' boundary point
			boundaries = append(boundaries, s.Boundary(bits, One, width))
			break
		}
	}
	return boundaries
}
