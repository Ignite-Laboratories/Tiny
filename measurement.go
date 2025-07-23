package tiny

import (
	"fmt"
	"strings"
	"tiny/travel"
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

// NewMeasurementOfBit creates a new Measurement of the provided bit-width consisting entirely of the provided Bit.
//
// Inward and outward travel directions are supported and work from the midpoint of the width, biased towards the west.
func NewMeasurementOfBit(w int, b Bit) Measurement {
	return NewMeasurementOfPattern(w, travel.Eastbound, b)
}

// NewMeasurementOfPattern creates a new Measurement of the provided bit-width consisting of the pattern emitted across it in the direction.Direction of travel.Travel.
//
// Inward and outward travel directions are supported and work from the midpoint of the width, biased towards the west.
func NewMeasurementOfPattern(w int, t travel.Travel, p ...Bit) Measurement {
	// TODO: Generate a random name
	if w <= 0 || len(p) == 0 {
		return Measurement{}
	}

	if t == travel.Northbound || t == travel.Southbound {
		panic(fmt.Sprintf("cannot take a latitudinal binary measurement [%v]", t.StringFull(true)))
	}

	printer := func(width int, tt travel.Travel) []Bit {
		bits := make([]Bit, width)
		patternI := 0
		for i := 0; i < width; i++ {
			ii := i
			if tt == travel.Westbound {
				ii = width - 1 - i
			}

			bits[ii] = p[patternI]
			patternI = (patternI + 1) % len(p)
		}
		return bits
	}

	if t == travel.Inbound || t == travel.Outbound {
		leftWidth := w / 2
		rightWidth := w - leftWidth

		if t == travel.Inbound {
			left := NewMeasurement(printer(leftWidth, travel.Eastbound)...)
			right := NewMeasurement(printer(rightWidth, travel.Westbound)...)
			return left.AppendMeasurements(right)
		}
		return NewMeasurement(printer(leftWidth, travel.Westbound)...).Append(printer(rightWidth, travel.Eastbound)...)
	}
	return NewMeasurement(printer(w, t)...)
}

// NewMeasurementOfZeros creates a new Measurement of the provided bit-width consisting entirely of 0s.
func NewMeasurementOfZeros(width int) Measurement {
	// TODO: Generate a random name
	return Measurement{
		Bytes: make([]byte, width/8),
		Bits:  make([]Bit, width%8),
	}.RollUp()
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
	return zeros.RollUp()
}

// NewMeasurement creates a new Measurement of the provided Bit slice.
func NewMeasurement(bits ...Bit) Measurement {
	// TODO: Generate a random name
	return Measurement{
		Bits: bits,
	}.RollUp()
}

// NewMeasurementOfBytes creates a new Measurement of the provided byte slice.
func NewMeasurementOfBytes(bytes ...byte) Measurement {
	// TODO: Generate a random name
	m := Measurement{}
	m = m.AppendBytes(bytes...)
	return m.RollUp()
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
			bits[7-i] = Bit((b >> i) & 1)
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
	return a.RollUp()
}

// AppendMeasurements places the provided measurement at the end of the measurement.
func (a Measurement) AppendMeasurements(m ...Measurement) Measurement {
	for _, mmt := range m {
		a = a.Append(mmt.GetAllBits()...)
	}
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

// PrependMeasurements places the provided measurement at the start of the measurement.
func (a Measurement) PrependMeasurements(m ...Measurement) Measurement {
	if len(m) == 0 {
		return a
	}

	result := m[len(m)-1]
	for i := len(m) - 2; i >= 0; i-- {
		result = m[i].AppendBytes(result.Bytes...).Append(result.Bits...)
	}
	result = result.AppendBytes(a.Bytes...).Append(a.Bits...)
	return result.RollUp()
}

// Reverse reverses the order of all bits in the measurement.
func (a Measurement) Reverse() Measurement {
	// TODO: reverse the measurement bit order
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

	bits := Emit(Expression{}, 8, a.Bytes[len(a.Bytes)-1])
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
		bits := Emit(Expression{}, 8, a.Bytes[0])
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
				b |= 1 << (7 - i)
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

// StringPretty returns a measurement-formatted string of the current binary information. Measurements
// are simply formatted with a single space between digits.
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
