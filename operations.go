package tiny

import (
	"log"
	"strconv"
)

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

// XORWithPattern takes the provided byte and XORs the provided bits against it, starting at the most significant bit.
func XORWithPattern(b byte, pattern []Bit) byte {
	if len(pattern) > 8 {
		log.Fatalf("input pattern should not be larger than a byte")
	}

	bits := FromByte(b)
	for i := 0; i < len(pattern); i++ {
		bits[i] ^= pattern[i]
	}

	return ToByte(bits)
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

// ToString converts a set of Bit values into a string.
func ToString[T SubByte](bits []T) string {
	output := ""
	for i := 0; i < len(bits); i++ {
		output += strconv.Itoa(int(bits[i]))
	}
	return output
}
