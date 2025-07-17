package tiny

import (
	"strings"
)

// Phrase represents a collection of raw binary measurements and their observed Endianness at the time of recording.
type Phrase struct {
	Name string
	Data []Measurement
	Endianness
}

// Logical represents a variable width phrase where the measurements are all similar fixed logical widths.
type Logical Phrase

// Index represents an implicitly fixed-width phrase of raw binary information.
type Index Phrase

// Natural represents a phrase holding a value belonging to the set of natural numbers, including zero.
//
// To those who think zero shouldn't be included in the set of natural numbers, I present a counter-argument:
// Base 1 has only one identifier, meaning it can only "represent" zero by -not- holding a value in an observable
// location.  Subsequently, all bases are built upon determining the size of a value through "identification" - in
// binary, through a series of zeros or ones, in decimal through the identifiers 0-9.
//
// Now here's where it gets tricky: a value doesn't even EXIST until it is given a place to exist within, meaning its
// existence directly implies a void which has now been filled - an identifiable "zero" state.  In fact, the very first
// identifier of all higher order bases (zero) specifically identifies this state!  Counting, itself, comes from the act of observing
// the general relativistic -presence- of anything - fingers, digits, different length squiggles, feelings - meaning to exclude
// zero attempts to redefine the very fundamental definition of identification itself: it's PERFECTLY reasonable to -naturally-
// count -zero- hairs on a magnificently bald head!
//
//	tl;dr - to count naturally involves identification, which implies accepting -non-existence- as a countable state
//
// I should note this entire system hinges on one fundamental flaw - this container technically holds one additional value beyond
// the 'natural' number set: nil!  I call this the "programmatic set" of numbers, and I can't stop you from setting your natural
// phrase to it, but I can empower you with awareness =)
type Natural Phrase

// Integer represents a phrase encoded as two measurements - a sign bit, and an arbitrary bit-width value.
//
// NOTE: The entire goal of tiny is to break away from the boundaries of overflow logic - if you explicitly
// require working with index-based overflow logic, please use an Index phrase.
type Integer Phrase

// Float32 represents a 32-bit phrase encoded as three measurements in accordance with IEEE 754 -
//
//	    Sign: 1 bit
//	Exponent: 8 bits
//	Mantissa: 23 bits
type Float32 Phrase

// Float64 represents a 64-bit phrase encoded as three measurements in accordance with IEEE 754 -
//
//	    Sign: 1 bit
//	Exponent: 11 bits
//	Mantissa: 52 bits
type Float64 Phrase

// Float128 represents a 128-bit phrase encoded as three measurements in accordance with IEEE 754 -
//
//	    Sign: 1 bit
//	Exponent: 15 bits
//	Mantissa: 112 bits
type Float128 Phrase

// Float256 represents a 256-bit phrase encoded as three measurements in accordance with IEEE 754 -
//
//	    Sign: 1 bit
//	Exponent: 19 bits
//	Mantissa: 236 bits
type Float256 Phrase

/**
New Functions
*/

// NewPhrase creates a named Phrase of the provided measurements, encoding scheme, and endianness.
func NewPhrase(name string, endianness Endianness, m ...Measurement) Phrase {
	return Phrase{
		Name:       name,
		Data:       m,
		Endianness: endianness,
	}
}

// NewPhraseFromBits creates a named Phrase of the provided bits, encoding scheme, and endianness.
func NewPhraseFromBits(name string, endianness Endianness, bits ...Bit) Phrase {
	return Phrase{
		Name:       name,
		Data:       []Measurement{NewMeasurement(bits...)},
		Endianness: endianness,
	}
}

// GetData returns the phrase's measurement data.  This is exposed as a method to guarantee
// the encoded accessors for any derived types are grouped together in your IDE's type-ahead.
func (a Phrase) GetData() []Measurement {
	return a.Data
}

// GetAllBits returns a slice of the Phrase's individual bits.
func (a Phrase) GetAllBits() []Bit {
	out := make([]Bit, 0, a.BitWidth())
	for _, m := range a.Data {
		out = append(out, m.GetAllBits()...)
	}
	return out
}

// BitWidth gets the total bit width of this Phrase's recorded data.
func (a Phrase) BitWidth() uint {
	total := uint(0)
	for _, m := range a.Data {
		total += m.BitWidth()
	}
	return uint(total)
}

// Append places the provided bits at the end of the Phrase.
func (a Phrase) Append(bits ...Bit) Phrase {
	if len(a.Data) == 0 {
		a.Data = append(a.Data, NewMeasurement(bits...))
		return a
	}

	last := len(a.Data) - 1
	a.Data[last] = a.Data[last].Append(bits...)
	return a.RollUp()
}

// AppendBytes places the provided bits at the end of the Phrase.
func (a Phrase) AppendBytes(bytes ...byte) Phrase {
	if len(a.Data) == 0 {
		a.Data = append(a.Data, NewMeasurementOfBytes(bytes...))
		return a
	}

	last := len(a.Data) - 1
	a.Data[last] = a.Data[last].AppendBytes(bytes...)
	return a.RollUp()
}

// AppendMeasurement places the provided measurement at the end of the Phrase.
func (a Phrase) AppendMeasurement(m ...Measurement) Phrase {
	a.Data = append(a.Data, m...)
	return a
}

// Prepend places the provided bits at the start of the Phrase.
func (a Phrase) Prepend(bits ...Bit) Phrase {
	if len(a.Data) == 0 {
		a.Data = append(a.Data, NewMeasurement(bits...))
		return a
	}

	a.Data[0] = a.Data[0].Prepend(bits...)
	return a.RollUp()
}

// PrependBytes places the provided bytes at the start of the Phrase.
func (a Phrase) PrependBytes(bytes ...byte) Phrase {
	if len(a.Data) == 0 {
		a.Data = append(a.Data, NewMeasurementOfBytes(bytes...))
		return a
	}

	a.Data[0] = a.Data[0].PrependBytes(bytes...)
	return a.RollUp()
}

// PrependMeasurement places the provided measurement at the start of the Phrase.
func (a Phrase) PrependMeasurement(m ...Measurement) Phrase {
	a.Data = append(m, a.Data...)
	return a
}

// Align ensures all Measurements are of the same width, with the last being smaller if measuring an uneven bit-width.
//
// If no width is provided, a standard alignment of 8-bits-per-byte will be used.
//
// For example -
//
//	let a = an un-aligned logical Phrase
//
//	| 0 1 - 0 1 0 - 0 1 1 0 1 0 0 0 - 1 0 1 1 0 - 0 0 1 0 0 0 0 1 |  ← Raw Bits
//	|  M0 -  M1   -  Measurement 2  -     M3    -  Measurement 4  |  ← Un-aligned Measurements
//
//	a.Align()
//
//	| 0 1 0 1 0 0 1 1 - 0 1 0 0 0 1 0 1 - 1 0 0 0 1 0 0 0 - 0 1 |  ← Raw Bits
//	|  Measurement 0  -  Measurement 1  -  Measurement 2  - M3  |  ← Aligned Measurements
//
//	a.Align(16)
//
//	| 0 1 0 1 0 0 1 1 0 1 0 0 0 1 0 1 - 1 0 0 0 1 0 0 0 0 1 |  ← Raw Bits
//	|          Measurement 0          -    Measurement 1    |  ← Aligned Measurements
func (a Phrase) Align(width ...uint) Phrase {
	w := 8
	if len(width) > 0 {
		w = int(width[0])
	}

	out := make([]Measurement, 0, int(a.BitWidth())/w)
	current := make([]Bit, 0, w)
	i := 0

	for _, m := range a.Data {
		for _, b := range m.GetAllBits() {
			current = append(current, b)
			i++
			if i == 8 {
				i = 0
				out = append(out, NewMeasurement(current...))
				current = make([]Bit, 0, w)
			}
		}
	}

	if len(current) > 0 {
		out = append(out, NewMeasurement(current...))
	}

	return Phrase{
		Data: out,
	}
}

// BleedLastBit returns the last bit of the phrase and a phrase missing that bit.
//
// NOTE: This is a destructive operation to the phrase's encoding scheme and returns a Raw Phrase.
func (a Phrase) BleedLastBit() (Bit, Phrase) {
	if a.BitWidth() == 0 {
		panic("cannot bleed the last bit of an empty phrase")
	}

	lastBit, lastMeasurement := a.Data[len(a.Data)-1].BleedLastBit()
	a.Data[len(a.Data)-1] = lastMeasurement
	return lastBit, a
}

// BleedFirstBit returns the first bit of the phrase and a phrase missing that bit.
//
// NOTE: This is a destructive operation to the phrase's encoding scheme and returns a Raw Phrase.
func (a Phrase) BleedFirstBit() (Bit, Phrase) {
	if a.BitWidth() == 0 {
		panic("cannot bleed the first bit of an empty phrase")
	}

	firstBit, firstMeasurement := a.Data[0].BleedFirstBit()
	a.Data[0] = firstMeasurement
	return firstBit, a
}

// RollUp calls Measurement.RollUp for every measurement in the phrase.
func (a Phrase) RollUp() Phrase {
	for i, m := range a.Data {
		a.Data[i] = m.RollUp()
	}
	return a
}

// Reverse reverses the order of all bits in the phrase.
func (a Phrase) Reverse() Phrase {
	// TODO: Reverse the measurement order
	return a
}

// String returns a string consisting entirely of 1s and 0s.
func (a Phrase) String() string {
	builder := strings.Builder{}
	builder.Grow(int(a.BitWidth()))

	for _, m := range a.Data {
		builder.WriteString(m.String())
	}

	return builder.String()
}

// StringPretty returns a phrase-formatted string of the current measurements.
//
// This means the bits will be placed between pipes and with dashes between measurements.
func (a Phrase) StringPretty() string {
	if len(a.Data) == 0 {
		return "||"
	}

	totalSize := 4 + (len(a.Data)-1)*3
	for _, m := range a.Data {
		totalSize += int(m.BitWidth()) * 2 // ← Approximately account for Measurement's StringPretty() spacing
	}

	builder := strings.Builder{}
	builder.Grow(totalSize)

	builder.WriteString("| ")

	builder.WriteString(a.Data[0].StringPretty())

	for _, m := range a.Data[1:] {
		builder.WriteString(" - ")
		builder.WriteString(m.StringPretty())
	}

	builder.WriteString(" | ")

	return builder.String()
}
