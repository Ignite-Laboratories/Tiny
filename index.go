package tiny

// Index represents an implicitly fixed-width phrase of raw binary information.
type Index struct {
	Phrase
}

func (a Index) GetData() []Measurement {
	return a.Phrase.GetData()
}

func (a Index) BitWidth() uint {
	return a.Phrase.BitWidth()
}
