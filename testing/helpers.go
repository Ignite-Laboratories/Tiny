package testing

import (
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func ShouldPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("Expected panic, but didn't get one")
		t.FailNow()
	}
}

func CompareValues[T comparable](a T, b T, t *testing.T) {
	if a != b {
		t.Errorf("Expected %v, got %v", a, b)
		t.FailNow()
	}
}

func CompareBitSlices(a []tiny.Bit, b []tiny.Bit, t *testing.T) {
	CompareByteSlices(tiny.Upcast(a), tiny.Upcast(b), t)
}

func CompareByteSlices(a []byte, b []byte, t *testing.T) {
	if len(a) != len(b) {
		t.Errorf("Slices are not the same length")
		t.FailNow()
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			t.Errorf("Expected %d at [%d], got %d", a[i], i, b[i])
			t.FailNow()
		}
	}
}
