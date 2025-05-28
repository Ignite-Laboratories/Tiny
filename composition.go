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
	Movements           map[string][]Passage
	Subdivisions        int
	SubdivisionBitWidth int
	Target              Phrase
	StartBitLength      int
}

func Distill(phrase Phrase, subdivisions int, subdivisionBitWidth int) *Composition {
	c := &Composition{
		Movements:           make(map[string][]Passage),
		Subdivisions:        subdivisions,
		SubdivisionBitWidth: subdivisionBitWidth,
	}
	c.Movements[MovementStart] = make([]Passage, 0)
	c.Movements[MovementPathway] = make([]Passage, 0)
	c.Movements[MovementSeed] = make([]Passage, 0)
	c.StartBitLength = phrase.AsBigInt().BitLen()
	c.Target = phrase
	c.distill(c.Target.AsBigInt())
	return c
}

// AddPassageToStart appends the provided passage to the end of the start movement.
func (c *Composition) AddPassageToStart(passage Passage) {
	c.Movements[MovementStart] = append(c.Movements[MovementStart], passage)
}

// AddPassageToPathway appends the provided passage to the end of the pathway movement.
func (c *Composition) AddPassageToPathway(passage Passage) {
	c.Movements[MovementPathway] = append(c.Movements[MovementPathway], passage)
}

// AddPassageToSeed appends the provided passage to the end of the seed movement.
func (c *Composition) AddPassageToSeed(passage Passage) {
	c.Movements[MovementSeed] = append(c.Movements[MovementSeed], passage)
}

func (c *Composition) distill(target *big.Int, height ...*big.Int) {
	passage := Passage{}
	if target.BitLen() <= 64 {
		return
	}

	var h *big.Int
	if len(height) > 0 {
		h = height[0]
	} else {
		upperBound := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(c.StartBitLength)), nil)
		h = new(big.Int).Div(upperBound, big.NewInt(int64(c.Subdivisions)))
	}
	nextH := new(big.Int).Div(h, big.NewInt(int64(c.Subdivisions)))

	multiplier := new(big.Int).Div(target, h)

	// "Box" in the target value
	low := new(big.Int).Mul(h, multiplier)
	high := new(big.Int).Add(low, h)

	// Get the deltas
	deltaLow := new(big.Int).Sub(target, low)
	deltaHigh := new(big.Int).Sub(high, target)

	// Find which is smaller
	switch deltaLow.Cmp(deltaHigh) {
	case -1:
		fallthrough
	case 0:
		// Low or equal condition
		passage = passage.Append(NewPhraseFromBits(0)) // Undershoot
		passage = passage.Append(NewPhraseFromBits(From.Number(int(multiplier.Int64()), c.SubdivisionBitWidth)...))
		c.AddPassageToPathway(passage)
		c.distill(deltaLow, nextH)
	case 1:
		// High condition
		passage = passage.Append(NewPhraseFromBits(1)) // Overshoot
		passage = passage.Append(NewPhraseFromBits(From.Number(int(multiplier.Int64()), c.SubdivisionBitWidth)...))
		c.AddPassageToPathway(passage)
		c.distill(deltaHigh, nextH)
	}
}
