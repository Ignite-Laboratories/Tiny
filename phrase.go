package tiny

import (
	"fmt"
	"github.com/ignite-laboratories/core/relatively"
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
func (a Phrase) AsBigInt() *big.Int {
	out := new(big.Int)
	out.SetString(a.StringBinary(), 2)
	return out
}

/**
Append/Prepend
*/

// AppendMeasurement appends the provided measurement to the phrase.
func (a Phrase) AppendMeasurement(m Measurement) Phrase {
	return append(a, m)
}

// AppendBits appends the provided bits to the end of the phrase.
//
// NOTE: This appends the bits as new measurements inside of a new phrase =)
func (a Phrase) AppendBits(bits ...Bit) Phrase {
	return append(a, NewPhraseFromBits(bits...)...)
}

// AppendBytes appends the provided bytes to the end of the phrase.
//
// NOTE: This appends the bits as new measurements inside of a new phrase =)
func (a Phrase) AppendBytes(bytes ...byte) Phrase {
	return append(a, NewPhraseFromBytesAndBits(bytes)...)
}

// AppendBitsAndBytes appends the provided bits, and then bytes, to the end of the phrase.
//
// NOTE: This appends the bits as new measurements inside of a new phrase =)
func (a Phrase) AppendBitsAndBytes(bits []Bit, bytes ...byte) Phrase {
	return append(a, NewPhraseFromBitsAndBytes(bits, bytes...)...)
}

// AppendBytesAndBits appends the provided bytes, and then bits, to the end of the phrase.
//
// NOTE: This appends the bits as new measurements inside of a new phrase =)
func (a Phrase) AppendBytesAndBits(bytes []byte, bits ...Bit) Phrase {
	return append(a, NewPhraseFromBytesAndBits(bytes, bits...)...)
}

// Append appends the provided phrase(s) to the end of the source phrase.
func (a Phrase) Append(p ...Phrase) Phrase {
	out := make(Phrase, 0, len(a))
	for _, item := range p {
		out = append(a, item...)
	}
	return out
}

// PrependMeasurement prepends the provided measurement to the phrase.
func (a Phrase) PrependMeasurement(m Measurement) Phrase {
	return append(Phrase{m}, a...)
}

// PrependBits prepends the provided bits to the beginning of the phrase.
//
// NOTE: This prepends the bits as new measurements inside of a new phrase =)
func (a Phrase) PrependBits(bits ...Bit) Phrase {
	return append(NewPhraseFromBits(bits...), a...)
}

// PrependBytes prepends the provided bytes to the beginning of the phrase.
//
// NOTE: This prepends the bits as new measurements inside of a new phrase =)
func (a Phrase) PrependBytes(bytes ...byte) Phrase {
	return append(NewPhraseFromBytesAndBits(bytes), a...)
}

// PrependBitsAndBytes prepends the provided bits, and then bytes, to the beginning of the phrase.
//
// NOTE: This prepends the bits as new measurements inside of a new phrase =)
func (a Phrase) PrependBitsAndBytes(bits []Bit, bytes ...byte) Phrase {
	return append(NewPhraseFromBitsAndBytes(bits, bytes...), a...)
}

// PrependBytesAndBits prepends the provided bytes, and then bits, to the beginning of the phrase.
//
// NOTE: This prepends the bits as new measurements inside of a new phrase =)
func (a Phrase) PrependBytesAndBits(bytes []byte, bits ...Bit) Phrase {
	return append(NewPhraseFromBytesAndBits(bytes, bits...), a...)
}

// Prepend prepends the provided phrase(s) to the beginning of the source phrase.
func (a Phrase) Prepend(p ...Phrase) Phrase {
	out := make(Phrase, 0, len(a))
	for _, item := range p {
		out = append(item, a...)
	}
	return out
}

// ToBytesAndBits converts its measurements into bytes and the remainder of bits.
func (a Phrase) ToBytesAndBits() ([]byte, []Bit) {
	out := make([]byte, 0, len(a))

	current := make([]Bit, 8)
	var i int
	for _, measure := range a {
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
func (a Phrase) BitLength() int {
	total := 0
	for _, m := range a {
		total += m.BitLength()
	}
	return total
}

// CountBelowThreshold counts any Measurement of the Phrase that's below the provided threshold value.
func (a Phrase) CountBelowThreshold(threshold int) int {
	var count int
	for _, m := range a {
		if m.Value() < threshold {
			count++
		}
	}
	return count
}

// AllBelowThreshold checks if every Measurement of the Phrase is below the provided threshold value.
func (a Phrase) AllBelowThreshold(threshold int) bool {
	for _, m := range a {
		if m.Value() > threshold {
			return false
		}
	}
	return true
}

// BreakMeasurementsApart breaks each Measurement of the Phrase apart at the provided index and returns
// the two resulting phrases.  The left phrase will contain the most significant bits, while the right
// phrase will contain the least significant bits.
func (a Phrase) BreakMeasurementsApart(index int) (left Phrase, right Phrase) {
	left = make(Phrase, len(a))
	right = make(Phrase, len(a))

	for i, m := range a {
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
func (a Phrase) AsInts() []int {
	out := make([]int, len(a))
	for i, m := range a {
		out[i] = m.Value()
	}
	return out
}

// AsBytes converts each Measurement of the Phrase into a byte.
func (a Phrase) AsBytes() []byte {
	out := make([]byte, len(a))
	for i, m := range a {
		out[i] = byte(m.Value())
	}
	return out
}

// Bits returns a slice of the phrase's underlying bits.
func (a Phrase) Bits() []Bit {
	out := make([]Bit, 0, a.BitLength())
	for _, m := range a {
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
func (a Phrase) Align(width ...int) Phrase {
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

	src := a
	out := make(Phrase, 0, len(src))
	for {
		measure, remainder, _ := src.ReadMeasurement(w)
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
// NOTE: If you request more bits than are available, the slices will only contain the available bits =)
//
// NOTE: This is intended for reading long stretches of bits.
// If you wish to read less than your architecture's bit width from the first measurement, using ReadMeasurement is a
// little easier to work with.
func (a Phrase) Read(length int) (read Phrase, remainder Phrase, err error) {
	l := length

	read = make(Phrase, 0, len(a))
	remainder = make(Phrase, 0, len(a))

	for _, m := range a {
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

	if read.BitLength() < l {
		err = ErrorEndOfBits
	}

	return read, remainder, err
}

// ReadFromEnd reads the provided number of bits from the end of the source phrase, followed by the remainder,
// as phrases. The returned bits remain in the same order as they originally existed from left-to-right, merely
// grouped into separate phrases.
//
// For example -
//
//		 let 𝑛 = 3
//			|         |←  𝑛  →|                  read ⬎         ⬐ remainder
//			[ 1 1 0 1   0 1 1 ].ReadFromEnd(𝑛) -> [ 0 1 1] [ 1 1 0 1]
//	                 ⬑  Bits are in same order ⬏
//
// NOTE: If you request more bits than are available, the slices will only contain the available bits =)
func (a Phrase) ReadFromEnd(length int) (read Phrase, remainder Phrase, err error) {
	remainder, read, _ = a.Read(a.BitLength() - length)
	if read.BitLength() < length {
		err = ErrorEndOfBits
	}
	return read, remainder, err
}

// ReadLastBit reads the last bit of the source phrase, followed by the remainder.
func (a Phrase) ReadLastBit() (last Bit, remainder Phrase, err error) {
	read, remainder, _ := a.ReadFromEnd(1)
	last, _, _ = read.ReadNextBit()
	if read.BitLength() < 1 {
		err = ErrorEndOfBits
	}
	return last, remainder, err
}

// ReadMeasurement reads the provided number of bits from the source phrase as a Measurement and provides the
// remainder as a Phrase.
//
// NOTE: This will panic if you attempt to read more than your architecture's bit width.
// For that, please use Read.
func (a Phrase) ReadMeasurement(length int) (read Measurement, remainder Phrase, err error) {
	if length > GetArchitectureBitWidth() {
		panic(errorMeasurementLimit)
	}

	read = NewMeasurement([]byte{})
	readMeasures, remainder, _ := a.Read(length)
	for _, m := range readMeasures {
		read.Append(m)
	}

	if read.BitLength() < length {
		err = ErrorEndOfBits
	}

	return read, remainder, err
}

// ReadNextBit reads a single bit from the source phrase and returns the remainder as a Phrase.
//
// NOTE: This kindly returns an ErrorMsgEndOfBits error if there are no more bits to read.
func (a Phrase) ReadNextBit() (read Bit, remainder Phrase, err error) {
	measure, remainder, err := a.ReadMeasurement(1)
	if measure.BitLength() <= 0 {
		return 0, remainder, ErrorEndOfBits
	}
	return measure.GetAllBits()[0], remainder, nil
}

// ReadUntilOne reads the source phrase until it reaches the first 1, then returns the zero count and remainder.
//
// If you'd like it to stop after a certain count, provide a limit.
func (a Phrase) ReadUntilOne(limit ...int) (zeros int, remainder Phrase) {
	l := -1
	if len(limit) > 0 {
		l = limit[0]
	}

	remainder = a
	for b, r, err := remainder.ReadNextBit(); err == nil; b, r, err = r.ReadNextBit() {
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
func (a Phrase) PadLeftToLength(overall int, char ...Bit) Phrase {
	c := Zero
	if len(char) > 0 {
		c = char[0]
	}

	toPad := overall - a.BitLength()
	if toPad <= 0 {
		return a
	}

	if c == Zero {
		return a.Prepend(Synthesize.Zeros(toPad))
	}
	return a.Prepend(Synthesize.Ones(toPad))
}

// PadRightToLength pads the phrase to the desired overall length with zeros on the right (least significant) side
// // of the bits.
//
// NOTE: If you'd prefer to pad with ones, please override the char parameter with tiny.One
func (a Phrase) PadRightToLength(overall int, char ...Bit) Phrase {
	c := Zero
	if len(char) > 0 {
		c = char[0]
	}

	toPad := overall - a.BitLength()
	if toPad <= 0 {
		return a
	}

	if c == Zero {
		return a.Append(Synthesize.Zeros(toPad))
	}
	return a.Append(Synthesize.Ones(toPad))
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
func (a Phrase) Trifurcate(startLen int, middleLen int) (start Phrase, middle Phrase, end Phrase, err error) {
	start, end, _ = a.Read(startLen)
	middle, end, _ = end.Read(middleLen)

	if start.BitLength()+middle.BitLength() < startLen+middleLen {
		err = ErrorEndOfBits
	}

	return start, middle, end, err
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
func (a Phrase) Bifurcate() (start Phrase, end Phrase, err error) {
	return a.Read(a.BitLength() / 2)
}

// WalkBits walks the bits of the source phrase at the provided stride and calls the
// provided function for each measurement step.
//
// NOTE: This will panic if given a stride greater than your architecture's bit width.
func (a Phrase) WalkBits(stride int, fn func(int, Measurement)) {
	if stride > GetArchitectureBitWidth() {
		panic(errorMeasurementLimit)
	}
	if stride <= 0 {
		panic("cannot walk at a stride of 0 or less")
	}

	remainder := a
	var bitM Measurement
	i := 0
	for bitM, remainder, _ = remainder.ReadMeasurement(stride); len(remainder) > 0; bitM, remainder, _ = remainder.ReadMeasurement(stride) {
		if bitM.BitLength() > 0 {
			fn(i, bitM)
			i++
		}
	}
	if bitM.BitLength() > 0 {
		fn(i, bitM)
	}
}

// NOT applies the logical operation `𝑎 ^ 1` for every bit of phrase `𝑎` in order to produce phrase `𝑏`.
// If the phrase bit lengths are uneven, the shorter phrase is left-padded with 0s to match the longer length.
// The results are guaranteed to always follow the below truth table -
//
//	"The NOT Truth Table"
//
//	        𝑎 | 𝑏
//	        0 | 1
//	        1 | 0
func (a Phrase) NOT() (b Phrase) {
	b = NewPhrase()

	var bit Bit
	var err error
	for bit, a, err = a.ReadNextBit(); err == nil; bit, a, err = a.ReadNextBit() {
		b = b.AppendBits(bit ^ 1)
	}
	return b.Align()
}

// AND applies the logical operation `𝑎 & 𝑏` for every bit of phrases `𝑎` and `𝑏` in order to produce phrase `𝑐`.
// If the phrase bit lengths are uneven, the shorter phrase is left-padded with 0s to match the longer length.
// The results are guaranteed to always follow the below truth table -
//
//	"The AND Truth Table"
//
//	     𝑎 | 𝑏 | 𝑐
//	     0 | 0 | 0
//	     0 | 1 | 0
//	     1 | 0 | 0
//	     1 | 1 | 1
func (a Phrase) AND(b Phrase) (c Phrase) {
	a, b = padToSameLength(a, b)

	c = NewPhrase()
	for bitA, bitB, pA, pB, err := readTwoPhrasesNextBit(a, b); err == nil; bitA, bitB, pA, pB, err = readTwoPhrasesNextBit(a, b) {
		a = pA
		b = pB
		result := bitA & bitB

		c = c.AppendBits(result)
	}
	return c.Align()
}

// OR applies the logical operation `𝑎 | 𝑏` for every bit of phrases `𝑎` and `𝑏` in order to produce phrase `𝑐`.
// If the phrase bit lengths are uneven, the shorter phrase is left-padded with 0s to match the longer length.
// The results are guaranteed to always follow the below truth table -
//
//	"The OR Truth Table"
//
//	     𝑎 | 𝑏 | 𝑐
//	     0 | 0 | 0
//	     0 | 1 | 1
//	     1 | 0 | 1
//	     1 | 1 | 1
func (a Phrase) OR(b Phrase) (c Phrase) {
	a, b = padToSameLength(a, b)

	c = NewPhrase()
	for bitA, bitB, pA, pB, err := readTwoPhrasesNextBit(a, b); err == nil; bitA, bitB, pA, pB, err = readTwoPhrasesNextBit(a, b) {
		a = pA
		b = pB
		result := bitA | bitB

		c = c.AppendBits(result)
	}
	return c.Align()
}

// XOR applies the logical operation `𝑎 ^ 𝑏` for every bit of phrases `𝑎` and `𝑏` in order to produce phrase `𝑐`.
// If the phrase bit lengths are uneven, the shorter phrase is left-padded with 0s to match the longer length.
// The results are guaranteed to always follow the below truth table -
//
//	"The XOR Truth Table"
//
//	     𝑎 | 𝑏 | 𝑐
//	     0 | 0 | 0
//	     0 | 1 | 1
//	     1 | 0 | 1
//	     1 | 1 | 0
func (a Phrase) XOR(b Phrase) (c Phrase) {
	a, b = padToSameLength(a, b)

	c = NewPhrase()
	for bitA, bitB, pA, pB, err := readTwoPhrasesNextBit(a, b); err == nil; bitA, bitB, pA, pB, err = readTwoPhrasesNextBit(a, b) {
		a = pA
		b = pB
		result := bitA ^ bitB

		c = c.AppendBits(result)
	}
	return c.Align()
}

// XNOR applies the logical operation `^(𝑎 ^ 𝑏)` for every bit of phrases `𝑎` and `𝑏` in order to produce phrase `𝑐`.
// If the phrase bit lengths are uneven, the shorter phrase is left-padded with 0s to match the longer length.
// The results are guaranteed to always follow the below truth table -
//
//	"The XNOR Truth Table"
//
//	     𝑎 | 𝑏 | 𝑐
//	     0 | 0 | 1
//	     0 | 1 | 0
//	     1 | 0 | 0
//	     1 | 1 | 1
func (a Phrase) XNOR(b Phrase) (c Phrase) {
	a, b = padToSameLength(a, b)

	c = NewPhrase()
	for bitA, bitB, pA, pB, err := readTwoPhrasesNextBit(a, b); err == nil; bitA, bitB, pA, pB, err = readTwoPhrasesNextBit(a, b) {
		a = pA
		b = pB
		result := 1 ^ (bitA ^ bitB)

		c = c.AppendBits(result)
	}
	return c.Align()
}

// NAND applies the logical operation `^(𝑎 & 𝑏)` for every bit of phrases `𝑎` and `𝑏` in order to produce phrase `𝑐`.
// If the phrase bit lengths are uneven, the shorter phrase is left-padded with 0s to match the longer length.
// The results are guaranteed to always follow the below truth table -
//
//	"The NAND Truth Table"
//
//	     𝑎 | 𝑏 | 𝑐
//	     0 | 0 | 1
//	     0 | 1 | 1
//	     1 | 0 | 1
//	     1 | 1 | 0
func (a Phrase) NAND(b Phrase) (c Phrase) {
	a, b = padToSameLength(a, b)

	c = NewPhrase()
	for bitA, bitB, pA, pB, err := readTwoPhrasesNextBit(a, b); err == nil; bitA, bitB, pA, pB, err = readTwoPhrasesNextBit(a, b) {
		a = pA
		b = pB
		result := 1 ^ (bitA & bitB)

		c = c.AppendBits(result)
	}
	return c.Align()
}

// NOR applies the logical operation `^(𝑎 | 𝑏)` for every bit of phrases `𝑎` and `𝑏` in order to produce phrase `𝑐`.
// If the phrase bit lengths are uneven, the shorter phrase is left-padded with 0s to match the longer length.
// The results are guaranteed to always follow the below truth table -
//
//	"The NOR Truth Table"
//
//	     𝑎 | 𝑏 | 𝑐
//	     0 | 0 | 1
//	     0 | 1 | 0
//	     1 | 0 | 0
//	     1 | 1 | 0
func (a Phrase) NOR(b Phrase) (c Phrase) {
	a, b = padToSameLength(a, b)

	c = NewPhrase()
	for bitA, bitB, pA, pB, err := readTwoPhrasesNextBit(a, b); err == nil; bitA, bitB, pA, pB, err = readTwoPhrasesNextBit(a, b) {
		a = pA
		b = pB
		result := 1 ^ (bitA | bitB)

		c = c.AppendBits(result)
	}
	return c.Align()
}

// Add performs binary addition between the source and provided phrases.
// The result will be at least as wide as the largest operand to be added.
func (a Phrase) Add(b Phrase) (c Phrase) {
	a, b = padToSameLength(a, b)

	carry := Zero
	c = NewPhrase()
	for bitA, bitB, pA, pB, err := readTwoPhrasesLastBit(a, b); err == nil; bitA, bitB, pA, pB, err = readTwoPhrasesLastBit(a, b) {
		a = pA
		b = pB

		ab := bitA + bitB + carry
		if carry == One {
			carry = Zero
		}
		if ab > 1 {
			ab -= 2
			carry = One
		}
		c = c.PrependBits(ab)
	}
	if carry == One {
		c = c.PrependBits(One)
	}
	return c.ToNumericForm().Align()
}

// Minus performs absolute binary subtraction between the source and provided phrases.
func (a Phrase) Minus(b Phrase) (c Phrase, negative bool) {
	a, b = padToSameLength(a, b)

	// We are performing -absolute- subtraction, so we the minuend must be greater than the subtrahend.
	if a.CompareTo(b) == relatively.Before {
		a, b = b, a
		negative = true
	}

	borrow := Zero
	c = NewPhrase()
	for bitA, bitB, pA, pB, err := readTwoPhrasesLastBit(a, b); err == nil; bitA, bitB, pA, pB, err = readTwoPhrasesLastBit(a, b) {
		a = pA
		b = pB

		if borrow == One {
			borrow = Zero
			bitA -= 1

			// Check for general byte underflow
			if bitA > 11 {
				borrow = One
				bitA += 2
			}
		}

		ab := bitA - bitB
		// Check for general byte underflow
		if ab > 11 {
			ab += 2
			borrow = One
		}
		c = c.PrependBits(Bit(ab))
	}
	if borrow == One {
		c = c.PrependBits(One)
	}
	return c.ToNumericForm().Align(), negative
}

// Times performs absolute binary multiplication between the source and provided phrases.
func (a Phrase) Times(b Phrase) (c Phrase) {
	c = NewPhrase()

	d := make([]Phrase, 0, b.BitLength())
	shift := 0

	for bitB, pB, errB := b.ReadLastBit(); errB == nil; bitB, pB, errB = b.ReadLastBit() {
		b = pB
		e := Synthesize.Zeros(shift)

		tempA := a
		for bitA, pA, errA := tempA.ReadLastBit(); errA == nil; bitA, pA, errA = tempA.ReadLastBit() {
			tempA = pA

			if bitA == 1 && bitB == 1 {
				e = e.PrependBits(1)
			} else {
				e = e.PrependBits(0)
			}
		}

		d = append(d, e)
		shift++
	}

	for _, p := range d {
		c = c.Add(p)
	}

	return c.ToNumericForm().Align()
}

// ToNumericForm removes any leading zeros from the phrase's bits, representing the digits in their
// smallest possible form.  This is the same as calling ReadUntilOne() and ignoring the zero count.
func (a Phrase) ToNumericForm() Phrase {
	_, remainder := a.ReadUntilOne()
	return remainder
}

// CompareTo determines if the numeric value of 𝑎 comes relatively.Before, relatively.Same as, or relatively.After the numeric value of 𝑏.
func (a Phrase) CompareTo(b Phrase) relatively.Relativity {
	a, b = padToSameLength(a, b)

	for bitA, bitB, pA, pB, err := readTwoPhrasesNextBit(a, b); err == nil; bitA, bitB, pA, pB, err = readTwoPhrasesNextBit(a, b) {
		a = pA
		b = pB

		if bitA > bitB {
			return relatively.After
		} else if bitA < bitB {
			return relatively.Before
		}
	}
	return relatively.Aligned
}

// Int returns a value up to your architecture's bit width from the source phrase and ignores
// the remaining bits entirely.
//
// For fully featured binary arithmetic operations, including easily converting between bases,
// please check out math.Big!  It's an EXCELLENT library with a rigorous amount of testing by
// far smarter individuals than myself.  You can convert between a phrase and a big.Int using
// Phrase.AsBigInt().
func (a Phrase) Int() int {
	bitWidth := GetArchitectureBitWidth()
	read, _, _ := a.Read(bitWidth)
	bits := read.Bits()
	return To.Number(len(bits), bits...)
}

// StringBinary returns the phrase's bits as a binary string of 1s and 0s.
func (a Phrase) StringBinary() string {
	out := ""
	for _, m := range a {
		out += m.StringBinary()
	}
	return out
}

func (a Phrase) String() string {
	out := make([]string, len(a))
	for i, m := range a {
		out[i] = m.StringBinary()
	}
	return fmt.Sprintf("%v", out)
}

/**
CONVENIENCE METHODS
*/

// readTwoPhrasesFromEnd is a convenience method to keep the code DRY.  It offers no gaurantees, by design.
// Please do not expose this method.
//
//	Functionality: This reads 'l' positions from the end of both phrases.
func readTwoPhrasesFromEnd(l int, a Phrase, b Phrase) (bitsA, bitsB []Bit, phraseA, phraseB Phrase) {
	var pA, pB Phrase
	pA, a, _ = a.ReadFromEnd(l)
	pB, b, _ = b.ReadFromEnd(l)
	return pA.Bits(), pB.Bits(), a, b
}

// readTwoPhrasesFromEnd is a convenience method to keep the code DRY.  It offers no gaurantees, by design.
// Please do not expose this method.
//
//	Functionality: This reads the rightmost bit from both phrases.
func readTwoPhrasesLastBit(a Phrase, b Phrase) (bitA, bitB Bit, phraseA, phraseB Phrase, err error) {
	bitsA, bitsB, a, b := readTwoPhrasesFromEnd(1, a, b)
	if len(bitsA) == 0 {
		return Zero, Zero, a, b, ErrorEndOfBits
	}
	return bitsA[0], bitsB[0], a, b, nil
}

// readTwoPhrases is a convenience method to keep the code DRY.  It offers no gaurantees, by design.
// Please do not expose this method.
//
//	Functionality: This reads 'l' positions from the start of both phrases.
func readTwoPhrases(l int, a Phrase, b Phrase) (bitsA, bitsB []Bit, phraseA, phraseB Phrase) {
	var pA, pB Phrase
	pA, a, _ = a.Read(l)
	pB, b, _ = b.Read(l)
	return pA.Bits(), pB.Bits(), a, b
}

// readTwoPhrasesFromEnd is a convenience method to keep the code DRY.  It offers no gaurantees, by design.
// Please do not expose this method.
//
//	Functionality: This reads the leftmost bit from both phrases.
func readTwoPhrasesNextBit(a Phrase, b Phrase) (bitA, bitB Bit, phraseA, phraseB Phrase, err error) {
	bitsA, bitsB, a, b := readTwoPhrases(1, a, b)
	if len(bitsA) == 0 {
		return Zero, Zero, a, b, ErrorEndOfBits
	}
	return bitsA[0], bitsB[0], a, b, nil
}

// padToSameLength is a convenience method to keep the code DRY.  It offers no gaurantees, by design.
// Please do not expose this method.
//
//	Functionality: This left pads the shorter phrase with 0s to match the length of the longer phrase.
func padToSameLength(a, b Phrase) (Phrase, Phrase) {
	aLen := a.BitLength()
	bLen := b.BitLength()
	length := int(math.Max(float64(aLen), float64(bLen)))

	if aLen < bLen {
		a = a.PadLeftToLength(length)
	} else {
		b = b.PadLeftToLength(length)
	}
	return a, b
}
