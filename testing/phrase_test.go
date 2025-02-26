package testing

import (
	"github.com/ignite-laboratories/support"
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Phrase_BitLength(t *testing.T) {
	phrase := tiny.NewPhrase(support.RandomBytes(32)...)
	length := phrase.BitLength()
	if length != 32*8 {
		t.Errorf("Expected %d, Got %d", 32*8, length)
	}
}

func Test_Phrase_AllBelowThreshold(t *testing.T) {
	below := tiny.NewPhrase(support.FixedBytes(32, 55)...)
	if !below.AllBelowThreshold(55) {
		t.Errorf("Input data was below threshold, but AllBelowThreshold returned false")
	}

	above := tiny.NewPhrase(support.FixedBytes(32, 77)...)
	if above.AllBelowThreshold(55) {
		t.Errorf("Input data was above threshold, but AllBelowThreshold returned true")
	}

	random := tiny.NewPhrase(support.RandomBytes(32)...)
	random[7] = tiny.NewMeasurement([]byte{77}) // ensure at least one is above threshold
	if above.AllBelowThreshold(55) {
		t.Errorf("Input data was above threshold, but AllBelowThreshold returned true")
	}
}
