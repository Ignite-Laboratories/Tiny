package testing

import (
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Create_Ones(t *testing.T) {
	for i := 0; i < 10; i++ {
		bits := tiny.Create.Ones(i)
		for ii := 0; ii < i; ii++ {
			if bits[ii] != 1 {
				t.Error("Expected all ones")
			}
		}
	}
}

func Test_Create_Zeros(t *testing.T) {
	for i := 0; i < 10; i++ {
		bits := tiny.Create.Zeros(i)
		for ii := 0; ii < i; ii++ {
			if bits[ii] != 0 {
				t.Error("Expected all zeros")
			}
		}
	}
}

func Test_Create_Repeating(t *testing.T) {
	patternTester := func(t *testing.T, pattern ...tiny.Bit) {
		for count := 0; count < 8; count++ {
			bits := tiny.Create.Repeating(count, pattern...)
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

func Test_Create_Random(t *testing.T) {
	for lengthI := 8; lengthI < 10; lengthI++ {
		for testI := 0; testI < 10; testI++ {
			random := tiny.Create.Random(lengthI)

			allZero := true
			allOne := true
			toggle0 := true
			toggle1 := true

			zeroOne := tiny.Zero
			oneZero := tiny.One
			for _, bit := range random {
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
			}
			if allOne {
				t.Error("Expected randomness, got all ones")
			}
			if toggle0 {
				t.Error("Expected randomness, got repeating 01s")
			}
			if toggle1 {
				t.Error("Expected randomness, got repeating 10s")
			}
		}
	}
}
