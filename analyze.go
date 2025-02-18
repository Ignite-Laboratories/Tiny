package tiny

import "bytes"

type _analyze int

// Shade gives heuristics around the distribution of 1s in the provided measure.
func (a _analyze) Shade(measure Measurement) BinaryShade {
	shade := a.ByteShade(measure.Bytes...)
	shade.combine(a.BitShade(measure.Bits...))
	return shade
}

// BitShade gives heuristics around the distribution of 1s in the provided bits.
func (_ _analyze) BitShade(bits ...Bit) BinaryShade {
	count := BinaryShade{}

	byteIndex := 0
	for i := 0; i < len(bits); i++ {
		if bits[i] == One {
			count.Distribution[byteIndex]++
			count.Ones++
		} else {
			count.Zeros++
		}
		count.Total++
		byteIndex++
		if byteIndex == 8 {
			byteIndex = 0
		}
	}
	count.calculate()

	return count
}

// ByteShade gives heuristics around the distribution of 1s in the provided bytes.
func (a _analyze) ByteShade(data ...byte) BinaryShade {
	count := BinaryShade{}

	for i := 0; i < len(data); i++ {
		c := a.BitShade(From.Byte(data[i])...)
		count.Ones += c.Ones
		count.Zeros += c.Zeros
		count.Total += c.Total
	}
	count.Distribution = a.OneDistribution(data...)
	count.calculate()

	return count
}

// HasPrefix checks if the source Bit slice begins with the provided Bit slice
func (_ _analyze) HasPrefix(data []Bit, prefix ...Bit) bool {
	return bytes.HasPrefix(Upcast(data), Upcast(prefix))
}

// OneDistribution counts how many ones occupy each position of the provided bytes.
func (_ _analyze) OneDistribution(data ...byte) [8]int {
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
