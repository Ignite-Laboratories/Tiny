package tiny

import "math/big"

type Movement struct {
	Signature  Phrase
	Delta      Phrase
	DeltaWidth int
}

// Perform uses the current movement information to re-build the original information.
func (m Movement) Perform() Phrase {
	signature := m.Signature
	bitLength := m.Signature.BitLength() + m.DeltaWidth
	delta := m.Delta.AsBigInt()

	for i := 0; i < bitLength; i++ {
		var sign Bit
		sign, signature = signature.ReadLastBit()

		midpoint := Synthesize.Midpoint(i + m.DeltaWidth + 1)

		if sign == One {
			delta = new(big.Int).Sub(midpoint.AsBigInt(), delta)
		} else {
			delta = new(big.Int).Add(midpoint.AsBigInt(), delta)
		}
	}

	return NewPhraseFromBigInt(delta)
}
