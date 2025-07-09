package tiny

import (
	"math/big"
	"strconv"
)

// Measurement is a variable-width slice of bits and is used to efficiently
// store them in operating memory.
// As most languages inherently require at least 8 bits to store custom types,
// storing each bit individually would need 8 times the size of every bit -
// thus, the measurement was born.
//
//	tl;dr: This holds bits in byte form, leaving anything less than a byte
//	       at the end of the binary information as a remainder of bits.
//
// NOTE: A measurement is limited to your architecture's bit-width wide by design.
// This allows you to easily grow or shrink bytes at the bit level and then capture the
// new value of each measurement as a standard 'int'.
//
// For longer stretches of binary information, string together measurements
// using a Phrase.
type Measurement struct {
	// Bytes holds complete byte data.
	Bytes []byte
	// Bits holds any remaining bits that didn't fit into a byte for whatever reason.
	Bits []Bit
}

// NewMeasurement constructs a Measurement, which represents a variable slice of bits.
func NewMeasurement(bytes []byte, bits ...Bit) Measurement {
	b := bytes
	ii := 0
	currentBits := make([]Bit, 0, 8)

	for _, bit := range bits {
		if ii > 7 {
			b = append(b, To.Byte(currentBits...))
			currentBits = make([]Bit, 0, 8)
			ii = 0
		}
		currentBits = append(currentBits, bit)
		ii++
	}

	if len(currentBits) == 8 {
		b = append(b, To.Byte(currentBits...))
		currentBits = make([]Bit, 0, 8)
	}

	return Measurement{
		Bytes: b,
		Bits:  currentBits,
	}
}

// NewMeasurementFromBits creates a new Measurement from the provided input bits.
//
// NOTE: This will panic if provided more bits than your architecture's bit width.
func NewMeasurementFromBits(bits ...Bit) Measurement {
	if len(bits) > GetArchitectureBitWidth() {
		panic(errorMeasurementLimit)
	}
	return NewMeasurement([]byte{}, bits...)
}

// NewMeasurementFromString creates a new Measurement from a binary string input.
//
// NOTE: This will panic if provided a string longer than your architecture's bit width.
func NewMeasurementFromString(s string) Measurement {
	if len(s) > GetArchitectureBitWidth() {
		panic(errorMeasurementLimit)
	}
	bits := make([]Bit, len(s))
	for i := 0; i < len(bits); i++ {
		bits[i] = Bit(s[i] & 1)
	}
	return NewMeasurement([]byte{}, bits...)
}

// NewMeasurementFromBigInt creates a new Measurement from a big.Int.
//
// NOTE: This will panic if provided a integer represented in base-2 longer than your architecture's bit width.
func NewMeasurementFromBigInt(b *big.Int) Measurement {
	return NewMeasurementFromString(b.Text(2))
}

// GetAllBits returns the measure in the form of a fully expanded Bit slice.
func (m *Measurement) GetAllBits() []Bit {
	byteBits := From.Bytes(m.Bytes...)
	return append(byteBits, m.Bits...)
}

// BitLength gets the total length of this Measurement's individual bits.
func (m *Measurement) BitLength() int {
	return m.ByteBitLength() + len(m.Bits)
}

// ByteBitLength gets the total length of this Measurement's byte's individual bits,
// ignoring the measurement's bits entirely.  This is helpful when attempting to find
// measurements that are not aligned to the width of a standard byte.
func (m *Measurement) ByteBitLength() int { return len(m.Bytes) * 8 }

// Value gets the integer value of the measure using its current bit representation.
// NOTE: Measures are limited to your architecture's bit width - intentionally limiting them to an int.
func (m *Measurement) Value() int {
	v, _ := strconv.ParseInt(To.String(m.GetAllBits()...), 2, GetArchitectureBitWidth())
	return int(v)
}

// Clear empties the Measurement of all bit information.
func (m *Measurement) Clear() {
	m.Bytes = []byte{}
	m.Bits = []Bit{}
}

// Toggle XORs every bit of each Measurement with 1.
func (m *Measurement) Toggle() {
	m.ForEachBit(func(_ int, bit Bit) Bit { return bit ^ One })
}

// ForEachBit calls the provided operation against every bit of the Measurement.
func (m *Measurement) ForEachBit(operation func(i int, bit Bit) Bit) {
	outBytes := make([]byte, len(m.Bytes))
	outBits := make([]Bit, len(m.Bits))
	bitI := 0
	for byteI, b := range m.Bytes {
		var newBits [8]Bit
		for subI, bit := range From.Byte(b) {
			newBits[subI] = operation(bitI, bit)
			bitI++
		}
		outBytes[byteI] = To.Byte(newBits[:]...)
	}
	for i, bit := range m.Bits {
		outBits[i] = operation(bitI, bit)
		bitI++
	}
	m.Bytes = outBytes
	m.Bits = outBits
}

// Read returns the individually addressed bits of the Measurement, ranged from the low
// index (inclusive) to the  high index (exclusive).  This intentionally follows standard
// Go slice [low:high] indexing, meaning it also fails the same if you reference beyond
// the measurable index boundaries.
func (m *Measurement) Read(low int, high int) []Bit {
	bits := m.GetAllBits()
	return bits[low:high]
}

// AppendBits places the provided bits at the end of the source Measurement.
//
// NOTE: A measurement can only hold up to your architecture's bit width!
func (m *Measurement) AppendBits(bits ...Bit) {
	if m.BitLength()+len(bits) > GetArchitectureBitWidth() {
		panic(errorMeasurementLimit)
	}
	m.Bits = append(m.Bits, bits...)          // Add the bits to the last remainder
	toAdd := To.Measure(m.Bits...)            // Convert that to byte form
	m.Bytes = append(m.Bytes, toAdd.Bytes...) // Combine the whole bytes
	m.Bits = toAdd.Bits                       // Bring forward the remaining bits
}

// AppendBytes places the provided bytes at the end of the source Measurement.
//
// NOTE: A measurement can only hold up to your architecture's bit width!
func (m *Measurement) AppendBytes(bytes ...byte) {
	if m.BitLength()+len(bytes)*8 > GetArchitectureBitWidth() {
		panic(errorMeasurementLimit)
	}
	lastBits := m.Bits
	for _, b := range bytes {
		byteBits := From.Byte(b)
		blended := append(lastBits, byteBits[:8-len(lastBits)]...)
		lastBits = byteBits[8-len(lastBits):]
		newByte := To.Byte(blended...)
		m.Bytes = append(m.Bytes, newByte)
	}
	m.Bits = lastBits
}

// Append places the provided Measurement at the end of the source Measurement.
// NOTE: A measurement can only hold up to your architecture's bit width!
func (m *Measurement) Append(measure Measurement) {
	m.AppendBytes(measure.Bytes...)
	m.AppendBits(measure.Bits...)
}

// PrependBits places the provided bits at the beginning of the source Measurement.
// NOTE: A measurement can only hold up to your architecture's bit width!
func (m *Measurement) PrependBits(bits ...Bit) {
	if m.BitLength()+len(bits) > GetArchitectureBitWidth() {
		panic(errorMeasurementLimit)
	}
	oldBits := m.Bits
	oldBytes := m.Bytes
	m.Bytes = []byte{}
	m.Bits = []Bit{}
	m.AppendBits(bits...)
	m.AppendBytes(oldBytes...)
	m.AppendBits(oldBits...)
}

// PrependBytes places the provided bytes at the beginning of the source Measurement.
// NOTE: A measurement can only hold up to your architecture's bit width!
func (m *Measurement) PrependBytes(bytes ...byte) {
	if m.BitLength()+len(bytes)*8 > GetArchitectureBitWidth() {
		panic(errorMeasurementLimit)
	}
	m.Bytes = append(bytes, m.Bytes...)
}

// Prepend places the provided Measurement at the beginning of the source Measurement.
// NOTE: A measurement can only hold up to your architecture's bit width!
func (m *Measurement) Prepend(measure Measurement) {
	m.PrependBits(measure.Bits...)   // First the ending bits get prepended
	m.PrependBytes(measure.Bytes...) // Then the starting bytes
}

// TrimStart removes the provided number of bits from the beginning of the Measurement.
func (m *Measurement) TrimStart(count int) {
	bits := m.GetAllBits()
	m.Clear()
	m.AppendBits(bits[count:]...)
}

// TrimEnd removes the provided number of bits from the end of the Measurement.
func (m *Measurement) TrimEnd(count int) {
	bits := m.GetAllBits()
	m.Clear()
	end := len(bits) - count - 1
	m.AppendBits(bits[:end]...)
}

// BreakApart splits the Measurement into two at the provided index and returns their results respectively.
//
// The first returned Measurement ("left") contains data from the start and up to (but not including) the index.
// The second returned Measurement ("right") contains data from the index to the end.
func (m *Measurement) BreakApart(index int) (Measurement, Measurement) {
	left := NewMeasurement([]byte{}, m.Read(0, index)...)
	right := NewMeasurement([]byte{}, m.Read(index, m.BitLength())...)
	return left, right
}

// Invert XORs every bit of the measurement against 1.
func (m *Measurement) Invert() {
	m.ForEachBit(func(_ int, bit Bit) Bit { return bit ^ One })
}

// StringBinary returns the measurement's bits as a binary string of 1s and 0s.
func (m *Measurement) StringBinary() string {
	return To.String(m.GetAllBits()...)
}

func (m *Measurement) String() string {
	return strconv.Itoa(To.Number(GetArchitectureBitWidth(), m.GetAllBits()...))
}
