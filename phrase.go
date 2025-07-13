package tiny

type Phrase []Measurement

// BitLength gets the total bit length of this Phrase's recorded data.
func (a Phrase) BitLength() int {
	total := 0
	for _, m := range a {
		total += m.BitLength()
	}
	return total
}
