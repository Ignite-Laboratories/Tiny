package tiny

import "crypto/rand"

// create is a way to synthesize binary slices from known parameters.
type create struct{}

// Ones creates a slice of '1's of the requested length.
func (_ create) Ones(count int) []Bit {
	var bits []Bit
	for i := 0; i < count; i++ {
		bits = append(bits, One)
	}
	return bits
}

// Zeros creates a slice of '0's of the requested length.
func (_ create) Zeros(count int) []Bit {
	return make([]Bit, count)
}

// Grey creates a slice of the provided pattern repeated the requested number of times.
func (_ create) Grey(count int, pattern ...Bit) []Bit {
	var bits []Bit
	for i := 0; i < count; i++ {
		bits = append(bits, pattern...)
	}
	return bits
}

// Random creates a random sequence of 1s and 0s.
func (_ create) Random(width int) []Bit {
	bits := make([]Bit, width)
	for i := 0; i < width; i++ {
		var b [1]byte
		_, _ = rand.Read(b[:])
		bits[i] = Bit(b[0] % 2)
	}
	return bits
}
