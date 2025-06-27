package tiny

import (
	"math/big"
)

// Approximation represents a synthetic approximation of a target value.
type Approximation struct {
	// Signature is where informational bits are placed during approximation
	Signature Phrase

	// Target is the target value to approximate.
	Target Phrase

	// Value is the current approximate value.
	Value *big.Int

	// Delta is a signed value representing the delta between the target and the closest boundary.
	Delta *big.Int

	// Height is the difference between the upper and lower boundaries.
	Height *big.Int
}

func (a *Approximation) CalculateBitDrop() int {
	length := a.Target.BitLength()
	dl := a.Delta.Text(2)
	sig := a.Signature.BitLength()
	return length - (sig + len(dl))
}

// Refine takes the current approximate value, divides its associated height in half, and then
// moves the approximation by that amount towards the target before calculating the new delta.
func (a *Approximation) Refine() {
	down := a.Delta.Sign() < 0
	a.Height = new(big.Int).Div(a.Height, big.NewInt(2))

	if down {
		a.Value = new(big.Int).Sub(a.Value, a.Height)
		a.Signature = a.Signature.AppendBits(1)
	} else {
		a.Value = new(big.Int).Add(a.Value, a.Height)
		a.Signature = a.Signature.AppendBits(0)
	}

	a.Delta = new(big.Int).Sub(a.Target.AsBigInt(), a.Value)
}
