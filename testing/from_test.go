package testing

import (
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_From_Number_NoWidth(t *testing.T) {
	bits := tiny.From.Number(10)
	expected := tiny.From.Bits(1, 0, 1, 0)
	CompareBitSlices(bits, expected, t)
}

func Test_From_Number_SameWidth(t *testing.T) {
	bits := tiny.From.Number(10, 4)
	expected := tiny.From.Bits(1, 0, 1, 0)
	CompareBitSlices(bits, expected, t)
}

func Test_From_Number_UnderWidth(t *testing.T) {
	bits := tiny.From.Number(10, 3)
	expected := tiny.From.Bits(1, 0, 1)
	CompareBitSlices(bits, expected, t)
}

func Test_From_Number_OverWidth(t *testing.T) {
	bits := tiny.From.Number(10, 5)
	expected := tiny.From.Bits(0, 1, 0, 1, 0)
	CompareBitSlices(bits, expected, t)
}
