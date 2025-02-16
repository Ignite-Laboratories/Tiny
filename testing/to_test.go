package testing

import (
	"ignitelabs.net/git/tiny"
	"testing"
)

func Test_To_Number_TooWide(t *testing.T) {
	number := tiny.To.Number(222, 1, 0, 1, 0)
	if number != 10 {
		t.Errorf("Expected %d, Got %d", 10, number)
	}
}

func Test_To_Number_SameWidth(t *testing.T) {
	number := tiny.To.Number(4, 1, 0, 1, 0)
	if number != 10 {
		t.Errorf("Expected %d, Got %d", 5, number)
	}
}

func Test_To_Number_UnderWide(t *testing.T) {
	number := tiny.To.Number(3, 1, 0, 1, 0)
	if number != 5 {
		t.Errorf("Expected %d, Got %d", 5, number)
	}
}

func Test_To_Number_LargeNumber(t *testing.T) {
	number := tiny.To.Number(32, 0, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 1, 1, 1, 0, 1, 0, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1, 1)
	if number != 1521080171 {
		t.Errorf("Expected %d, Got %d", 5, 1521080171)
	}
}

func Test_To_String(t *testing.T) {
	random := tiny.Synthesize.Random(22)
	bits := random.GetAllBits()
	str := tiny.To.String(bits...)
	for i := 0; i < len(str); i++ {
		stringValue := str[i] - '0'
		if stringValue != uint8(bits[i]) {
			t.Errorf("Expected %d, Got %d", bits[i], stringValue)
		}
	}
}

func Test_To_Measure(t *testing.T) {
	bytes := []byte{255, 77, 0}
	byteBits := tiny.From.Bytes(bytes...)
	bits := tiny.From.Bits(0, 1, 1)
	combined := append(byteBits, bits...)

	measure := tiny.To.Measure(combined...)
	CompareByteSlices(measure.Bytes, bytes, t)
	CompareBitSlices(measure.Bits, bits, t)
}
