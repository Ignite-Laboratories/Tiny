package tiny

/**
ZLE
*/

type ZLEScheme interface {
	KeyBitWidth() int
	Read(Phrase) (value int, remainder Phrase)
	Encode(int) (key Phrase, projection Phrase)
}

/**
Shade
*/

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

// BinaryShade is a count of the number of 0s and 1s within binary data.
type BinaryShade struct {
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

func (c *BinaryShade) combine(other BinaryShade) {
	c.Zeros += other.Zeros
	c.Ones += other.Ones
	c.Total += other.Total
	for i := range c.Distribution {
		c.Distribution[i] += other.Distribution[i]
	}
	c.calculate()
}

func (c *BinaryShade) calculate() {
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

// Upcast is a convenience method to upcast slices of SubByte to byte.
func Upcast[TIn SubByte](data []TIn) []byte {
	out := make([]byte, len(data))
	for i, bit := range data {
		out[i] = byte(bit)
	}
	return out
}

// Bit represents one binary place. [0 - 1]
type Bit byte

// Crumb represents two binary places. [0-3]
type Crumb byte

// Note represents three binary places. [0-7]
type Note byte

// Nibble represents four binary places. [0-15]
type Nibble byte

// Flake represents five binary places. [0-31]
type Flake byte

// Morsel represents six binary places. [0-63]
type Morsel byte

// Shred represents seven binary places. [0-127]
type Shred byte

/**
Over Byte Types
*/

// Scale represents twelve binary places - See MaxScale.
type Scale int16

// Motif represents sixteen binary places - See MaxMotif.
type Motif int16

// Riff represents twenty-four binary places - See MaxRiff.
type Riff int32

// Cadence represents thirty-two binary places - See MaxCadence.
type Cadence int32

// Hook represents forty-eight binary places - See MaxHook.
type Hook int64

// Melody represents sixty-four binary places - See MaxMelody.
type Melody int64

// Verse represents one-hundred & twenty-eight binary places - See MaxVerse.
type Verse [2]int64

/**
Bits()
*/

// Bits uses the provided value to build a 1 Bit slice.
func (v Bit) Bits() []Bit {
	return From.Number(int(v))
}

// Bits uses the provided value to build a 2 Bit Crumb slice.
func (v Crumb) Bits() []Bit {
	return From.Crumb(v)
}

// Bits uses the provided value to build a 3 Bit Note slice.
func (v Note) Bits() []Bit {
	return From.Note(v)
}

// Bits uses the provided value to build a 4 Bit Nibble slice.
func (v Nibble) Bits() []Bit {
	return From.Nibble(v)
}

// Bits uses the provided value to build a 5 Bit Flake slice.
func (v Flake) Bits() []Bit {
	return From.Flake(v)
}

// Bits uses the provided value to build a 6 Bit Morsel slice.
func (v Morsel) Bits() []Bit {
	return From.Morsel(v)
}

// Bits uses the provided value to build a 7 Bit Shred slice.
func (v Shred) Bits() []Bit {
	return From.Shred(v)
}

// Bits uses the provided value to build a 12 Bit Scale slice.
func (v Scale) Bits() []Bit {
	return From.Scale(v)
}

// Bits uses the provided value to build a 16 Bit Motif slice.
func (v Motif) Bits() []Bit {
	return From.Motif(v)
}

// Bits uses the provided value to build a 24 Bit Riff slice.
func (v Riff) Bits() []Bit {
	return From.Riff(v)
}

// Bits uses the provided value to build a 32 Bit Cadence slice.
func (v Cadence) Bits() []Bit {
	return From.Cadence(v)
}

// Bits uses the provided value to build a 48 Bit Hook slice.
func (v Hook) Bits() []Bit {
	return From.Hook(v)
}

// Bits uses the provided value to build a 64 Bit Melody slice.
func (v Melody) Bits() []Bit {
	return From.Melody(v)
}

// Bits uses the provided value to build a 128 Bit Verse slice.
func (v Verse) Bits() []Bit {
	return From.Verse(v)
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

// String converts a Scale to a 12-bit string.
func (v Scale) String() string {
	return To.String(v.Bits()...)
}

// String converts a Motif to a 16-bit string.
func (v Motif) String() string {
	return To.String(v.Bits()...)
}

// String converts a Riff to a 24-bit string.
func (v Riff) String() string {
	return To.String(v.Bits()...)
}

// String converts a Cadence to a 32-bit string.
func (v Cadence) String() string {
	return To.String(v.Bits()...)
}

// String converts a Hook to a 48-bit string.
func (v Hook) String() string {
	return To.String(v.Bits()...)
}

// String converts a Melody to a 64-bit string.
func (v Melody) String() string {
	return To.String(v.Bits()...)
}

// String converts a Verse to a 128-bit string.
func (v Verse) String() string {
	return To.String(v.Bits()...)
}
