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
func (a Approximation) CalculateBitDrop() int {
	length := a.Target.BitLength()
	dl := a.Delta.Text(2)
	sig := a.Signature.BitLength()
	return length - (sig + len(dl))
}

// Modulate calls the provided modulation function for every pattern of the source approximation
// and returns the resulting pattern information.
func (a Approximation) Modulate(fn ModulationFunc) Approximation {
	bitLen := a.Value.BitLength()
	l := bitLen / a.BitDepth
	if bitLen%a.BitDepth > 0 {
		l += 1
	}

	result := make([]Bit, 0, a.Value.BitLength())
	i := 0
	for remainder := a.Value; remainder.BitLength() > 0; {
		var current Measurement
		current, remainder = remainder.ReadMeasurement(a.BitDepth)
		replacement := fn(i, l, current)
		result = append(result, replacement.GetAllBits()...)
		i++
	}

	a.Value = NewPhraseFromBits(result...)
	a.Delta = new(big.Int).Sub(a.TargetBigInt, a.Value.AsBigInt())

	return a
}
