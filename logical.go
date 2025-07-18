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
