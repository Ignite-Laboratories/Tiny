package tiny

import (
	"strings"
)

// Measurement is a variable-width slice of bits and is used to efficiently store them in operating memory.
// As most languages inherently require at least 8 bits to store custom types, storing each bit individually
// would need 8 times the size of every bit - thus, the measurement was born.
type Measurement struct {
	// Name represents the name of this measurement.
	Name string
	// Bytes holds complete byte data.
	Bytes []byte
	// Bits holds any remaining bits.
	Bits []Bit
}

// NewMeasurementOfDigit creates a new Measurement of the provided bit-width consisting entirely of the provided digit.
func NewMeasurementOfDigit(width int, digit Bit) Measurement {
	// TODO: Generate a random name
	if digit == One {
		return NewMeasurementOfOnes(width)
	}
	return NewMeasurementOfZeros(width)
}

// NewMeasurementOfZeros creates a new Measurement of the provided bit-width consisting entirely of 0s.
func NewMeasurementOfZeros(width int) Measurement {
	// TODO: Generate a random name
	return Measurement{
		Bytes: make([]byte, width/8),
		Bits:  make([]Bit, width%8),
	}
}

// NewMeasurementOfOnes creates a new Measurement of the provided bit-width consisting entirely of 1s.
func NewMeasurementOfOnes(width int) Measurement {
	// TODO: Generate a random name
	zeros := NewMeasurementOfZeros(width)
	for i := range zeros.Bytes {
		zeros.Bytes[i] = 255
	}
	for i := range zeros.Bits {
		zeros.Bits[i] = 1
	}
	return zeros
}

// NewMeasurement creates a new Measurement of the provided Bit slice.
func NewMeasurement(bits ...Bit) Measurement {
	// TODO: Generate a random name
	return Measurement{
		Bits: bits,
	}
}

// NewMeasurementOfBytes creates a new Measurement of the provided byte slice.
func NewMeasurementOfBytes(bytes ...byte) Measurement {
	// TODO: Generate a random name
	m := Measurement{}
	m = m.AppendBytes(bytes...)
	return m
}

// BitWidth gets the total bit width of this Measurement's recorded data.
func (a Measurement) BitWidth() uint {
	return uint((len(a.Bytes) * 8) + len(a.Bits))
}

// GetAllBits returns a slice of the Measurement's individual bits.
func (a Measurement) GetAllBits() []Bit {
	a = a.sanityCheck()
	var byteBits []Bit
	for _, b := range a.Bytes {
		bits := make([]Bit, 8)
		for i := 7; i >= 0; i-- {
			bits[i] = Bit((b >> i) & 1)
		}
		byteBits = append(byteBits, bits...)
	}
	return append(byteBits, a.Bits...)
}

// Append places the provided bits at the end of the Measurement.
func (a Measurement) Append(bits ...Bit) Measurement {
	a = a.sanityCheck(bits...)

	a.Bits = append(a.Bits, bits...)
	return a.RollUp()
}

// AppendBytes places the provided bits at the end of the Measurement.
func (a Measurement) AppendBytes(bytes ...byte) Measurement {
	a = a.sanityCheck()

	lastBits := a.Bits
	for _, b := range bytes {
		bits := make([]Bit, 8)

		for i := byte(7); i < 8; i-- {
			bits[7-i] = Bit((b >> i) & 1)
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
	return a.RollUp()
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
	return a.RollUp()
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
	return a.RollUp()
}

// Reverse reverses the order of all bits in the measurement.
func (a Measurement) Reverse() Measurement {
	// TODO: reverse the measurement
	return a
}

// BleedLastBit returns the last bit of the measurement and a measurement missing that bit.
func (a Measurement) BleedLastBit() (Bit, Measurement) {
	if a.BitWidth() == 0 {
		panic("cannot bleed the last bit of an empty measurement")
	}

	if len(a.Bits) >= 0 {
		a.Bits = a.Bits[:len(a.Bits)-1]
		return a.Bits[len(a.Bits)-1], a.RollUp()
	}

	bits := Emit(Bits.All(), 8, a.Bytes[len(a.Bytes)-1])
	last := bits[7]
	bits = bits[:7]
	a.Bits = append(bits, a.Bits...)
	return last, a.RollUp()
}

// BleedFirstBit returns the first bit of the measurement and a measurement missing that bit.
func (a Measurement) BleedFirstBit() (Bit, Measurement) {
	if a.BitWidth() == 0 {
		panic("cannot bleed the first bit of an empty measurement")
	}

	if len(a.Bytes) >= 0 {
		bits := Emit(Bits.All(), 8, a.Bytes[0])
		first := bits[0]
		bits = bits[1:]
		a.Bytes = a.Bytes[1:]
		a = a.Prepend(bits...)
		return first, a.RollUp()
	} else {
		bit := a.Bits[0]
		a.Bits = a.Bits[1:]
		return bit, a.RollUp()
	}
}

// RollUp combines the currently measured bits into the measured bytes if there is enough recorded.
func (a Measurement) RollUp() Measurement {
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

/**
Utilities
*/

// sanityCheck ensures the provided bits are all 1s and 0s and rolls the currently measured bits into bytes, if possible.
func (a Measurement) sanityCheck(bits ...Bit) Measurement {
	if a.Bytes == nil {
		a.Bytes = []byte{}
	}
	if a.Bits == nil {
		a.Bits = []Bit{}
	}
	SanityCheck(bits...)
	return a.RollUp()
}

// String converts the measurement to a binary string entirely consisting of 1s and 0s.
func (a Measurement) String() string {
	bits := a.GetAllBits()

	builder := strings.Builder{}
	builder.Grow(len(bits))
	for _, b := range bits {
		builder.WriteString(b.String())
	}
	return builder.String()
}

// StringPretty converts the measurement to a binary string with a single space character
// between each 1 and 0.
func (a Measurement) StringPretty() string {
	bits := a.GetAllBits()

	if len(bits) == 0 {
		return ""
	}

	builder := strings.Builder{}
	builder.Grow(len(bits)*2 - 1)

	builder.WriteString(bits[0].String())

	if len(bits) > 1 {
		for _, bit := range bits[1:] {
			builder.WriteString(" ")
			builder.WriteString(bit.String())
		}
	}

	return builder.String()
}
