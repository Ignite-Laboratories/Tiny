package testing

import (
	"github.com/ignite-laboratories/support"
	"github.com/ignite-laboratories/support/test"
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Phrase_NewPhraseFromBits(t *testing.T) {
	bits := []tiny.Bit{0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	p := tiny.NewPhraseFromBits(bits...)
	r, _ := p.Read(p.BitLength())
	test.CompareSlices(bits, r.Bits(), t)
}

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
	ComparePhrases(aligned, expected, t)
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
	ComparePhrases(aligned, expected, t)
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
	ComparePhrases(aligned, expected, t)
}

func Test_Phrase_Align_Simple(t *testing.T) {
	m1 := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0)
	m2 := tiny.NewMeasurement([]byte{}, 1, 1, 0, 1)
	phrase := tiny.Phrase{m1, m2}

	aligned := phrase.Align()
	expected := tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 1, 0, 0, 1, 1, 0, 1)}
	ComparePhrases(aligned, expected, t)
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
	phrase := tiny.NewPhrase(11, 33, 55, 99, 77, 22, 33)
	phrase.Align(tiny.GetArchitectureBitWidth() + 1)
}

/**
Read
*/

func Test_Phrase_Read(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder := phrase.Read(4)

	left := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0)
	ComparePhrases(read, tiny.Phrase{left}, t)

	right := tiny.NewMeasurement([]byte{}, 1, 1, 0, 1)
	ComparePhrases(remainder, tiny.Phrase{right}, t)
}

func Test_Phrase_Read_Zero(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder := phrase.Read(0)

	ComparePhrases(read, tiny.Phrase{}, t)
	ComparePhrases(remainder, tiny.Phrase{tiny.NewMeasurement([]byte{77})}, t)
}

func Test_Phrase_Read_Negative(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder := phrase.Read(-5)

	ComparePhrases(read, tiny.Phrase{}, t)
	ComparePhrases(remainder, tiny.Phrase{tiny.NewMeasurement([]byte{77})}, t)
}

func Test_Phrase_Read_AcrossMeasurements(t *testing.T) {
	phrase := tiny.NewPhrase(77, 22)

	read, remainder := phrase.Read(10)

	left1 := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0, 1, 1, 0, 1)
	left2 := tiny.NewMeasurement([]byte{}, 0, 0)
	ComparePhrases(read, tiny.Phrase{left1, left2}, t)

	right := tiny.NewMeasurement([]byte{}, 0, 1, 0, 1, 1, 0)
	ComparePhrases(remainder, tiny.Phrase{right}, t)
}

/**
ReadMeasurement
*/

func Test_Phrase_ReadMeasurement(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder := phrase.ReadMeasurement(4)

	left := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0)
	CompareMeasurements(read, left, t)

	right := tiny.NewMeasurement([]byte{}, 1, 1, 0, 1)
	ComparePhrases(remainder, tiny.Phrase{right}, t)
}

func Test_Phrase_ReadMeasurement_Zero(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder := phrase.ReadMeasurement(0)

	CompareMeasurements(read, tiny.NewMeasurement([]byte{}), t)
	ComparePhrases(remainder, tiny.Phrase{tiny.NewMeasurement([]byte{77})}, t)
}

func Test_Phrase_ReadMeasurement_Negative(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder := phrase.ReadMeasurement(-5)

	CompareMeasurements(read, tiny.NewMeasurement([]byte{}), t)
	ComparePhrases(remainder, tiny.Phrase{tiny.NewMeasurement([]byte{77})}, t)
}

func Test_Phrase_ReadMeasurement_OverByte(t *testing.T) {
	phrase := tiny.NewPhrase(77, 22, 33)

	read, remainder := phrase.ReadMeasurement(10)

	left := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0)
	CompareMeasurements(read, left, t)

	right1 := tiny.NewMeasurement([]byte{}, 0, 1, 0, 1, 1, 0)
	right2 := tiny.NewMeasurement([]byte{33})
	ComparePhrases(remainder, tiny.Phrase{right1, right2}, t)
}

func Test_Phrase_ReadMeasurement_ShouldPanicIfOverArchitectureBitWidth(t *testing.T) {
	defer test.ShouldPanic(t)
	tiny.NewPhrase().ReadMeasurement(tiny.GetArchitectureBitWidth() + 1)
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
	ComparePhrases(s, tiny.NewPhrase(77), t)
	ComparePhrases(m, tiny.NewPhrase(22), t)
	ComparePhrases(e, tiny.NewPhrase(33), t)
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

	ComparePhrases(s, eStart, t)
	ComparePhrases(m, eMiddle, t)
	ComparePhrases(e, eEnd, t)
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

	ComparePhrases(s, eStart, t)
	ComparePhrases(m, eMiddle, t)
	ComparePhrases(e, eEnd, t)
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

	ComparePhrases(s, eStart, t)
	ComparePhrases(m, eMiddle, t)
	ComparePhrases(e, eEnd, t)
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

	ComparePhrases(s, eStart, t)
	ComparePhrases(m, eMiddle, t)
	ComparePhrases(e, eEnd, t)
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

	ComparePhrases(s, eStart, t)
	ComparePhrases(m, eMiddle, t)
	ComparePhrases(e, eEnd, t)
}

/**
NOTE: FuzzyRead tests are located in fuzzy_test.go
*/

/*
*
Focus
*/

func Test_Phrase_Focus(t *testing.T) {
	phrase := tiny.Synthesize.RandomPhrase(1024, 32)
	length := phrase.BitLength()
	l, r := phrase.Focus()
	lb := l.BitLength()
	rb := r.BitLength()

	elb := length / 2
	erb := length - elb

	if lb != elb {
		t.Errorf("Expected length of left focus to be %d, got %d", elb, lb)
	}
	if rb != erb {
		t.Errorf("Expected length of right focus to be %d, got %d", erb, rb)
	}
}

func Test_Phrase_Focus_UnevenLength(t *testing.T) {
	phrase := tiny.Synthesize.RandomPhrase(1024, 32)
	phrase = append(tiny.NewPhraseFromBits(1), phrase...)
	length := phrase.BitLength()
	l, r := phrase.Focus()
	lb := l.BitLength()
	rb := r.BitLength()

	elb := length / 2
	erb := length - elb

	if lb != elb {
		t.Errorf("Expected length of left focus to be %d, got %d", elb, lb)
	}
	if rb != erb {
		t.Errorf("Expected length of right focus to be %d, got %d", erb, rb)
	}
}

func Test_Phrase_Focus_Recursion(t *testing.T) {
	phrase := tiny.Synthesize.RandomPhrase(1024, 32)
	length := phrase.BitLength()
	l, r := phrase.Focus(5)
	lb := l.BitLength()
	rb := r.BitLength()

	elb := length / 2 / 2 / 2 / 2 / 2
	erb := length - elb

	if lb != elb {
		t.Errorf("Expected length of left focus to be %d, got %d", elb, lb)
	}
	if rb != erb {
		t.Errorf("Expected length of right focus to be %d, got %d", erb, rb)
	}
}

func Test_Phrase_Focus_RecursionUnevenLength(t *testing.T) {
	phrase := tiny.Synthesize.RandomPhrase(1024, 32)
	phrase = append(tiny.NewPhraseFromBits(1), phrase...)
	length := phrase.BitLength()
	l, r := phrase.Focus(5)
	lb := l.BitLength()
	rb := r.BitLength()

	elb := length / 2 / 2 / 2 / 2 / 2
	erb := length - elb

	if lb != elb {
		t.Errorf("Expected length of left focus to be %d, got %d", elb, lb)
	}
	if rb != erb {
		t.Errorf("Expected length of right focus to be %d, got %d", erb, rb)
	}
}

/**
WalkBits
*/

func Test_Phrase_WalkBits(t *testing.T) {
	remainder := tiny.Synthesize.RandomPhrase(4, 32)
	bits := remainder.Bits()

	remainder.WalkBits(1, func(i int, m tiny.Measurement) {
		if bits[i] != m.GetAllBits()[0] {
			t.Errorf("Expected bit %d to be %d, got %d", i, bits[i], m.GetAllBits()[0])
		}
	})
}

func Test_Phrase_WalkBits_AtStride(t *testing.T) {
	remainder := tiny.Synthesize.RandomPhrase(4, 32)
	bits := remainder.Bits()
	i := 0

	remainder.WalkBits(3, func(_ int, m tiny.Measurement) {
		m.ForEachBit(func(_ int, b tiny.Bit) tiny.Bit {
			if bits[i] != b {
				t.Errorf("Expected bit %d to be %d, got %d", i, bits[i], b)
			}
			i++
			return b
		})
	})
}

func Test_Phrase_WalkBits_ShouldPanicIfStrideTooLarge(t *testing.T) {
	defer test.ShouldPanic(t)

	remainder := tiny.Synthesize.RandomPhrase(4, 8)
	remainder.WalkBits(tiny.GetArchitectureBitWidth()+1, func(i int, m tiny.Measurement) {})
}

func Test_Phrase_WalkBits_ShouldPanicIfStrideIsNegative(t *testing.T) {
	defer test.ShouldPanic(t)

	remainder := tiny.Synthesize.RandomPhrase(4, 8)
	remainder.WalkBits(-1, func(i int, m tiny.Measurement) {})
}

func Test_Phrase_WalkBits_ShouldPanicIfStrideIsZero(t *testing.T) {
	defer test.ShouldPanic(t)

	remainder := tiny.Synthesize.RandomPhrase(4, 8)
	remainder.WalkBits(0, func(i int, m tiny.Measurement) {})
}

/**
Read ZLE Tests

NOTE: These are nearly identical to their fuzzy_test.go counterpart functions,
      These just test the convenience method on Phrase
*/

func Test_Phrase_InterpretScaledZLE(t *testing.T) {
	tester := func(input int) {
		passage := tiny.NewZLEScaledPassage(input)
		result := tiny.Fuzzy.InterpretZLEScaled(passage)
		if input != result {
			t.Errorf("Expected %d, got %d", input, result)
		}
	}

	tester(0)
	tester(2)
	tester(5)
	tester(77)
	tester(333)
	tester(65536)
}

func Test_Phrase_ReadScaledZLE(t *testing.T) {
	data := tiny.NewPhraseFromBytesAndBits([]byte{77, 22, 33, 11, 77, 22, 33, 11}, 0, 1)

	tester := func(length int, eKey tiny.Measurement, data tiny.Phrase) {
		eProjection, eRemainder := data.Read(length)
		phrase := append(tiny.Phrase{eKey}, data...)
		key, projection, remainder := phrase.ReadZLEScaled()
		CompareMeasurements(key, eKey, t)
		ComparePhrases(projection, eProjection, t)
		ComparePhrases(remainder, eRemainder, t)
	}

	tester(2, tiny.NewMeasurement([]byte{}, 1), data)
	tester(3, tiny.NewMeasurement([]byte{}, 0, 1), data)
	tester(8, tiny.NewMeasurement([]byte{}, 0, 0, 1), data)
	tester(16, tiny.NewMeasurement([]byte{}, 0, 0, 0, 0), data)
	tester(64, tiny.NewMeasurement([]byte{}, 0, 0, 0, 1), data)
}

func Test_Phrase_Read64BitZLE(t *testing.T) {
	data := tiny.NewPhraseFromBytesAndBits([]byte{77, 22, 33, 11, 77, 22, 33, 11}, 0, 1)

	tester := func(length int, eKey tiny.Measurement, data tiny.Phrase) {
		eProjection, eRemainder := data.Read(length)
		phrase := append(tiny.Phrase{eKey}, data...)
		key, projection, remainder := phrase.ReadZLE64()
		CompareMeasurements(key, eKey, t)
		ComparePhrases(projection, eProjection, t)
		ComparePhrases(remainder, eRemainder, t)
	}

	tester(4, tiny.NewMeasurement([]byte{}, 1), data)
	tester(8, tiny.NewMeasurement([]byte{}, 0, 1), data)
	tester(16, tiny.NewMeasurement([]byte{}, 0, 0, 1), data)
	tester(32, tiny.NewMeasurement([]byte{}, 0, 0, 0, 0), data)
	tester(64, tiny.NewMeasurement([]byte{}, 0, 0, 0, 1), data)
}

func Test_Phrase_Read5BitZLE(t *testing.T) {
	data := tiny.NewPhraseFromBytesAndBits([]byte{77, 22, 33, 11, 77, 22, 33, 11}, 0, 1)

	tester := func(length int, eKey tiny.Measurement, data tiny.Phrase) {
		eProjection, eRemainder := data.Read(length)
		phrase := append(tiny.Phrase{eKey}, data...)
		key, projection, remainder := phrase.ReadZLE5()
		CompareMeasurements(key, eKey, t)
		ComparePhrases(projection, eProjection, t)
		ComparePhrases(remainder, eRemainder, t)
	}

	tester(1, tiny.NewMeasurement([]byte{}, 1), data)
	tester(2, tiny.NewMeasurement([]byte{}, 0, 1), data)
	tester(3, tiny.NewMeasurement([]byte{}, 0, 0, 1), data)
	tester(4, tiny.NewMeasurement([]byte{}, 0, 0, 0, 0), data)
	tester(5, tiny.NewMeasurement([]byte{}, 0, 0, 0, 1), data)
}

func Test_Phrase_ReadZLE(t *testing.T) {
	data := tiny.NewPhraseFromBytesAndBits([]byte{77, 22, 33, 11, 77, 22, 33, 11}, 0, 1)

	tester := func(length int, eKey tiny.Measurement, data tiny.Phrase) {
		eProjection, eRemainder := data.Read(length)
		phrase := append(tiny.Phrase{eKey}, data...)
		key, projection, remainder := phrase.ReadZLE(-1)
		CompareMeasurements(key, eKey, t)
		ComparePhrases(projection, eProjection, t)
		ComparePhrases(remainder, eRemainder, t)
	}

	tester(0, tiny.NewMeasurement([]byte{}, 1), data)
	tester(2, tiny.NewMeasurement([]byte{}, 0, 1), data)
	tester(4, tiny.NewMeasurement([]byte{}, 0, 0, 1), data)
	tester(8, tiny.NewMeasurement([]byte{}, 0, 0, 0, 1), data)
	tester(16, tiny.NewMeasurement([]byte{}, 0, 0, 0, 0, 1), data)
	tester(32, tiny.NewMeasurement([]byte{}, 0, 0, 0, 0, 0, 1), data)
}

/**
Invert
*/

func Test_Phrase_Invert(t *testing.T) {
	expected := tiny.NewPhraseFromBytesAndBits([]byte{178, 233, 222}, 0, 1, 1, 0)
	phrase := tiny.NewPhraseFromBytesAndBits([]byte{77, 22, 33}, 1, 0, 0, 1)
	// |        77       |         22      |        33       |    9    | <- Input Values
	// | 0 1 0 0 1 1 0 1 | 0 0 0 1 0 1 1 0 | 0 0 1 0 0 0 0 1 | 1 0 0 1 | <- Input
	// | 1 0 1 1 0 0 1 0 | 1 1 1 0 1 0 0 1 | 1 1 0 1 1 1 1 0 | 0 1 1 0 | <- Inverted
	// |       178       |        233      |       222       |    6    | <- Inverted Values
	phrase = phrase.Invert()
	ComparePhrases(phrase, expected, t)
}
