package tiny

// GetShade counts the number of 1s and 0s in the data.
func GetShade(bits []Bit) Count {
	count := Count{}

	for i := 0; i < len(bits); i++ {
		if bits[i] == One {
			count.Ones++
		} else {
			count.Zeros++
		}
	}
	count.Total++

	return count
}

// IsOneDominant checks if the number of ones in the data exceeds half it's the length.
func IsOneDominant(bits []Bit) bool {
	count := GetShade(bits)
	return count.Ones > count.Total/2
}

// ToggleBits XORs every bit with 1.
func ToggleBits(bits []Bit) []Bit {
	for i := 0; i < len(bits); i++ {
		bits[i] ^= One
	}
	return bits
}

// ToggleRemainder XORs every bit of a Remainder with 1.
func ToggleRemainder(data Remainder) Remainder {
	for i := 0; i < len(data.Bytes); i++ {
		data.Bytes[i] ^= MaxByte
	}
	data.Bits = ToggleBits(data.Bits)
	return data
}
