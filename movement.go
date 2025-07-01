package tiny

import (
	"fmt"
	"math/big"
)

type Movement struct {
	Signature    Phrase
	Delta        Phrase
	DeltaWidth   int
	InitialWidth int
}

// Perform uses the current movement information to re-build the original information.
func (m Movement) Perform() Phrase {
	signature := m.Signature
	delta := m.Delta.AsBigInt()
	bitLength := m.InitialWidth

	for i := m.DeltaWidth; i < bitLength+1; i++ {
		var sign Bit
		sign, signature = signature.ReadLastBit()

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
