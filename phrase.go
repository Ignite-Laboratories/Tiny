package tiny

import "strings"

// Phrase represents a collection of measurements, plus their Encoding scheme.
type Phrase struct {
	Data []Measurement
	Encoding
}

// NewLogicalPhrase creates a new Logical Phrase of the provided measurements.
func NewLogicalPhrase(m ...Measurement) Phrase {
	return Phrase{
		Data:     m,
		Encoding: Logical,
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
		a.Data = append(a.Data, NewMeasurementFromBytes(bytes...))
		return a
	}

	last := len(a.Data) - 1
	a.Data[last] = a.Data[last].AppendBytes(bytes...)
	return a.RollUp()
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
		a.Data = append(a.Data, NewMeasurementFromBytes(bytes...))
		return a
	}

	a.Data[0] = a.Data[0].PrependBytes(bytes...)
	return a.RollUp()
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
	current := make([]Bit, 0, 8)
	i := 0

	for _, m := range a.Data {
		for _, b := range m.GetAllBits() {
			current = append(current, b)
			i++
			if i == 8 {
				i = 0
				out = append(out, NewMeasurement(current...))
			}
		}
	}

	if len(current) > 0 {
		out = append(out, NewMeasurement(current...))
	}

	return Phrase{
		Data:     out,
		Encoding: a.Encoding,
	}
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
	// TODO: Reverse the phrase
	return Phrase{}
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

	builder.WriteString("| ")

	return builder.String()
}
