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

// BitLength gets the total bit length of this Measurement's recorded data.
func (a Measurement) BitLength() int {
	return (len(a.Bytes) * 8) + len(a.Bits)
}

// GetAllBits returns a slice of the Measurement's individual bits.
func (m *Measurement) GetAllBits() []Bit {
	var byteBits []Bit
	for _, v := range m.Bytes {
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
	return a.rollupBits()
}

// AppendBytes places the provided bits at the end of the Measurement.
func (a Measurement) AppendBytes(bytes ...byte) Measurement {
	a = a.sanityCheckBytes(bytes...)

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
	return a.rollupBits()
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
	return a.rollupBits()
}

// PrependBytes places the provided bytes at the start of the Measurement.
func (a Measurement) PrependBytes(bytes ...byte) Measurement {
	a = a.sanityCheckBytes(bytes...)

	oldBits := a.Bits
	oldBytes := a.Bytes
	a.Bytes = bytes
	a.Bits = []Bit{}
	a = a.AppendBytes(oldBytes...)
	a = a.Append(oldBits...)
	return a.rollupBits()
}

/**
Utilities
*/

// sanityCheck does three things:
//
// 1. Ensures the provided bits are all 1s and 0s.
//
// 2. Ensures the resulting bit length will not exceed the WordWidth.
//
// 3. Rolls up the currently measured bits into bytes, if possible.
func (a Measurement) sanityCheck(bits ...Bit) Measurement {
	SanityCheck(bits...)
	if a.BitLength()+len(bits) > WordWidth {
		panic(ErrorMeasurementLimit)
	}
	return a.rollupBits()
}

// sanityCheckBytes does two things:
//
// 1. Ensures the resulting bit length will not exceed the WordWidth.
//
// 2. Rolls up the currently measured bits to bytes, if possible.
func (a Measurement) sanityCheckBytes(bytes ...byte) Measurement {
	if a.BitLength()+(len(bytes)*8) > WordWidth {
		panic(ErrorMeasurementLimit)
	}
	return a.rollupBits()
}

// rollupBits rolls the currently measured bits into bytes, if possible.
func (a Measurement) rollupBits() Measurement {
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
