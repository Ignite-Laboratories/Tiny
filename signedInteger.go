package tiny

// SignedInteger represents a phrase where the first measurement is a sign and the remaining bits are the data.
type SignedInteger Phrase

// NewInteger creates a new SignedInteger.
func NewInteger(sign Bit, data Phrase) SignedInteger {
	return SignedInteger(NewPhraseFromBits(sign).Append(data))
}

// NewIntegerFromBool creates a new SignedInteger using a boolean for the sign.
func NewIntegerFromBool(sign bool, data Phrase) SignedInteger {
	if sign {
		return SignedInteger(NewPhraseFromBits(1).Append(data))
	}
	return SignedInteger(NewPhraseFromBits(0).Append(data))
}

// GetSign returns the first bit as the sign bit.
func (a SignedInteger) GetSign() Bit {
	if len(a) > 0 {
		sign, _, _ := Phrase(a).ReadNextBit()
		return sign
	}
	return 0
}

// GetValue returns everything after the initial sign bit.
func (a SignedInteger) GetValue() Phrase {
	if len(a) > 0 {
		return Phrase(a[1:])
	}
	return NewPhrase()
}
