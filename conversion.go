package tiny

import (
	"log"
	"strconv"
)

// ToBytes takes in binary data and returns it in Remainder form.
func ToBytes(bits []Bit) Remainder {
	// The resulting slice of bytes
	var bytes []byte

	// Process bits in groups of 8 (1 byte)
	for i := 0; i+7 < len(bits); i += 8 {
		var b byte
		for j := 0; j < 8; j++ {
			b = (b << 1) | byte(bits[i+j]) // Shift left and add the next bit
		}
		bytes = append(bytes, b) // Add the full byte to the result
	}

	// Extract any leftover bits that don't form a complete byte
	remainingBits := bits[len(bytes)*8:]

	return Remainder{bytes, remainingBits}
}

// ToByte takes in binary data and returns its byte.
func ToByte(bits []Bit) byte {
	if len(bits) > 8 {
		log.Fatalf("input must contain less than 8 bits")
	}

	result := byte(0)
	padding := 8 - len(bits) // Calculate left padding for smaller slices

	for i, bit := range bits {
		// Shift each bit to its correct position considering the padding
		result |= byte(bit) << (7 - (i + padding))
	}
	return result
}

// BytesToBits takes a slice of bytes and returns a slice of all of its individual bits.
func BytesToBits(data []byte) []Bit {
	bits := make([]Bit, 0, len(data)*8)
	for _, b := range data {
		bits = append(bits, ToBitsFixedWidth(int(b), 8)...)
	}
	return bits
}

// ToBits takes an int and returns its minimum constituent bits.
func ToBits(value int) []Bit {
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

	return bits
}

// ToBitsFixedWidth takes an int and returns its constituent bits, prepended with 0 to the desired width.
func ToBitsFixedWidth(value int, width int) []Bit {
	bits := ToBits(value)
	result := make([]Bit, width-len(bits))
	result = append(result, bits...)
	return result
}

// Bits uses the provided value to build a 1 Bit slice.
func (v Bit) Bits() []Bit {
	return ToBitsFixedWidth(int(v), WidthBit)
}

// Bits uses the provided value to build a 2 Bit slice.
func (v Crumb) Bits() []Bit {
	return ToBitsFixedWidth(int(v), WidthCrumb)
}

// Bits uses the provided value to build a 3 Bit slice.
func (v Note) Bits() []Bit {
	return ToBitsFixedWidth(int(v), WidthNote)
}

// Bits uses the provided value to build a 4 Bit slice.
func (v Nibble) Bits() []Bit {
	return ToBitsFixedWidth(int(v), WidthNibble)
}

// Bits uses the provided value to build a 5 Bit slice.
func (v Flake) Bits() []Bit {
	return ToBitsFixedWidth(int(v), WidthFlake)
}

// Bits uses the provided value to build a 6 Bit slice.
func (v Morsel) Bits() []Bit {
	return ToBitsFixedWidth(int(v), WidthMorsel)
}

// Bits uses the provided value to build a 7 Bit slice.
func (v Shred) Bits() []Bit {
	return ToBitsFixedWidth(int(v), WidthShred)
}

// FromByte returns a byte's constituent bits.
func FromByte(b byte) []Bit {
	return ToBitsFixedWidth(int(b), WidthByte)
}

// ToString converts a set of Bit values into a string.
func ToString[T SubByte](bits []T) string {
	output := ""
	for i := 0; i < len(bits); i++ {
		output += strconv.Itoa(int(bits[i]))
	}
	return output
}

/**
String()
*/

// String converts a Bit to a 1-bit string.
func (v Bit) String() string {
	return ToString([]Bit{v})
}

// String converts a Crumb to a 2-bit string.
func (v Crumb) String() string {
	return ToString(ToBitsFixedWidth(int(v), WidthCrumb))
}

// String converts a Note to a 3-bit string.
func (v Note) String() string {
	return ToString(ToBitsFixedWidth(int(v), WidthNote))
}

// String converts a Nibble to a 4-bit string.
func (v Nibble) String() string {
	return ToString(ToBitsFixedWidth(int(v), WidthNibble))
}

// String converts a Flake to a 5-bit string.
func (v Flake) String() string {
	return ToString(ToBitsFixedWidth(int(v), WidthFlake))
}

// String converts a Morsel to a 6-bit string.
func (v Morsel) String() string {
	return ToString(ToBitsFixedWidth(int(v), WidthMorsel))
}

// String converts a Shred to a 7-bit string.
func (v Shred) String() string {
	return ToString(ToBitsFixedWidth(int(v), WidthShred))
}
