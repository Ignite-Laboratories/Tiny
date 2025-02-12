package testing

import (
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func CompareBitSlices(a []tiny.Bit, b []tiny.Bit, t *testing.T) {
	CompareByteSlices(tiny.Upcast(a), tiny.Upcast(b), t)
}

func CompareByteSlices(a []byte, b []byte, t *testing.T) {
	if len(a) != len(b) {
		t.Errorf("Slices are not the same length")
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			t.Errorf("Expected %d at [%d], got %d", a[i], i, b[i])
		}
	}
}
