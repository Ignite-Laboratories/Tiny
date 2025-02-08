package tiny

// A SubByte is any type representable in less than 8 bits.
type SubByte interface {
	Bit | Crumb | Note | Nibble | Flake | Morsel | Shred
	String() string
}

// A Bit represents one binary value. [0 - 1]
type Bit byte

// A Crumb represents two binary values. [0-3]
type Crumb byte

// A Note represents three binary values. [0-7]
type Note byte

// A Nibble represents four binary values. [0-15]
type Nibble byte

// A Flake represents five binary values. [0-31]
type Flake byte

// A Morsel represents six binary values. [0-63]
type Morsel byte

// A Shred represents seven binary values. [0-127]
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

// A Remainder is used to efficiently store Bits in operating memory.  In Go, all types are sized around
// 8-bits (a byte) - thus, every instance of the Bit type takes up 8 bits of operational memory. Because of
// this, we only operate at the Bit level when necessary. The Bytes field holds the majority of the
// information, while the Bits field holds the remaining bits that didn't fit into a byte.
type Remainder struct {
	Bytes []byte
	Bits  []Bit
}

// Count represents the count of 1s and 0s within binary data.
type Count struct {
	Zeros             int
	Ones              int
	Total             int
	Shade             Shade
	PredominantlyDark bool
	Distribution      [8]int
}

// Calculate fills in a Count's Shade information.
func (c *Count) Calculate() {
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
