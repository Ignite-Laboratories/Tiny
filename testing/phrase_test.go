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

/**
BreakMeasurementsApart and RecombineMeasurements
*/

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

	recombined := tiny.RecombineMeasurements(l, r)

	for i, m := range recombined {
		if m.Value() != int(expected[i]) {
			t.Errorf("Expected %d, Got %d", i, m.Value())
		}
	}
}

/**
Align
*/

func Test_Phrase_AlignOnByteWidth(t *testing.T) {
	// Test logic:
	//
	// Input:
	//   bits -> 0 1
	//     77 -> 0 1 0 0 1 1 0 1
	//     22 -> 0 0 0 1 0 1 1 0
	//     33 -> 0 0 1 0 0 0 0 1
	//
	// Output:
	//      |       77        |       22        |       33        |  <- "Unaligned"
	//   0 1 0 1 0 0 1 1 - 0 1 0 0 0 1 0 1 - 1 0 0 0 1 0 0 0 - 0 1   <- Raw Bits
	//  |      83        |       69        |      136        |       <- "Aligned"

	// Build the phrase
	phrase := append(tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 1)}, tiny.NewPhrase(77, 22, 33)...)

	// Align it
	aligned := phrase.Align()

	// Test the result
	expected := tiny.NewPhrase(83, 69, 136)
	expected = append(expected, tiny.NewMeasurement([]byte{}, 0, 1))
	test.ComparePhrases(aligned, expected, t)
}

func Test_Phrase_AlignOnSmallerWidth(t *testing.T) {
	// Test logic:
	//
	// Input:
	//   bits -> 0 1
	//     77 -> 0 1 0 0 1 1 0 1
	//     22 -> 0 0 0 1 0 1 1 0
	//     33 -> 0 0 1 0 0 0 0 1
	//
	// Output:
	//      |       77        |       22        |       33        |  <- "Unaligned"
	//   0 1 0 1 0 0 1 1 - 0 1 0 0 0 1 0 1 - 1 0 0 0 1 0 0 0 - 0 1   <- Raw Bits
	//  |   5   |   3    |    4   |   5    |    8   |   8    |  1    <- "Aligned"

	// Build the phrase
	phrase := append(tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 1)}, tiny.NewPhrase(77, 22, 33)...)

	// Align it
	aligned := phrase.Align(4)

	// Test the result
	one := tiny.NewMeasurement([]byte{}, 0, 1)
	three := tiny.NewMeasurement([]byte{}, 0, 0, 1, 1)
	four := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0)
	five := tiny.NewMeasurement([]byte{}, 0, 1, 0, 1)
	eight := tiny.NewMeasurement([]byte{}, 1, 0, 0, 0)
	expected := tiny.Phrase{five, three, four, five, eight, eight, one}
	test.ComparePhrases(aligned, expected, t)
}

func Test_Phrase_AlignOnLargerWidth(t *testing.T) {
	// Test logic:
	//
	// Input:
	//   bits -> 0 1
	//     77 -> 0 1 0 0 1 1 0 1
	//     22 -> 0 0 0 1 0 1 1 0
	//     33 -> 0 0 1 0 0 0 0 1
	//
	// Output:
	//      |       77        |       22        |       33        |  <- "Unaligned"
	//   0 1 0 1 0 0 1 1 - 0 1 0 0 0 1 0 1 - 1 0 0 0 1 0 0 0 - 0 1   <- Raw Bits
	//  |        333          |        88           |     33         <- "Aligned"

	phrase := append(tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 1)}, tiny.NewPhrase(77, 22, 33)...)
	aligned := phrase.Align(10)

	m1 := tiny.NewMeasurement([]byte{}, 0, 1, 0, 1, 0, 0, 1, 1, 0, 1)
	m2 := tiny.NewMeasurement([]byte{}, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0)
	m3 := tiny.NewMeasurement([]byte{}, 1, 0, 0, 0, 0, 1)
	expected := tiny.Phrase{m1, m2, m3}
	test.ComparePhrases(aligned, expected, t)
}

func Test_Phrase_Align_Simple(t *testing.T) {
	m1 := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0)
	m2 := tiny.NewMeasurement([]byte{}, 1, 1, 0, 1)
	phrase := tiny.Phrase{m1, m2}

	aligned := phrase.Align()
	expected := tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 1, 0, 0, 1, 1, 0, 1)}
	test.ComparePhrases(aligned, expected, t)
}

func Test_Phrase_Align_PanicIfZeroWidth(t *testing.T) {
	defer test.ShouldPanic(t)
	phrase := tiny.NewPhrase(77, 22, 33)
	phrase.Align(0)
}

func Test_Phrase_Align_PanicIfNegativeWidth(t *testing.T) {
	defer test.ShouldPanic(t)
	phrase := tiny.NewPhrase(77, 22, 33)
	phrase.Align(-1)
}

func Test_Phrase_Align_PanicIfWidthTooLarge(t *testing.T) {
	defer test.ShouldPanic(t)
	phrase := tiny.NewPhrase(77, 22, 33)
	phrase.Align(33)
}

/**
Read
*/

func Test_Phrase_Read(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder := phrase.Read(4)

	left := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0)
	test.ComparePhrases(read, tiny.Phrase{left}, t)

	right := tiny.NewMeasurement([]byte{}, 1, 1, 0, 1)
	test.ComparePhrases(remainder, tiny.Phrase{right}, t)
}

func Test_Phrase_Read_Zero(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder := phrase.Read(0)

	test.ComparePhrases(read, tiny.Phrase{}, t)
	test.ComparePhrases(remainder, tiny.Phrase{tiny.NewMeasurement([]byte{77})}, t)
}

func Test_Phrase_Read_Negative(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder := phrase.Read(-5)

	test.ComparePhrases(read, tiny.Phrase{}, t)
	test.ComparePhrases(remainder, tiny.Phrase{tiny.NewMeasurement([]byte{77})}, t)
}

func Test_Phrase_Read_AcrossMeasurements(t *testing.T) {
	phrase := tiny.NewPhrase(77, 22)

	read, remainder := phrase.Read(10)

	left1 := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0, 1, 1, 0, 1)
	left2 := tiny.NewMeasurement([]byte{}, 0, 0)
	test.ComparePhrases(read, tiny.Phrase{left1, left2}, t)

	right := tiny.NewMeasurement([]byte{}, 0, 1, 0, 1, 1, 0)
	test.ComparePhrases(remainder, tiny.Phrase{right}, t)
}

/**
ReadMeasurement
*/

func Test_Phrase_ReadMeasurement(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder := phrase.ReadMeasurement(4)

	left := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0)
	test.CompareMeasurements(read, left, t)

	right := tiny.NewMeasurement([]byte{}, 1, 1, 0, 1)
	test.ComparePhrases(remainder, tiny.Phrase{right}, t)
}

func Test_Phrase_ReadMeasurement_Zero(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder := phrase.ReadMeasurement(0)

	test.CompareMeasurements(read, tiny.NewMeasurement([]byte{}), t)
	test.ComparePhrases(remainder, tiny.Phrase{tiny.NewMeasurement([]byte{77})}, t)
}

func Test_Phrase_ReadMeasurement_Negative(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder := phrase.ReadMeasurement(-5)

	test.CompareMeasurements(read, tiny.NewMeasurement([]byte{}), t)
	test.ComparePhrases(remainder, tiny.Phrase{tiny.NewMeasurement([]byte{77})}, t)
}

func Test_Phrase_ReadMeasurement_OverByte(t *testing.T) {
	phrase := tiny.NewPhrase(77, 22, 33)

	read, remainder := phrase.ReadMeasurement(10)

	left := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0)
	test.CompareMeasurements(read, left, t)

	right1 := tiny.NewMeasurement([]byte{}, 0, 1, 0, 1, 1, 0)
	right2 := tiny.NewMeasurement([]byte{33})
	test.ComparePhrases(remainder, tiny.Phrase{right1, right2}, t)
}

func Test_Phrase_ReadMeasurement_ShouldPanicIfOver32(t *testing.T) {
	defer test.ShouldPanic(t)
	tiny.NewPhrase().ReadMeasurement(33)
}

/**
Trifurcate
*/

func Test_Phrase_Trifurcate(t *testing.T) {
	// Test logic:
	//
	// Input:
	//     77 -> 0 1 0 0 1 1 0 1
	//     22 -> 0 0 0 1 0 1 1 0
	//     33 -> 0 0 1 0 0 0 0 1
	//
	// Output:
	//   0 1 0 0 1 1 0 1 - 0 0 0 1 0 1 1 0 - 0 0 1 0 0 0 0 1  <- Raw Bits
	//  |     Start      |     Middle      |      End       | <- "Trifurcated"
	phrase := tiny.NewPhrase(77, 22, 33)

	s, m, e := phrase.Trifurcate(8, 8)
	test.ComparePhrases(s, tiny.NewPhrase(77), t)
	test.ComparePhrases(m, tiny.NewPhrase(22), t)
	test.ComparePhrases(e, tiny.NewPhrase(33), t)
}

func Test_Phrase_Trifurcate_OddSize(t *testing.T) {
	// Test logic:
	//
	// Input:
	//     77 -> 0 1 0 0 1 1 0 1
	//     22 -> 0 0 0 1 0 1 1 0
	//     33 -> 0 0 1 0 0 0 0 1
	//
	// Output:
	//   0 1 0 0 - 1 1 0 1 - 0 0 0 1 0 1 1 0 - 0 0 1 0 - 0 0 0 1  <- Raw Bits
	//  | Start  | Middle1 |     Middle2     | Middle3 |   End  | <- "Trifurcated"
	phrase := tiny.NewPhrase(77, 22, 33)

	s, m, e := phrase.Trifurcate(4, 16)

	eStart := tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 1, 0, 0)}

	eMiddle1 := tiny.NewMeasurement([]byte{}, 1, 1, 0, 1)
	eMiddle2 := tiny.NewMeasurement([]byte{}, 0, 0, 0, 1, 0, 1, 1, 0)
	eMiddle3 := tiny.NewMeasurement([]byte{}, 0, 0, 1, 0)
	eMiddle := tiny.Phrase{eMiddle1, eMiddle2, eMiddle3}

	eEnd := tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 0, 0, 1)}

	test.ComparePhrases(s, eStart, t)
	test.ComparePhrases(m, eMiddle, t)
	test.ComparePhrases(e, eEnd, t)
}

func Test_Phrase_Trifurcate_ExcessiveMiddleLength(t *testing.T) {
	// Test logic:
	//
	// Input:
	//     77 -> 0 1 0 0 1 1 0 1
	//
	// Output:
	//   0 1 - 0 0 1 1 0 1      <- Raw Bits
	//  | S  |   Middle   | E | <- "Trifurcated"
	phrase := tiny.NewPhrase(77)

	s, m, e := phrase.Trifurcate(2, 8)

	eStart := tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 1)}
	eMiddle := tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 0, 1, 1, 0, 1)}
	eEnd := tiny.Phrase{}

	test.ComparePhrases(s, eStart, t)
	test.ComparePhrases(m, eMiddle, t)
	test.ComparePhrases(e, eEnd, t)
}

func Test_Phrase_Trifurcate_ExcessiveStartLength(t *testing.T) {
	// Test logic:
	//
	// Input:
	//     77 -> 0 1 0 0 1 1 0 1
	//
	// Output:
	//   0 1 0 0 1 1 0 1          <- Raw Bits
	//  |     Start     | M | E | <- "Trifurcated"
	phrase := tiny.NewPhrase(77)

	s, m, e := phrase.Trifurcate(10, 8)

	eStart := tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 1, 0, 0, 1, 1, 0, 1)}
	eMiddle := tiny.Phrase{}
	eEnd := tiny.Phrase{}

	test.ComparePhrases(s, eStart, t)
	test.ComparePhrases(m, eMiddle, t)
	test.ComparePhrases(e, eEnd, t)
}

func Test_Phrase_Trifurcate_ZeroStartLength(t *testing.T) {
	// Test logic:
	//
	// Input:
	//     77 -> 0 1 0 0 1 1 0 1
	//
	// Output:
	//       0 1 0 0 - 1 1 0 1  <- Raw Bits
	//  | S | Middle |  End   | <- "Trifurcated"
	phrase := tiny.NewPhrase(77)

	s, m, e := phrase.Trifurcate(0, 4)

	eStart := tiny.Phrase{}
	eMiddle := tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 1, 0, 0)}
	eEnd := tiny.Phrase{tiny.NewMeasurement([]byte{}, 1, 1, 0, 1)}

	test.ComparePhrases(s, eStart, t)
	test.ComparePhrases(m, eMiddle, t)
	test.ComparePhrases(e, eEnd, t)
}

func Test_Phrase_Trifurcate_ZeroStartLengthAndNoEnd(t *testing.T) {
	// Test logic:
	//
	// Input:
	//     77 -> 0 1 0 0 1 1 0 1
	//
	// Output:
	//       0 1 0 0 1 1 0 1      <- Raw Bits
	//  | S |    Middle     | E | <- "Trifurcated"
	phrase := tiny.NewPhrase(77)

	s, m, e := phrase.Trifurcate(0, 10)

	eStart := tiny.Phrase{}
	eMiddle := tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 1, 0, 0, 1, 1, 0, 1)}
	eEnd := tiny.Phrase{}

	test.ComparePhrases(s, eStart, t)
	test.ComparePhrases(m, eMiddle, t)
	test.ComparePhrases(e, eEnd, t)
}
