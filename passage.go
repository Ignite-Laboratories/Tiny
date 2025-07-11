package tiny

import (
	"fmt"
	"math/big"
)

// Passage represents a conversion between encoded and traditional values.
type Passage struct {
	Signature    Phrase
	Delta        Phrase
	DeltaWidth   int
	InitialWidth int
}

// Perform uses the current passage information to re-build the original information.
func (p Passage) Perform() Phrase {
	signature := p.Signature
	delta := p.Delta.AsBigInt()
	bitLength := p.InitialWidth

	for i := p.DeltaWidth; i < bitLength+1; i++ {
		var sign Bit
		sign, signature, _ = signature.ReadLastBit()

		midpoint := Synthesize.Midpoint(i)
		fmt.Println(midpoint)

		if sign == One {
			delta = new(big.Int).Sub(midpoint.AsBigInt(), delta)
		} else {
			delta = new(big.Int).Add(midpoint.AsBigInt(), delta)
		}
		fmt.Println(delta.Text(2))
	}

	return NewPhraseFromBigInt(delta)
}

// AsPhrase returns the passage as a Phrase aligned to the provided alignment.
//
// NOTE: If no alignment is provided, a standard of 8 bits is used.
func (p Passage) AsPhrase(alignment ...int) Phrase {
	a := 8
	if len(alignment) > 0 {
		a = alignment[0]
		if a < 1 {
			a = 1
		}
	}
	return append(p.Signature, p.Delta...).Align(a)
}
