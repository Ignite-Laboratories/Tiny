package tiny

// Remainder is used to efficiently store Bits in operating memory, as the
// language inherently requires at least 8 bits to store any custom type.
type Remainder struct {
	// Bytes holds complete byte data.
	Bytes []byte
	// Bits holds any remaining bits that didn't fit into a byte for whatever reason.
	Bits []Bit
}

// NewRemainder constructs a Remainder, which represents an uneven amount of binary bits.
func NewRemainder(bytes []byte, bits ...Bit) Remainder {
	return Remainder{
		Bytes: bytes,
		Bits:  bits,
	}
}

// AppendBits places the provided bits at the end of the source Remainder.
func (source *Remainder) AppendBits(bits ...Bit) {
	source.Bits = append(source.Bits, bits...)          // Add the bits to the last remainder
	toAdd := To.Bytes(source.Bits...)                   // Convert that to byte form
	source.Bytes = append(source.Bytes, toAdd.Bytes...) // Combine the whole bytes
	source.Bits = toAdd.Bits                            // Bring forward the remaining bits
}

// AppendBytes places the provided bytes at the end of the source Remainder.
func (source *Remainder) AppendBytes(bytes ...byte) {
	remainderLength := len(source.Bits)
	lastRemainder := source.Bits
	for _, lastByte := range bytes {
		bits := From.Byte(lastByte)
		newBits := append(lastRemainder, bits[:remainderLength]...)
		lastRemainder = bits[remainderLength:]
		source.Bytes = append(source.Bytes, To.Byte(newBits...))
	}
	source.Bits = lastRemainder
}

// AppendRemainder places the provided Remainder at the end of the source Remainder.
func (source *Remainder) AppendRemainder(remainder Remainder) {
	source.AppendBytes(remainder.Bytes...)
	source.AppendBits(remainder.Bits...)
}

// PrependBits places the provided bits at the beginning of the source Remainder.
func (source *Remainder) PrependBits(bits ...Bit) {
	prependLength := len(bits)
	currentBits := bits
	newBytes := make([]byte, 0)
	for _, nextByte := range source.Bytes {
		bits = From.Byte(nextByte)                                // Get the next byte
		newBits := append(currentBits, bits[:8-prependLength]...) // Grab the bits to complete the current byte out of it
		currentBits = bits[8-prependLength:]                      // Grab the remaining bits as what to prepend next
		newBytes = append(newBytes, To.Byte(newBits...))          // Add the newly formed bit structure to the bytes
	}
	source.Bytes = newBytes
	source.Bits = currentBits // This is now the last few bits of the original last byte
}

// PrependBytes places the provided bytes at the beginning of the source Remainder.
func (source *Remainder) PrependBytes(bytes ...byte) {
	source.Bytes = append(bytes, source.Bytes...)
}

// PrependRemainder places the provided Remainder at the beginning of the source Remainder.
func (source *Remainder) PrependRemainder(remainder Remainder) {
	source.PrependBits(remainder.Bits...)   // First the ending bits get prepended
	source.PrependBytes(remainder.Bytes...) // Then the starting bytes
}
