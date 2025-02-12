package testing

import (
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Modify_XORByteWithByte(t *testing.T) {
	pattern := tiny.From.Bits(0, 1)
	bytes := []byte{155, 255, 128, 127, 77}
	expected := []byte{219, 191, 192, 63, 13}
	xor := tiny.Modify.XORBytesWithBits(bytes, pattern...)

	for i, b := range xor {
		if b != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], b)
		}
	}
}

func Test_Modify_XORByteWithBits(t *testing.T) {
	pattern := tiny.From.Bits(0, 1)
	xor := tiny.Modify.XORByteWithBits(byte(155), pattern...)
	if xor != 219 {
		t.Errorf("Expected %d, got %d", 206, xor)
	}
}

func Test_Modify_ToggleBytes(t *testing.T) {
	bytes := []byte{255, 0, 128, 127, 77}
	inverse := []byte{0, 255, 127, 128, 178}

	toggled := tiny.Modify.ToggleBytes(bytes...)
	for i, b := range toggled {
		if b != inverse[i] {
			t.Errorf("Expected %d, got %d", inverse[i], b)
		}
	}
}

func Test_Modify_ToggleBits(t *testing.T) {
	data := tiny.From.Bits(1, 0, 0, 1, 1, 0, 1, 1)
	inverse := tiny.From.Bits(0, 1, 1, 0, 0, 1, 0, 0)

	toggled := tiny.Modify.ToggleBits(data...)
	for i, bit := range toggled {
		if bit != inverse[i] {
			t.Errorf("Expected %d, got %d", inverse[i], bit)
		}
	}
}

func Test_Modify_DropMostSignificantBit(t *testing.T) {
	remainder := tiny.Modify.DropMostSignificantBit(byte(155))
	if len(remainder.Bits) != 7 {
		t.Errorf("Expected 7 bits, got %d", len(remainder.Bits))
	}
	expected := tiny.From.Bits(0, 0, 1, 1, 0, 1, 1)
	for i, bit := range expected {
		if remainder.Bits[i] != bit {
			t.Errorf("Expected %d, got %d", bit, remainder.Bits[i])
		}
	}
}

func Test_Modify_DropMostSignificantBits(t *testing.T) {
	for count := 1; count < 8; count++ {
		expected155 := tiny.From.Bits(1, 0, 0, 1, 1, 0, 1, 1)[count:]
		expected255 := tiny.From.Bits(1, 1, 1, 1, 1, 1, 1, 1)[count:]
		expected33 := tiny.From.Bits(0, 0, 1, 0, 0, 0, 0, 1)[count:]
		expected127 := tiny.From.Bits(0, 1, 1, 1, 1, 1, 1, 1)[count:]
		expected0 := tiny.From.Bits(0, 0, 0, 0, 0, 0, 0, 0)[count:]
		expectedBits := append(expected155, expected255...)
		expectedBits = append(expectedBits, expected33...)
		expectedBits = append(expectedBits, expected127...)
		expectedBits = append(expectedBits, expected0...)

		remainder := tiny.Modify.DropMostSignificantBits(count, byte(155), byte(255), byte(33), byte(127), byte(0))
		count := 0
		for _, b := range remainder.Bytes {
			bits := tiny.From.Byte(b)
			for _, bit := range bits {
				if bit != expectedBits[count] {
					t.Errorf("Expected %d, got %d", expectedBits[count], bit)
				}
				count++
			}
		}
		for i := 0; i < len(remainder.Bits); i++ {
			if remainder.Bits[i] != expectedBits[count] {
				t.Errorf("Expected 0, got %d", remainder.Bits[i])
			}

			count++
		}
	}
}
