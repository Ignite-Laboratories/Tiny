package tiny

// Real represents an Operable Phrase that represents a "real number" - real numbers are held in four measurements:
//
//	  ⬐ The Sign              ⬐ The Fractional Part
//	| 1 - 1 0 1 1 0 - 1 0 0 0 1 0 1 0 1 1 - 1 1 0 | ( -22.5556̅ )
//	          ⬑ The Whole Part                ⬑ The Periodic Part
//
// All parts (except the sign) can grow in arbitrary widths to accommodate whatever size number you can imagine.  To put
// this into perspective, a 256 bit number can hold up to the value 1.1579208923731619542357098500869e+77!  A single
// gigabyte of memory can hold up to 3.2e-8 256-bit real numbers - meaning you've got plenty of storage in the modern age
// to work with =)
//
// By default, real numbers are given a maximum combined decimal precision bit width of 256 bits - but you may override that if desired.
//
// After every arithmetic operation, a check is performed to see if the periodic part is missing and if the fractional part
// fills the entire allotted precision - if so, the real number is deemed to be "irrational".
//
// The REASON to work with a type like this is to ensure that all arithmetic is done as -Math- intended, not within the bounds of
// tight computational memory spaces.  By knowing exactly where the decimal point is located, all reals can be aligned implicitly
// by the matrix engine without first performing what's called 'type coercion' (or defining how to switch between numeric encoding
// schemes on the fly).  The need to differentiate between floating point and integer numbers is entirely a computer science issue
// born from ancient memory requirements and NOT one that a mathematician should have to bear while exploring their theories.
//
// See Natural, Complex, Index, and Binary
type Real struct {
	// Precision represents the maximum combined bit-width of any part of the real number beyond the decimal place.
	//
	// NOTE: This defaults to 256 bits.
	Precision int // Defaults to 256

	// Irrational is true when the number continues on indefinitely with no observed repetition up to the defined precision.
	Irrational bool

	// Negative represents the sign of the real number - with true representing negative.
	Negative bool

	// Whole represents the whole part of the real number.
	Whole Natural

	// Fractional represents the decimal portion of the real number.
	Fractional Natural

	// Periodic represents the periodic end of the fractional portion of the real number and may or may not be present.
	Periodic Natural
}
