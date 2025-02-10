package tiny

import "bytes"

// analyze is a way to glean information about existing binary information.
type analyze struct{}

// BitShade counts the number of 1s and 0s in the data.
func (_ analyze) BitShade(bits []Bit) BinaryCount {
	count := BinaryCount{}

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

// DataShade checks if the number of ones in the data exceeds half it's the length.
func (a analyze) DataShade(data []byte) BinaryCount {
	count := BinaryCount{}

	for i := 0; i < len(data); i++ {
		c := a.BitShade(From.Byte(data[i]))
		count.Ones += c.Ones
		count.Zeros += c.Zeros
		count.Total += c.Total
	}
	count.Distribution = a.OneDistribution(data)
	count.Calculate()

	return count
}

// HasPrefix upcasts the input slices to bytes and then calls bytes.HasPrefix.
func (_ analyze) HasPrefix(data []Bit, pattern []Bit) bool {
	return bytes.HasPrefix(subToByte(data), subToByte(pattern))
}

// OneDistribution counts how many ones occupy each position of the provided slice of bytes.
func (_ analyze) OneDistribution(data []byte) [8]int {
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
