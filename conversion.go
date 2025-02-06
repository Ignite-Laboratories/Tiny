package tiny

import "log"

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

// ToBitsFixedWidth takes an int and returns its constituent bits, prepended with 0 to the desired width.
func ToBitsFixedWidth(value int, width int) []Bit {
	bits := ToBits(value)
	result := make([]Bit, width-len(bits))
	result = append(result, bits...)
	return result
}

// ToBits takes an int and returns its constituent bits.
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

// FromByte takes a byte and returns its constituent bits.
func FromByte(b byte) []Bit {
	// Yes, this is a shorthand convenience method - sue me =)
	return ToBitsFixedWidth(int(b), 8)
}

// BytesToBits takes a slice of bytes and returns a slice of all of its individual bits.
func BytesToBits(data []byte) []Bit {
	bits := make([]Bit, 0, len(data)*8)
	for _, b := range data {
		bits = append(bits, ToBitsFixedWidth(int(b), 8)...)
	}
	return bits
}
