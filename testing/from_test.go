package testing

import (
	"github.com/ignite-laboratories/support/test"
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_From_Byte(t *testing.T) {
	bits170 := []tiny.Bit{1, 0, 1, 0, 1, 0, 1, 0}
	bits85 := []tiny.Bit{0, 1, 0, 1, 0, 1, 0, 1}
	expected1 := append(bits170, bits85...)
	bits1 := tiny.From.Bytes([]byte{170, 85}...)
	test.CompareSlices(bits1, expected1, t)

	bits2 := tiny.From.Bytes([]byte{170}...)
	test.CompareSlices(bits2, bits170, t)

	bits3 := tiny.From.Bytes([]byte{}...)
	test.CompareSlices(bits3, []tiny.Bit{}, t)
}

func Test_From_Number_NoWidth(t *testing.T) {
	bits := tiny.From.Number(10)
	expected := tiny.From.Bits(1, 0, 1, 0)
	test.CompareSlices(bits, expected, t)
}

func Test_From_Number_SameWidth(t *testing.T) {
	bits := tiny.From.Number(10, 4)
	expected := tiny.From.Bits(1, 0, 1, 0)
	test.CompareSlices(bits, expected, t)
}

func Test_From_Number_UnderWidth(t *testing.T) {
	bits := tiny.From.Number(10, 3)
	expected := tiny.From.Bits(1, 0, 1)
	test.CompareSlices(bits, expected, t)
}

func Test_From_Number_OverWidth(t *testing.T) {
	bits := tiny.From.Number(10, 5)
	expected := tiny.From.Bits(0, 1, 0, 1, 0)
	test.CompareSlices(bits, expected, t)
}
