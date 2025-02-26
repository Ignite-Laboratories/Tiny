package tiny

type _modify int

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

// QuarterSplit calls Measurement.QuarterSplit on each element of the provided slice.
// NOTE: QuarterSplit is a mutable operation against the underlying measurements.
func (_ _modify) QuarterSplit(data ...Measurement) {
	for _, d := range data {
		d.QuarterSplit()
	}
}
