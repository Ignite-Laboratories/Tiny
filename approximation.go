package tiny

import (
	"math/big"
)

// Approximation represents a synthetic approximation of a target value.
type Approximation struct {
	// Signature is where informational bits are placed during approximation.
	Signature Phrase

	// Target is the target value to approximate.
	Target Phrase

	// TargetBigInt is the target value as a big.Int.
	//
	// NOTE: This is purely an efficiency gain to avoid reconverting a static value in loop logic.
	TargetBigInt *big.Int

	// Value is the current approximate value.
	Value Phrase

	// Delta is the difference between Target and Value.
	Delta *big.Int

	// BitDepth is the width of each pattern.
	BitDepth int

	// IndexWidth is the initial bit width of the data index that holds the target value.
	IndexWidth int
}

// CalculateBitDrop calculates the difference between the target bit length and the
// approximation bit length using the following formula:
//
//	target - (signature + delta)
func (a *Approximation) CalculateBitDrop() int {
	tLength := a.Target.BitLength()
	dLength := a.Delta.Text(2)
	sLength := a.Signature.BitLength()
	return tLength - (sLength + len(dLength))
}

// Refine finds the closest synthetic value to the delta and applies it to the approximation
// using the following encoding scheme and should be called iteratively until enough bits
// are gained to 'bailout'.
//
//	 Dark ⬎    Light ⬎    Sign ⬎
//		╭ [ ⁰⁄₁ ... ] [ ⁰⁄₁ ... ] ⁰⁄₁ ╮
//		╰────────←───LOOP────←────────╯
//
// The three major parts of this scheme are the 'cursor' position, the 'bailout' condition,
// and the synthesis of a value relative to the initial index size.
//
// The starting cursor position is optionally provided as an input parameter and indicates
// how far in from the left side of the target index to start calculating from.  This represents
// a point from which ones can be synthesized until an ideal dark value is found.  At each
// step an appropriate light value is found, which indicates how many zeros to synthesize
// on the right side of the dark value.  Once the best refinement is found, the approximation
// is updated and the dark stride is returned to be fed into the next iteration's position value.
//
// The bailout condition is indicated when the dark ZLE value is '0' and should be added to the
// signature by the -calling- function.
//
// For example, let's synthesize a few values:
//
//		Position: 2
//		    Dark: 3
//		   Light: 1
//
//	   2   |   3   |         | 1  <- Inputs
//		0 1 0 | 1 0 1 | 0 1 0 1 | 0  <- Approximation
//		      |       | 1 1 1 1 | 0  <- Synthesized value
//
//		Position: 0
//		    Dark: 4
//		   Light: 5
//
//	  0  |    4    |     |     5     <- Inputs
//		    | 0 1 0 1 | 0 1 | 0 1 0 1 0 <- Approximation
//		    |         | 1 1 | 0 0 0 0 0 <- Synthesized value
//
//		Position: 1
//		    Dark: 4
//		   Light: 0
//
//	  1  |    4    |               | 0 <- Inputs
//		 0  | 1 0 1 0 | 1 0 1 0 1 0 1 |   <- Approximation
//		    |         | 1 1 1 1 1 1 1 |   <- Synthesized value
func (a *Approximation) Refine(position ...int) (stride int) {
	sign := Zero
	if a.Delta.Sign() < 0 {
		sign = One
	}
	a.Signature = a.Signature.AppendBits(sign)

	p := 0
	if len(position) > 0 {
		p = position[0]
		if p < 0 {
			p = 0
		}
	}

	offset := 0
	for i := a.IndexWidth - 1; i >= 0; i-- {
		//ones := Synthesize.Ones(a.IndexWidth - p - offset)

		offset++
	}
	return -1
}
