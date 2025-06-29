package testing

import (
	"github.com/ignite-laboratories/support/test"
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Modulate_Toggle(t *testing.T) {
	depth := 3
	d := tiny.Synthesize.Ones(128)

	tester := func(width int, startHigh bool) {
		if width <= 0 {
			width = 1
		}
		a := tiny.Synthesize.Approximation(d, depth)
		one := tiny.From.Number((1<<depth)-1, depth)
		zero := tiny.From.Number(0, depth)
		a = tiny.Modulate.Approximation(a, tiny.Modulate.Toggle(width, startHigh, zero...))

		ii := 0
		high := startHigh
		for remainder := a.Value; remainder.BitLength() > 0; {
			if ii >= width {
				ii = 0
				high = !high
			}

			var current tiny.Measurement
			current, remainder = remainder.ReadMeasurement(depth)

			bits := current.GetAllBits()
			if high {
				test.CompareSlices(bits, zero[:len(bits)], t)
			} else {
				test.CompareSlices(bits, one[:len(bits)], t)
			}

			ii++
		}
	}

	for i := 0; i < 32; i++ {
		tester(i, false)
		tester(i, true)
	}
}
