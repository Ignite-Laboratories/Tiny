package tiny

// GetShade counts the number of 1s and 0s in the data.
func GetShade(data Remainder) Count {
	count := Count{}

	// Define how to count
	increment := func(bits []Bit) {
		for i := 0; i < len(bits); i++ {
			if bits[i] == One {
				count.Ones++
			} else {
				count.Zeros++
			}
		}
		count.Total++
	}

	// Walk each byte and count
	for i := 0; i < len(data.Bytes); i++ {
		increment(FromByte(data.Bytes[i]))
	}
	// Walk each remaining bit and count
	increment(data.Bits)

	return count
}

// IsOneDominant checks if the number of ones in the data exceeds half it's the length.
func IsOneDominant(data Remainder) bool {
	count := GetShade(data)
	return count.Ones > count.Total/2
}

// ToggleBits XORs every bit with 1.
func ToggleBits(bits []Bit) []Bit {
	for i := 0; i < len(bits); i++ {
		bits[i] ^= One
	}
	return bits
}

// ToggleData XORs every bit of a Remainder with 1.
func ToggleData(data Remainder) Remainder {
	for i := 0; i < len(data.Bytes); i++ {
		data.Bytes[i] ^= MaxByte
	}
	data.Bits = ToggleBits(data.Bits)
	return data
}
