package tiny

// NewPhrase calls NewMeasurement for each input byte and returns a slice of the results.
func NewPhrase(data ...byte) Phrase {
	out := make(Phrase, len(data))
	for i, d := range data {
		out[i] = NewMeasurement([]byte{d})
	}
	return out
}

// QuarterSplit calls Measurement.QuarterSplit for each Measurement of the Phrase.
func (phrase Phrase) QuarterSplit() {
	for i, m := range phrase {
		m.QuarterSplit()
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

// AllBelowThreshold checks if every Measurement of the Phrase is below the provided threshold value.
func (phrase Phrase) AllBelowThreshold(threshold int) bool {
	for _, m := range phrase {
		if m.Value() > threshold {
			return false
		}
	}
	return true
}
