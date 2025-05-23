package testing

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/ignite-laboratories/support/test"
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Synthesize_ForEach(t *testing.T) {
	script := tiny.From.Bytes(170, 85)
	script = append(script, tiny.From.Bits(0, 1, 1, 0, 1, 0)...)
	measure := tiny.Synthesize.ForEach(22, func(i int) tiny.Bit {
		return script[i]
	})
	test.CompareSlices(measure.GetAllBits(), script, t)
}

func Test_Synthesize_Ones(t *testing.T) {
	for i := 0; i < 10; i++ {
		measure := tiny.Synthesize.Ones(i)
		bits := measure.GetAllBits()
		for ii := 0; ii < i; ii++ {
			if bits[ii] != 1 {
				t.Error("Expected all ones")
			}
		}
	}
}

func Test_Synthesize_Zeros(t *testing.T) {
	for i := 0; i < 10; i++ {
		measure := tiny.Synthesize.Zeros(i)
		bits := measure.GetAllBits()
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
			measure := tiny.Synthesize.Repeating(count, pattern...)
			bits := measure.GetAllBits()
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
	unevenM := tiny.Synthesize.Pattern(8, tiny.From.Bits(0, 1, 1)...)
	uneven := unevenM.GetAllBits()
	expectedUneven := tiny.From.Bits(0, 1, 1, 0, 1, 1, 0, 1)
	evenM := tiny.Synthesize.Pattern(9, tiny.From.Bits(0, 1, 1)...)
	even := evenM.GetAllBits()
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
			measure := tiny.Synthesize.Random(lengthI, g)
			bits := measure.GetAllBits()

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
