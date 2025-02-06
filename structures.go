package tiny

const Zero Bit = 0
const One Bit = 1
const ZeroZero Crumb = 0
const ZeroOne Crumb = 1
const OneZero Crumb = 2
const OneOne Crumb = 3

const MaxCrumb = 3
const MaxNote = 7
const MaxNibble = 15
const MaxFlake = 31
const MaxMorsel = 63
const MaxShred = 127
const MaxByte = 255

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

// A Remainder is used to efficiently store Bits in operating memory.  In Golang, all types are
// sized around 8-bits (a byte) - thus, every instance of the Bit type takes up 8 bits of operational memory.
// Because of this, we only operate at the Bit level when necessary. The Bytes field holds the majority of the
// information, while the Bits field holds the remaining bits that didn't fit into a multiple of 8 in size.
type Remainder struct {
	Bytes []byte
	Bits  []Bit
}

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

// Count represents the count of 1s and 0s within binary data.
type Count struct {
	Zeros             uint64
	Ones              uint64
	Total             uint64
	Shade             Shade
	PredominantlyDark bool
}

func shadeCount(count Count) Count {
	count.PredominantlyDark = count.Ones > count.Total/2

	if count.Zeros == 0 && count.Ones >= 0 {
		count.Shade = Light
	} else if count.Zeros > 0 && count.Ones == 0 {
		count.Shade = Dark
	} else {
		count.Shade = Grey
	}

	return count
}
