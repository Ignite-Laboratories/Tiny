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
// Measurement - A measurement is a variable width slice of bits up to your architecture's bit width.
type Composition struct {
	Movements   map[string][]Passage
	Target      *big.Int
	TargetFloat *big.Float
	Fuzzy       *big.Int

	pathway []Bit
}

func Distill(phrase Phrase) *Composition {
	c := &Composition{
		Movements: make(map[string][]Passage),
	}
	c.Movements[MovementStart] = make([]Passage, 0)
	c.Movements[MovementPathway] = make([]Passage, 0)
	c.Movements[MovementSeed] = make([]Passage, 0)
	c.pathway = make([]Bit, 0)
	c.Target = phrase.AsBigInt()
	c.TargetFloat = new(big.Float).SetInt(c.Target)
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
