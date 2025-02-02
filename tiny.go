package tiny

// Tiny is a package devoted to making it easier to interface with small bit ranges in Golang. Some of the bit
// ranges already have defined terms, but I've filled in the gaps with as fitting of terms as I could find.
//
// Bit Range Terminology:
// Width  Range  Name
//   1      0-1    Bit
//   2      0-3    Crumb
//   3      0-7    Note
//   4     0-15    Nibble
//   5     0-31    Flake
//   6     0-63    Morsel
//   7    0-127    Shred

// A Bit represents one binary value [0 - 1]
type Bit byte

// A Crumb represents two binary values [0-3]
type Crumb byte

// A Note represents three binary values [0-7]
type Note byte

// A Nibble represents four binary values [0-15]
type Nibble byte

// A Flake represents five binary values [0-31]
type Flake byte

// A Morsel represents six binary values [0-63]
type Morsel byte

// A Shred represents seven binary values [0-127]
type Shred byte

// Zero -> 0
const Zero Bit = 0

// One -> 1
const One Bit = 1

// ZeroZero -> 00
const ZeroZero Crumb = 0

// ZeroOne -> 01
const ZeroOne Crumb = 1

// OneZero -> 10
const OneZero Crumb = 2

// OneOne -> 11
const OneOne Crumb = 3

// CrumbMax -> 3
const CrumbMax = 3

// NoteMax -> 7
const NoteMax = 7

// NibbleMax -> 15
const NibbleMax = 15

// FlakeMax -> 31
const FlakeMax = 31

// MorselMax -> 63
const MorselMax = 63

// ShredMax -> 127
const ShredMax = 127

// ByteMax -> 255
const ByteMax = 255

// A Remainder is used to efficiently store Bits in operating memory.  In Golang, all types are
// sized around 8-bits (a byte) - thus, every instance of the Bit type takes up 8 bits of operational memory.
// Because of this, we only operate at the Bit level when necessary. The Bytes field holds the majority of the
// information, while the Bits field holds the remaining bits that didn't fit into a multiple of 8 in size.
type Remainder struct {
	Bytes []byte
	Bits  []Bit
}

// NewRemainder initializes a new instance of a Remainder type with empty slices.
func NewRemainder() Remainder {
	return Remainder{[]byte{}, []Bit{}}
}

// ToBytes takes in binary data and returns it in Remainder form.
func ToBytes(bits []Bit) Remainder {
	// The resulting slice of bytes
	var bytes []byte

	// Process bits in groups of 8 (1 byte)
	for i := 0; i+7 < len(bits); i += 8 {
		var b byte
		for j := 0; j < 8; j++ {
			b = (b << 1) | byte(bits[i+j]) // Shift left and add the next bit
		}
		bytes = append(bytes, b) // Add the full byte to the result
	}

	// Extract any leftover bits that don't form a complete byte
	remainingBits := bits[len(bytes)*8:]

	return Remainder{bytes, remainingBits}
}

// ToBits takes a generic type and returns its constituent bits.
func ToBits(value int) []Bit {
	if value == 0 {
		return []Bit{Bit(0)}
	}

	var bits []Bit
	for value > 0 {
		bit := Bit(value % 2)    // Get the least significant Bit
		bits = append(bits, bit) // Append the Bit
		value /= 2               // Shift right by dividing by 2
	}

	// Reverse the slice
	for left, right := 0, len(bits)-1; left < right; left, right = left+1, right-1 {
		bits[left], bits[right] = bits[right], bits[left]
	}

	return bits
}

// BytesToBits takes a slice of bytes and returns a slice of all of its individual bits
func BytesToBits(data []byte) []Bit {
	bits := make([]Bit, 0, len(data)*8)
	for _, b := range data {
		bits = append(bits, ToBits(int(b))...)
	}
	return bits
}
