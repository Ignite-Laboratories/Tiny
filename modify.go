package tiny

import "log"

// modify is a way to alter existing binary information.
type modify struct{}

// XORWithPattern XORs a pattern of bits against a byte, starting at the most significant bit.
func (_ modify) XORWithPattern(b byte, pattern []Bit) byte {
	if len(pattern) > 8 {
		log.Fatalf("input pattern should not be larger than a byte")
	}

	bits := From.Byte(b)
	for i := 0; i < len(pattern); i++ {
		bits[i] ^= pattern[i]
	}

	return To.Byte(bits...)
}

// XORDataWithPattern XORs a pattern of bits against every byte, starting at the most significant bit of each.
func (m modify) XORDataWithPattern(data []byte, pattern []Bit) []byte {
	for i := 0; i < len(data); i++ {
		data[i] = m.XORWithPattern(data[i], pattern)
	}
	return data
}

// ToggleBits XORs every bit with 1.
func (_ modify) ToggleBits(bits []Bit) []Bit {
	for i := 0; i < len(bits); i++ {
		bits[i] ^= One
	}
	return bits
}

// ToggleData XORs every bit of a byte with 1.
func (_ modify) ToggleData(data []byte) []byte {
	for i := 0; i < len(data); i++ {
		data[i] ^= MaxByte
	}
	return data
}
