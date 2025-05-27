package tiny

import (
	"fmt"
	"math"
)

// FuzzyHandler is a factory for creating or referencing fuzzy projection functions.
type FuzzyHandler int

// Count returns a function that will return true the requested number of times.
func (_ FuzzyHandler) Count(value int) func(Bit) bool {
	i := 0
	return func(b Bit) bool {
		i++
		return i < value
	}
}

// WhileZero returns true until the value of 1 is reached.
func (_ FuzzyHandler) WhileZero(b Bit) bool {
	return b == Zero
}

// WhileOne returns true until the value of 0 is reached.
func (_ FuzzyHandler) WhileOne(b Bit) bool {
	return b == Zero
}

// ZLEKey reads up to four bits or until a value of 1 is reached.
// This will yield a Zero Length Encoding key that can be parsed using FuzzyHandler.ParseZLE64
//
// If you would like to read a ZLE key longer than 4 bits, you may provide an upper limit.
//
// If you wish for no upper limit (just read until EOD or a 1) then provide <= 0 as the upper limit..
func (_ FuzzyHandler) ZLEKey(upperLimit ...int) func(Bit) bool {
	limit := 4
	if len(upperLimit) > 0 {
		limit = upperLimit[0]
	}

	i := 0
	return func(b Bit) bool {
		i++
		if limit <= 0 {
			return b == Zero
		}
		return b == Zero && i < limit
	}
}

// ParseZLEScaled uses the provided Zero Length Encoding key to calculate how many more bits to read.
//
// This particular flavor of ZLE will yield an addressable range up to 64 bits wide while prioritizing
// a minimal number of bits under that length.
//
// NOTE: This function merely gives the bit ranges for each key entry.
//
//	ZLE Key | Bit Range | Value Range
//	      1 |     2     |   0-3
//	    0 1 |     3     |   0-2³ + 4 (4-11)
//	  0 0 1 |     8     |   0-2⁸ + 12 (12-267)
//	0 0 0 0 |    16     |   0-2¹⁶
//	0 0 0 1 |    64     |   0-2⁶⁴
func (_ FuzzyHandler) ParseZLEScaled(key Measurement) int {
	switch bits := key.Bits; {
	case len(bits) == 1 && key.Value() == 1:
		return 2
	case len(bits) == 2 && key.Value() == 1:
		return 3
	case len(bits) == 3 && key.Value() == 1:
		return 8
	case len(bits) == 4 && key.Value() == 0:
		return 16
	case len(bits) == 4 && key.Value() == 1:
		return 64
	default:
		panic(fmt.Sprintf("invalid scaled ZLE key: %v", key.Bits))
	}
}

// InterpretZLEScaled returns the -interpreted- value of a scaled ZLE phrase.
//
//	ZLE Key | Bit Range | Value Range
//	      1 |     2     |   0-3
//	    0 1 |     3     |   0-2³ + 4 (4-11)
//	  0 0 1 |     8     |   0-2⁸ + 12 (12-267)
//	0 0 0 0 |    16     |   0-2¹⁶
//	0 0 0 1 |    64     |   0-2⁶⁴
func (_ FuzzyHandler) InterpretZLEScaled(passage Passage) int {
	key := passage[0][0]
	projection := passage[1]
	switch bits := key.Bits; {
	case len(bits) == 1 && key.Value() == 1:
		return To.Number(2, projection.Bits()...)
	case len(bits) == 2 && key.Value() == 1:

		return To.Number(3, projection.Bits()...) + 4
	case len(bits) == 3 && key.Value() == 1:
		return To.Number(8, projection.Bits()...) + 12
	case len(bits) == 4 && key.Value() == 0:
		return To.Number(16, projection.Bits()...)
	case len(bits) == 4 && key.Value() == 1:
		return int(projection.AsBigInt().Int64())
	default:
		panic(fmt.Sprintf("invalid scaled ZLE key: %v", key.Bits))
	}
}

// ParseZLE64 uses the provided Zero Length Encoding key to calculate how many more bits to read.
//
// This particular flavor of ZLE will yield an addressable range up to 64 bits wide.
//
//	ZLE Key | Bit Range
//	      1 | 4
//	    0 1 | 8
//	  0 0 1 | 16
//	0 0 0 0 | 32
//	0 0 0 1 | 64
func (_ FuzzyHandler) ParseZLE64(key Measurement) int {
	switch bits := key.Bits; {
	case len(bits) == 1 && key.Value() == 1:
		return 4
	case len(bits) == 2 && key.Value() == 1:
		return 8
	case len(bits) == 3 && key.Value() == 1:
		return 16
	case len(bits) == 4 && key.Value() == 0:
		return 32
	case len(bits) == 4 && key.Value() == 1:
		return 64
	default:
		panic(fmt.Sprintf("invalid ZLE key: %v", key.Bits))
	}
}

// ParseZLE5 uses the provided Zero Length Encoding key to calculate how many more bits to read.
//
// This particular flavor of ZLE will yield an addressable range up to 5 bits wide.
//
//	ZLE Key | Bit Range
//	      1 | 1
//	    0 1 | 2
//	  0 0 1 | 3
//	0 0 0 0 | 4
//	0 0 0 1 | 5
func (_ FuzzyHandler) ParseZLE5(key Measurement) int {
	switch bits := key.Bits; {
	case len(bits) == 1 && key.Value() == 1:
		return 1
	case len(bits) == 2 && key.Value() == 1:
		return 2
	case len(bits) == 3 && key.Value() == 1:
		return 3
	case len(bits) == 4 && key.Value() == 0:
		return 4
	case len(bits) == 4 && key.Value() == 1:
		return 5
	default:
		panic(fmt.Sprintf("invalid micro ZLE key: %v", key.Bits))
	}
}

// ParseZLE uses the provided Zero Length Encoding key to calculate how many more bits to read.
//
// This returns back 2ⁿ - where 𝑛 is the number of zeros found.
//
// NOTE: This will overflow if you let it read too far =)
//
//		ZLE Key | Bit Range
//		      1 | 0
//		    0 1 | 2
//		  0 0 1 | 4
//		0 0 0 1 | 8
//	           ...
//	      𝑛   1 | 2ⁿ
func (_ FuzzyHandler) ParseZLE(key Measurement) int {
	count := 0
	for _, b := range key.Bits {
		if b == Zero {
			count++
		} else {
			break
		}
	}
	if count == 0 {
		return 0
	}
	return int(math.Pow(2, float64(count)))
}

// SixtyFour uses the key Measurement value to calculate a bit range of up to six bits, yielding 64 unique values.
//
// NOTE: This will still return a bit length of 6 if provided a key value greater than 3.
//
//	Key | Bit Range
//	  0 | 0
//	  1 | 2
//	  2 | 4
//	  3 | 6
func (_ FuzzyHandler) SixtyFour(key Measurement) int {
	switch v := key.Value(); v {
	case 0:
		return 0
	case 1:
		return 2
	case 2:
		return 4
	case 3:
		return 6
	default:
		return 6
	}
}

// Window creates a fuzzy projection function using a window width.
// The window width is multiplied against the value of the key measurement (plus one) to build
// the continuation projection range.
//
// NOTE: This will panic if provided a window width <= 0.
//
// For example:
//
//		 FuzzyRead(2, tiny.Fuzzy.Window(3))
//
//	Value-> 1  |   1                                             <- Window Occurances
//	     | 0 0 | 0 0 1 | 1 0 1 0 0 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | <- Raw bits
//	     | Key | Cont  |            Remainder                  | <- Fuzzy read
//
//	Value-> 2  |   1       2                                       <- Window Occurances
//	     | 0 1 | 0 0 1 - 1 0 1 | 0 0 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | <- Raw bits
//	     | Key |  Continuation |         Remainder               | <- Fuzzy read
//
//	Value-> 3  |   1       2       3                                 <- Window Occurances
//	     | 1 0 | 0 0 1 - 1 0 1 - 0 0 0 | 1 0 1 1 0 0 0 1 0 0 0 0 1 | <- Raw bits
//	     | Key |     Continuation      |        Remainder          | <- Fuzzy read
//
//	Value-> 4  |   1       2       3       4                           <- Window Occurances
//	     | 1 1 | 0 0 1 - 1 0 1 - 0 0 0 - 1 0 1 | 1 0 0 0 1 0 0 0 0 1 | <- Raw bits
//	     | Key |          Continuation         |      Remainder      | <- Fuzzy read
func (_ FuzzyHandler) Window(windowWidth int) func(Measurement) int {
	if windowWidth <= 0 {
		panic("fuzzy.Window: window width must be greater than zero")
	}

	return func(key Measurement) int {
		return (key.Value() + 1) * windowWidth
	}
}

// PowerWindow creates a fuzzy projection function using a window width.
// The window width is multiplied against the value of the key measurement to build
// the continuation projection range.
// This differs from a standard Window operation in that the key value is considered to be a
// 'power of 2', providing further projection into the data (with 0 being a single occurrance.)
//
// NOTE: This will panic if provided a window width <= 0.
//
// For example:
//
//		 FuzzyRead(2, tiny.Fuzzy.PowerWindow(2))
//
//	Value-> 1  |  1                                              <- Window Occurances
//	     | 0 0 | 0 0 | 1 1 0 1 0 0 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | <- Raw bits
//	     | Key |  C  |             Remainder                   | <- Fuzzy read
//
//	Value-> 2  |  1     2                                          <- Window Occurances
//	     | 0 1 | 0 0 - 1 1 | 0 1 0 0 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | <- Raw bits
//	     | Key |   Cont    |           Remainder                 | <- Fuzzy read
//
//	Value-> 4  |  1     2     3     4                                  <- Window Occurances
//	     | 1 0 | 0 0 - 1 1 - 0 1 - 0 0 | 0 1 0 1 1 0 0 0 1 0 0 0 0 1 | <- Raw bits
//	     | Key |      Continuation     |          Remainder          | <- Fuzzy read
//
//	Value-> 8  |  1     2     3     4     5     6     7     8                  <- Window Occurances
//	     | 1 1 | 0 0 - 1 1 - 0 1 - 0 0 - 0 1 - 0 1 - 1 0 - 0 0 | 1 0 0 0 0 1 | <- Raw bits
//	     | Key |                  Continuation                 |  Remainder  | <- Fuzzy read
//
// NOTE: The output continuation phrase will be aligned to its source form, but a call to Phrase.Align(windowWidth)
// will yield even measurements as demonstrated above.
func (_ FuzzyHandler) PowerWindow(windowWidth int) func(Measurement) int {
	if windowWidth <= 0 {
		panic("fuzzy.Window: window width must be greater than zero")
	}

	return func(key Measurement) int {
		v := key.Value()
		if v == 0 {
			return windowWidth
		}
		return int(math.Pow(2, float64(v))) * windowWidth
	}
}
