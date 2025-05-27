package tiny

import (
	"math/big"
)

// Composition represents a working structure for synthesizing binary information.
//
// A composition consists of Movements which can be referenced and updated while distilling binary information.
//
// Passage - A passage is a collection of related phrases.
// This allows for looping phrases to be clustered together.
//
// Phrase - A phrase is a way of grouping together longer runs of bits.
// Since a measurement is limited in width, this allows any arbitrary length of bits to be grouped together.
//
// Measurement - A measurement is a variable width slice of bits up to MaxMeasurementBitLength.
type Composition struct {
	Movements map[string]Passage
}

func Distill(phrase Phrase) *Composition {
	c := &Composition{
		Movements: make(map[string]Passage),
	}
	c.Movements[MovementStart] = make(Passage, 0)
	c.Movements[MovementPathway] = make(Passage, 0)
	c.Movements[MovementSeed] = make(Passage, 0)
	startLen := phrase.BitLength()
	target := phrase.AsBigInt()
	endLen := target.BitLen()
	c.shrink(target, startLen-endLen)
	return c
}

func (c *Composition) AddPassageToStart(passage Passage) {
	c.Movements[MovementStart] = append(c.Movements[MovementStart], passage...)
}

func (c *Composition) AddPassageToPathway(passage Passage) {
	c.Movements[MovementPathway] = append(c.Movements[MovementPathway], passage...)
}

func (c *Composition) AddPassageToSeed(passage Passage) {
	c.Movements[MovementSeed] = append(c.Movements[MovementSeed], passage...)
}

func (c *Composition) shrink(target *big.Int, delta int) {
	// 0 - Encode the bit length delta as a new passage
	passage := NewZLE64PassageInt(delta)
	bitLength := big.NewInt(int64(target.BitLen()))

	// 1 - Calculate the subdivision height and index
	upperBound := new(big.Int).Exp(big.NewInt(2), bitLength, nil)
	height := upperBound.Div(upperBound, big.NewInt(8))
	index := new(big.Int).Div(target, height)

	// 2 - Encode the index onto the current passage
	passage = append(passage, NewPhraseFromBigInt(index))

	// 3 - Calculate the difference from the target to the closest index
	difference := new(big.Int).Sub(target, index.Mul(index, height))

	// 4 - Add the current passage to the pathway
	c.AddPassageToPathway(passage)

	// 5 - Recurse on the difference if longer than 64 bits
	if difference.BitLen() > 64 {
		c.shrink(difference, target.BitLen()-difference.BitLen())
	}
}
