package tiny

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
