package tiny

type _modify int

// XORWithPattern creates a repeating pattern of the provided bits and XORs it against
// the entire length of the Measure.
func (_ _modify) XORWithPattern(measure Measure, pattern ...Bit) {
	patternI := 0
	measure.ForEachBit(func(i int, bit Bit) Bit {
		bit = bit ^ pattern[patternI]
		patternI++
		if patternI >= len(pattern) {
			patternI = 0
		}
		return bit
	})
}

// XORWithBits walks the provided pattern and XORs every bit with the source Measure's
// bits, starting from the most significant bit.
func (_ _modify) XORWithBits(measure Measure, bits ...Bit) {
	measure.ForEachBit(func(i int, bit Bit) Bit {
		if i > len(bits) {
			return bit
		}
		return bit ^ bits[i]
	})
}

// XORByteWithBits XORs a fixed range of bits against a byte, starting at the most significant bit.
func (_ _modify) XORByteWithBits(b byte, bits ...Bit) byte {
	if len(bits) > 8 {
		bits = bits[:8]
	}

	byteBits := From.Byte(b)
	for i := 0; i < len(bits); i++ {
		byteBits[i] ^= bits[i]
	}

	return To.Byte(byteBits...)
}

// XORBytesWithBits XORs a fixed range of bits against every byte, starting at the most significant bit of each.
func (m _modify) XORBytesWithBits(data []byte, bits ...Bit) []byte {
	for i := 0; i < len(data); i++ {
		data[i] = m.XORByteWithBits(data[i], bits...)
	}
	return data
}

// Toggle XORs every bit of each Measure with 1.
func (_ _modify) Toggle(measures ...Measure) {
	for _, measure := range measures {
		measure.ForEachBit(func(_ int, bit Bit) Bit { return bit ^ One })
	}
}

// ToggleBits XORs every bit with 1.
func (_ _modify) ToggleBits(bits ...Bit) []Bit {
	for i := 0; i < len(bits); i++ {
		bits[i] ^= One
	}
	return bits
}

// ToggleBytes XORs every bit of a byte with 1.
func (_ _modify) ToggleBytes(data ...byte) []byte {
	for i := 0; i < len(data); i++ {
		data[i] ^= MaxByte
	}
	return data
}

// DropMostSignificantBit removes the '128' bit from the input bytes and returns a Measure as
// it may not fit into a standard slice of bytes.
// This should be called when your bytes are naturally all below the 128 threshold.
func (m _modify) DropMostSignificantBit(data ...byte) Measure {
	return m.DropMostSignificantBits(1, data...)
}

// DropMostSignificantBits removes the provided number of most significant bits from the input bytes
// and returns a Measure from this operation, as it may not fit back into a standard slice of bytes.
func (_ _modify) DropMostSignificantBits(count int, data ...byte) Measure {
	remainder := Measure{}
	for _, b := range data {
		remainder.AppendBits(From.Byte(b)[count:]...)
	}
	return remainder
}
