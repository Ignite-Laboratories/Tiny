package tiny

// Measure is used to efficiently store Bits in operating memory, as most
// languages inherently requires at least 8 bits to store any custom type.
// Conceptually, this is a musical "measure" - facilitating rhythmic changes
// to binary information; however, a "measurement" is, abstractly, a unit
// defining the presence of something.  Linguistics are cool, folks!
// I highly encourage you consider the alternative meanings of every term
// in your preferred programming languages, they become far more intuitive =)
//
// TL;DR: This holds bits in byte form, leaving anything less than a byte
// at the end of the binary information as a remainder of Bits.
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

// Read returns the individually addressed bits of the Measure, ranged from the low
// index (inclusive) to the  high index (exclusive).  This intentionally follows standard
// Go slice [low:high] indexing, meaning it also fails the same if you reference beyond
// the measurable index boundaries.
func (m *Measure) Read(low uint64, high uint64) []Bit {
	// Step 1: Are we looking entirely in the Bits section?
	if low >= m.ByteBitLength() {
		// This only concerns the bits, so we can simply drop the byte's bit length
		low -= m.ByteBitLength()
		high -= m.ByteBitLength()
		return m.Bits[low:high]
	}

	var foundBytes []byte
	var foundBits []Bit
	var output []Bit

	lowByteIndex := low / 8    // This is the byte index
	lowByteSubIndex := low % 8 // This is the bit index of that byte
	highByteIndex := high / 8
	highByteSubIndex := high % 8

	// Step 2: Are we split across both the Bits and Bytes?
	if high > m.ByteBitLength() {
		// Yes?  Grab all the bytes from the starting byte...
		foundBytes = m.Bytes[lowByteIndex:]
		high -= m.ByteBitLength()
		foundBits = m.Bits[:high]
	} else {
		// No?  Grab
		foundBytes = m.Bytes[lowByteIndex:highByteIndex]
	}

	// Step 3: Split the found bytes apart
	for i, b := range foundBytes {
		if i == 0 {
			// We need to use the low side's sub index
			output = append(output, From.Byte(b)[lowByteSubIndex:]...)
		} else if i == len(foundBytes) {
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

// BitLength gets the total length of this Measure's individual bits.
func (m *Measure) BitLength() uint64 {
	return m.ByteBitLength() + uint64(len(m.Bits))
}

// ByteBitLength gets the total length of this Measure's byte's individual bits.
func (m *Measure) ByteBitLength() uint64 { return uint64(len(m.Bytes)) * 8 }

// AppendBits places the provided bits at the end of the source Measure.
func (m *Measure) AppendBits(bits ...Bit) {
	m.Bits = append(m.Bits, bits...)          // Add the bits to the last remainder
	toAdd := To.Bytes(m.Bits...)              // Convert that to byte form
	m.Bytes = append(m.Bytes, toAdd.Bytes...) // Combine the whole bytes
	m.Bits = toAdd.Bits                       // Bring forward the remaining bits
}

// AppendBytes places the provided bytes at the end of the source Measure.
func (m *Measure) AppendBytes(bytes ...byte) {
	remainderLength := len(m.Bits)
	lastRemainder := m.Bits
	for _, lastByte := range bytes {
		bits := From.Byte(lastByte)
		newBits := append(lastRemainder, bits[:remainderLength]...)
		lastRemainder = bits[remainderLength:]
		m.Bytes = append(m.Bytes, To.Byte(newBits...))
	}
	m.Bits = lastRemainder
}

// AppendRemainder places the provided Measure at the end of the source Measure.
func (m *Measure) AppendRemainder(measure Measure) {
	m.AppendBytes(measure.Bytes...)
	m.AppendBits(measure.Bits...)
}

// PrependBits places the provided bits at the beginning of the source Measure.
func (m *Measure) PrependBits(bits ...Bit) {
	prependLength := len(bits)
	currentBits := bits
	newBytes := make([]byte, 0)
	for _, nextByte := range m.Bytes {
		bits = From.Byte(nextByte)                                // Get the next byte
		newBits := append(currentBits, bits[:8-prependLength]...) // Grab the bits to complete the current byte out of it
		currentBits = bits[8-prependLength:]                      // Grab the remaining bits as what to prepend next
		newBytes = append(newBytes, To.Byte(newBits...))          // Add the newly formed bit structure to the bytes
	}
	m.Bytes = newBytes
	m.Bits = currentBits // This is now the last few bits of the original last byte
}

// PrependBytes places the provided bytes at the beginning of the source Measure.
func (m *Measure) PrependBytes(bytes ...byte) {
	m.Bytes = append(bytes, m.Bytes...)
}

// PrependRemainder places the provided Measure at the beginning of the source Measure.
func (m *Measure) PrependRemainder(remainder Measure) {
	m.PrependBits(remainder.Bits...)   // First the ending bits get prepended
	m.PrependBytes(remainder.Bytes...) // Then the starting bytes
}
