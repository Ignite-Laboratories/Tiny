package tiny

// Natural represents a phrase holding a single measurement value belonging to the set of natural numbers, or all positive whole numbers (including zero.)
//
// To those who think zero shouldn't be included in the set of natural numbers, I present a counter-argument:
// Base 1 has only one identifier, meaning it can only "represent" zero by -not- holding a value in an observable
// location.  Subsequently, all bases are built upon determining the size of a value through "identification" - in
// binary, through a series of zeros or ones, in decimal through the digits 0-9.
//
// Now here's where it gets tricky: a value cannot even EXIST until it is given a place to exist within, meaning its
// existence directly implies a void which has now been filled - an identifiable "zero" state.  In fact, the very first
// identifier of all higher order bases (zero) specifically identifies this state!  Counting, itself, comes from the act of observing
// the general relativistic -presence- of anything - fingers, digits, different length squiggles, feelings - meaning to exclude
// zero attempts to redefine the very fundamental definition of identification itself: it's PERFECTLY reasonable to -naturally-
// count -zero- hairs on a magnificently bald head!
//
//	tl;dr - to count naturally involves identification, including identifying -non-existence- as a countable state!
//
// I should note this entire system hinges on one fundamental flaw - this container technically holds one additional value beyond
// the 'natural' number set: nil! Technically, until a number occupies a location, that space holds a 'nil' value in all bases
// above base 1, which observes it as the value 'zero'.  When factoring this trait in, I call it the "programmatic set" of
// numbers.  I can't stop you from setting your natural phrase to it - but I can empower you with awareness of it =)
type Natural struct {
	Phrase
}

// The below is for reference when converting to back to an IEEE 754 float:

// Float32 represents a 32-bit phrase encoded as three measurements in accordance with IEEE 754 -
//
//	    Sign: 1 bit
//	Exponent: 8 bits
//	Mantissa: 23 bits
//type Float32 Phrase

// Float64 represents a 64-bit phrase encoded as three measurements in accordance with IEEE 754 -
//
//	    Sign: 1 bit
//	Exponent: 11 bits
//	Mantissa: 52 bits
//type Float64 Phrase

// Float128 represents a 128-bit phrase encoded as three measurements in accordance with IEEE 754 -
//
//	    Sign: 1 bit
//	Exponent: 15 bits
//	Mantissa: 112 bits
//type Float128 Phrase

// Float256 represents a 256-bit phrase encoded as three measurements in accordance with IEEE 754 -
//
//	    Sign: 1 bit
//	Exponent: 19 bits
//	Mantissa: 236 bits
//type Float256 Phrase
