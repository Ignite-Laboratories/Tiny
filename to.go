package tiny

import "strconv"

type _to int

// toNumeric takes in a set of bits and converts it to an numeric value up to the
// specified width.  If the amount of provided bits exceeds that width, the excess
// bits dropped entirely.  If the width exceeds 32 (a maximum for an int), they
// are also dropped.
// For example: If 4 is provided, a Nibble value of [0-15] is returned even if 8 bits are provided.
func toNumeric(width int, bits ...Bit) int {
	if width > 31 {
		bits = bits[:32]
	}

	if len(bits) > width {
		bits = bits[:width]
	}

	result := 0
	padding := width - len(bits)

	for i, bit := range bits {
		result |= int(bit) << ((width - 1) - (i + padding))
	}
	return result
}

// Crumb converts the first 2 bits of the Bit slice to a Crumb and ignores the rest.
func (_ _to) Crumb(bits ...Bit) Crumb {
	return Crumb(toNumeric(2, bits...))
}

// Note converts the first 3 bits of the Bit slice to a Note and ignores the rest.
func (_ _to) Note(bits ...Bit) Note {
	return Note(toNumeric(3, bits...))
}

// Nibble converts the first 4 bits of the Bit slice to a Nibble and ignores the rest.
func (_ _to) Nibble(bits ...Bit) Nibble {
	return Nibble(toNumeric(4, bits...))
}

// Flake converts the first 5 bits of the Bit slice to a Flake and ignores the rest.
func (_ _to) Flake(bits ...Bit) Flake {
	return Flake(toNumeric(5, bits...))
}

// Morsel converts the first 6 bits of the Bit slice to a Morsel and ignores the rest.
func (_ _to) Morsel(bits ...Bit) Morsel {
	return Morsel(toNumeric(6, bits...))
}

// Shred converts the first 7 bits of the Bit slice to a Shred and ignores the rest.
func (_ _to) Shred(bits ...Bit) Shred {
	return Shred(toNumeric(7, bits...))
}

// Byte converts the first 8 bits of the Bit slice to a byte and ignores the rest.
func (_ _to) Byte(bits ...Bit) byte {
	return byte(toNumeric(8, bits...))
}

// Bytes converts a Bit slice to a Measure.
func (t _to) Bytes(bits ...Bit) Measure {
	var bytes []byte
	for i := 0; i+7 < len(bits); i += 8 {
		bytes = append(bytes, t.Byte(bits[i:i+8]...))
	}
	remainingBits := bits[len(bytes)*8:]
	return Measure{Bytes: bytes, Bits: remainingBits}
}

// String creates a slice of mixed 1s and 0s from the provided Bit slice
func (_ _to) String(bits ...Bit) string {
	output := ""
	for i := 0; i < len(bits); i++ {
		output += strconv.Itoa(int(bits[i]))
	}
	return output
}
