package tiny

// Shade is a descriptor of whether the binary data is Light, Dark, or Grey.
type Shade int

const (
	// Light represents all 0s.
	Light Shade = iota

	// Dark represents all 1s.
	Dark

	// Grey represents a mixture of 1s and 0s.
	Grey
)

// BinaryCount is a count of the number of 0s and 1s within binary data.
type BinaryCount struct {
	// Zeros are a count of the number of 0s in the target.
	Zeros int
	// Ones are a count of the number of 1s in the target.
	Ones int
	// Total is a tally of all 1s and 0s.
	Total int
	// Shade details the general nature of the 1s and 0s.
	Shade Shade
	// PredominantlyDark is true if more than half of the binary data is a 1.
	PredominantlyDark bool
	// Distribution is a count of how many ones are in each position of every byte of the target.
	Distribution [8]int
}

func (c *BinaryCount) combine(other BinaryCount) {
	c.Zeros += other.Zeros
	c.Ones += other.Ones
	c.Total += other.Total
	for i := range c.Distribution {
		c.Distribution[i] += other.Distribution[i]
	}
	c.calculate()
}

func (c *BinaryCount) calculate() {
	c.PredominantlyDark = c.Ones > c.Total/2

	if c.Zeros > 0 && c.Ones == 0 {
		// It's all zeros
		c.Shade = Light
	} else if c.Zeros == 0 && c.Ones >= 0 {
		// It's all ones
		c.Shade = Dark
	} else {
		// It's a mixture
		c.Shade = Grey
	}
}

/**
Sub Byte Types
*/

// SubByte is any type representable in less than 8 bits.
type SubByte interface {
	Bit | Crumb | Note | Nibble | Flake | Morsel | Shred
	String() string
}

// upcast is a convenience method to upcast slices of SubByte to byte.
func upcast[TIn SubByte](data []TIn) []byte {
	out := make([]byte, len(data))
	for i, bit := range data {
		out[i] = byte(bit)
	}
	return out
}

// Bit represents one binary value. [0 - 1]
type Bit byte

// Crumb represents two binary values. [0-3]
type Crumb byte

// Note represents three binary values. [0-7]
type Note byte

// Nibble represents four binary values. [0-15]
type Nibble byte

// Flake represents five binary values. [0-31]
type Flake byte

// Morsel represents six binary values. [0-63]
type Morsel byte

// Shred represents seven binary values. [0-127]
type Shred byte

/**
Bits()
*/

// Bits uses the provided value to build a 1 Bit slice.
func (v Bit) Bits() []Bit {
	return From.Number(int(v))
}

// Bits uses the provided value to build a 2 Bit slice.
func (v Crumb) Bits() []Bit {
	return From.Crumb(v)
}

// Bits uses the provided value to build a 3 Bit slice.
func (v Note) Bits() []Bit {
	return From.Note(v)
}

// Bits uses the provided value to build a 4 Bit slice.
func (v Nibble) Bits() []Bit {
	return From.Nibble(v)
}

// Bits uses the provided value to build a 5 Bit slice.
func (v Flake) Bits() []Bit {
	return From.Flake(v)
}

// Bits uses the provided value to build a 6 Bit slice.
func (v Morsel) Bits() []Bit {
	return From.Morsel(v)
}

// Bits uses the provided value to build a 7 Bit slice.
func (v Shred) Bits() []Bit {
	return From.Shred(v)
}

/**
String()
*/

// String converts a Bit to a 1-bit string.
func (v Bit) String() string {
	return To.String(v)
}

// String converts a Crumb to a 2-bit string.
func (v Crumb) String() string {
	return To.String(v.Bits()...)
}

// String converts a Note to a 3-bit string.
func (v Note) String() string {
	return To.String(v.Bits()...)
}

// String converts a Nibble to a 4-bit string.
func (v Nibble) String() string {
	return To.String(v.Bits()...)
}

// String converts a Flake to a 5-bit string.
func (v Flake) String() string {
	return To.String(v.Bits()...)
}

// String converts a Morsel to a 6-bit string.
func (v Morsel) String() string {
	return To.String(v.Bits()...)
}

// String converts a Shred to a 7-bit string.
func (v Shred) String() string {
	return To.String(v.Bits()...)
}
