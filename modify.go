package tiny

type _modify int

// TODO: XOR with PATTERN of bits!

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
