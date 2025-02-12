package testing

import (
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Measure_Toggle(t *testing.T) {
	bytes := []byte{255, 0, 128, 127, 77}
	inverseBytes := []byte{0, 255, 127, 128, 178}
	bits := tiny.From.Bits(0, 1, 0, 1)
	inverseBits := tiny.From.Bits(1, 0, 1, 0)
	m := tiny.NewMeasure(bytes, bits...)
	m.Toggle()

	for i, b := range m.Bytes {
		if b != inverseBytes[i] {
			t.Errorf("Expected %d, got %d", inverseBytes[i], b)
		}
	}
	for i, b := range m.Bits {
		if b != inverseBits[i] {
			t.Errorf("Expected %d, got %d", inverseBits[i], b)
		}
	}
}
