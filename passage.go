package tiny

import (
	"fmt"
	"math"
)

// Passage represents a slice of Phrases.
type Passage []Phrase

// NewPassage implicitly casts the provided phrase slice into a passage.
func NewPassage(phrases ...Phrase) Passage {
	return phrases
}

// AppendBitsAsPhrase appends the provided bits to the end of the passage as a new phrase.
func (passage Passage) AppendBitsAsPhrase(bits ...Bit) Passage {
	return append(passage, NewPhraseFromBits(bits...))
}

// Append appends the provided phrases to the end of the passage.
func (passage Passage) Append(phrase ...Phrase) Passage {
	return append(passage, phrase...)
}

// Prepend prepends the provided phrases to the beginning of the passage.
func (passage Passage) Prepend(phrase ...Phrase) Passage {
	return append(phrase, passage...)
}

// NewZLEPassageInt stores the provided number using the below ZLE scheme, where the number of leading zeros
// indicates the power of two worth of bits to read.
// This returns a passage where the first phrase is the key and the second is the value.
//
// NOTE: For 2⁰, a bit range of 0 is returned - not 1.
//
//		ZLE Key | Bit Range To Store
//		      1 | 0
//		    0 1 | 2
//		  0 0 1 | 4
//		0 0 0 1 | 8
//	           ...
//	      𝑛   1 | 2ⁿ
func NewZLEPassageInt(input int) Passage {
	return NewZLEPassage(NewPhraseFromBits(From.Number(input)...))
}

// NewZLEPassage stores the provided bits using the below ZLE scheme, where the number of leading zeros
// indicates the power of two worth of bits to read.
// This returns a passage where the first phrase is the key and the second is the value.
//
// NOTE: For 2⁰, a bit range of 0 is returned - not 1.
//
//		ZLE Key | Bit Range To Store
//		      1 | 0
//		    0 1 | 2
//		  0 0 1 | 4
//		0 0 0 1 | 8
//	           ...
//	      𝑛   1 | 2ⁿ
func NewZLEPassage(input Phrase) Passage {
	zle := NewPhraseFromBits(1)
	bitLen := input.BitLength()

	if bitLen == 0 {
		return Passage{zle, Phrase{}}
	}

	power := 1
	zle = zle.PrependBits(0).Align()

	for math.Pow(2, float64(power)) < float64(bitLen) {
		power++
		zle = zle.PrependBits(0).Align()
	}

	zeroCount := int(math.Pow(2, float64(power)))
	value := make([]Bit, zeroCount-bitLen)           // Pad in the leading zeros
	value = append(value, input.Bits()...)           // Append the input bits
	return Passage{zle, NewPhraseFromBits(value...)} // Combine the ZLE key and the input value
}

// NewZLE5PassageInt stores the provided number using the below ZLE scheme.
// This returns a passage where the first phrase is the key and the second is the value.
//
//	ZLE Key | Bit Range To Store
//	      1 | 1
//	    0 1 | 2
//	  0 0 1 | 3
//	0 0 0 0 | 4
//	0 0 0 1 | 5
func NewZLE5PassageInt(input int) Passage {
	return NewZLE5Passage(NewMeasurementFromBits(From.Number(input)...))
}

// NewZLE5Passage stores the provided measurement using the below ZLE scheme.
// This returns a passage where the first phrase is the key and the second is the value.
//
//	ZLE Key | Bit Range To Store
//	      1 | 1
//	    0 1 | 2
//	  0 0 1 | 3
//	0 0 0 0 | 4
//	0 0 0 1 | 5
func NewZLE5Passage(input Measurement) Passage {
	switch bitLen := input.BitLength(); bitLen {
	case 0:
		return Passage{NewPhraseFromBits(1), NewPhraseFromBits(0)}
	case 1:
		return Passage{NewPhraseFromBits(1), Phrase{input}}
	case 2:
		return Passage{NewPhraseFromBits(0, 1), Phrase{input}}
	case 3:
		return Passage{NewPhraseFromBits(0, 0, 1), Phrase{input}}
	case 4:
		return Passage{NewPhraseFromBits(0, 0, 0, 0), Phrase{input}}
	case 5:
		return Passage{NewPhraseFromBits(0, 0, 0, 1), Phrase{input}}
	}
	panic(fmt.Sprintf("invalid 5-bit ZLE key: %v", input.Bits))
}

// NewZLE64PassageInt stores the provided number using the below ZLE scheme.
// This returns a passage where the first phrase is the key and the second is the value.
//
//	ZLE Key | Bit Range To Store
//	      1 | 4
//	    0 1 | 8
//	  0 0 1 | 16
//	0 0 0 0 | 32
//	0 0 0 1 | 64
func NewZLE64PassageInt(input int) Passage {
	return NewZLE64Passage(NewPhraseFromBits(From.Number(input)...))
}

// NewZLE64Passage stores the provided phrase using the below ZLE scheme.
// This returns a passage where the first phrase is the key and the second is the value.
//
//	ZLE Key | Bit Range To Store
//	      1 | 4
//	    0 1 | 8
//	  0 0 1 | 16
//	0 0 0 0 | 32
//	0 0 0 1 | 64
func NewZLE64Passage(input Phrase) Passage {
	bitLen := input.BitLength()

	if bitLen <= 4 {
		value := make([]Bit, 4-bitLen)
		value = append(value, input.Bits()...)
		return Passage{NewPhraseFromBits(1), NewPhraseFromBits(value...)}
	}
	if bitLen <= 8 {
		value := make([]Bit, 8-bitLen)
		value = append(value, input.Bits()...)
		return Passage{NewPhraseFromBits(0, 1), NewPhraseFromBits(value...)}
	}
	if bitLen <= 16 {
		value := make([]Bit, 16-bitLen)
		value = append(value, input.Bits()...)
		return Passage{NewPhraseFromBits(0, 0, 1), NewPhraseFromBits(value...)}
	}
	if bitLen <= 32 {
		value := make([]Bit, 32-bitLen)
		value = append(value, input.Bits()...)
		return Passage{NewPhraseFromBits(0, 0, 0, 0), NewPhraseFromBits(value...)}
	}
	if bitLen <= 64 {
		value := make([]Bit, 64-bitLen)
		value = append(value, input.Bits()...)
		return Passage{NewPhraseFromBits(0, 0, 0, 1), NewPhraseFromBits(value...)}
	}
	panic(fmt.Sprintf("invalid 64-bit ZLE key: %v", input.Bits))
}

// NewZLEScaledPassage stores the provided phrase using the below ZLE scheme.
// This returns a passage where the first phrase is the key and the second is the value.
//
// This ZLE is unique in that the values are -interpreted- differently than stored
// for the second and third keys (though, practically speaking, they are stored no differently).
//
//	ZLE Key | Bit Range | Value Range
//	      1 |     2     |   0-3
//	    0 1 |     3     |   0-2³ + 4 (4-11)
//	  0 0 1 |     8     |   0-2⁸ + 12 (12-267)
//	0 0 0 0 |    16     |   0-2¹⁶
//	0 0 0 1 |    64     |   0-2⁶⁴
func NewZLEScaledPassage(input int) Passage {
	switch x := input; {
	case x < 4:
		return newZLEScaledPassage(NewPhraseFromBits(From.Number(x, 2)...))
	case x < 12:
		x -= 4
		return newZLEScaledPassage(NewPhraseFromBits(From.Number(x, 3)...))
	case x < 268:
		x -= 12
		return newZLEScaledPassage(NewPhraseFromBits(From.Number(x, 8)...))
	case x < 65536:
		return newZLEScaledPassage(NewPhraseFromBits(From.Number(x, 16)...))
	default:
		return newZLEScaledPassage(NewPhraseFromBits(From.Number(x, 64)...))
	}
}

// newZLEScaledPassage stores the provided phrase using the below ZLE scheme.
// This returns a passage where the first phrase is the key and the second is the value.
//
// This ZLE is unique in that the values are -interpreted- differently than stored
// for the second and third keys (though, practically speaking, they are stored no differently).
//
//	ZLE Key | Bit Range | Value Range
//	      1 |     2     |   0-3
//	    0 1 |     3     |   0-2³ + 4 (4-11)
//	  0 0 1 |     8     |   0-2⁸ + 12 (12-267)
//	0 0 0 0 |    16     |   0-2¹⁶
//	0 0 0 1 |    64     |   0-2⁶⁴
func newZLEScaledPassage(input Phrase) Passage {
	switch bitLen := input.BitLength(); {
	case bitLen <= 2:
		value := make([]Bit, 2-bitLen)
		value = append(value, input.Bits()...)
		return Passage{NewPhraseFromBits(1), NewPhraseFromBits(value...)}
	case bitLen <= 3:
		value := make([]Bit, 3-bitLen)
		value = append(value, input.Bits()...)
		return Passage{NewPhraseFromBits(0, 1), NewPhraseFromBits(value...)}
	case bitLen <= 8:
		value := make([]Bit, 8-bitLen)
		value = append(value, input.Bits()...)
		return Passage{NewPhraseFromBits(0, 0, 1), NewPhraseFromBits(value...)}
	case bitLen <= 16:
		value := make([]Bit, 16-bitLen)
		value = append(value, input.Bits()...)
		return Passage{NewPhraseFromBits(0, 0, 0, 0), NewPhraseFromBits(value...)}
	case bitLen <= 64:
		value := make([]Bit, 64-bitLen)
		value = append(value, input.Bits()...)
		return Passage{NewPhraseFromBits(0, 0, 0, 1), NewPhraseFromBits(value...)}
	default:
		panic(fmt.Sprintf("invalid scaled ZLE key: %v", input.Bits))
	}
}

func (passage Passage) String() string {
	out := make([]string, len(passage))
	for i, p := range passage {
		out[i] = p.String()
	}
	return fmt.Sprintf("%v", out)
}
