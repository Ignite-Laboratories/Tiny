package tiny

import "math/big"

/**
This is a darling of mine - It has absolutely no functional use in this library,
but it demonstrates how one can "modulate" literally any repeating pattern.

The concepts of binary synthesis echo across so many different fields because all
abstractly rely upon the same fundamental principles of logic and reason that drive
this very library's existence.

--Alex Petz
*/

type _modulate struct{}

// ModulationFunc is a type of function that is called for each instance of a binary
// pattern from an Approximation.
//
// Parameters:
//   - i: The current pattern's index.
//   - l: The total number of patterns in the source approximation.
//   - m: The current pattern.
//
// The returned measurement replaces the current pattern in the source approximation.
type ModulationFunc func(i int, l int, m Measurement) Measurement

// Approximation calls the provided modulation function for every pattern of the source approximation
// and returns the resulting pattern information.
func (m _modulate) Approximation(a Approximation, fn ModulationFunc) Approximation {
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

// Inject(count int) - goes high for one instance and then skips the provided count

// Halve(times int) - goes high when the appropriate half is reached

// Layer(skip int) - goes high after the provided count

// Sine(amplitude int, period int, start Measurement) - Generates a sine wave from the starting measurement

// Cosine(amplitude int, period int, start Measurement) - Generates a cosine wave from the starting measurement

// Tangent(amplitude int, period int, start Measurement) - Generates a tangent wave from the starting measurement

// Toggle returns either the source (low) or the provided bit pattern (high) for the provided number of intervals.
func (m _modulate) Toggle(width int, startHigh bool, pattern ...Bit) ModulationFunc {
	// NOTE: There's no difference between '0' and '1' width, so treat anything less than '1' as '1'
	if width <= 0 {
		width = 1
	}

	high := startHigh
	w := 0
	return func(i int, l int, m Measurement) Measurement {
		if w >= width {
			w = 0
			high = !high
		}
		w++

		if high {
			return NewMeasurementFromBits(pattern[:m.BitLength()]...)
		} else {
			return m
		}
	}
}
