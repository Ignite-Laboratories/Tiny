package testing

import (
	"github.com/ignite-laboratories/tiny"
	"testing"
)

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
