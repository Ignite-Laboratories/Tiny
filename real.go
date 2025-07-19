package tiny

// Real represents a phrase that holds a real number.  A real number is held in four parts:
//
//	  ⬐ The Sign              ⬐ The Fractional Part
//	| 1 - 1 0 1 1 0 - 1 0 0 0 1 0 1 0 1 1 - 1 1 0 | ( -22.5556̅ )
//	          ⬑ The Whole Part                ⬑ The Periodic Part
//
// All parts (except the sign) can grow in arbitrary widths to accommodate whatever size number you can imagine.  To put
// this into perspective, a 256 bit number can hold up to the value 1.1579208923731619542357098500869e+77!  A single
// gigabyte of memory can hold up to 3.2e-8 of those sized real numbers - you've got plenty of storage in the modern age
// to work with =)
//
// By default, real numbers are given a maximum -fractional- precision bit width of 256 bits - but you may override that if desired.
//
// The REASON to work with a type like this to ensure that all arithmetic is done as -Math- intended, not within the bounds of
// tight computational memory spaces.  By knowing exactly where the decimal point is located, all reals can be aligned implicitly
// by the matrix engine without first performing what's called 'type coercion' (or defining how to switch between numeric encoding
// schemes on the fly).
//
// After every arithmetic operation a check is performed to see if the periodic part is missing and if the fractional part
// fills the entire allotted maximum precision - if so, the real number is deemed to be "irrational".
type Real struct {
	Phrase
	Precision  uint
	Irrational bool
}
