package tiny

import "fmt"

// Measurement is a variable-width slice of bits and is used to efficiently store them in operating memory.
// As most languages inherently require at least 8 bits to store custom types, storing each bit individually
// would need 8 times the size of every bit - thus, the measurement was born.
//
//	 tl;dr: This holds bits in byte form, leaving anything less than a byte
//		       at the end of the binary information as a remainder of bits.
//
// NOTE: A measurement can only hold up to a WordWidth of bits and will panic with ErrorMeasurementLimit if you
// attempt to store beyond that limit.  For working with longer stretches of binary information, please see the Phrase.
type Measurement struct {
	// Bytes holds complete byte data.
	Bytes []byte
	// Bits holds any remaining bits.
	Bits []Bit
}

// NewMeasurement creates a new Measurement of the provided Bit slice.
func NewMeasurement(bits ...Bit) Measurement {
	return Measurement{
		Bits: bits,
	}
}

// BitLength gets the total bit length of this Measurement's recorded data.
func (a Measurement) BitLength() int {
	return (len(a.Bytes) * 8) + len(a.Bits)
}

// GetAllBits returns a slice of the Measurement's individual bits.
func (m *Measurement) GetAllBits() []Bit {
	var byteBits []Bit
	for _, b := range m.Bytes {
		bits := make([]Bit, 8)
		for i := byte(7); i < 8; i-- {
			bits[i] = Bit((b >> i) & 1)
		}
		byteBits = append(byteBits, bits...)
	}
	return append(byteBits, m.Bits...)
}

// Append places the provided bits at the end of the Measurement.
func (a Measurement) Append(bits ...Bit) Measurement {
	a = a.sanityCheck(bits...)

	a.Bits = append(a.Bits, bits...)
	return a.rollup()
}

// AppendBytes places the provided bits at the end of the Measurement.
func (a Measurement) AppendBytes(bytes ...byte) Measurement {
	a = a.sanityCheck()

	lastBits := a.Bits
	for _, b := range bytes {
		bits := make([]Bit, 8)

		for i := byte(7); i < 8; i-- {
			bits[i] = Bit((b >> i) & 1)
		}

		blended := append(lastBits, bits[:8-len(lastBits)]...)
		lastBits = bits[8-len(lastBits):]

		var newByte byte
		for i := byte(7); i < 8; i-- {
			newByte |= byte(blended[i]) << i
		}

		a.Bytes = append(a.Bytes, newByte)
	}

	a.Bits = lastBits
	return a.rollup()
}

// Prepend places the provided bits at the start of the Measurement.
func (a Measurement) Prepend(bits ...Bit) Measurement {
	a = a.sanityCheck(bits...)

	oldBits := a.Bits
	oldBytes := a.Bytes
	a.Bytes = []byte{}
	a.Bits = []Bit{}
	a = a.Append(bits...)
	a = a.AppendBytes(oldBytes...)
	a = a.Append(oldBits...)
	return a.rollup()
}

// PrependBytes places the provided bytes at the start of the Measurement.
func (a Measurement) PrependBytes(bytes ...byte) Measurement {
	a = a.sanityCheck()

	oldBits := a.Bits
	oldBytes := a.Bytes
	a.Bytes = bytes
	a.Bits = []Bit{}
	a = a.AppendBytes(oldBytes...)
	a = a.Append(oldBits...)
	return a.rollup()
}

/**
Utilities
*/

// sanityCheck ensures the provided bits are all 1s and 0s and rolls the currently measured bits into bytes, if possible.
func (a Measurement) sanityCheck(bits ...Bit) Measurement {
	SanityCheck(bits...)
	return a.rollup()
}

// rollup combines the currently measured bits into bytes, if possible.
func (a Measurement) rollup() Measurement {
	for len(a.Bits) >= 8 {
		var b byte
		for i := byte(7); i < 8; i-- {
			if a.Bits[i] == 1 {
				b |= 1 << i
			}
		}
		a.Bits = a.Bits[8:]
		a.Bytes = append(a.Bytes, b)
	}
	return a
}

func (a Measurement) String() string {
	return fmt.Sprintf("%v", m.StringBinary())
}
