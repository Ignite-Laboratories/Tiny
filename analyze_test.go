package tiny

import (
	"testing"
)

func TestOneDistribution_FullSpectrum(t *testing.T) {
	// With this arrangement each byte has one more one than the last in each index
	// This means the count of 1s per index should be i+1
	bytes := []byte{0, 1, 3, 7, 15, 31, 63, 127, 255}
	ones := Analyze.OneDistribution(bytes...)
	for i, index := range ones {
		if index != i+1 {
			t.Errorf("Invalid one count")
		}
	}
}

func TestOneDistribution_Light(t *testing.T) {
	bytes := []byte{0, 0, 0, 0, 0}
	ones := Analyze.OneDistribution(bytes...)
	for _, index := range ones {
		if index != 0 {
			t.Errorf("Should not have any ones")
		}
	}
}

func TestOneDistribution_Dark(t *testing.T) {
	bytes := []byte{255, 255, 255, 255, 255}
	ones := Analyze.OneDistribution(bytes...)
	for _, index := range ones {
		if index != 5 {
			t.Errorf("Should have exactly 5 ones")
		}
	}
}

func TestOneDistribution_Grey(t *testing.T) {
	bytes := []byte{42, 22, 88, 222, 133}
	expected := []int{2, 2, 1, 3, 3, 3, 3, 1}
	ones := Analyze.OneDistribution(bytes...)
	for i := 0; i < len(expected); i++ {
		if ones[i] != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], ones[i])
		}
	}
}
