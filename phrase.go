package tiny

import "fmt"

// Phrase represents a Measurement slice and provides easy clustered measurement functionality.
type Phrase []Measurement

// NewPhrase calls NewMeasurement for each input byte and returns a Phrase of the results.
func NewPhrase(data ...byte) Phrase {
	out := make(Phrase, len(data))
	for i, d := range data {
		out[i] = NewMeasurement([]byte{d})
	}
	return out
}

// NewPhraseFromBits creates a phrase of the provided bits at a standard 8-bits-per-byte measurement interval.
func NewPhraseFromBits(data ...Bit) Phrase {
	out := make(Phrase, 0)
	current := make([]Bit, 0, 8)
	ii := 0

	for _, d := range data {

		if ii > 7 {
			ii = 0
			out = append(out, NewMeasurement([]byte{}, current...))
			current = make([]Bit, 0, 8)
		}
		current = append(current, d)
		ii++
	}
	if len(current) > 0 {
		out = append(out, NewMeasurement([]byte{}, current...))
	}
	return out
}

// NewPhraseFromBitsAndBytes creates a phrase by combining NewPhraseFromBits(bits) and then NewPhrase(bytes).
func NewPhraseFromBitsAndBytes(bits []Bit, bytes ...byte) Phrase {
	p := NewPhraseFromBits(bits...)
	p = append(p, NewPhrase(bytes...)...)
	return p
}

// NewPhraseFromBytesAndBits creates a phrase by combining NewPhrase(bytes) and then NewPhraseFromBits(bits).
func NewPhraseFromBytesAndBits(bytes []byte, bits ...Bit) Phrase {
	p := NewPhrase(bytes...)
	p = append(p, NewPhraseFromBits(bits...)...)
	return p
}

// NewPhraseFromString creates a new phrase from a binary string input.
func NewPhraseFromString(s string) Phrase {
	bits := make([]Bit, len(s))
	for i := 0; i < len(bits); i++ {
		bits[i] = Bit(s[i] & 1)
	}
	return NewPhraseFromBits(bits...)
}

// ToBytesAndBits converts its measurements into bytes and the remainder of bits.
func (phrase Phrase) ToBytesAndBits() ([]byte, []Bit) {
	out := make([]byte, 0, len(phrase))

	current := make([]Bit, 8)
	var i int
	for _, measure := range phrase {
		for _, bit := range measure.GetAllBits() {
			current[i] = bit
			i++

			if i == 8 {
				out = append(out, To.Byte(current...))
				current = make([]Bit, 8)
				i = 0
			}
		}
	}
	current = current[:i]
	return out, current
}

// QuarterSplit quarter splits each Measurement of the Phrase.
func (phrase Phrase) QuarterSplit() {
	for i, m := range phrase {
		m.QuarterSplit()
		phrase[i] = m
	}
}

// UnQuarterSplit reverses a quarter split operation on each Measurement of the Phrase.
func (phrase Phrase) UnQuarterSplit() {
	for i, m := range phrase {
		m.UnQuarterSplit()
		phrase[i] = m
	}
}

// BitLength returns the total length of all bits in each Measurement of the Phrase.
func (phrase Phrase) BitLength() int {
	total := 0
	for _, m := range phrase {
		total += m.BitLength()
	}
	return total
}

// CountBelowThreshold counts any Measurement of the Phrase that's below the provided threshold value.
func (phrase Phrase) CountBelowThreshold(threshold int) int {
	var count int
	for _, m := range phrase {
		if m.Value() < threshold {
			count++
		}
	}
	return count
}

// AllBelowThreshold checks if every Measurement of the Phrase is below the provided threshold value.
func (phrase Phrase) AllBelowThreshold(threshold int) bool {
	for _, m := range phrase {
		if m.Value() > threshold {
			return false
		}
	}
	return true
}

// BreakMeasurementsApart breaks each Measurement of the Phrase apart at the provided index and returns
// the two resulting phrases.  The left phrase will contain the most significant bits, while the right
// phrase will contain the least significant bits.
func (phrase Phrase) BreakMeasurementsApart(index int) (left Phrase, right Phrase) {
	left = make(Phrase, len(phrase))
	right = make(Phrase, len(phrase))

	for i, m := range phrase {
		l, r := m.BreakApart(index)
		left[i] = l
		right[i] = r
	}

	return left, right
}

// RecombineMeasurements recombines the two provided measurement phrases into a single phrase.
// The left phrase should contain the most significant bits, while the right phrase should contain
// the least significant bits.
//
// NOTE: The two phrases must be the same length.  If they are not, this will panic.
func RecombineMeasurements(left Phrase, right Phrase) Phrase {
	if len(left) != len(right) {
		panic("left and right must be the same length")
	}

	out := make(Phrase, len(left))
	for i := 0; i < len(left); i++ {
		// NOTE: We create a new measurement since Append is a pointer operation
		m := NewMeasurement(left[i].Bytes, left[i].Bits...)
		m.Append(right[i])
		out[i] = m
	}

	return out
}

// AsInts converts each Measurement of the Phrase into an int.
func (phrase Phrase) AsInts() []int {
	out := make([]int, len(phrase))
	for i, m := range phrase {
		out[i] = m.Value()
	}
	return out
}

// AsBytes converts each Measurement of the Phrase into a byte.
func (phrase Phrase) AsBytes() []byte {
	out := make([]byte, len(phrase))
	for i, m := range phrase {
		out[i] = byte(m.Value())
	}
	return out
}

// Bits returns a slice of the phrase's underlying bits.
func (phrase Phrase) Bits() []Bit {
	out := make([]Bit, 0, phrase.BitLength())
	for _, m := range phrase {
		out = append(out, m.GetAllBits()...)
	}
	return out
}

// Align ensures all but the final Measurement of the source phrase are of the provided width.
//
// If no width is provided, a standard alignment of 8-bits-per-byte will be used.
//
// A Phrase is considered "aligned" if all except the -final- measurement are of the same width.
//
// For example:
//
//		0 1 | 0 1 0 | 0 1 1 0 1 0 0 0 | 1 0 1 1 0 | 0 0 1 0 0 0 0 1 |  <- Raw Bits
//		 M1 |  M2   |  Measurement 3  |     M4    |  Measurement 5  |  <- "Unaligned" Phrase
//
//	 Align()
//
//		0 1 0 1 0 0 1 1 | 0 1 0 0 0 1 0 1 | 1 0 0 0 1 0 0 0 | 0 1 |  <- Raw Bits
//		 Measurement1   |  Measurement 2  |  Measurement 3  | M4  |  <- "Aligned" Phrase
//
// NOTE: This will panic if you provide a width greater than the maximum width of a Measurement (32 bits),
// or if you provide a width of <= 0.
func (phrase Phrase) Align(width ...int) Phrase {
	w := 8
	if len(width) > 0 {
		w = width[0]
	}
	if w > MaxMeasurementBitLength {
		panic(errorMeasurementLimit)
	}
	if w <= 0 {
		panic(fmt.Sprintf("cannot align at a %d bit width", width))
	}

	src := phrase
	out := make(Phrase, 0, len(src))
	for {
		measure, remainder := src.ReadMeasurement(w)
		if len(remainder) == 0 {
			if measure.BitLength() > 0 {
				out = append(out, measure)
			}
			break
		}

		out = append(out, measure)
		src = remainder
	}

	return out
}

// Read reads the provided number of bits from the source phrase, plus the remainder, as phrases.
//
// NOTE: If you provide a length in excess of the phrase bit-length, only the available bits will be read
// and the remainder will be empty.
//
// NOTE: This is intended for reading long stretches of bits.
// If you wish to read less than 32 bits from the first measurement, using ReadMeasurement is a
// little easier to work with.
func (phrase Phrase) Read(length int) (read Phrase, remainder Phrase) {
	read = make(Phrase, 0, len(phrase))
	remainder = make(Phrase, 0, len(phrase))

	for _, m := range phrase {
		if length <= 0 {
			remainder = append(remainder, m)
			continue
		}

		bitLen := m.BitLength()
		if bitLen <= length {
			read = append(read, m)
		} else {
			bits := m.GetAllBits()
			read = append(read, NewMeasurement([]byte{}, bits[0:length]...))
			remainder = append(remainder, NewMeasurement([]byte{}, bits[length:]...))
		}

		length -= bitLen
	}

	return read, remainder
}

// FuzzyRead reads the provided number of bits from the source phrase and then passes them to the provided fuzzy
// projection function.
// The projection function should return the number of bits to continue reading, based upon the found "key" bits.
//
// Finally, the key Measurement, continuation Phrase, and remainder Phrase are returned.
//
// NOTE: The most common fuzzy projection functions are accessible from the tiny.Fuzzy instance of tiny.FuzzyReader.
//
// For example:
//
//		 FuzzyRead(2, tiny.Fuzzy.SixtyFour)
//
//	 value-> 0    𝧘 Resulting Continuation Size
//	      | 0 0 | | 0 0 1 1 0 1 0 0 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | <- Raw bits
//	      | Key |C|                Remainder                    | <- Fuzzy read
//
//	 value-> 1      𝧘 Resulting Continuation Size
//	      | 0 1 | 0 0 | 1 1 0 1 0 0 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | <- Raw bits
//	      | Key |  C  |               Remainder                 | <- Fuzzy read
//
//	 value-> 2        𝧘 Resulting Continuation Size
//	      | 1 0 | 0 0 1 1 | 0 1 0 0 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | <- Raw bits
//	      | Key |  Cont   |             Remainder               | <- Fuzzy read
//
//	 value-> 3          𝧘 Resulting Continuation Size
//	      | 1 1 | 0 0 1 1 0 1 | 0 0 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | <- Raw bits
//	      | Key |   Continue  |            Remainder            | <- Fuzzy read
func (phrase Phrase) FuzzyRead(length int, fuzzyProjection func(Measurement) int) (key Measurement, continuation Phrase, remainder Phrase) {
	key, remainder = phrase.ReadMeasurement(length)
	continuation, remainder = remainder.Read(fuzzyProjection(key))
	return key, continuation, remainder
}

// ReadMeasurement reads the provided number of bits from the source phrase as a Measurement and provides the
// remainder as a Phrase.
//
// NOTE: This will panic if you attempt to read more than 32 bits as it cannot contain the result in a single
// measurement.
// If you wish to read more than 32 bits, please use Read.
func (phrase Phrase) ReadMeasurement(length int) (read Measurement, remainder Phrase) {
	if length > MaxMeasurementBitLength {
		panic(errorMeasurementLimit)
	}

	read = NewMeasurement([]byte{})
	readMeasures, remainder := phrase.Read(length)
	for _, m := range readMeasures {
		read.Append(m)
	}

	return read, remainder
}

// Trifurcate takes the source phrase and subdivides it in thrice - start, middle, and end.
//
// For example:
//
//		tiny.Phrase{ 77, 22, 33 }
//
//		|        77       |        22       |        33       |  <- Bytes
//		| 0 1 0 0 1 1 0 1 | 0 0 0 1 0 1 1 0 | 0 0 1 0 0 0 0 1 |  <- Raw Bits
//		|  Measurement 1  |  Measurement 2  |  Measurement 3  |  <- Source Phrase
//
//		Trifurcate(4,16)
//
//		|    4    |                  16                 |           <- Trifurcation lengths
//		| 0 1 0 0 | 1 1 0 1 - 0 0 0 1 0 1 1 0 - 0 0 1 0 | 0 0 0 1 | <- Raw Bits
//		|  Start  |               Middle                |   End   | <- Trifurcated Phrases
//		|  Start  | Middle1 |     Middle2     | Middle3 |   End   | <- Phrase Measurements
//
//	 (Optional) Align() each phrase
//
//		| 0 1 0 0 | 1 1 0 1 0 0 0 1 - 0 1 1 0 0 0 1 0 | 0 0 0 1 | <- Raw Bits
//		|  Start  |     Middle1     |     Middle2     |   End   | <- Aligned Phrase Measurements
func (phrase Phrase) Trifurcate(startLen int, middleLen int) (start Phrase, middle Phrase, end Phrase) {
	start, end = phrase.Read(startLen)
	middle, end = end.Read(middleLen)
	return start, middle, end
}

// Focus is used to limit the width of eminently relevant measurements.
// It finds the midpoint of the phrase (using flooring) to split it in twain.
// Because of the floored split point, the right phrase will be larger if the data is odd in length.
//
// You may optionally provide a 'times' parameter that indicates how many times to "focus" into the
// eminent measurements of the phrase recursively.
// This will continue to bisect the left phrase and grow the right by prepending it with the remainder.
func (phrase Phrase) Focus(times ...int) (left Phrase, right Phrase) {
	t := 1
	if len(times) > 0 {
		t = times[0]
		if t < 1 {
			// If provided a negative or zero value, just bisect once
			t = 1
		}
	}

	length := phrase.BitLength()
	midpoint := length / 2
	left, right = phrase.Read(midpoint)

	if t > 1 {
		ll, rr := left.Focus(t - 1)
		left = ll
		right = append(rr, right...)
	}

	return left, right
}

// WalkBits walks the bits of the source phrase at the provided stride and calls the
// provided function for each measurement step.
func (phrase Phrase) WalkBits(stride int, fn func(int, Measurement)) {
	if stride > MaxMeasurementBitLength {
		panic(errorMeasurementLimit)
	}
	if stride <= 0 {
		panic("cannot walk at a stride of 0 or less")
	}

	remainder := phrase
	var bitM Measurement
	i := 0
	for bitM, remainder = remainder.ReadMeasurement(stride); len(remainder) > 0; bitM, remainder = remainder.ReadMeasurement(stride) {
		if bitM.BitLength() > 0 {
			fn(i, bitM)
			i++
		}
	}
	if bitM.BitLength() > 0 {
		fn(i, bitM)
	}
}

// StringBinary returns the phrase's bits as a binary string of 1s and 0s.
func (phrase Phrase) StringBinary() string {
	out := ""
	for _, m := range phrase {
		out += m.StringBinary()
	}
	return out
}
