package tiny

import "strconv"

// to is a way to convert binary slices to other forms.
type to struct{}

// Byte converts the first 8 indices of the tiny.Bit slice to a byte and ignores the rest.
func (_ to) Byte(bits ...Bit) byte {
	if len(bits) > 8 {
		bits = bits[:8]
	}

	result := byte(0)
	padding := 8 - len(bits)

	for i, bit := range bits {
		result |= byte(bit) << (7 - (i + padding))
	}
	return result
}

// Bytes converts a tiny.Bit slice to a tiny.Remainder.
func (t to) Bytes(bits ...Bit) Remainder {
	var bytes []byte
	for i := 0; i+7 < len(bits); i += 8 {
		bytes = append(bytes, t.Byte(bits[i:i+8]...))
	}
	remainingBits := bits[len(bytes)*8:]
	return Remainder{Bytes: bytes, Bits: remainingBits}
}

// String creates a slice of mixed 1s and 0s from the provided tiny.Bit slice
func (_ to) String(bits ...Bit) string {
	output := ""
	for i := 0; i < len(bits); i++ {
		output += strconv.Itoa(int(bits[i]))
	}
	return output
}
