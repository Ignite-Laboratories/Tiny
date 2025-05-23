package testing

import (
	"github.com/ignite-laboratories/support/test"
	"github.com/ignite-laboratories/tiny"
	"testing"
)

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
		k, c, r := phrase.FuzzyRead(2, tiny.Fuzzy.SixtyFour)

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
	_, c, _ := phrase.FuzzyRead(12, tiny.Fuzzy.SixtyFour)
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
		k, c, r := phrase.FuzzyRead(2, tiny.Fuzzy.Window(3))
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
	phrase.FuzzyRead(2, tiny.Fuzzy.Window(0))
}

func Test_Phrase_FuzzyRead_Window_NegativeWidth(t *testing.T) {
	defer test.ShouldPanic(t)
	phrase := tiny.NewPhrase(13, 22, 33)
	phrase.FuzzyRead(2, tiny.Fuzzy.Window(-1))
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
		k, c, r := phrase.FuzzyRead(2, tiny.Fuzzy.PowerWindow(3))
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
	phrase.FuzzyRead(2, tiny.Fuzzy.PowerWindow(0))
}

func Test_Phrase_FuzzyRead_PowerWindow_NegativeWidth(t *testing.T) {
	defer test.ShouldPanic(t)
	phrase := tiny.NewPhrase(13, 22, 33)
	phrase.FuzzyRead(2, tiny.Fuzzy.PowerWindow(-1))
}
