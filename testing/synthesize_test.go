package testing

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/ignite-laboratories/support/test"
	"github.com/ignite-laboratories/tiny"
	"math/big"
	"testing"
)

func Test_Synthesize_ForEach(t *testing.T) {
	script := tiny.From.Bytes(170, 85)
	script = append(script, tiny.From.Bits(0, 1, 1, 0, 1, 0)...)
	phrase := tiny.Synthesize.ForEach(22, func(i int) tiny.Bit {
		return script[i]
	})
	test.CompareSlices(phrase.Bits(), script, t)
}

func Test_Synthesize_Ones(t *testing.T) {
	for i := 0; i < 10; i++ {
		phrase := tiny.Synthesize.Ones(i)
		bits := phrase.Bits()
		for ii := 0; ii < i; ii++ {
			if bits[ii] != 1 {
				t.Error("Expected all ones")
			}
		}
	}
}

func Test_Synthesize_Zeros(t *testing.T) {
	for i := 0; i < 10; i++ {
		phrase := tiny.Synthesize.Zeros(i)
		bits := phrase.Bits()
		for ii := 0; ii < i; ii++ {
			if bits[ii] != 0 {
				t.Error("Expected all zeros")
			}
		}
	}
}

func Test_Synthesize_Repeating(t *testing.T) {
	patternTester := func(t *testing.T, pattern ...tiny.Bit) {
		for count := 0; count < 8; count++ {
			phrase := tiny.Synthesize.Repeating(count, pattern...)
			bits := phrase.Bits()
			for i := 0; i < count; i++ {
				offset := i * len(pattern)
				for patternI := 0; patternI < len(pattern); patternI++ {
					if bits[offset+patternI] != pattern[patternI] {
						t.Error("Expected repeating pattern")
					}
				}
			}
		}
	}

	patternTester(t, tiny.From.Bits()...)
	patternTester(t, tiny.From.Bits(0)...)
	patternTester(t, tiny.From.Bits(1)...)
	patternTester(t, tiny.From.Bits(0, 1, 1)...)
	patternTester(t, tiny.From.Bits(1, 0, 0)...)
	patternTester(t, tiny.From.Bits(1, 0, 1, 0, 0)...)
}

func Test_Synthesize_Pattern(t *testing.T) {
	unevenP := tiny.Synthesize.Pattern(8, tiny.From.Bits(0, 1, 1)...)
	uneven := unevenP.Bits()
	expectedUneven := tiny.From.Bits(0, 1, 1, 0, 1, 1, 0, 1)
	evenP := tiny.Synthesize.Pattern(9, tiny.From.Bits(0, 1, 1)...)
	even := evenP.Bits()
	expectedEven := tiny.From.Bits(0, 1, 1, 0, 1, 1, 0, 1, 1)

	for i, bit := range uneven {
		if bit != expectedUneven[i] {
			t.Errorf("Expected %d, got %d", expectedUneven[i], bit)
		}
	}
	for i, bit := range even {
		if bit != expectedEven[i] {
			t.Errorf("Expected %d, got %d", expectedEven[i], bit)
		}
	}
}

func Test_Synthesize_RandomPhrase_ShouldPanicWithLargeMeasurementWidth(t *testing.T) {
	defer test.ShouldPanic(t)
	remainder := tiny.Synthesize.RandomPhrase(1, tiny.MaxMeasurementBitLength+1)
	remainder.WalkBits(3, func(i int, m tiny.Measurement) {})
}

func Test_Synthesize_RandomPhrase_ShouldPanicWithNegativeMeasurementWidth(t *testing.T) {
	defer test.ShouldPanic(t)
	remainder := tiny.Synthesize.RandomPhrase(1, -1)
	remainder.WalkBits(3, func(i int, m tiny.Measurement) {})
}

func Test_Synthesize_RandomPhrase_ShouldPanicWith0MeasurementWidth(t *testing.T) {
	defer test.ShouldPanic(t)
	remainder := tiny.Synthesize.RandomPhrase(1, 0)
	remainder.WalkBits(3, func(i int, m tiny.Measurement) {})
}

func Test_Synthesize_Random_StressTest(t *testing.T) {
	for i := 0; i < 10000; i++ {
		Test_Synthesize_Random(t)
	}
}

func Test_Synthesize_Random_CustomGenerator(t *testing.T) {
	synthesize_random(t, func(i int) tiny.Bit {
		var v uint64
		binary.Read(rand.Reader, binary.BigEndian, &v)
		return tiny.Bit(v % 2)
	})
}

func Test_Synthesize_Random(t *testing.T) {
	synthesize_random(t, nil)
}

func synthesize_random(t *testing.T, g func(int) tiny.Bit) {
	// There is no way to "test" that 1 or 2 digit binary sets are "random"...it's only four possible values =)
	for lengthI := 3; lengthI < 10; lengthI++ {
		for testI := 0; testI < 10; testI++ {
			phrase := tiny.Synthesize.Random(lengthI, g)
			bits := phrase.Bits()

			allZero := true
			allOne := true
			toggle0 := true
			toggle1 := true

			zeroOne := tiny.Zero
			oneZero := tiny.One
			for _, bit := range bits {
				if bit == 0 {
					allOne = false
				}
				if bit == 1 {
					allZero = false
				}
				if zeroOne != bit {
					toggle0 = false
				}
				if oneZero != bit {
					toggle1 = false
				}

				zeroOne ^= 1
				oneZero ^= 1
			}
			if allZero {
				t.Error("Expected randomness, got all zeros")
				t.FailNow()
			}
			if allOne {
				t.Error("Expected randomness, got all ones")
				t.FailNow()
			}
			if toggle0 {
				t.Error("Expected randomness, got repeating 01s")
				t.FailNow()
			}
			if toggle1 {
				t.Error("Expected randomness, got repeating 10s")
				t.FailNow()
			}
		}
	}
}

func Test_Synthesize_Subdivided_ShouldAcceptNegativeIndex(t *testing.T) {
	result := tiny.Synthesize.Subdivided(8, -1, 7)
	test.CompareSlices(result, tiny.From.Number(0, 8), t)
}

func Test_Synthesize_Subdivided_ShouldAcceptOversizedIndex(t *testing.T) {
	result := tiny.Synthesize.Subdivided(8, 8, 7)
	test.CompareSlices(result, tiny.From.Number(255, 8), t)
}

func Test_Synthesize_Subdivided_Byte(t *testing.T) {
	expected := [][]tiny.Bit{
		tiny.From.Number(0, 8),
		tiny.From.Number(36, 8),
		tiny.From.Number(72, 8),
		tiny.From.Number(109, 8),
		tiny.From.Number(145, 8),
		tiny.From.Number(182, 8),
		tiny.From.Number(218, 8),
		tiny.From.Number(255, 8),
	}

	for i := 0; i < 8; i++ {
		result := tiny.Synthesize.Subdivided(8, i, 7)
		test.CompareSlices(result, expected[i], t)
	}
}

func Test_Synthesize_Subdivided_Int16(t *testing.T) {
	expected := [][]tiny.Bit{
		tiny.From.Number(0, 16),
		tiny.From.Number(9362, 16),
		tiny.From.Number(18724, 16),
		tiny.From.Number(28086, 16),
		tiny.From.Number(37448, 16),
		tiny.From.Number(46810, 16),
		tiny.From.Number(56172, 16),
		tiny.From.Number(65535, 16),
	}

	for i := 0; i < 8; i++ {
		result := tiny.Synthesize.Subdivided(16, i, 7)
		test.CompareSlices(result, expected[i], t)
	}
}

func Test_Synthesize_Approximate_42(t *testing.T) {
	expectedBits := tiny.Synthesize.Pattern(42, 1, 0, 0)
	expectedIndex := 4
	result, index := tiny.Synthesize.Approximate(tiny.Synthesize.Pattern(42, 1, 0).AsBigInt(), 7)
	test.CompareSlices(result, expectedBits.Bits(), t)
	if index != expectedIndex {
		t.Fatalf("Expected index %d, got %d", expectedIndex, index)
	}
}

func Test_Synthesize_Approximate_Byte(t *testing.T) {
	expected := [][]tiny.Bit{
		tiny.From.Number(0, 8),
		tiny.From.Number(36, 8),
		tiny.From.Number(72, 8),
		tiny.From.Number(109, 8),
		tiny.From.Number(145, 8),
		tiny.From.Number(182, 8),
		tiny.From.Number(218, 8),
		tiny.From.Number(255, 8),
	}

	expectedIndex := -1
	for i := int64(0); i < 256; i++ {
		expectedIndex = 0
		if i > 36 {
			expectedIndex = 1
		}
		if i > 72 {
			expectedIndex = 2
		}
		if i > 109 {
			expectedIndex = 3
		}
		if i > 145 {
			expectedIndex = 4
		}
		if i > 182 {
			expectedIndex = 5
		}
		if i > 218 {
			expectedIndex = 6
		}
		if i == 255 {
			expectedIndex = 7
		}

		result, index := tiny.Synthesize.Approximate(big.NewInt(i), 7, 8)

		if index != expectedIndex {
			t.Fatalf("Expected index %d, got %d", expectedIndex, index)
		}
		test.CompareSlices(result, expected[index], t)
	}
}
