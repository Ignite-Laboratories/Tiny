package tiny

import (
	"fmt"
	"math"
	"math/big"
)

// Phrase represents a Measurement slice and provides clustered measurement functionality.
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

// AppendBits appends the provided bits to the end of the phrase.
//
// NOTE: This appends the bits as new measurements inside of a new phrase =)
func (phrase Phrase) AppendBits(bits ...Bit) Phrase {
	return append(phrase, NewPhraseFromBits(bits...)...)
}

// AppendBytes appends the provided bytes to the end of the phrase.
//
// NOTE: This appends the bits as new measurements inside of a new phrase =)
func (phrase Phrase) AppendBytes(bytes ...byte) Phrase {
	return append(phrase, NewPhraseFromBytesAndBits(bytes)...)
}

// AppendBitsAndBytes appends the provided bits, and then bytes, to the end of the phrase.
//
// NOTE: This appends the bits as new measurements inside of a new phrase =)
func (phrase Phrase) AppendBitsAndBytes(bits []Bit, bytes ...byte) Phrase {
	return append(phrase, NewPhraseFromBitsAndBytes(bits, bytes...)...)
}

// AppendBytesAndBits appends the provided bytes, and then bits, to the end of the phrase.
//
// NOTE: This appends the bits as new measurements inside of a new phrase =)
func (phrase Phrase) AppendBytesAndBits(bytes []byte, bits ...Bit) Phrase {
	return append(phrase, NewPhraseFromBytesAndBits(bytes, bits...)...)
}

// Append appends the provided phrase(s) to the end of the source phrase.
func (phrase Phrase) Append(p ...Phrase) Phrase {
	out := make(Phrase, 0, len(phrase))
	for _, item := range p {
		out = append(phrase, item...)
	}
	return out
}

// PrependMeasurement prepends the provided measurement to the phrase.
func (phrase Phrase) PrependMeasurement(m Measurement) Phrase {
	return append(Phrase{m}, phrase...)
}

// PrependBits prepends the provided bits to the beginning of the phrase.
//
// NOTE: This prepends the bits as new measurements inside of a new phrase =)
func (phrase Phrase) PrependBits(bits ...Bit) Phrase {
	return append(NewPhraseFromBits(bits...), phrase...)
}

// PrependBytes prepends the provided bytes to the beginning of the phrase.
//
// NOTE: This prepends the bits as new measurements inside of a new phrase =)
func (phrase Phrase) PrependBytes(bytes ...byte) Phrase {
	return append(NewPhraseFromBytesAndBits(bytes), phrase...)
}

// PrependBitsAndBytes prepends the provided bits, and then bytes, to the beginning of the phrase.
//
// NOTE: This prepends the bits as new measurements inside of a new phrase =)
func (phrase Phrase) PrependBitsAndBytes(bits []Bit, bytes ...byte) Phrase {
	return append(NewPhraseFromBitsAndBytes(bits, bytes...), phrase...)
}

// PrependBytesAndBits prepends the provided bytes, and then bits, to the beginning of the phrase.
//
// NOTE: This prepends the bits as new measurements inside of a new phrase =)
func (phrase Phrase) PrependBytesAndBits(bytes []byte, bits ...Bit) Phrase {
	return append(NewPhraseFromBytesAndBits(bytes, bits...), phrase...)
}

// Prepend prepends the provided phrase(s) to the beginning of the source phrase.
func (phrase Phrase) Prepend(p ...Phrase) Phrase {
	out := make(Phrase, 0, len(phrase))
	for _, item := range p {
		out = append(item, phrase...)
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
// @formatter:off
//
// For example-
//
//	Starting Phrase:
//
//	| 0 1 | 0 1 0 | 0 1 1 0 1 0 0 0 | 1 0 1 1 0 | 0 0 1 0 0 0 0 1 |  ← Raw Bits
//	|  M0 -  M1   -  Measurement 2  -     M3    -  Measurement 4  |  ← "Unaligned" Phrase
//
//	Align()
//
//	| 0 1 0 1 0 0 1 1 | 0 1 0 0 0 1 0 1 | 1 0 0 0 1 0 0 0 | 0 1 |  ← Raw Bits
//	|  Measurement 0  -  Measurement 1  -  Measurement 2  - M3  |  ← "Aligned" Phrase
//
// NOTE: This will panic if you provide a width greater than your architecture's bit width, or if
// given a width of <= 0.
//
// @formatter:on
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

// Read reads the provided number of bits from the source phrase, followed by the remainder, as phrases.
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

// ReadFromEnd reads the provided number of bits from the end of the source phrase, followed by the remainder, as phrases.
// The returned bits remain in logical order from left-to-right - for example -
//
//	                                  read ⬎         ⬐ remainder
//	[ 1 0 0 1 0 1 1 ].ReadFromEnd(3) -> [ 0 1 1] [ 1 0 0 1]
//
// NOTE: If you provide a length in excess of the phrase bit-length, only the available bits will be read
// and the remainder will be empty.
func (phrase Phrase) ReadFromEnd(length int) (read Phrase, remainder Phrase) {
	remainder, read = phrase.Read(phrase.BitLength() - length)
	return read, remainder
}

// ReadLastBit reads the last bit of the source phrase, followed by the remainder.
func (phrase Phrase) ReadLastBit() (last Bit, remainder Phrase) {
	end, remainder := phrase.ReadFromEnd(1)
	last, _, _ = end.ReadBit()
	return last, remainder
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
//
// NOTE: This returns an error if there are no more bits to read.
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

/**
Padding
*/

// PadLeftToLength pads the phrase to the desired overall length with zeros on the left (most significant) side
// of the bits.
//
// NOTE: If you'd prefer to pad with ones, please override the char parameter with tiny.One
func (phrase Phrase) PadLeftToLength(overall int, char ...Bit) Phrase {
	c := Zero
	if len(char) > 0 {
		c = char[0]
	}

	toPad := overall - phrase.BitLength()
	if toPad <= 0 {
		return phrase
	}

	if c == Zero {
		return phrase.Prepend(Synthesize.Zeros(toPad))
	}
	return phrase.Prepend(Synthesize.Ones(toPad))
}

// PadRightToLength pads the phrase to the desired overall length with zeros on the right (least significant) side
// // of the bits.
//
// NOTE: If you'd prefer to pad with ones, please override the char parameter with tiny.One
func (phrase Phrase) PadRightToLength(overall int, char ...Bit) Phrase {
	c := Zero
	if len(char) > 0 {
		c = char[0]
	}

	toPad := overall - phrase.BitLength()
	if toPad <= 0 {
		return phrase
	}

	if c == Zero {
		return phrase.Append(Synthesize.Zeros(toPad))
	}
	return phrase.Append(Synthesize.Ones(toPad))
}

// Trifurcate takes the source phrase and subdivides it in thrice - start, middle, and end.
//
// @formatter:off
//
// For example:
//
//		tiny.Phrase{ 77, 22, 33 }
//
//		|        77       |        22       |        33       |  ← Bytes
//		| 0 1 0 0 1 1 0 1 | 0 0 0 1 0 1 1 0 | 0 0 1 0 0 0 0 1 |  ← Raw Bits
//		|  Measurement 0  |  Measurement 1  |  Measurement 2  |  ← Source Phrase
//
//		Trifurcate(4,16)
//
//		|    4    |                  16                 |           ← Trifurcation lengths
//		| 0 1 0 0 | 1 1 0 1 - 0 0 0 1 0 1 1 0 - 0 0 1 0 | 0 0 0 1 | ← Raw Bits
//		|  Start  |               Middle                |   End   | ← Trifurcated Phrases
//		|  Start  | Middle0 -     Middle1     - Middle2 |   End   | ← Phrase Measurements
//
//	 (Optional) Align() each phrase
//
//		| 0 1 0 0 | 1 1 0 1 0 0 0 1 - 0 1 1 0 0 0 1 0 | 0 0 0 1 | ← Raw Bits
//		|  Start  |     Middle0     -     Middle1     |   End   | ← Aligned Phrase Measurements
//
// @formatter:on
func (phrase Phrase) Trifurcate(startLen int, middleLen int) (start Phrase, middle Phrase, end Phrase) {
	start, end = phrase.Read(startLen)
	middle, end = end.Read(middleLen)
	return start, middle, end
}

// Bifurcate takes the source phrase and subdivides it in twain - start and end.
//
// The ending bits will contain any odd bits from the splitting operation.
//
// @formatter:off
//
// For example:
//
//	tiny.Phrase{ 77, 22, 33 }
//	tiny.AppendBits(1, 0, 0)
//
//	|        77       |        22       |        33       |   5   |  ← Values
//	| 0 1 0 0 1 1 0 1 | 0 0 0 1 0 1 1 0 | 0 0 1 0 0 0 0 1 | 1 0 0 |  ← Raw Bits
//	|  Measurement 0  -  Measurement 1  -  Measurement 2  -   M3  |  ← Source Phrase Measurements
//
//	Bifurcate()
//
//	|        77       |         22        |        33       |   5   |  ← Values
//	| 0 1 0 0 1 1 0 1 | 0 0 0 1 0 - 1 1 0 | 0 0 1 0 0 0 0 1 | 1 0 0 |  ← Raw Bits
//	|             Start           |               End               |  ← Bifurcated Phrases
//	|     Start 0     -  Start 1  | End 0 -      End 1      - End 2 |  ← Bifurcated Phrase Measurements
//
//	(Optional) Align() each phrase
//
//	| 0 1 0 0 1 1 0 1 | 0 0 0 1 0 - 1 1 0 | 0 0 1 0 0 - 0 0 1 | 1 0 0 |  ← Raw Bits
//	|     Start 0     -  Start 1  |       End 0       -     End 1     |  ← Aligned Measurements
//
// @formatter:on
func (phrase Phrase) Bifurcate() (start Phrase, end Phrase) {
	return phrase.Read(phrase.BitLength() / 2)
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

// Add performs binary addition between the source phrase and the provided phrase.
// The result will be at least as wide as the largest operand to be added.
func (phrase Phrase) Add(b Phrase) Phrase {
	aLen := phrase.BitLength()
	bLen := b.BitLength()
	length := int(math.Max(float64(aLen), float64(bLen)))

	if aLen < bLen {
		phrase = phrase.PadLeftToLength(length)
	} else {
		b = b.PadLeftToLength(length)
	}

	reader := func() (pA, pB Phrase) {
		pA, phrase = phrase.ReadFromEnd(1)
		pB, b = b.ReadFromEnd(1)
		return pA, pB
	}

	carry := Zero
	out := NewPhrase()
	for pA, pB := reader(); pA.BitLength() > 0; pA, pB = reader() {
		bitA := pA.Bits()[0]
		bitB := pB.Bits()[0]

		c := bitA + bitB + carry
		if carry == One {
			carry = Zero
		}
		if c > 1 {
			c -= 2
			carry = One
		}
		out = out.PrependBits(c)
	}
	if carry == One {
		out = out.PrependBits(One)
	}
	return out
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
