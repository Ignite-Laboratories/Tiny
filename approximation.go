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
}

// CalculateBitDrop calculates the difference between the target bit length and the
// approximation bit length using the following formula:
//
//	target - (signature + delta)
func (a *Approximation) CalculateBitDrop() int {
	length := a.Target.BitLength()
	dl := a.Delta.Text(2)
	sig := a.Signature.BitLength()
	return length - (sig + len(dl))
}

// Refine finds the closest synthetic value to the delta and applies it to the approximation
// using the following encoding scheme:
//
//  Sign ⬎    ⬐ Dark ZLE  ⬐ Light ZLE
//      ⁰⁄₁ [ ⁰⁄₁ ... ] [ ⁰⁄₁ ... ] ╮
//       ↑                          │
//       ╰───────────LOOP───────────╯
func (a *Approximation) Refine(retain ...int) (stride int) {

}
