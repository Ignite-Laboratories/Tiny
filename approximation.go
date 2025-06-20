package tiny

import "math/big"

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
	BitDepth
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

}
