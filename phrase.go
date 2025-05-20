package tiny

import "fmt"

// Phrase represents a Measurement slice
type Phrase []Measurement

// NewPhrase calls NewMeasurement for each input byte and returns a Phrase of the results.
func NewPhrase(data ...byte) Phrase {
	out := make(Phrase, len(data))
	for i, d := range data {
		out[i] = NewMeasurement([]byte{d})
	}
	return out
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
//	 Align(8)
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
	if w > 32 {
		panic(errorMeasureLimit)
	}
	if w <= 0 {
		panic(fmt.Sprintf("cannot read at a %d bit width", width))
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
// NOTE: If you provide a length in excess of the phrase bit-length only available bits will be read.
//
// NOTE: This is intended for reading long stretches of bits.
// If you wish to read less than 32 bits from the first measurement, using Phrase.ReadMeasurement is a
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

// ReadMeasurement reads the provided number of bits from the source phrase as a Measurement and the
// remainder as a phrase.
//
// This will panic if you attempt to read more than 32 bits as it cannot contain the result in a single
// measurement.
// If you wish to read more than 32 bits, please use Phrase.Read.
func (phrase Phrase) ReadMeasurement(length int) (read Measurement, remainder Phrase) {
	if length > 32 {
		panic(errorMeasureLimit)
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
//	 (Optional) Start.Align(), Middle.Align(), End.Align()
//
//		| 0 1 0 0 | 1 1 0 1 0 0 0 1 - 0 1 1 0 0 0 1 0 | 0 0 0 1 | <- Raw Bits
//		|  Start  |     Middle1     |     Middle2     |   End   | <- Aligned Phrase Measurements
func (phrase Phrase) Trifurcate(startLen int, middleLen int) (start Phrase, middle Phrase, end Phrase) {
	start, end = phrase.Read(startLen)
	middle, end = end.Read(middleLen)
	return start, middle, end
}
