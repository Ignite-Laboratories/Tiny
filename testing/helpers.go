package testing

import (
	"github.com/ignite-laboratories/support/test"
	"github.com/ignite-laboratories/tiny"
	"testing"
)

// ComparePhrases fails the test if the two provided phrases are not equal.
func ComparePhrases(a tiny.Phrase, b tiny.Phrase, t *testing.T) {
	test.CompareSlicesFunc(a, b, CompareMeasurements, t)
}

// CompareUnalignedPhrases fails the test if the two provided phrases are not equal, but aligns them prior to testing.
func CompareUnalignedPhrases(a tiny.Phrase, b tiny.Phrase, t *testing.T) {
	a = a.Align()
	b = b.Align()
	test.CompareSlicesFunc(a, b, CompareMeasurements, t)
}

// CompareMeasurements fails the test if the two provided measurements are not equal.
func CompareMeasurements(a tiny.Measurement, b tiny.Measurement, t *testing.T) {
	test.CompareSlices(a.GetAllBits(), b.GetAllBits(), t)
}
