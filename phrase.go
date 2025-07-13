package tiny

import "fmt"

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

// GetBits returns a Bit slice of all the Phrase's underlying bits.
//
// If you'd prefer to Read specific measurements, you may provide an Expression.
func (a Phrase) GetBits(expr ...Expression) []Bit {
	bits := make([]Bit, a.BitLength())

	var measurements []Measurement
	if len(expr) == 0 {
		measurements = Express(Read.All(), a.Data)
	} else {
		measurements = Express(expr[0], a.Data)
	}

	for _, m := range measurements {
		for ii, b := range m.GetAllBits() {
			bits[ii] = b
		}
	}

	return bits
}

// BitLength gets the total bit length of this Phrase's recorded data.
func (a Phrase) BitLength() int {
	total := 0
	for _, m := range a.Data {
		total += m.BitLength()
	}
	return total
}

// Align ensures all Measurements are of the same width, with the last being smaller if measuring an uneven bit-width.
//
// If no width is provided, a standard alignment of 8-bits-per-byte will be used.
//
// For example-
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

	out := make([]Measurement, 0, a.BitLength()/8)
	for _, m := range a.Data {
		out = append(out, m)
	}

	return Phrase{
		Data:     out,
		Encoding: a.Encoding,
	}
}
