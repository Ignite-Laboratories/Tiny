package tiny

// Logical represents a variable width phrase where the measurements are all similar fixed logical widths.
type Logical struct {
	Phrase
}

// AsLogical converts the phrase into a logically aligned phrase of 8-bits-per-measurement.
//
// If you'd prefer a different logical width, you may provide it.
func (a Phrase) AsLogical(alignment ...uint) Logical {
	// TODO: Walk the bits of the phrase and actually split it into a single measurement per logical unit
	if len(alignment) > 0 {
		return Logical{a.Align(alignment[0])}
	}
	return Logical{a.Align()}
}

// StringPretty returns a phrase-formatted string of the current measurements.
//
// This means the bits will be placed between pipes and with dashes between measurements.
func (a Logical) StringPretty() string {
	return a.Phrase.StringPretty()
}

/**
IPhrase Passthrough Functions
*/

func (a Logical) GetData() []Measurement {
	return a.Phrase.GetData()
}

func (a Logical) BitWidth() uint {
	return a.Phrase.BitWidth()
}

func (a Logical) BleedLastBit() (Bit, any) {
	b, p := a.Phrase.BleedLastBit()
	a.Phrase = p
	return b, a
}

func (a Logical) BleedFirstBit() (Bit, any) {
	b, p := a.Phrase.BleedFirstBit()
	a.Phrase = p
	return b, a
}

func (a Logical) RollUp() any {
	p := a.Phrase.RollUp()
	a.Phrase = p
	return a
}

func (a Logical) Reverse() any {
	p := a.Phrase.Reverse()
	a.Phrase = p
	return a
}

func (a Logical) Append(bits ...Bit) any {
	p := a.Phrase.Append()
	a.Phrase = p
	return a
}

func (a Logical) AppendBytes(bytes ...byte) any {
	p := a.Phrase.AppendBytes()
	a.Phrase = p
	return a
}

func (a Logical) AppendMeasurement(m ...Measurement) any {
	p := a.Phrase.AppendMeasurement()
	a.Phrase = p
	return a
}

func (a Logical) Prepend(bits ...Bit) any {
	p := a.Phrase.Prepend()
	a.Phrase = p
	return a
}

func (a Logical) PrependBytes(bytes ...byte) any {
	p := a.Phrase.PrependBytes()
	a.Phrase = p
	return a
}

func (a Logical) PrependMeasurement(m ...Measurement) any {
	p := a.Phrase.PrependMeasurement()
	a.Phrase = p
	return a
}

func (a Logical) Align(width ...uint) any {
	p := a.Phrase.Align()
	a.Phrase = p
	return a
}
