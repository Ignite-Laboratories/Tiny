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
		count.Total++
	}
	count.Calculate()

	return count
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
	count.Distribution = GetDistributionOfOnes(data)
	count.Calculate()

	return count
}

// ToggleBits XORs every bit with 1.
func ToggleBits(bits []Bit) []Bit {
	for i := 0; i < len(bits); i++ {
		bits[i] ^= One
	}
	return bits
}

// ToggleData XORs every bit of a byte with 1.
func ToggleData(data []byte) []byte {
	for i := 0; i < len(data); i++ {
		data[i] ^= MaxByte
	}
	return data
}

// GetDistributionOfOnes counts how many ones occupy each position of the provided slice of bytes.
func GetDistributionOfOnes(data []byte) [8]int {
	output := [8]int{}
	for _, b := range data {
		for i := 0; i < 8; i++ {
			if (b & (1 << (7 - i))) != 0 { // Extract the i-th bit (starting from MSB)
				output[i]++
			}
		}
	}
	return output
}
