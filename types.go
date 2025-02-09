package tiny

// SubByte is any type representable in less than 8 bits.
type SubByte interface {
	Bit | Crumb | Note | Nibble | Flake | Morsel | Shred
	String() string
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

// Remainder is used to efficiently store Bits in operating memory.
type Remainder struct {
	// Bytes holds complete byte data.
	Bytes []byte
	// Bits holds any remaining bits that didn't fit into a byte for whatever reason.
	Bits []Bit
}

// BinaryCount is a count of the number of 0s and 1s within the target.
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

// Calculate fills in a BinaryCount's Shade information.
func (c *BinaryCount) Calculate() {
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
