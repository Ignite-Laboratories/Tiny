package tiny

// Build is a way to construct binary slices.
type Build int

// Ones creates a slice of '1's of the requested length.
func (b Build) Ones(count int) []Bit {
	var bits []Bit
	for i := 0; i < count; i++ {
		bits = append(bits, One)
	}
	return bits
}

// Zeros creates a slice of '0's of the requested length.
func (b Build) Zeros(count int) []Bit {
	return make([]Bit, count)
}

// Grey creates a slice of the provided pattern repeated the requested number of times.
func (b Build) Grey(count int, pattern ...Bit) []Bit {
	var bits []Bit
	for i := 0; i < count; i++ {
		bits = append(bits, pattern...)
	}
	return bits
}

// FromBits uses the provided ones and zeros to build a Bit slice.
func (b Build) FromBits(bits ...Bit) []Bit {
	return append([]Bit{}, bits...)
}

// FromInt uses the provided int to build a Bit slice. If no width is provided, the result is
// given in its smallest possible width.  Otherwise, the data is MSB padded with 0s to the
// specified width.
func (b Build) FromInt(value int, width ...int) []Bit {
	if value == 0 {
		return []Bit{Bit(0)}
	}

	var bits []Bit
	for value > 0 {
		bit := Bit(value % 2)    // Get the least significant Bit
		bits = append(bits, bit) // Append the Bit
		value /= 2               // Shift right by dividing by 2
	}

	// Reverse the slice
	for left, right := 0, len(bits)-1; left < right; left, right = left+1, right-1 {
		bits[left], bits[right] = bits[right], bits[left]
	}

	if len(width) > 0 {
		result := make([]Bit, width[0]-len(bits))
		result = append(result, bits...)
		return result
	}

	return bits
}

// FromBytes uses the provided slice of bytes to build a Bit slice.
func (b Build) FromBytes(bytes ...byte) []Bit {
	var output []Bit
	for _, v := range bytes {
		output = append(output, b.FromInt(int(v), 8)...)
	}
	return output
}
