package tiny

// GetBitShade counts the number of 1s and 0s in the data.
func GetBitShade(bits []Bit) Count {
	count := Count{}

	for i := 0; i < len(bits); i++ {
		if bits[i] == One {
			count.Ones++
		} else {
			count.Zeros++
		}
	}
	count.Total++

	return shadeCount(count)
}

// GetDataShade checks if the number of ones in the data exceeds half it's the length.
func GetDataShade(data []byte) Count {
	count := Count{}

	for i := 0; i < len(data); i++ {
		c := GetBitShade(FromByte(data[i]))
		count.Ones += c.Ones
		count.Zeros += c.Zeros
		count.Total += c.Total
	}

	return shadeCount(count)
}

// ToggleBits XORs every bit with 1.
func ToggleBits(bits []Bit) []Bit {
	for i := 0; i < len(bits); i++ {
		bits[i] ^= One
	}
	return bits
}

// ToggleData XORs every bit of a Remainder with 1.
func ToggleData(data []byte) []byte {
	for i := 0; i < len(data); i++ {
		data[i] ^= MaxByte
	}
	return data
}
