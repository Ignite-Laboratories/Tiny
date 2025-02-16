package tiny

import "strconv"

// Measure is used to efficiently store Bits in operating memory, as most
// languages inherently requires at least 8 bits to store any custom type.
// Conceptually, this is a musical "measure" - facilitating rhythmic changes
// to binary information; however, a "measurement" is, abstractly, any unit
// defining the presence of something.  Linguistics are cool, folks!
// I highly encourage you consider the alternative meanings of every term
// in your preferred programming languages, they become far more intuitive =)
//
// TL;DR: This holds bits in byte form, leaving anything less than a byte
// at the end of the binary information as a remainder of Bits.
//
// NOTE: A measure is limited to 32 bits wide by design.  This allows you
// to easily grow or shrink bits at the byte level and then capture the
// new value of each individual measure.
type Measure struct {
	// Bytes holds complete byte data.
	Bytes []byte
	// Bits holds any remaining bits that didn't fit into a byte for whatever reason.
	Bits []Bit
}

// NewMeasure constructs a Measure, which represents an uneven amount of binary bits.
func NewMeasure(bytes []byte, bits ...Bit) Measure {
	return Measure{
		Bytes: bytes,
		Bits:  bits,
	}
}

// GetAllBits returns the measure in the form of a fully expanded Bit slice.
func (m *Measure) GetAllBits() []Bit {
	byteBits := From.Bytes(m.Bytes...)
	return append(byteBits, m.Bits...)
}

// BitLength gets the total length of this Measure's individual bits.
func (m *Measure) BitLength() int {
	return m.ByteBitLength() + len(m.Bits)
}

// ByteBitLength gets the total length of this Measure's byte's individual bits.
func (m *Measure) ByteBitLength() int { return len(m.Bytes) * 8 }

// Value gets the integer value of the measure using its current bit representation.
// NOTE: Measures are limited to 32 bits wide - intentionally limiting them to an int.
func (m *Measure) Value() int {
	v, _ := strconv.ParseInt(To.String(m.GetAllBits()...), 2, 32)
	return int(v)
}

// Toggle XORs every bit of each Measure with 1.
func (m *Measure) Toggle() {
	m.ForEachBit(func(_ int, bit Bit) Bit { return bit ^ One })
}

// ForEachBit calls the provided operation against every bit of the Measure.
func (m *Measure) ForEachBit(operation func(i int, bit Bit) Bit) {
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

// Read returns the individually addressed bits of the Measure, ranged from the low
// index (inclusive) to the  high index (exclusive).  This intentionally follows standard
// Go slice [low:high] indexing, meaning it also fails the same if you reference beyond
// the measurable index boundaries.
func (m *Measure) Read(low int, high int) []Bit {
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

	// Step 2: Are we split across both the Bits and Measure?
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

// AppendBits places the provided bits at the end of the source Measure.
func (m *Measure) AppendBits(bits ...Bit) {
	if m.BitLength()+len(bits) > 32 {
		panic(errorMeasureLimit)
	}
	m.Bits = append(m.Bits, bits...)          // Add the bits to the last remainder
	toAdd := To.Measure(m.Bits...)            // Convert that to byte form
	m.Bytes = append(m.Bytes, toAdd.Bytes...) // Combine the whole bytes
	m.Bits = toAdd.Bits                       // Bring forward the remaining bits
}

// AppendBytes places the provided bytes at the end of the source Measure.
func (m *Measure) AppendBytes(bytes ...byte) {
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

// Append places the provided Measure at the end of the source Measure.
func (m *Measure) Append(measure Measure) {
	m.AppendBytes(measure.Bytes...)
	m.AppendBits(measure.Bits...)
}

// PrependBits places the provided bits at the beginning of the source Measure.
func (m *Measure) PrependBits(bits ...Bit) {
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

// PrependBytes places the provided bytes at the beginning of the source Measure.
func (m *Measure) PrependBytes(bytes ...byte) {
	if m.BitLength()+len(bytes)*8 > 32 {
		panic(errorMeasureLimit)
	}
	m.Bytes = append(bytes, m.Bytes...)
}

// Prepend places the provided Measure at the beginning of the source Measure.
func (m *Measure) Prepend(measure Measure) {
	m.PrependBits(measure.Bits...)   // First the ending bits get prepended
	m.PrependBytes(measure.Bytes...) // Then the starting bytes
}
