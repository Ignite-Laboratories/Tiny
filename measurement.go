package tiny

import (
	"strconv"
)

// Measurement is used to efficiently store bits in operating memory.  As most
// languages inherently require at least 8 bits to store custom types, storing
// each bit individually would need 8 times the size of every bit - thus, the
// measurement was born.
//
// Conceptually, this is a musical "measure" - facilitating rhythmic changes
// to binary information; however, a "measurement" is, abstractly, any unit
// recording the presence of something.  Linguistics are cool, folks!
// I highly encourage you consider the alternative meanings of every term
// in your preferred programming languages, they become far more intuitive =)
//
// TL;DR: This holds bits in byte form, leaving anything less than a byte
// at the end of the binary information as a remainder of bits.
//
// NOTE: A measurement is limited to 32 bits wide by design.  This allows you
// to easily grow or shrink bits at the byte level and then capture the
// new value of each individual measurement.
type Measurement struct {
	// Bytes holds complete byte data.
	Bytes []byte
	// Bits holds any remaining bits that didn't fit into a byte for whatever reason.
	Bits []Bit
}

// NewMeasurement constructs a Measurement, which represents an uneven amount of binary bits.
func NewMeasurement(bytes []byte, bits ...Bit) Measurement {
	return Measurement{
		Bytes: bytes,
		Bits:  bits,
	}
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
// NOTE: Measures are limited to 32 bits wide - intentionally limiting them to an int.
func (m *Measurement) Value() int {
	v, _ := strconv.ParseInt(To.String(m.GetAllBits()...), 2, 32)
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
	// Step 1: Are we looking entirely in the Bits section?
	if low >= m.ByteBitLength() {
		// This only concerns the bits, so we can simply drop the byte's bit length
		low -= m.ByteBitLength()
		high -= m.ByteBitLength()
		return m.Bits[low:high]
	}

	lowByteIndex := low / 8    // This is the byte index
	lowByteSubIndex := low % 8 // This is the bit index of that byte
	highByteIndex := high / 8
	highByteSubIndex := high % 8

	var output []Bit
	var foundBytes []byte
	var foundBits []Bit
	var isSplit bool

	// Step 2: Are we split across both the Bits and Measurement?
	if high > m.ByteBitLength() {
		// Yes?  Grab all the bytes from the starting byte...
		foundBytes = m.Bytes[lowByteIndex:]
		foundBits = m.Bits[:highByteSubIndex]
		isSplit = true
	} else {
		// No?  Grab
		foundBytes = m.Bytes[lowByteIndex : highByteIndex+1]
	}

	// Step 3: Split the found bytes apart
	for i, b := range foundBytes {
		if i == 0 {
			// We need to use the low side's sub index
			output = append(output, From.Byte(b)[lowByteSubIndex:]...)
		} else if !isSplit && i == len(foundBytes)-1 {
			// We need to use the high side's sub index
			output = append(output, From.Byte(b)[:highByteSubIndex]...)
		} else {
			// Get the full byte
			output = append(output, From.Byte(b)...)
		}
	}

	// Step 4: Combine the bytes and bits and return
	output = append(output, foundBits...)
	return output
}

// AppendBits places the provided bits at the end of the source Measurement.
// NOTE: A measurement can only hold up to 32 bits!
func (m *Measurement) AppendBits(bits ...Bit) {
	if m.BitLength()+len(bits) > 32 {
		panic(errorMeasureLimit)
	}
	m.Bits = append(m.Bits, bits...)          // Add the bits to the last remainder
	toAdd := To.Measure(m.Bits...)            // Convert that to byte form
	m.Bytes = append(m.Bytes, toAdd.Bytes...) // Combine the whole bytes
	m.Bits = toAdd.Bits                       // Bring forward the remaining bits
}

// AppendBytes places the provided bytes at the end of the source Measurement.
// NOTE: A measurement can only hold up to 32 bits!
func (m *Measurement) AppendBytes(bytes ...byte) {
	if m.BitLength()+len(bytes)*8 > 32 {
		panic(errorMeasureLimit)
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
// NOTE: A measurement can only hold up to 32 bits!
func (m *Measurement) Append(measure Measurement) {
	m.AppendBytes(measure.Bytes...)
	m.AppendBits(measure.Bits...)
}

// PrependBits places the provided bits at the beginning of the source Measurement.
// NOTE: A measurement can only hold up to 32 bits!
func (m *Measurement) PrependBits(bits ...Bit) {
	if m.BitLength()+len(bits) > 32 {
		panic(errorMeasureLimit)
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
// NOTE: A measurement can only hold up to 32 bits!
func (m *Measurement) PrependBytes(bytes ...byte) {
	if m.BitLength()+len(bytes)*8 > 32 {
		panic(errorMeasureLimit)
	}
	m.Bytes = append(bytes, m.Bytes...)
}

// Prepend places the provided Measurement at the beginning of the source Measurement.
// NOTE: A measurement can only hold up to 32 bits!
func (m *Measurement) Prepend(measure Measurement) {
	m.PrependBits(measure.Bits...)   // First the ending bits get prepended
	m.PrependBytes(measure.Bytes...) // Then the starting bytes
}

// QuarterSplit is a unique operation to a byte.  The two most significant bits in a byte
// represent 128 and 64, which means that those two bits cover three quarters of the entire
// address space.  The act of quarter splitting exploits this to reduce the size of bytes
// under a value of 64, while keeping no change in bit length for 64-127, but taking at 1-bit
// hit on anything 128+.  In doing so, a self describing bit scheme can be used for readability.
//
// The first 1-2 bits describe how to read the next bits:
//
//	 0: 6 more bits (the value was under 64)
//	10: 6 more bits (the value was 64-127 and has had 64 subtracted from it)
//	11: 7 more bits (the value was 128+ and has had 128 subtracted from it)
func (m *Measurement) QuarterSplit() {
	// Get the measurement's value and clear it out
	value := m.Value()
	valueWidth := 6
	m.Clear()

	if value < 64 {
		m.AppendBits(0)
	} else if value < 128 {
		m.AppendBits(1, 0)
		value -= 64
	} else {
		m.AppendBits(1, 1)
		value -= 128
		valueWidth = 7
	}
	m.AppendBits(From.Number(value, valueWidth)...)
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

// SplitAtIndex splits the Measurement into two at the provided index and returns their results respectively.
//
// The first returned Measurement ("left") contains data from the start and up to (but not including) the index.
// The second returned Measurement ("right") contains data from the index to the end.
func (m *Measurement) SplitAtIndex(index int) (Measurement, Measurement) {
	left := NewMeasurement([]byte{}, m.Read(0, index)...)
	right := NewMeasurement([]byte{}, m.Read(index, m.BitLength())...)
	return left, right
}
