package testing

import (
	"fmt"
	"github.com/ignite-laboratories/core/relatively"
	"github.com/ignite-laboratories/tiny"
	"math"
	"math/big"
	"testing"
)

func Test_Phrase_NewPhraseFromBits(t *testing.T) {
	bits := []tiny.Bit{0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 1}
	p := tiny.NewPhraseFromBits(bits...)
	r, _, _ := p.Read(p.BitLength())
	CompareSlices(bits, r.Bits(), t)
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

	CompareSlices(bytes, expectedBytes, t)
	CompareSlices(bits, expectedBits, t)
}

func Test_Phrase_BitLength(t *testing.T) {
	phrase := tiny.Synthesize.RandomPhrase(32)
	length := phrase.BitLength()
	if length != 32*8 {
		t.Errorf("Expected %d, Got %d", 32*8, length)
	}
}

func Test_Phrase_AllBelowThreshold(t *testing.T) {
	below := tiny.NewPhrase(32, 55)
	if !below.AllBelowThreshold(55) {
		t.Errorf("Input data was below threshold, but AllBelowThreshold returned false")
	}

	above := tiny.NewPhrase(32, 77)
	if above.AllBelowThreshold(55) {
		t.Errorf("Input data was above threshold, but AllBelowThreshold returned true")
	}

	random := tiny.Synthesize.RandomPhrase(32)
	random[7] = tiny.NewMeasurement([]byte{77}) // ensure at least one is above threshold
	if above.AllBelowThreshold(55) {
		t.Errorf("Input data was above threshold, but AllBelowThreshold returned true")
	}
}

func Test_Phrase_CountBelowThreshold(t *testing.T) {
	threshold := 55

	below := tiny.Synthesize.Repeating(32, tiny.From.Byte(33)...)
	belowCount := below.CountBelowThreshold(threshold)
	if belowCount != 32 {
		t.Errorf("Expected 32 below a threshold of %d, Got %d", threshold, belowCount)
	}

	above := tiny.Synthesize.Repeating(32, tiny.From.Byte(77)...)
	aboveCount := above.CountBelowThreshold(threshold)
	if aboveCount != 0 {
		t.Errorf("Expected 0 below a threshold of %d, Got %d", threshold, aboveCount)
	}

	random := tiny.Synthesize.RandomPhrase(32)
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
	defer ShouldPanic(t)

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
	//      |       77        |       22        |       33        |  ← "Unaligned"
	//   0 1 0 1 0 0 1 1 - 0 1 0 0 0 1 0 1 - 1 0 0 0 1 0 0 0 - 0 1   ← Raw Bits
	//  |      83        |       69        |      136        |       ← "Aligned"

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
	//      |       77        |       22        |       33        |  ← "Unaligned"
	//   0 1 0 1 0 0 1 1 - 0 1 0 0 0 1 0 1 - 1 0 0 0 1 0 0 0 - 0 1   ← Raw Bits
	//  |   5   |   3    |    4   |   5    |    8   |   8    |  1    ← "Aligned"

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
	//      |       77        |       22        |       33        |  ← "Unaligned"
	//   0 1 0 1 0 0 1 1 - 0 1 0 0 0 1 0 1 - 1 0 0 0 1 0 0 0 - 0 1   ← Raw Bits
	//  |        333          |        88           |     33         ← "Aligned"

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
	defer ShouldPanic(t)
	phrase := tiny.NewPhrase(77, 22, 33)
	phrase.Align(0)
}

func Test_Phrase_Align_PanicIfNegativeWidth(t *testing.T) {
	defer ShouldPanic(t)
	phrase := tiny.NewPhrase(77, 22, 33)
	phrase.Align(-1)
}

func Test_Phrase_Align_PanicIfWidthTooLarge(t *testing.T) {
	defer ShouldPanic(t)
	phrase := tiny.NewPhrase(11, 33, 55, 99, 77, 22, 33)
	phrase.Align(tiny.GetArchitectureBitWidth() + 1)
}

/**
Read
*/

func Test_Phrase_Read(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder, _ := phrase.Read(4)

	left := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0)
	ComparePhrases(read, tiny.Phrase{left}, t)

	right := tiny.NewMeasurement([]byte{}, 1, 1, 0, 1)
	ComparePhrases(remainder, tiny.Phrase{right}, t)
}

func Test_Phrase_Read_Zero(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder, _ := phrase.Read(0)

	ComparePhrases(read, tiny.Phrase{}, t)
	ComparePhrases(remainder, tiny.Phrase{tiny.NewMeasurement([]byte{77})}, t)
}

func Test_Phrase_Read_Negative(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder, _ := phrase.Read(-5)

	ComparePhrases(read, tiny.Phrase{}, t)
	ComparePhrases(remainder, tiny.Phrase{tiny.NewMeasurement([]byte{77})}, t)
}

func Test_Phrase_Read_AcrossMeasurements(t *testing.T) {
	phrase := tiny.NewPhrase(77, 22)

	read, remainder, _ := phrase.Read(10)

	left1 := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0, 1, 1, 0, 1)
	left2 := tiny.NewMeasurement([]byte{}, 0, 0)
	ComparePhrases(read, tiny.Phrase{left1, left2}, t)

	right := tiny.NewMeasurement([]byte{}, 0, 1, 0, 1, 1, 0)
	ComparePhrases(remainder, tiny.Phrase{right}, t)
}

func Test_Phrase_Read_ReadPhraseOperationsErrorAppropriately(t *testing.T) {
	phrase := tiny.NewPhrase(77)
	length := phrase.BitLength()

	tester := func(operation func(int) (tiny.Phrase, tiny.Phrase, error)) {
		d := 1
		r, rr, err := operation(d)
		fmt.Println(r, rr)
		if err != nil {
			t.Fatalf("Did not expect an error reading %d bits into an %d bit phrase", d, length)
		}
		_, _, err = operation(length + d)
		if err == nil {
			t.Fatalf("Expected an error reading %d bits into an %d bit phrase", length+d, length)
		}
	}

	tester(phrase.Read)
	tester(phrase.ReadFromEnd)
}

func Test_Phrase_Read_ReadMeasurementErrorsAppropriately(t *testing.T) {
	phrase := tiny.NewPhrase(77)
	length := phrase.BitLength()

	d := 1
	r, rr, err := phrase.ReadMeasurement(d)
	fmt.Println(r, rr)
	if err != nil {
		t.Fatalf("Did not expect an error reading %d bits into an %d bit phrase", d, length)
	}
	_, _, err = phrase.ReadMeasurement(length + d)
	if err == nil {
		t.Fatalf("Expected an error reading %d bits into an %d bit phrase", length+d, length)
	}
}

func Test_Phrase_Read_ReadBitOperationsErrorAppropriately(t *testing.T) {
	phrase := tiny.NewPhrase(77)
	_, _, err := phrase.ReadLastBit()
	if err != nil {
		t.Fatalf("Did not expect an error reading the last bit of a 1 bit phrase")
	}

	phrase = tiny.NewPhrase()
	_, _, err = phrase.ReadLastBit()
	if err == nil {
		t.Fatalf("Expected an error reading the last bit of a 1 bit phrase")
	}

	phrase = tiny.NewPhrase(77)
	_, _, err = phrase.ReadNextBit()
	if err != nil {
		t.Fatalf("Did not expect an error reading the next bit of a 1 bit phrase")
	}

	phrase = tiny.NewPhrase()
	_, _, err = phrase.ReadNextBit()
	if err == nil {
		t.Fatalf("Expected an error reading the next bit of a 1 bit phrase")
	}
}

/**
ReadFromEnd
*/

func Test_Phrase_ReadFromEnd(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder, _ := phrase.ReadFromEnd(4)

	ComparePhrases(read, tiny.NewPhraseFromBits(1, 1, 0, 1), t)
	ComparePhrases(remainder, tiny.NewPhraseFromBits(0, 1, 0, 0), t)
}

func Test_Phrase_ReadFromEnd_NoData(t *testing.T) {
	phrase := tiny.NewPhrase()

	read, remainder, _ := phrase.ReadFromEnd(4)

	ComparePhrases(read, tiny.NewPhrase(), t)
	ComparePhrases(remainder, tiny.NewPhrase(), t)
}

func Test_Phrase_ReadFromEnd_UndersizedData(t *testing.T) {
	phrase := tiny.NewPhraseFromBits(1, 1)

	read, remainder, _ := phrase.ReadFromEnd(4)

	ComparePhrases(read, tiny.NewPhraseFromBits(1, 1), t)
	ComparePhrases(remainder, tiny.NewPhrase(), t)
}

/**
ReadLastBit
*/

func Test_Phrase_ReadLastBit_ErrorIfEmptyPhrase(t *testing.T) {
	phrase := tiny.NewPhrase()
	read, remainder, err := phrase.ReadLastBit()
	if err == nil {
		t.Fatalf("ReadLastBit should have returned an error with no bits left to read")
	}
	if read != tiny.Zero {
		t.Fatalf("Expected the read bit to be zero when no bits are present")
	}
	if remainder.BitLength() > 0 {
		t.Fatalf("Expected the remainder to be empty when no bits are present")
	}
}

func Test_Phrase_ReadLastBit(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	bit, remainder, _ := phrase.ReadLastBit()

	if bit != 1 {
		t.Errorf("Expected bit to be 1, got %d", bit)
	}

	ComparePhrases(remainder, tiny.NewPhraseFromBits(0, 1, 0, 0, 1, 1, 0), t)
}

func Test_Phrase_ReadLastBit_NoData(t *testing.T) {
	phrase := tiny.NewPhrase()

	bit, remainder, _ := phrase.ReadLastBit()

	if bit != 0 {
		t.Errorf("Expected bit to be 0, got %d", bit)
	}

	if remainder.BitLength() > 0 {
		t.Errorf("Expected remainder to be empty, got %s", remainder)
	}
}

func Test_Phrase_ReadLastBit_OneBit(t *testing.T) {
	phrase := tiny.NewPhraseFromBits(1)

	bit, remainder, _ := phrase.ReadLastBit()

	if bit != 1 {
		t.Errorf("Expected bit to be 0, got %d", bit)
	}

	if remainder.BitLength() > 0 {
		t.Errorf("Expected remainder to be empty, got %s", remainder)
	}
}

/**
ReadMeasurement
*/

func Test_Phrase_ReadMeasurement(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder, _ := phrase.ReadMeasurement(4)

	left := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0)
	CompareMeasurements(read, left, t)

	right := tiny.NewMeasurement([]byte{}, 1, 1, 0, 1)
	ComparePhrases(remainder, tiny.Phrase{right}, t)
}

func Test_Phrase_ReadMeasurement_Zero(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder, _ := phrase.ReadMeasurement(0)

	CompareMeasurements(read, tiny.NewMeasurement([]byte{}), t)
	ComparePhrases(remainder, tiny.Phrase{tiny.NewMeasurement([]byte{77})}, t)
}

func Test_Phrase_ReadMeasurement_Negative(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	read, remainder, _ := phrase.ReadMeasurement(-5)

	CompareMeasurements(read, tiny.NewMeasurement([]byte{}), t)
	ComparePhrases(remainder, tiny.Phrase{tiny.NewMeasurement([]byte{77})}, t)
}

func Test_Phrase_ReadMeasurement_OverByte(t *testing.T) {
	phrase := tiny.NewPhrase(77, 22, 33)

	read, remainder, _ := phrase.ReadMeasurement(10)

	left := tiny.NewMeasurement([]byte{}, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0)
	CompareMeasurements(read, left, t)

	right1 := tiny.NewMeasurement([]byte{}, 0, 1, 0, 1, 1, 0)
	right2 := tiny.NewMeasurement([]byte{33})
	ComparePhrases(remainder, tiny.Phrase{right1, right2}, t)
}

func Test_Phrase_ReadMeasurement_ShouldPanicIfOverArchitectureBitWidth(t *testing.T) {
	defer ShouldPanic(t)
	tiny.NewPhrase().ReadMeasurement(tiny.GetArchitectureBitWidth() + 1)
}

/**
ReadNextBit
*/

func Test_Phrase_ReadNextBit(t *testing.T) {
	for x := 0; x < 256; x++ {
		for y := 0; y < 256; y++ {
			phrase := tiny.NewPhrase(byte(x), byte(y))
			bit, remainder, err := phrase.ReadNextBit()
			remainder = remainder.Align()

			expected := tiny.From.Number(x, 8)
			eBit := expected[0]
			eRemainder := expected[1:]
			eRemainder = append(eRemainder, tiny.From.Number(y, 8)...)

			if err != nil {
				t.Errorf("Expected no error, got %s", err)
			}
			if bit != eBit {
				t.Errorf("Expected bit to be %d, got %d", eBit, bit)
			}
			ComparePhrases(remainder, tiny.NewPhraseFromBits(eRemainder...), t)
		}
	}
}

func Test_Phrase_ReadNextBit_ShouldErrorWhenEndOfPhrase(t *testing.T) {
	phrase := tiny.NewPhrase(33, 22)

	for i := 0; i <= phrase.BitLength(); i++ {
		_, remainder, err := phrase.ReadNextBit()
		phrase = remainder
		if i == phrase.BitLength() && err == nil {
			t.Fatalf("Expected an error when reading beyond the end of the phrase, got nil")
		}
	}
}

/**
ReadUntilOne
*/

func Test_Phrase_ReadUntilOne(t *testing.T) {
	for i := 0; i < 16; i++ {
		input := tiny.Synthesize.Zeros(i)
		data := tiny.Synthesize.RandomPhrase(8)
		data = data.PrependBits(1) // Ensure the data starts with a 1

		input = append(input, data...)
		zeros, remainder := input.ReadUntilOne()
		if zeros != i {
			t.Errorf("Expected %d zeros, got %d", i, zeros)
		}
		CompareUnalignedPhrases(remainder, data, t)
	}
}

func Test_Phrase_ReadUntilOne_WithLimit(t *testing.T) {
	input := tiny.NewPhraseFromBits(0, 0, 0, 0, 0, 1, 0, 0, 1)
	zeros, remainder := input.ReadUntilOne(4)
	remainder = remainder.Align()
	if zeros != 4 {
		t.Errorf("Expected %d zeros, got %d", 4, zeros)
	}
	ComparePhrases(remainder, tiny.NewPhraseFromBits(0, 1, 0, 0, 1), t)
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
	//   0 1 0 0 1 1 0 1 - 0 0 0 1 0 1 1 0 - 0 0 1 0 0 0 0 1  ← Raw Bits
	//  |     Start      |     Middle      |      End       | ← "Trifurcated"
	phrase := tiny.NewPhrase(77, 22, 33)

	s, m, e, _ := phrase.Trifurcate(8, 8)
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
	//   0 1 0 0 - 1 1 0 1 - 0 0 0 1 0 1 1 0 - 0 0 1 0 - 0 0 0 1  ← Raw Bits
	//  | Start  | Middle1 |     Middle2     | Middle3 |   End  | ← "Trifurcated"
	phrase := tiny.NewPhrase(77, 22, 33)

	s, m, e, _ := phrase.Trifurcate(4, 16)

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
	//   0 1 - 0 0 1 1 0 1      ← Raw Bits
	//  | S  |   Middle   | E | ← "Trifurcated"
	phrase := tiny.NewPhrase(77)

	s, m, e, _ := phrase.Trifurcate(2, 8)

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
	//   0 1 0 0 1 1 0 1          ← Raw Bits
	//  |     Start     | M | E | ← "Trifurcated"
	phrase := tiny.NewPhrase(77)

	s, m, e, _ := phrase.Trifurcate(10, 8)

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
	//       0 1 0 0 - 1 1 0 1  ← Raw Bits
	//  | S | Middle |  End   | ← "Trifurcated"
	phrase := tiny.NewPhrase(77)

	s, m, e, _ := phrase.Trifurcate(0, 4)

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
	//       0 1 0 0 1 1 0 1      ← Raw Bits
	//  | S |    Middle     | E | ← "Trifurcated"
	phrase := tiny.NewPhrase(77)

	s, m, e, _ := phrase.Trifurcate(0, 10)

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
	defer ShouldPanic(t)

	remainder := tiny.Synthesize.RandomPhrase(4, 8)
	remainder.WalkBits(tiny.GetArchitectureBitWidth()+1, func(i int, m tiny.Measurement) {})
}

func Test_Phrase_WalkBits_ShouldPanicIfStrideIsNegative(t *testing.T) {
	defer ShouldPanic(t)

	remainder := tiny.Synthesize.RandomPhrase(4, 8)
	remainder.WalkBits(-1, func(i int, m tiny.Measurement) {})
}

func Test_Phrase_WalkBits_ShouldPanicIfStrideIsZero(t *testing.T) {
	defer ShouldPanic(t)

	remainder := tiny.Synthesize.RandomPhrase(4, 8)
	remainder.WalkBits(0, func(i int, m tiny.Measurement) {})
}

/**
Padding
*/

func Test_Phrase_PadLeftToLength(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	// Pad to over the target index

	paddedAA := phrase.PadLeftToLength(20)
	paddedAB := phrase.PadLeftToLength(20, tiny.One)
	expectedAA := tiny.NewPhrase(0)
	expectedAA = expectedAA.AppendBits(0, 0, 0, 0)
	expectedAA = expectedAA.AppendBytes(77)

	expectedAB := tiny.NewPhrase(255)
	expectedAB = expectedAB.AppendBits(1, 1, 1, 1)
	expectedAB = expectedAB.AppendBytes(77)

	ComparePhrases(paddedAA, expectedAA, t)
	ComparePhrases(paddedAB, expectedAB, t)

	// Pad to under the target index

	expectedUndersized := tiny.NewPhraseFromBits(0, 1, 0, 0, 1, 1, 0, 1)

	paddedBA := phrase.PadLeftToLength(5)
	paddedBB := phrase.PadLeftToLength(5, tiny.One)

	ComparePhrases(paddedBA, expectedUndersized, t)
	ComparePhrases(paddedBB, expectedUndersized, t)

	paddedCA := phrase.PadLeftToLength(0)
	paddedCB := phrase.PadLeftToLength(0, tiny.One)

	ComparePhrases(paddedCA, expectedUndersized, t)
	ComparePhrases(paddedCB, expectedUndersized, t)

	paddedDA := phrase.PadLeftToLength(-7)
	paddedDB := phrase.PadLeftToLength(7, tiny.One)

	ComparePhrases(paddedDA, expectedUndersized, t)
	ComparePhrases(paddedDB, expectedUndersized, t)
}

func Test_Phrase_PadRightToLength(t *testing.T) {
	phrase := tiny.NewPhrase(77)

	// Pad to over the target index

	paddedAA := phrase.PadRightToLength(20)
	paddedAB := phrase.PadRightToLength(20, tiny.One)

	expectedAA := tiny.NewPhrase(77)
	expectedAA = expectedAA.AppendBytes(0)
	expectedAA = expectedAA.AppendBits(0, 0, 0, 0)

	expectedAB := tiny.NewPhrase(77)
	expectedAB = expectedAB.AppendBytes(255)
	expectedAB = expectedAB.AppendBits(1, 1, 1, 1)

	ComparePhrases(paddedAA, expectedAA, t)
	ComparePhrases(paddedAB, expectedAB, t)

	// Pad to under the target index

	expectedUndersized := tiny.NewPhraseFromBits(0, 1, 0, 0, 1, 1, 0, 1)

	paddedBA := phrase.PadRightToLength(5)
	paddedBB := phrase.PadRightToLength(5, tiny.One)

	ComparePhrases(paddedBA, expectedUndersized, t)
	ComparePhrases(paddedBB, expectedUndersized, t)

	paddedCA := phrase.PadRightToLength(0)
	paddedCB := phrase.PadRightToLength(0, tiny.One)

	ComparePhrases(paddedCA, expectedUndersized, t)
	ComparePhrases(paddedCB, expectedUndersized, t)

	paddedDA := phrase.PadRightToLength(-7)
	paddedDB := phrase.PadRightToLength(7, tiny.One)

	ComparePhrases(paddedDA, expectedUndersized, t)
	ComparePhrases(paddedDB, expectedUndersized, t)
}

/**
Arithmetic and Logic Gates
*/

func Test_Phrase_Add_StressTest(t *testing.T) {
	for i := 0; i < 1<<14; i++ {
		a := tiny.Synthesize.RandomBits(11)
		b := tiny.Synthesize.RandomBits(11)

		c := a.Add(b)
		cStr := c.StringBinary()

		cBigInt := new(big.Int).Add(a.AsBigInt(), b.AsBigInt())
		cBigIntStr := cBigInt.Text(2)

		if cStr != cBigIntStr {
			t.Errorf("Expected %s + %s = %s, got %s", a.StringBinary(), b.StringBinary(), cBigIntStr, cStr)
		}
	}
}

func Test_Phrase_Add_Multiple_StressTest(t *testing.T) {
	for i := 0; i < 1<<14; i++ {
		a := tiny.Synthesize.RandomBits(11)
		b := tiny.Synthesize.RandomBits(11)
		c := tiny.Synthesize.RandomBits(11)
		d := tiny.Synthesize.RandomBits(11)

		e := a.Add(b).Add(c).Add(d)
		eStr := e.StringBinary()

		eBigInt := new(big.Int).Add(a.AsBigInt(), b.AsBigInt())
		eBigInt = new(big.Int).Add(eBigInt, c.AsBigInt())
		eBigInt = new(big.Int).Add(eBigInt, d.AsBigInt())
		eBigIntStr := eBigInt.Text(2)

		if eStr != eBigIntStr {
			t.Errorf("Expected %s, got %s", eBigIntStr, eStr)
		}
	}
}

func Test_Phrase_Minus_StressTest(t *testing.T) {
	for i := 0; i < 1<<16; i++ {
		a := tiny.Synthesize.RandomBits(11)
		b := tiny.Synthesize.RandomBits(11)

		c := a.Minus(b)
		data := c.GetData()
		sign := c.GetSign()
		cStr := data.StringBinary()

		if sign == 1 {
			cStr = "-" + cStr
		}

		if len(cStr) <= 0 {
			cStr = "0"
		}

		cBigInt := new(big.Int).Sub(a.AsBigInt(), b.AsBigInt())
		cBigIntStr := cBigInt.Text(2)

		if cStr != cBigIntStr {
			t.Errorf("Expected %s - %s = %s, got %s", a.StringBinary(), b.StringBinary(), cBigIntStr, cStr)
		}
	}
}

func Test_Phrase_Times_StressTest(t *testing.T) {
	for i := 0; i < 1<<13; i++ {
		a := tiny.Synthesize.RandomBits(11)
		b := tiny.Synthesize.RandomBits(11)

		c := a.Times(b)
		cStr := c.StringBinary()

		cBigInt := new(big.Int).Mul(a.AsBigInt(), b.AsBigInt())
		cBigIntStr := cBigInt.Text(2)

		if cStr != cBigIntStr {
			t.Errorf("Expected %s * %s = %s, got %s", a.StringBinary(), b.StringBinary(), cBigIntStr, cStr)
		}
	}
}

func Test_Phrase_DividedBy_StressTest(t *testing.T) {
	for i := 0; i < 1<<15; i++ {
		a := tiny.Synthesize.RandomBits(11)
		b := tiny.Synthesize.RandomBits(3)

		c := a.DividedBy(b)
		cStr := c.StringBinary()

		if len(cStr) == 0 {
			cStr = "0"
		}

		cBigInt := new(big.Int).Div(a.AsBigInt(), b.AsBigInt())
		cBigIntStr := cBigInt.Text(2)

		if cStr != cBigIntStr {
			t.Errorf("Expected %s / %s = %s, got %s", a.StringBinary(), b.StringBinary(), cBigIntStr, cStr)
		}
	}
}

func Test_Phrase_Modulo_StressTest(t *testing.T) {
	for i := 0; i < 1<<16; i++ {
		a := tiny.Synthesize.RandomBits(11)
		b := tiny.Synthesize.RandomBits(3)

		c := a.Modulo(b)
		cStr := c.StringBinary()

		if len(cStr) == 0 {
			cStr = "0"
		}

		cBigInt := new(big.Int).Mod(a.AsBigInt(), b.AsBigInt())
		cBigIntStr := cBigInt.Text(2)

		if cStr != cBigIntStr {
			t.Errorf("Expected %s %% %s = %s, got %s", a.StringBinary(), b.StringBinary(), cBigIntStr, cStr)
		}
	}
}

//func Test_Phrase_ToThePowerOf_StressTest(t *testing.T) {
//	for i := 0; i < 1<<16; i++ {
//		a := tiny.Synthesize.RandomBits(11)
//		b := tiny.Synthesize.RandomBits(3)
//
//		c, _ := a.ToThePowerOf(false, b)
//		cStr := c.StringBinary()
//
//		if len(cStr) == 0 {
//			cStr = "0"
//		}
//
//		cBigInt := new(big.Int).Exp(a.AsBigInt(), b.AsBigInt(), nil)
//		cBigIntStr := cBigInt.Text(2)
//
//		if cStr != cBigIntStr {
//			t.Errorf("Expected %s %% %s = %s, got %s", a.StringBinary(), b.StringBinary(), cBigIntStr, cStr)
//		}
//	}
//}

func Test_Phrase_LogicGates(t *testing.T) {
	// Test logic:
	//     We will test that a and b yield the following identities
	//
	// | 0 1 0 0 1 1 0 1 |  (77) ← a
	// | 0 1 0 1 1 0 0 0 |  (88) ← b
	//
	// Unary Identities -
	// | 1 0 1 1 0 0 1 0 | (178) ← NOT a
	// | 1 0 1 0 0 1 1 1 | (167) ← NOT b
	//
	// Binary Identities -
	// | 0 1 0 0 1 0 0 0 |  (72) ← a  AND b
	// | 0 1 0 1 1 1 0 1 |  (93) ← a   OR b
	// | 0 0 0 1 0 1 0 1 |  (21) ← a  XOR b
	// | 1 1 1 0 1 0 1 0 | (234) ← a XNOR b
	// | 1 0 1 1 0 1 1 1 | (183) ← a NAND b
	// | 1 0 1 0 0 0 1 0 | (162) ← a  NOR b

	a := tiny.NewPhrase(77)
	b := tiny.NewPhrase(88)

	notA := a.NOT().Int()
	if notA != 178 {
		t.Errorf("NOT a - Expected 178, got %d", notA)
	}

	notB := b.NOT().Int()
	if notB != 167 {
		t.Errorf("NOT b - Expected 167, got %d", notB)
	}

	and := a.AND(b).Int()
	if and != 72 {
		t.Errorf("a AND b - Expected 72, got %d", and)
	}

	or := a.OR(b).Int()
	if or != 93 {
		t.Errorf("a OR b - Expected 93, got %d", or)
	}

	xor := a.XOR(b).Int()
	if xor != 21 {
		t.Errorf("a XOR b - Expected 21, got %d", xor)
	}

	xnor := a.XNOR(b).Int()
	if xnor != 234 {
		t.Errorf("a XNOR b - Expected 234, got %d", xnor)
	}

	nand := a.NAND(b).Int()
	if nand != 183 {
		t.Errorf("a NAND b - Expected 183, got %d", nand)
	}

	nor := a.NOR(b).Int()
	if nor != 162 {
		t.Errorf("a NOR b - Expected 162, got %d", nor)
	}
}

func Test_Phrase_CompareTo(t *testing.T) {
	a := tiny.NewPhrase(66)
	b := tiny.NewPhrase(77)
	c := tiny.NewPhrase(88)

	aa := a.CompareTo(a)
	if aa != relatively.Aligned {
		t.Errorf("Expected %d, got %d", relatively.Aligned, aa)
	}

	ab := a.CompareTo(b)
	if ab != relatively.Before {
		t.Errorf("Expected %d, got %d", relatively.Before, ab)
	}

	ac := a.CompareTo(c)
	if ac != relatively.Before {
		t.Errorf("Expected %d, got %d", relatively.Before, ac)
	}

	ba := b.CompareTo(a)
	if ba != relatively.After {
		t.Errorf("Expected %d, got %d", relatively.After, ba)
	}

	bb := b.CompareTo(b)
	if bb != relatively.Aligned {
		t.Errorf("Expected %d, got %d", relatively.Aligned, bb)
	}

	bc := b.CompareTo(c)
	if bc != relatively.Before {
		t.Errorf("Expected %d, got %d", relatively.Before, bc)
	}

	ca := c.CompareTo(a)
	if ca != relatively.After {
		t.Errorf("Expected %d, got %d", relatively.After, ca)
	}

	cb := c.CompareTo(b)
	if cb != relatively.After {
		t.Errorf("Expected %d, got %d", relatively.After, cb)
	}

	cc := c.CompareTo(c)
	if cc != relatively.Aligned {
		t.Errorf("Expected %d, got %d", relatively.Aligned, cc)
	}
}

func Test_Phrase_Int_IgnoresBitsAboveArchitectureBitWidth(t *testing.T) {
	data := tiny.Synthesize.RandomBits(tiny.GetArchitectureBitWidth())
	valBefore := data.Int()

	data.AppendBytes(77)
	valAfter := data.Int()

	if valBefore != valAfter {
		t.Errorf("Expected %d, got %d", valBefore, valAfter)
	}
}

func Test_Phrase_Int_StressTest(t *testing.T) {
	for i := 0; i <= tiny.GetArchitectureBitWidth(); i++ {
		oversize := i - tiny.GetArchitectureBitWidth()
		if oversize < 0 {
			oversize = 0
		}

		for ii := 0; ii < 1<<10; ii++ {
			data := tiny.Synthesize.RandomBits(int(math.Min(float64(tiny.GetArchitectureBitWidth()), float64(i))))
			data = append(data, tiny.Synthesize.RandomBits(oversize)...)

			val := data.Int()
			valBigInt := int(data.AsBigInt().Int64())

			if val != valBigInt {
				t.Errorf("Expected %d, got %d", valBigInt, val)
			}
		}
	}

	data := tiny.NewPhraseFromBytesAndBits([]byte{77, 33}, 0, 1)
	val := data.Int()
	if val != 78981 {
		t.Errorf("Expected 77, got %d", val)
	}
}
