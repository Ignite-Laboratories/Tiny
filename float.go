package tiny

// TODO: Change this to a single "Float" phrase type that can be converted TO the desired 754 bit width. Internally, it should track its own precision.

// Float32 represents a 32-bit phrase encoded as three measurements in accordance with IEEE 754 -
//
//	    Sign: 1 bit
//	Exponent: 8 bits
//	Mantissa: 23 bits
type Float32 tiny.Phrase

// Float64 represents a 64-bit phrase encoded as three measurements in accordance with IEEE 754 -
//
//	    Sign: 1 bit
//	Exponent: 11 bits
//	Mantissa: 52 bits
type Float64 tiny.Phrase

// Float128 represents a 128-bit phrase encoded as three measurements in accordance with IEEE 754 -
//
//	    Sign: 1 bit
//	Exponent: 15 bits
//	Mantissa: 112 bits
type Float128 tiny.Phrase

// Float256 represents a 256-bit phrase encoded as three measurements in accordance with IEEE 754 -
//
//	    Sign: 1 bit
//	Exponent: 19 bits
//	Mantissa: 236 bits
type Float256 tiny.Phrase
