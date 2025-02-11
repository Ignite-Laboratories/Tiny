package tiny

import "crypto/rand"

type _create struct{}

// Ones creates a slice of '1's of the requested length.
func (_ _create) Ones(count int) []Bit {
	var bits []Bit
	for i := 0; i < count; i++ {
		bits = append(bits, One)
	}
	return bits
}

// Zeros creates a slice of '0's of the requested length.
func (_ _create) Zeros(count int) []Bit {
	return make([]Bit, count)
}

// Repeating creates a slice of the provided pattern repeated the requested number of times.
func (_ _create) Repeating(count int, pattern ...Bit) []Bit {
	var bits []Bit
	for i := 0; i < count; i++ {
		bits = append(bits, pattern...)
	}
	return bits
}

// Pattern creates a slice of the provided pattern repeated and trimmed to the provided length.
func (_ _create) Pattern(length uint64, pattern ...Bit) []Bit {
	count := (length / uint64(len(pattern))) + 1
	repeated := _create{}.Repeating(int(count), pattern...)
	return repeated[:length]
}

// Random creates a random sequence of 1s and 0s.
func (_ _create) Random(width int) []Bit {
	bits := make([]Bit, width)
	for i := 0; i < width; i++ {
		var b [1]byte
		_, _ = rand.Read(b[:])
		bits[i] = Bit(b[0] % 2)
	}
	return bits
}
