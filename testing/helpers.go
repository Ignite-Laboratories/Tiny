package testing

import (
	"github.com/ignite-laboratories/tiny"
	"testing"
)

// ComparePhrases fails the test if the two provided phrases are not equal.
func ComparePhrases(a tiny.Phrase, b tiny.Phrase, t *testing.T) {
	CompareSlicesFunc(a, b, CompareMeasurements, t)
}

// CompareUnalignedPhrases fails the test if the two provided phrases are not equal, but aligns them prior to testing.
func CompareUnalignedPhrases(a tiny.Phrase, b tiny.Phrase, t *testing.T) {
	a = a.Align()
	b = b.Align()
	CompareSlicesFunc(a, b, CompareMeasurements, t)
}

// CompareMeasurements fails the test if the two provided measurements are not equal.
func CompareMeasurements(a tiny.Measurement, b tiny.Measurement, t *testing.T) {
	CompareSlices(a.GetAllBits(), b.GetAllBits(), t)
}

// ShouldPanic fails the test if the test did not panic.
// It should be called at the start of your test with:
//
//	defer test.ShouldPanic(t)
func ShouldPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("Expected panic, but didn't get one")
		t.FailNow()
	}
}

// CompareValues fails the test if the two provided values are not equal.
func CompareValues[T comparable](a T, b T, t *testing.T) {
	if a != b {
		t.Errorf("Expected %v, got %v", a, b)
		t.FailNow()
	}
}

// CompareSlices fails the test if the two slices are unequal in length, or if the
// elements are not equal for every index.
func CompareSlices[T comparable](a []T, b []T, t *testing.T) {
	if len(a) != len(b) {
		t.Errorf("Slices are not the same length")
		t.FailNow()
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			t.Errorf("Expected %v at [%d], got %v", a[i], i, b[i])
			t.FailNow()
		}
	}
}

// CompareSlicesFunc fails the test if the two slices are unequal in length, or if the
// elements are not equal for every index.
//
// NOTE: This differs from CompareSlices in that it allows you to provide a custom comparison function.
func CompareSlicesFunc[T any](a []T, b []T, compare func(T, T, *testing.T), t *testing.T) {
	if len(a) != len(b) {
		t.Errorf("Slices are not the same length")
		t.FailNow()
	}
	for i := 0; i < len(a); i++ {
		compare(a[i], b[i], t)
	}
}
