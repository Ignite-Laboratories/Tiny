package tiny

import "log"

// modify is a way to alter existing binary information.
type modify struct{}

// XORBitsWithPattern XORs a pattern of bits against a byte, starting at the most significant bit.
func (_ modify) XORBitsWithPattern(b byte, pattern ...Bit) byte {
	if len(pattern) > 8 {
		log.Fatalf("input pattern should not be larger than a byte")
	}

	bits := From.Byte(b)
	for i := 0; i < len(pattern); i++ {
		bits[i] ^= pattern[i]
	}

	return To.Byte(bits...)
}

// XORBytesWithPattern XORs a pattern of bits against every byte, starting at the most significant bit of each.
func (m modify) XORBytesWithPattern(data []byte, pattern ...Bit) []byte {
	for i := 0; i < len(data); i++ {
		data[i] = m.XORBitsWithPattern(data[i], pattern...)
	}
	return data
}

// ToggleBits XORs every bit with 1.
func (_ modify) ToggleBits(bits ...Bit) []Bit {
	for i := 0; i < len(bits); i++ {
		bits[i] ^= One
	}
	return bits
}

// ToggleBytes XORs every bit of a byte with 1.
func (_ modify) ToggleBytes(data ...byte) []byte {
	for i := 0; i < len(data); i++ {
		data[i] ^= MaxByte
	}
	return data
}

// DropMostSignificantBit removes the '128' bit from the input bytes and returns a Remainder as
// it may not fit into a standard slice of bytes.
// This should be called when your bytes are naturally all below the 128 threshold.
func (m modify) DropMostSignificantBit(data ...byte) Remainder {
	return m.DropMostSignificantBits(1, data...)
}

// DropMostSignificantBits removes the provided number of most significant bits from the input bytes
// and returns the remainder from this operation, as it may not fit back into a standard slice of bytes.
func (_ modify) DropMostSignificantBits(count int, data ...byte) Remainder {
	remainder := Remainder{}
	for _, b := range data {
		remainder.AppendBits(From.Byte(b)[count:]...)
	}
	return remainder
}
