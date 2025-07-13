package tiny

type Phrase struct {
	Data []Measurement
	Encoding
}

// NewLogicalPhrase creates a new Logical Phrase of the provided measurements.
func NewLogicalPhrase(m ...Measurement) Phrase {
	return Phrase{
		Data: m,
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
