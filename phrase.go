package tiny

// NewPhrase calls NewMeasurement for each input byte and returns a slice of the results.
func NewPhrase(data ...byte) Phrase {
	out := make(Phrase, len(data))
	for i, d := range data {
		out[i] = NewMeasurement([]byte{d})
	}
	return out
}

// ToBytesAndBits converts its measurements into bytes and a remainder of bits.
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
