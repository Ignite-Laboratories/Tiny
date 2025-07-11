package tiny

import "math/big"

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

// Run uses the provided value to build a 10-bit slice.
func (f _from) Run(value Run) []Bit {
	return f.Number(int(value), WidthRun)
}

// Scale uses the provided value to build a 12-bit slice.
func (f _from) Scale(value Scale) []Bit {
	return f.Number(int(value), WidthScale)
}

// Motif uses the provided value to build a 16-bit slice.
func (f _from) Motif(value Motif) []Bit {
	return f.Number(int(value), WidthMotif)
}

// Riff uses the provided value to build a 24-bit slice.
func (f _from) Riff(value Riff) []Bit {
	return f.Number(int(value), WidthRiff)
}

// Cadence uses the provided value to build a 32-bit slice.
func (f _from) Cadence(value Cadence) []Bit {
	return f.Number(int(value), WidthCadence)
}

// Hook uses the provided value to build a 48-bit slice.
func (f _from) Hook(value Hook) []Bit {
	return f.Number(int(value), WidthHook)
}

// Bytes uses the provided slice of bytes to build a Bit slice.
func (f _from) Bytes(bytes ...byte) []Bit {
	var output []Bit
	for _, v := range bytes {
		output = append(output, f.Number(int(v), 8)...)
	}
	return output
}

// BigInt returns the bits of the provided big.Int padded to the specified width with zeros.
//
// If no width is provided, the result is given in its smallest possible width.
func (_ _from) BigInt(value *big.Int, width ...int) []Bit {
	str := value.Text(2)
	out := make([]Bit, len(str))

	for i := 0; i < len(str); i++ {
		if str[i] == '1' {
			out[i] = One
		} else {
			out[i] = Zero
		}
	}

	if len(width) > 0 {
		pad := make([]Bit, width[0]-len(out))
		out = append(pad, out...)
	}
	return out
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
		if width[0] < len(bits) {
			return bits[:width[0]]
		} else {
			result := make([]Bit, width[0]-len(bits))
			result = append(result, bits...)
			return result
		}
	}

	return bits
}
