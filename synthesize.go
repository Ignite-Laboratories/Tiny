package tiny

import "crypto/rand"

type _synthesize struct{}

// ForEach calls the provided function the desired number of times and then builds
// a Measure from the collected results of all invocations.
// For Example, to synthesize a string of 5 ones:
// Synthesize.ForEach(5, func(i int) Bit { return 1 })
func (_ _synthesize) ForEach(count int, f func(int) Bit) Measure {
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
	return NewMeasure(bytes, bits...)
}

// Ones creates a slice of '1's of the requested length.
func (s _synthesize) Ones(count int) Measure {
	return s.ForEach(count, func(i int) Bit { return One })
}

// Zeros creates a slice of '0's of the requested length.
func (s _synthesize) Zeros(count int) Measure {
	return s.ForEach(count, func(i int) Bit { return Zero })
}

// Repeating repeats the provided pattern the desired number of times.
// Use Repeating when you want the entire pattern emitted a fixed number of times.
// Use Pattern when you want the pattern to fit within a specified length.
func (s _synthesize) Repeating(count int, pattern ...Bit) Measure {
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
func (s _synthesize) Pattern(length int, pattern ...Bit) Measure {
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

// Random creates a random sequence of 1s and 0s.
func (s _synthesize) Random(length int) Measure {
	return s.ForEach(length, func(_ int) Bit {
		var b [1]byte
		_, _ = rand.Read(b[:])
		return Bit(b[0] % 2)
	})
}
