package testing

import (
	"crypto/rand"
	"encoding/binary"
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Synthesize_ForEach(t *testing.T) {
	script := tiny.From.Bytes(170, 85)
	script = append(script, tiny.From.Bits(0, 1, 1, 0, 1, 0)...)
	phrase := tiny.Synthesize.ForEach(22, func(i int) tiny.Bit {
		return script[i]
	})
	CompareSlices(phrase.Bits(), script, t)
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

func Test_Synthesize_TrailingZeros(t *testing.T) {
	bitLength := 10

	for i := 0; i <= bitLength; i++ {
		phrase := tiny.Synthesize.TrailingZeros(bitLength, i)
		bits := phrase.Bits()

		expected1s := bitLength - i
		expected0s := i
		found1s := 0
		found0s := 0

		gotOnes := false

		for _, b := range bits {
			if b == 1 {
				gotOnes = true
				found1s++
			} else {
				if !gotOnes && expected1s > 0 {
					t.Errorf("Expected ones followed by zeros, got %v", bits)
				}
				found0s++
			}
		}

		if expected0s != found0s {
			t.Errorf("Expected %d zeros, got %d", expected0s, found0s)
		}

		if expected1s != found1s {
			t.Errorf("Expected %d ones, got %d", expected1s, found1s)
		}
	}
}

func Test_Synthesize_Midpoint(t *testing.T) {
	for i := 1; i < 1<<10; i++ {
		phrase := tiny.Synthesize.Midpoint(i)
		bits := phrase.Bits()

		if bits[0] != 1 {
			t.Error("Expected a one in the first position")
		}
		for ii := 1; ii < i; ii++ {
			if bits[ii] != 0 {
				t.Error("Expected all zeros after the initial one.")
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
	defer ShouldPanic(t)
	remainder := tiny.Synthesize.RandomPhrase(1, tiny.GetArchitectureBitWidth()+1)
	remainder.WalkBits(3, func(i int, m tiny.Measurement) {})
}

func Test_Synthesize_RandomPhrase_ShouldPanicWithNegativeMeasurementWidth(t *testing.T) {
	defer ShouldPanic(t)
	remainder := tiny.Synthesize.RandomPhrase(1, -1)
	remainder.WalkBits(3, func(i int, m tiny.Measurement) {})
}

func Test_Synthesize_RandomPhrase_ShouldPanicWith0MeasurementWidth(t *testing.T) {
	defer ShouldPanic(t)
	remainder := tiny.Synthesize.RandomPhrase(1, 0)
	remainder.WalkBits(3, func(i int, m tiny.Measurement) {})
}

func Test_Synthesize_Random_StressTest(t *testing.T) {
	for i := 0; i < 10000; i++ {
		Test_Synthesize_Random(t)
	}
}

func Test_Synthesize_Random_CustomGenerator(t *testing.T) {
	called := false
	synthesize_random(t, func(i int) tiny.Bit {
		called = true
		var v uint64
		binary.Read(rand.Reader, binary.BigEndian, &v)
		return tiny.Bit(v % 2)
	})
	if !called {
		t.Error("Expected custom generator to be called")
	}
}

func Test_Synthesize_Random(t *testing.T) {
	synthesize_random(t, nil)
}

func synthesize_random(t *testing.T, g func(int) tiny.Bit) {
	// There is no way to "test" that 1 or 2 digit binary sets are "random"...it's only four possible values =)
	for lengthI := 3; lengthI < 10; lengthI++ {
		for testI := 0; testI < 10; testI++ {
			phrase := tiny.Synthesize.RandomBits(lengthI, g)
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

func Test_Synthesize_Boundary(t *testing.T) {
	limit := 16
	for i := 0; i <= limit; i++ {
		value := i
		repetend := tiny.Zero
		if i == limit {
			value = limit - 1
			repetend = tiny.One
		}

		bits := tiny.From.Number(value, 3)
		repeating := tiny.Synthesize.Repeating(5, repetend)
		out := tiny.Synthesize.Boundary(bits, repetend, 8)
		expected := tiny.NewPhraseFromBits(bits...).Append(repeating).Align()
		ComparePhrases(out, expected, t)
	}
}

func Test_Synthesize_AllBoundaries(t *testing.T) {
	// Test Logic: We know all the 3-bit boundaries, so just test that it outputs the right ones
	// Genuinely, we'd be doing the same exact thing if we actually validated each bit individually =)
	expected := []tiny.Phrase{
		tiny.NewPhrase(0),
		tiny.NewPhrase(32),
		tiny.NewPhrase(64),
		tiny.NewPhrase(96),
		tiny.NewPhrase(128),
		tiny.NewPhrase(160),
		tiny.NewPhrase(192),
		tiny.NewPhrase(224),
		tiny.NewPhrase(255),
	}
	out := tiny.Synthesize.AllBoundaries(3, 8)
	for i, p := range expected {
		ComparePhrases(p, out[i], t)
	}
}

func Test_Synthesize_AllBoundaries_ExceptDark(t *testing.T) {
	// Test Logic: We know all the 3-bit boundaries, so just test that it outputs the right ones
	// Genuinely, we'd be doing the same exact thing if we actually validated each bit individually =)
	expected := []tiny.Phrase{
		tiny.NewPhrase(0),
		tiny.NewPhrase(32),
		tiny.NewPhrase(64),
		tiny.NewPhrase(96),
		tiny.NewPhrase(128),
		tiny.NewPhrase(160),
		tiny.NewPhrase(192),
		tiny.NewPhrase(224),
	}
	out := tiny.Synthesize.AllBoundaries(3, 8, false)
	for i, p := range expected {
		ComparePhrases(p, out[i], t)
	}
}

func Test_Synthesize_AllBoundaries_ShouldPanicWithNegativeDepth(t *testing.T) {
	defer ShouldPanic(t)
	tiny.Synthesize.AllBoundaries(-1, 8)
}

func Test_Synthesize_AllBoundaries_ShouldPanicWithNegativeWidth(t *testing.T) {
	defer ShouldPanic(t)
	tiny.Synthesize.AllBoundaries(3, -1)
}
