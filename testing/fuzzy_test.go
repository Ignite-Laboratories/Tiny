package testing

import (
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
	//  |         77          |        22       |         33      |
	//  | 0 1 | 0 0 | 1 1 0 1 - 0 0 0 1 0 1 1 0 - 0 0 1 0 0 0 0 1 |  <- Raw Bits
	//  | Key |  C  |                 Remainder                   | <- Fuzzy Read

	test := func(phrase tiny.Phrase, keyBits []tiny.Bit, continuationBits []tiny.Bit, remainderBits []tiny.Bit) {
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

	test(tiny.NewPhrase(13, 22, 33), []tiny.Bit{0, 0}, []tiny.Bit{}, []tiny.Bit{0, 0, 1, 1, 0, 1})
	test(tiny.NewPhrase(77, 22, 33), []tiny.Bit{0, 1}, []tiny.Bit{0, 0}, []tiny.Bit{1, 1, 0, 1})
	test(tiny.NewPhrase(141, 22, 33), []tiny.Bit{1, 0}, []tiny.Bit{0, 0, 1, 1}, []tiny.Bit{0, 1})
	test(tiny.NewPhrase(205, 22, 33), []tiny.Bit{1, 1}, []tiny.Bit{0, 0, 1, 1, 0, 1}, []tiny.Bit{})
}
