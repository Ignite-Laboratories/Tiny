package testing

import (
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Modify_XORByteWithBits(t *testing.T) {
	pattern := tiny.From.Bits(0, 1)
	xor := tiny.Modify.XORByteWithBits(byte(155), pattern...)
	if xor != 219 {
		t.Errorf("Expected %d, got %d", 206, xor)
	}
}

func Test_Modify_XORBytesWithBits(t *testing.T) {
	pattern := tiny.From.Bits(0, 1)
	bytes := []byte{155, 255, 128, 127, 77}
	expected := []byte{219, 191, 192, 63, 13}
	xor := tiny.Modify.XORBytesWithBits(bytes, pattern...)
	CompareByteSlices(xor, expected, t)
}

func Test_Modify_ToggleBytes(t *testing.T) {
	bytes := []byte{255, 0, 128, 127, 77}
	inverse := []byte{0, 255, 127, 128, 178}

	toggled := tiny.Modify.ToggleBytes(bytes...)
	CompareByteSlices(toggled, inverse, t)
}

func Test_Modify_ToggleBits(t *testing.T) {
	data := tiny.From.Bits(1, 0, 0, 1, 1, 0, 1, 1)
	inverse := tiny.From.Bits(0, 1, 1, 0, 0, 1, 0, 0)

	toggled := tiny.Modify.ToggleBits(data...)
	CompareBitSlices(toggled, inverse, t)
}

func Test_Modify_DropMostSignificantBit(t *testing.T) {
	for _, measure := range tiny.Modify.DropMostSignificantBit(byte(155)) {
		if len(measure.Bits) != 7 {
			t.Errorf("Expected 7 bits, got %d", len(measure.Bits))
		}
		expected := tiny.From.Bits(0, 0, 1, 1, 0, 1, 1)
		CompareBitSlices(measure.Bits, expected, t)
	}
}

func Test_Modify_DropMostSignificantBits(t *testing.T) {
	for dropCount := 0; dropCount <= 8; dropCount++ {
		expected155 := tiny.From.Bits(1, 0, 0, 1, 1, 0, 1, 1)[dropCount:]
		expected255 := tiny.From.Bits(1, 1, 1, 1, 1, 1, 1, 1)[dropCount:]
		expected33 := tiny.From.Bits(0, 0, 1, 0, 0, 0, 0, 1)[dropCount:]
		expected127 := tiny.From.Bits(0, 1, 1, 1, 1, 1, 1, 1)[dropCount:]
		expected0 := tiny.From.Bits(0, 0, 0, 0, 0, 0, 0, 0)[dropCount:]
		expectedBits := [][]tiny.Bit{
			expected155,
			expected255,
			expected33,
			expected127,
			expected0,
		}

		for i, measure := range tiny.Modify.DropMostSignificantBits(dropCount, byte(155), byte(255), byte(33), byte(127), byte(0)) {
			if dropCount == 0 {
				// We didn't drop any bits, treat it as a byte
				expectedByte := tiny.To.Byte(expectedBits[i]...)
				if expectedByte != measure.Bytes[0] {
					t.Errorf("Expected %d, got %d", expectedByte, measure.Bytes[0])
				}
			} else {
				if len(measure.Bytes) > 0 {
					t.Errorf("Expected only bits, got %d bytes + %d bits", len(measure.Bytes), len(measure.Bits))
				}
				CompareBitSlices(measure.Bits, expectedBits[i], t)
			}
		}
	}
}
