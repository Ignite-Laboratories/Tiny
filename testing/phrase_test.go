package testing

import (
	"github.com/ignite-laboratories/core/test"
	"github.com/ignite-laboratories/support"
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Phrase_ToBytesAndBits(t *testing.T) {
	// Test logic:
	//
	// Input:
	//     77 -> 0 1 0 0 1 1 0 1
	//   bits -> 0 1 0
	//   bits -> 1 0 1
	//     22 -> 0 0 0 1 0 1 1 0
	//     33 -> 0 0 1 0 0 0 0 1
	//
	// Output:
	//   0 1 0 0 1 1 0 1 - 0 1 0 1 0 1 0 0 - 0 1 0 1 1 0 0 0 - 1 0 0 0 0 1
	//         77                84                88           remainder

	b1 := tiny.NewMeasurement([]byte{77}, 0, 1, 0)
	b2 := tiny.NewMeasurement([]byte{}, 1, 0, 1)
	b3 := tiny.NewMeasurement([]byte{22, 33})

	phrase := tiny.Phrase{b1, b2, b3}
	bytes, bits := phrase.ToBytesAndBits()

	expectedBytes := []byte{77, 84, 88}
	expectedBits := []tiny.Bit{1, 0, 0, 0, 0, 1}

	test.CompareSlices(bytes, expectedBytes, t)
	test.CompareSlices(bits, expectedBits, t)
}

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

func Test_Phrase_CountBelowThreshold(t *testing.T) {
	threshold := 55

	below := tiny.NewPhrase(support.FixedBytes(32, 33)...)
	belowCount := below.CountBelowThreshold(threshold)
	if belowCount != 32 {
		t.Errorf("Expected 32 below a threshold of %d, Got %d", threshold, belowCount)
	}

	above := tiny.NewPhrase(support.FixedBytes(32, 77)...)
	aboveCount := above.CountBelowThreshold(threshold)
	if aboveCount != 0 {
		t.Errorf("Expected 0 below a threshold of %d, Got %d", threshold, aboveCount)
	}

	random := tiny.NewPhrase(support.RandomBytes(32)...)
	randomCount := random.CountBelowThreshold(threshold)
	var count int
	for _, b := range random {
		if b.Value() < threshold {
			count++
		}
	}

	if randomCount != count {
		t.Errorf("Expected %d below a threshold of %d, Got %d", count, threshold, randomCount)
	}
}

func Test_Phrase_BreakMeasurementsApart(t *testing.T) {
	// The data in this test intentionally increments the left and right regions of each measure by 1 per measure.
	data := tiny.NewPhrase(0, 65, 130, 195)
	l, r := data.BreakMeasurementsApart(2)

	for i := 0; i < 4; i++ {
		if l[i].Value() != i {
			t.Errorf("Expected %d, Got %d", i, l[i].Value())
		}
		if r[i].Value() != i {
			t.Errorf("Expected %d, Got %d", i, r[i].Value())
		}
	}
}

func Test_Phrase_BreakMeasurementsApart_PanicBeyondBounds(t *testing.T) {
	defer test.ShouldPanic(t)

	data := tiny.NewPhrase(0, 65, 130, 195)
	data.BreakMeasurementsApart(9)
}

func Test_Phrase_BreakMeasurementsApart_EmptyData(t *testing.T) {
	data := tiny.NewPhrase()
	data.BreakMeasurementsApart(2)
}

func Test_Phrase_RecombineMeasures(t *testing.T) {
	expected := []byte{0, 65, 130, 195}

	data := tiny.NewPhrase(expected...)
	l, r := data.BreakMeasurementsApart(2)

	recombined := tiny.RecombinePhrases(l, r)

	for i, m := range recombined {
		if m.Value() != int(expected[i]) {
			t.Errorf("Expected %d, Got %d", i, m.Value())
		}
	}
}
