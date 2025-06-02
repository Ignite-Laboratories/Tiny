package testing

import (
	"fmt"
	"github.com/ignite-laboratories/support/test"
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Phrase_FuzzyRead_ZLEKeyReading(t *testing.T) {
	tester := func(expected tiny.Measurement) {
		phrase := tiny.NewPhraseFromBitsAndBytes(expected.Bits, 77, 22, 33)
		phrase.FuzzyRead(tiny.Fuzzy.ZLEKey(), func(key tiny.Measurement) int {
			CompareMeasurements(expected, key, t)
			return 0
		})
	}

	tester(tiny.NewMeasurement([]byte{}, 1))
	tester(tiny.NewMeasurement([]byte{}, 0, 1))
	tester(tiny.NewMeasurement([]byte{}, 0, 0, 1))
	tester(tiny.NewMeasurement([]byte{}, 0, 0, 0, 0))
	tester(tiny.NewMeasurement([]byte{}, 0, 0, 0, 1))
}

func Test_Phrase_FuzzyRead_64BitZLE(t *testing.T) {
	data := tiny.NewPhraseFromBytesAndBits([]byte{77, 22, 33, 11, 77, 22, 33, 11}, 0, 1)

	tester := func(length int, eK ...tiny.Bit) {
		eP, eR := data.Read(length)
		bits := data.Bits()
		newBits := append(eK, bits...)
		phrase := tiny.NewPhraseFromBits(newBits...)
		k, p, r := phrase.FuzzyRead(tiny.Fuzzy.ZLEKey(), tiny.Fuzzy.ParseZLE64)

		test.CompareSlices(eK, k.Bits, t)
		test.CompareSlices(eP.Bits(), p.Bits(), t)
		test.CompareSlices(eR.Bits(), r.Bits(), t)
	}

	tester(4, 1)
	tester(8, 0, 1)
	tester(16, 0, 0, 1)
	tester(32, 0, 0, 0, 0)
	tester(64, 0, 0, 0, 1)
}

func Test_Phrase_FuzzyRead_5BitZLE(t *testing.T) {
	data := tiny.NewPhraseFromBytesAndBits([]byte{77, 22, 33, 11, 77, 22, 33, 11}, 0, 1)

	tester := func(length int, eK ...tiny.Bit) {
		eP, eR := data.Read(length)
		bits := data.Bits()
		newBits := append(eK, bits...)
		phrase := tiny.NewPhraseFromBits(newBits...)
		k, p, r := phrase.FuzzyRead(tiny.Fuzzy.ZLEKey(), tiny.Fuzzy.ParseZLE5)

		test.CompareSlices(eK, k.Bits, t)
		test.CompareSlices(eP.Bits(), p.Bits(), t)
		test.CompareSlices(eR.Bits(), r.Bits(), t)
	}

	tester(1, 1)
	tester(2, 0, 1)
	tester(3, 0, 0, 1)
	tester(4, 0, 0, 0, 0)
	tester(5, 0, 0, 0, 1)
}

func Test_Phrase_FuzzyRead_ZLE(t *testing.T) {
	data := tiny.NewPhraseFromBytesAndBits([]byte{77, 22, 33, 11, 77, 22, 33, 11}, 0, 1)

	tester := func(length int, eK ...tiny.Bit) {
		eP, eR := data.Read(length)
		bits := data.Bits()
		newBits := append(eK, bits...)
		phrase := tiny.NewPhraseFromBits(newBits...)
		k, p, r := phrase.FuzzyRead(tiny.Fuzzy.ZLEKey(-1), tiny.Fuzzy.ParseZLE)

		test.CompareSlices(eK, k.Bits, t)
		test.CompareSlices(eP.Bits(), p.Bits(), t)
		test.CompareSlices(eR.Bits(), r.Bits(), t)
	}

	tester(0, 1)
	tester(2, 0, 1)
	tester(4, 0, 0, 1)
	tester(8, 0, 0, 0, 1)
	tester(16, 0, 0, 0, 0, 1)
	tester(32, 0, 0, 0, 0, 0, 1)
}

func Test_Phrase_FuzzyRead_SixtyFour(t *testing.T) {
	// Test logic:
	//
	// Input Keys:
	//     13 -> 0 0 - 0 0 1 1 0 1
	//     77 -> 0 1 - 0 0 1 1 0 1
	//    141 -> 1 0 - 0 0 1 1 0 1
	//    205 -> 1 1 - 0 0 1 1 0 1
	//
	// Remainder Values:
	//     22 -> 0 0 0 1 0 1 1 0
	//     33 -> 0 0 1 0 0 0 0 1
	//
	// Example Output:
	//  |        Input        |        22       |         33      |
	//  | 0 1 | 0 0 | 1 1 0 1 - 0 0 0 1 0 1 1 0 - 0 0 1 0 0 0 0 1 |  <- Raw Bits
	//  | Key |  C  |                 Remainder                   | <- Fuzzy Read

	tester := func(phrase tiny.Phrase, keyBits []tiny.Bit, continuationBits []tiny.Bit, remainderBits []tiny.Bit) {
		k, c, r := phrase.FuzzyRead(tiny.Fuzzy.Count(2), tiny.Fuzzy.SixtyFour)

		eKey := tiny.NewMeasurement([]byte{}, keyBits...)
		eContinuation := tiny.Phrase{}
		if len(continuationBits) > 0 {
			eContinuation = append(eContinuation, tiny.NewMeasurement([]byte{}, continuationBits...))
		}
		eRemainder := tiny.Phrase{}
		if len(remainderBits) > 0 {
			eRemainder = append(eRemainder, tiny.NewMeasurement([]byte{}, remainderBits...))
		}
		eRemainder = append(eRemainder, tiny.NewMeasurement([]byte{22}), tiny.NewMeasurement([]byte{33}))

		CompareMeasurements(k, eKey, t)
		ComparePhrases(c, eContinuation, t)
		ComparePhrases(r, eRemainder, t)
	}

	tester(tiny.NewPhrase(13, 22, 33), []tiny.Bit{0, 0}, []tiny.Bit{}, []tiny.Bit{0, 0, 1, 1, 0, 1})
	tester(tiny.NewPhrase(77, 22, 33), []tiny.Bit{0, 1}, []tiny.Bit{0, 0}, []tiny.Bit{1, 1, 0, 1})
	tester(tiny.NewPhrase(141, 22, 33), []tiny.Bit{1, 0}, []tiny.Bit{0, 0, 1, 1}, []tiny.Bit{0, 1})
	tester(tiny.NewPhrase(205, 22, 33), []tiny.Bit{1, 1}, []tiny.Bit{0, 0, 1, 1, 0, 1}, []tiny.Bit{})
}

func Test_Phrase_FuzzyRead_SixtyFour_MaxSixBits(t *testing.T) {
	phrase := tiny.NewPhrase(13, 22, 33)
	_, c, _ := phrase.FuzzyRead(tiny.Fuzzy.Count(12), tiny.Fuzzy.SixtyFour)
	if c.BitLength() != 6 {
		t.Errorf("Expected 6 bits, got %d", c.BitLength())
	}
}

func Test_Phrase_FuzzyRead_Window(t *testing.T) {
	// Test logic:
	//
	// Input Keys:
	//     13 -> 0 0 - 0 0 1 1 0 1
	//     77 -> 0 1 - 0 0 1 1 0 1
	//    141 -> 1 0 - 0 0 1 1 0 1
	//    205 -> 1 1 - 0 0 1 1 0 1
	//
	// Remainder Values:
	//     22 -> 0 0 0 1 0 1 1 0
	//     33 -> 0 0 1 0 0 0 0 1
	//
	// Example Output:
	//
	//		 FuzzyRead(2, tiny.Fuzzy.Window(3))
	//
	//       |        Input        |        22       |         33      | <- Raw Values
	//       | 0 0 | 0 0 1 | 1 0 1 - 0 0 0 1 0 1 1 0 - 0 0 1 0 0 0 0 1 | <- Raw bits
	//       | Key | Cont  |              Remainder                    | <- Fuzzy read
	//
	//       |        Input      |        22       |         33      | <- Raw Values
	//       | 0 1 | 0 0 1 1 0 1 | 0 0 0 1 0 1 1 0 - 0 0 1 0 0 0 0 1 | <- Raw bits
	//       | Key |     Cont    |          Remainder                | <- Fuzzy read
	//
	//       |        Input      |         22        |         33      | <- Raw Values
	//       | 1 0 | 0 0 1 1 0 1 - 0 0 0 | 1 0 1 1 0 - 0 0 1 0 0 0 0 1 | <- Raw bits
	//       | Key |    Continuation     |         Remainder           | <- Fuzzy read
	//
	//       |        Input      |         22        |         33      | <- Raw Values
	//       | 1 1 | 0 0 1 1 0 1 - 0 0 0 1 0 1 | 1 0 - 0 0 1 0 0 0 0 1 | <- Raw bits
	//       | Key |        Continuation       |       Remainder       | <- Fuzzy read

	tester := func(phrase tiny.Phrase, keyBits []tiny.Bit, continuationBits tiny.Phrase, remainder tiny.Phrase) {
		k, c, r := phrase.FuzzyRead(tiny.Fuzzy.Count(2), tiny.Fuzzy.Window(3))
		eKey := tiny.NewMeasurement([]byte{}, keyBits...)
		CompareMeasurements(k, eKey, t)
		ComparePhrases(c, continuationBits, t)
		ComparePhrases(r, remainder, t)
	}

	tester(tiny.NewPhrase(13, 22, 33), []tiny.Bit{0, 0}, tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 0, 1)}, tiny.Phrase{tiny.NewMeasurement([]byte{}, 1, 0, 1), tiny.NewMeasurement([]byte{22}), tiny.NewMeasurement([]byte{33})})
	tester(tiny.NewPhrase(77, 22, 33), []tiny.Bit{0, 1}, tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 0, 1, 1, 0, 1)}, tiny.Phrase{tiny.NewMeasurement([]byte{22}), tiny.NewMeasurement([]byte{33})})
	tester(tiny.NewPhrase(141, 22, 33), []tiny.Bit{1, 0}, tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 0, 1, 1, 0, 1), tiny.NewMeasurement([]byte{}, 0, 0, 0)}, tiny.Phrase{tiny.NewMeasurement([]byte{}, 1, 0, 1, 1, 0), tiny.NewMeasurement([]byte{33})})
	tester(tiny.NewPhrase(205, 22, 33), []tiny.Bit{1, 1}, tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 0, 1, 1, 0, 1), tiny.NewMeasurement([]byte{}, 0, 0, 0, 1, 0, 1)}, tiny.Phrase{tiny.NewMeasurement([]byte{}, 1, 0), tiny.NewMeasurement([]byte{33})})
}

func Test_Phrase_FuzzyRead_Window_ZeroWidth(t *testing.T) {
	defer test.ShouldPanic(t)
	phrase := tiny.NewPhrase(13, 22, 33)
	phrase.FuzzyRead(tiny.Fuzzy.Count(2), tiny.Fuzzy.Window(0))
}

func Test_Phrase_FuzzyRead_Window_NegativeWidth(t *testing.T) {
	defer test.ShouldPanic(t)
	phrase := tiny.NewPhrase(13, 22, 33)
	phrase.FuzzyRead(tiny.Fuzzy.Count(2), tiny.Fuzzy.Window(-1))
}

func Test_Phrase_FuzzyRead_PowerWindow(t *testing.T) {
	// Test logic:
	//
	// Input Keys:
	//     13 -> 0 0 - 0 0 1 1 0 1
	//     77 -> 0 1 - 0 0 1 1 0 1
	//    141 -> 1 0 - 0 0 1 1 0 1
	//    205 -> 1 1 - 0 0 1 1 0 1
	//
	// Remainder Values:
	//     22 -> 0 0 0 1 0 1 1 0
	//     33 -> 0 0 1 0 0 0 0 1
	//     55 -> 0 0 1 1 0 1 1 1
	//
	// Example Output:
	//
	//		 FuzzyRead(2, tiny.Fuzzy.Window(3))
	//
	//       |        Input        |        22       |         33      |         55      | <- Raw Values
	//       | 0 0 | 0 0 1 | 1 0 1 - 0 0 0 1 0 1 1 0 - 0 0 1 0 0 0 0 1 - 0 0 1 1 0 1 1 1 | <- Raw bits
	//       | Key | Cont  |                        Remainder                            | <- Fuzzy read
	//
	//       |       Input       |        22       |         33      |         55      | <- Raw Values
	//       | 0 1 | 0 0 1 1 0 1 | 0 0 0 1 0 1 1 0 - 0 0 1 0 0 0 0 1 - 0 0 1 1 0 1 1 1 | <- Raw bits
	//       | Key |    Cont     |                    Remainder                        | <- Fuzzy read
	//
	//       |       Input       |         22        |         33      |         55      | <- Raw Values
	//       | 1 0 | 0 0 1 1 0 1 - 0 0 0 1 0 1 | 1 0 - 0 0 1 0 0 0 0 1 - 0 0 1 1 0 1 1 1 | <- Raw bits
	//       | Key |       Continuation        |                Remainder                | <- Fuzzy read
	//
	//       |       Input       |        22       |         33      |         55        | <- Raw Values
	//       | 1 1 | 0 0 1 1 0 1 - 0 0 0 1 0 1 1 0 - 0 0 1 0 0 0 0 1 - 0 0 | 1 1 0 1 1 1 | <- Raw bits
	//       | Key |                      Continuation                     | Remainder   | <- Fuzzy read

	tester := func(phrase tiny.Phrase, keyBits []tiny.Bit, continuationBits tiny.Phrase, remainder tiny.Phrase) {
		k, c, r := phrase.FuzzyRead(tiny.Fuzzy.Count(2), tiny.Fuzzy.PowerWindow(3))
		eKey := tiny.NewMeasurement([]byte{}, keyBits...)
		CompareMeasurements(k, eKey, t)
		ComparePhrases(c, continuationBits, t)
		ComparePhrases(r, remainder, t)
	}

	tester(tiny.NewPhrase(13, 22, 33, 55), []tiny.Bit{0, 0}, tiny.NewPhraseFromBits(0, 0, 1), tiny.NewPhraseFromBitsAndBytes([]tiny.Bit{1, 0, 1}, 22, 33, 55))
	tester(tiny.NewPhrase(77, 22, 33, 55), []tiny.Bit{0, 1}, tiny.NewPhraseFromBits(0, 0, 1, 1, 0, 1), tiny.NewPhrase(22, 33, 55))
	tester(tiny.NewPhrase(141, 22, 33, 55), []tiny.Bit{1, 0}, tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 0, 1, 1, 0, 1), tiny.NewMeasurement([]byte{}, 0, 0, 0, 1, 0, 1)}, tiny.NewPhraseFromBitsAndBytes([]tiny.Bit{1, 0}, 33, 55))
	tester(tiny.NewPhrase(205, 22, 33, 55), []tiny.Bit{1, 1}, tiny.Phrase{tiny.NewMeasurement([]byte{}, 0, 0, 1, 1, 0, 1), tiny.NewMeasurement([]byte{22}), tiny.NewMeasurement([]byte{33}), tiny.NewMeasurement([]byte{}, 0, 0)}, tiny.NewPhraseFromBits(1, 1, 0, 1, 1, 1))
}

func Test_Phrase_FuzzyRead_PowerWindow_ZeroWidth(t *testing.T) {
	defer test.ShouldPanic(t)
	phrase := tiny.NewPhrase(13, 22, 33)
	phrase.FuzzyRead(tiny.Fuzzy.Count(2), tiny.Fuzzy.PowerWindow(0))
}

func Test_Phrase_FuzzyRead_PowerWindow_NegativeWidth(t *testing.T) {
	defer test.ShouldPanic(t)
	phrase := tiny.NewPhrase(13, 22, 33)
	phrase.FuzzyRead(tiny.Fuzzy.Count(2), tiny.Fuzzy.PowerWindow(-1))
}

/**
Encoding
*/

func Test_Passage_FuzzyWrite_ZLE(t *testing.T) {
	tester := func(input tiny.Phrase, expectedKey tiny.Phrase, expectedValue tiny.Phrase) {
		passage := tiny.NewZLEPassage(input)
		ComparePhrases(passage[0], expectedKey, t)
		ComparePhrases(passage[1], expectedValue, t)
	}

	tester(tiny.NewPhrase(77, 22, 33), tiny.NewPhraseFromBits(0, 0, 0, 0, 0, 1), tiny.NewPhrase(0, 77, 22, 33))
	tester(tiny.NewPhrase(77, 22), tiny.NewPhraseFromBits(0, 0, 0, 0, 1), tiny.NewPhrase(77, 22))
	tester(tiny.NewPhrase(77), tiny.NewPhraseFromBits(0, 0, 0, 1), tiny.NewPhrase(77))
	tester(tiny.NewPhraseFromBits(1, 0, 1), tiny.NewPhraseFromBits(0, 0, 1), tiny.NewPhraseFromBits(0, 1, 0, 1))
	tester(tiny.NewPhraseFromBits(0, 1), tiny.NewPhraseFromBits(0, 1), tiny.NewPhraseFromBits(0, 1))
	tester(tiny.NewPhraseFromBits(1), tiny.NewPhraseFromBits(0, 1), tiny.NewPhraseFromBits(0, 1))
	tester(tiny.NewPhraseFromBits(), tiny.NewPhraseFromBits(1), tiny.NewPhraseFromBits())
}

func Test_Passage_FuzzyWrite_ZLE5(t *testing.T) {
	tester := func(input tiny.Measurement, expectedKey tiny.Phrase, expectedValue tiny.Phrase) {
		passage := tiny.NewZLE5Passage(input)
		ComparePhrases(passage[0], expectedKey, t)
		ComparePhrases(passage[1], expectedValue, t)
	}

	tester(tiny.NewMeasurementFromBits(0), tiny.NewPhraseFromBits(1), tiny.Phrase{tiny.NewMeasurementFromBits(0)})
	tester(tiny.NewMeasurementFromBits(1), tiny.NewPhraseFromBits(1), tiny.Phrase{tiny.NewMeasurementFromBits(1)})

	for i := 0; i < 4; i++ {
		v := tiny.NewMeasurementFromBits(tiny.From.Number(i, 2)...)
		tester(v, tiny.NewPhraseFromBits(0, 1), tiny.Phrase{v})
	}

	for i := 0; i < 8; i++ {
		v := tiny.NewMeasurementFromBits(tiny.From.Number(i, 3)...)
		tester(v, tiny.NewPhraseFromBits(0, 0, 1), tiny.Phrase{v})
	}

	for i := 0; i < 16; i++ {
		v := tiny.NewMeasurementFromBits(tiny.From.Number(i, 4)...)
		tester(v, tiny.NewPhraseFromBits(0, 0, 0, 0), tiny.Phrase{v})
	}

	for i := 0; i < 32; i++ {
		v := tiny.NewMeasurementFromBits(tiny.From.Number(i, 5)...)
		tester(v, tiny.NewPhraseFromBits(0, 0, 0, 1), tiny.Phrase{v})
	}
}

func Test_Passage_FuzzyWrite_ZLE5_ShouldPanicAbove5BitInput(t *testing.T) {
	defer test.ShouldPanic(t)
	tiny.NewZLE5Passage(tiny.NewMeasurementFromBits(0, 1, 0, 1, 0, 1))
}

func Test_Passage_FuzzyWrite_ZLE64(t *testing.T) {
	tester := func(input tiny.Phrase, expectedKey tiny.Phrase, expectedValue tiny.Phrase) {
		passage := tiny.NewZLE64Passage(input)
		ComparePhrases(passage[0], expectedKey, t)
		ComparePhrases(passage[1], expectedValue, t)
	}

	tester(tiny.Phrase{}, tiny.NewPhraseFromBits(1), tiny.NewPhraseFromBits(0, 0, 0, 0))
	tester(tiny.NewPhraseFromBits(0, 1, 0), tiny.NewPhraseFromBits(1), tiny.NewPhraseFromBits(0, 0, 1, 0))
	tester(tiny.NewPhrase(77), tiny.NewPhraseFromBits(0, 1), tiny.NewPhrase(77))
	tester(tiny.NewPhrase(77, 22), tiny.NewPhraseFromBits(0, 0, 1), tiny.NewPhrase(77, 22))
	tester(tiny.NewPhrase(77, 22, 33), tiny.NewPhraseFromBits(0, 0, 0, 0), tiny.NewPhrase(0, 77, 22, 33))
	tester(tiny.NewPhrase(77, 22, 33, 55), tiny.NewPhraseFromBits(0, 0, 0, 0), tiny.NewPhrase(77, 22, 33, 55))
	tester(tiny.NewPhrase(77, 22, 33, 55, 11), tiny.NewPhraseFromBits(0, 0, 0, 1), tiny.NewPhrase(0, 0, 0, 77, 22, 33, 55, 11))
	tester(tiny.NewPhrase(77, 22, 33, 55, 11, 88), tiny.NewPhraseFromBits(0, 0, 0, 1), tiny.NewPhrase(0, 0, 77, 22, 33, 55, 11, 88))
	tester(tiny.NewPhrase(77, 22, 33, 55, 11, 88, 44), tiny.NewPhraseFromBits(0, 0, 0, 1), tiny.NewPhrase(0, 77, 22, 33, 55, 11, 88, 44))
	tester(tiny.NewPhrase(77, 22, 33, 55, 11, 88, 44, 99), tiny.NewPhraseFromBits(0, 0, 0, 1), tiny.NewPhrase(77, 22, 33, 55, 11, 88, 44, 99))
}

func Test_Passage_FuzzyWrite_ZLE64_ShouldPanicAbove64BitInput(t *testing.T) {
	defer test.ShouldPanic(t)
	tiny.NewZLE64Passage(tiny.NewPhrase(77, 22, 33, 55, 11, 88, 44, 99, 66))
}

func Test_Passage_FuzzyWrite_ZLEScaled(t *testing.T) {
	tester := func(input int, expectedKey tiny.Phrase, expectedValue tiny.Phrase) {
		passage := tiny.NewZLEScaledPassage(input)
		ComparePhrases(passage[0], expectedKey, t)
		ComparePhrases(passage[1], expectedValue, t)
	}

	tester(0, tiny.NewPhraseFromBits(1), tiny.NewPhraseFromBits(0, 0))
	tester(2, tiny.NewPhraseFromBits(1), tiny.NewPhraseFromBits(1, 0))
	tester(5, tiny.NewPhraseFromBits(0, 1), tiny.NewPhraseFromBits(0, 0, 1))
	tester(77, tiny.NewPhraseFromBits(0, 0, 1), tiny.NewPhrase(65))
	tester(333, tiny.NewPhraseFromBits(0, 0, 0, 0), tiny.NewPhrase(1, 77))
	tester(65536, tiny.NewPhraseFromBits(0, 0, 0, 1), tiny.NewPhrase(0, 0, 0, 0, 0, 1, 0, 0))
}

func Test_Synthesize_FuzzyApproximate(t *testing.T) {
	data := tiny.Synthesize.RandomPhrase(32)
	_, data, _ = data.ReadBit()
	data = data.PrependBits(1)

	approx := tiny.Fuzzy.Approximation(data.AsBigInt(), 3)

	fmt.Println(approx.Indices)
	fmt.Println(approx.Target.Text(2))
	fmt.Println(approx.Value.Text(2))
	fmt.Println(approx.Delta.Text(2))

	fmt.Println(approx.Target)
	fmt.Println(approx.Value)
	fmt.Println(approx.Delta)
	fmt.Println(approx.Relativity)
}

func Test_Synthesize_FuzzyApproximate3(t *testing.T) {
	counts := make([]int, 8)

	for i := 0; i < 256; i++ {
		results := approximate()

		for ii := 0; ii < 8; ii++ {
			counts[ii] += results[ii]
		}
	}

	for ii := 0; ii < 8; ii++ {
		counts[ii] /= 256
		fmt.Printf("[%d] %d\n", ii, counts[ii])
	}
}

func approximate() []int {
	out := make([]int, 8)
	data := tiny.Synthesize.RandomPhrase(32)
	_, data, _ = data.ReadBit()
	data = data.PrependBits(1)

	target := data.AsBigInt()
	bitLen := target.BitLen()

	approx := tiny.Fuzzy.Approximation(target, 1)
	out[0] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 2)
	out[1] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 3)
	out[2] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 4)
	out[3] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 5)
	out[4] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 6)
	out[5] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 7)
	out[6] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 8)
	out[7] = bitLen - approx.Delta.BitLen()
	return out
}
