package tiny

import (
	"fmt"
	"math/big"
)

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

// NewPhraseFromMeasurement creates a phrase from a single measurement.
func NewPhraseFromMeasurement(m Measurement) Phrase {
	return Phrase{m}
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

// NewPhraseFromString creates a new Phrase from a binary string input.
func NewPhraseFromString(s string) Phrase {
	bits := make([]Bit, len(s))
	for i := 0; i < len(bits); i++ {
		bits[i] = Bit(s[i] & 1)
	}
	return NewPhraseFromBits(bits...)
}

// NewPhraseFromBigInt creates a new Phrase from a big.Int.
func NewPhraseFromBigInt(b *big.Int) Phrase {
	return NewPhraseFromString(b.Text(2))
}

// AsBigInt converts the tiny.Phrase into a big.Int.
func (phrase Phrase) AsBigInt() *big.Int {
	out := new(big.Int)
	out.SetString(phrase.StringBinary(), 2)
	return out
}

/**
Append/Prepend
*/

// AppendMeasurement appends the provided measurement to the phrase.
func (phrase Phrase) AppendMeasurement(m Measurement) Phrase {
	return append(phrase, m)
}

// PrependMeasurement prepends the provided measurement to the phrase.
func (phrase Phrase) PrependMeasurement(m Measurement) Phrase {
	return append(Phrase{m}, phrase...)
}

// AppendBits appends the provided bits to the end of the phrase.
func (phrase Phrase) AppendBits(bits ...Bit) Phrase {
	return append(phrase, NewPhraseFromBits(bits...)...)
}

// PrependBits prepends the provided bits to the beginning of the phrase.
func (phrase Phrase) PrependBits(bits ...Bit) Phrase {
	return append(NewPhraseFromBits(bits...), phrase...)
}

// AppendBytes appends the provided bytes to the end of the phrase.
func (phrase Phrase) AppendBytes(bytes ...byte) Phrase {
	return append(phrase, NewPhraseFromBytesAndBits(bytes)...)
}

// PrependBytes prepends the provided bytes to the beginning of the phrase.
func (phrase Phrase) PrependBytes(bytes ...byte) Phrase {
	return append(NewPhraseFromBytesAndBits(bytes), phrase...)
}

// AppendBitsAndBytes appends the provided bits, and then bytes, to the end of the phrase.
func (phrase Phrase) AppendBitsAndBytes(bits []Bit, bytes ...byte) Phrase {
	return append(phrase, NewPhraseFromBitsAndBytes(bits, bytes...)...)
}

// PrependBitsAndBytes prepends the provided bits, and then bytes, to the beginning of the phrase.
func (phrase Phrase) PrependBitsAndBytes(bits []Bit, bytes ...byte) Phrase {
	return append(NewPhraseFromBitsAndBytes(bits, bytes...), phrase...)
}

// AppendBytesAndBits appends the provided bytes, and then bits, to the end of the phrase.
func (phrase Phrase) AppendBytesAndBits(bytes []byte, bits ...Bit) Phrase {
	return append(phrase, NewPhraseFromBytesAndBits(bytes, bits...)...)
}

// PrependBytesAndBits prepends the provided bytes, and then bits, to the beginning of the phrase.
func (phrase Phrase) PrependBytesAndBits(bytes []byte, bits ...Bit) Phrase {
	return append(NewPhraseFromBytesAndBits(bytes, bits...), phrase...)
}

// Append appends the provided phrase to the end of the source phrase.
func (phrase Phrase) Append(p Phrase) Phrase {
	return append(phrase, p...)
}

// Prepend prepends the provided phrase to the beginning of the source phrase.
func (phrase Phrase) Prepend(p Phrase) Phrase {
	return append(p, phrase...)
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

// BitLength returns the total length of all bits in each Measurement of the Phrase.
func (phrase Phrase) BitLength() int {
	total := 0
	for _, m := range phrase {
		total += m.BitLength()
	}
	return total
}

// BitLengthAsBigInt returns the total length of all bits in each Measurement of the Phrase as a big.Int.
func (phrase Phrase) BitLengthAsBigInt() *big.Int {
	return big.NewInt(int64(phrase.BitLength()))
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
// NOTE: This will panic if you provide a width greater than your architecture's bit width, or if
// given a width of <= 0.
func (phrase Phrase) Align(width ...int) Phrase {
	w := 8
	if len(width) > 0 {
		w = width[0]
	}
	if w > GetArchitectureBitWidth() {
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
// If you wish to read less than your architecture's bit width from the first measurement, using ReadMeasurement is a
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

// ReadMeasurement reads the provided number of bits from the source phrase as a Measurement and provides the
// remainder as a Phrase.
//
// NOTE: This will panic if you attempt to read more than your architecture's bit width.
// For that, please use Read.
func (phrase Phrase) ReadMeasurement(length int) (read Measurement, remainder Phrase) {
	if length > GetArchitectureBitWidth() {
		panic(errorMeasurementLimit)
	}

	read = NewMeasurement([]byte{})
	readMeasures, remainder := phrase.Read(length)
	for _, m := range readMeasures {
		read.Append(m)
	}

	return read, remainder
}

// ReadBit reads a single bit from the source phrase and returns the remainder as a Phrase.
func (phrase Phrase) ReadBit() (read Bit, remainder Phrase, err error) {
	measure, remainder := phrase.ReadMeasurement(1)
	if measure.BitLength() == 0 {
		return 0, nil, fmt.Errorf("no more bits to read")
	}
	return measure.GetAllBits()[0], remainder, nil
}

// ReadUntilOne reads the source phrase until it reaches the first 1, then returns the zero count and remainder.
//
// If you'd like it to stop after a certain count, provide a limit.
func (phrase Phrase) ReadUntilOne(limit ...int) (zeros int, remainder Phrase) {
	l := -1
	if len(limit) > 0 {
		l = limit[0]
	}

	remainder = phrase
	for b, r, err := remainder.ReadBit(); err == nil; b, r, err = r.ReadBit() {
		if l >= 0 && zeros >= l {
			break
		}
		if b == 1 {
			return zeros, remainder
		}
		zeros++
		remainder = r
	}
	return zeros, remainder
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

// Bifurcate takes the source phrase and subdivides it in twain - start and end.
//
// The ending bits will contain any odd bits from the splitting operation.
//
// For example:
//
//			tiny.Phrase{ 77, 22, 33 }
//	     tiny.AppendBits(1, 0, 0)
//
//			|        77       |        22       |        33       |   5   |  <- Values
//			| 0 1 0 0 1 1 0 1 | 0 0 0 1 0 1 1 0 | 0 0 1 0 0 0 0 1 | 1 0 0 |  <- Raw Bits
//			|  Measurement 1  |  Measurement 2  |  Measurement 3  |   M4  |  <- Source Phrase Measurements
//
//			Bifurcate()
//
//			|        77       |         22        |        33       |   5   |  <- Values
//			| 0 1 0 0 1 1 0 1 | 0 0 0 1 0 - 1 1 0 | 0 0 1 0 0 0 0 1 | 1 0 0 |  <- Raw Bits
//			|             Start           |               End               |  <- Bifurcated Phrases
//			|     Start 1     |  Start 2  | End 1 |      End 2      | End 3 |  <- Bifurcated Phrase Measurements
//
//		 (Optional) Align() each phrase
//
//			| 0 1 0 0 1 1 0 1 | 0 0 0 1 0 - 1 1 0 | 0 0 1 0 0 - 0 0 1 | 1 0 0 |  <- Raw Bits
//			|     Start 1     |  Start 2  |       End 1       |     End 2     |  <- Aligned Measurements
func (phrase Phrase) Bifurcate() (start Phrase, end Phrase) {
	return phrase.Read(phrase.BitLength() / 2)
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
//
// NOTE: This will panic if given a stride greater than your architecture's bit width.
func (phrase Phrase) WalkBits(stride int, fn func(int, Measurement)) {
	if stride > GetArchitectureBitWidth() {
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

// Invert XORs every bit of every measurement against 1.
//
// NOTE: This does so iteratively, bit-by-bit.
func (phrase Phrase) Invert() Phrase {
	out := make(Phrase, len(phrase))
	for i, m := range phrase {
		m.Invert()
		out[i] = m
	}
	return out
}

// DeltaEncode uses the source phrase value to encode the delta from a known relative boundary point.
//
// Don't worry, it's not that complex - here's an index of data:
//
//	[ 1 1 1 1 ... 1 1 1 ] <- Dark Boundary
//	↕   Fourth Index    ↕
//	[ 1 1 0 0 ... 0 0 0 ] <- Upper Quarter Boundary
//	↕    Third Index    ↕
//	[ 1 0 0 0 ... 0 0 0 ] <- Mid Boundary
//	↕   Second Index    ↕
//	[ 0 1 0 0 ... 0 0 0 ] <- Lower Quarter Boundary
//	↕    First Index    ↕
//	[ 0 0 0 0 ... 0 0 0 ] <- Light Boundary
//
// A signature dictates where the target exists relative to the final bit width, while the delta (Δ) tells us
// how many bits were removed from the source phrase.
// The remainder tells us how far our target is from the signature's vantage point.
//
// The signature is a very simple measurement - the first bit tells you if the data existed in the lower
// or upper half of the address space, with each next bit telling you if its closer to the top or bottom of
// each subsequent halving of the address space.
// The standard depth to walk is three bits, but you can optionally override this at your discretion.
//
//	NOTE: This will panic if you give it a depth of 0 or less
//
// So, for example:
//
//	              ⬐ Remainder
//	[ 1 0 1 ] [ 1 0 1 ]
//	    ⬑ Key Signature
//
// The key signature tells us it's in the upper eighth of the lower quandrant of the upper half.
// Because it's in the UPPER eighth, the remainder value is considered to be subtracted from the
// synthetic vantage point.
//
// So, with an 8 bit Δ, you'd synthesize 11 bits (including the remainder) as such:
//
//	| 1 1 0 0 0 0 0 0 0 0 0 | <- Synthesized value at the appropriate vantage point (the 3rd quarter, here)
//	|                 1 0 1 | <- Remainder to subtract
//	| 1 0 1 1 1 1 1 1 0 1 1 | <- Resulting value
func (phrase Phrase) DeltaEncode(depth ...int) (signature Measurement, delta int, remainder Phrase) {
	return Measurement{}, -1, nil
}

func (phrase Phrase) DeltaDecode(signature Measurement, delta int, remainder Phrase) Phrase {
	return nil
}

// StringBinary returns the phrase's bits as a binary string of 1s and 0s.
func (phrase Phrase) StringBinary() string {
	out := ""
	for _, m := range phrase {
		out += m.StringBinary()
	}
	return out
}

func (phrase Phrase) String() string {
	out := make([]string, len(phrase))
	for i, m := range phrase {
		out[i] = m.StringBinary()
	}
	return fmt.Sprintf("%v", out)
}
