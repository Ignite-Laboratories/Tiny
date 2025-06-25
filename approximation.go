package tiny

import (
	"math/big"
)

// Approximation represents a synthetic approximation of a target value.
type Approximation struct {
	// Signature is where informational bits are placed during approximation
	Signature Phrase

	// Target is the target value to approximate
	Target Phrase

	// Delta is a signed value representing the delta between the target and value
	Delta *big.Int

	// Upper is the next largest approximation of the target
	Upper *big.Int

	// Lower is the next smaller approximation of the target
	Lower *big.Int

	// BitDepth represents the width of each emitted signature measurement.
	BitDepth int

	// Cycles represents the number of times an approximation was made to yield the current delta.
	// It is incremented each time an approximation is made, as well as a refinement.
	Cycles int
}

// GetClosestValue returns the closest approximate value to the target.
func (a *Approximation) GetClosestValue() *big.Int {
	if a.Delta.Sign() < 0 {
		return a.Upper
	}
	return a.Lower
}

// Refine logarithmically subdivides the range between the lower and upper approximation values,
// emitting a signature measurement of the provided bit depth along the way.
func (a *Approximation) Refine() {
	maxBitDepth := 1 << a.BitDepth
	difference := new(big.Int).Sub(a.Upper, a.Lower)
	height := new(big.Int).Div(difference, big.NewInt(int64(maxBitDepth-1)))

	best := a.Delta
	bestI := -1
	values := make([]*big.Int, maxBitDepth)

	for i := 0; i < maxBitDepth; i++ {
		offset := new(big.Int).Mul(big.NewInt(int64(i)), height)
		value := new(big.Int).Add(a.Lower, offset)
		values[i] = value

		delta := new(big.Int).Sub(a.Target.AsBigInt(), value)
		if delta.CmpAbs(best) < 0 {
			best = delta
			bestI = i
		}
	}

	if bestI < 0 {
		return
	}

	if best.Sign() < 0 {
		a.Upper = values[bestI]
		a.Lower = values[bestI-1]
	} else {
		a.Upper = values[bestI+1]
		a.Lower = values[bestI]
	}

	a.Delta = best
	a.Cycles++
}
