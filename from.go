package tiny

type _from int

// Bits uses the provided ones and zeros to build a Bit slice.
func (_ _from) Bits(bits ...Bit) []Bit {
	return append([]Bit{}, bits...)
}

// Crumb uses the provided value to build a 2-bit slice.
func (f _from) Crumb(value Crumb) []Bit {
	return f.Number(int(value), WidthCrumb)
}

// Note uses the provided value to build a 3-bit slice.
func (f _from) Note(value Note) []Bit {
	return f.Number(int(value), WidthNote)
}

// Nibble uses the provided value to build a 4-bit slice.
func (f _from) Nibble(value Nibble) []Bit {
	return f.Number(int(value), WidthNibble)
}

// Flake uses the provided value to build a 5-bit slice.
func (f _from) Flake(value Flake) []Bit {
	return f.Number(int(value), WidthFlake)
}

// Morsel uses the provided value to build a 6-bit slice.
func (f _from) Morsel(value Morsel) []Bit {
	return f.Number(int(value), WidthMorsel)
}

// Shred uses the provided value to build a 7-bit slice.
func (f _from) Shred(value Shred) []Bit {
	return f.Number(int(value), WidthShred)
}

// Byte uses the provided value to build a 8-bit slice.
func (f _from) Byte(value byte) []Bit {
	return f.Number(int(value), WidthByte)
}

// Bytes uses the provided slice of bytes to build a Bit slice.
func (f _from) Bytes(bytes ...byte) []Bit {
	var output []Bit
	for _, v := range bytes {
		output = append(output, f.Number(int(v), 8)...)
	}
	return output
}

// Number uses the provided int to build a tiny.Bit slice. If no width is provided, the result is
// given in its smallest possible width.  Otherwise, the data is MSB padded with 0s to the
// specified width.
func (_ _from) Number(value int, width ...int) []Bit {
	if value == 0 {
		if width == nil {
			width = []int{1}
		}
		return make([]Bit, width[0])
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
