package tiny

// Complex represents a phrase encoded as two measurements - a real number, and an imaginary number.
//
// Both components of the complex number can be any numeric type, but they will always be like typed.
type Complex struct {
	Phrase
}

func (a Complex) GetData() []Measurement {
	return a.Phrase.GetData()
}

func (a Complex) BitWidth() uint {
	return a.Phrase.BitWidth()
}
