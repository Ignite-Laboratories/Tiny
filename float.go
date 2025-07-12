package tiny

// Float represents a phrase where the first measurement is a sign, the next measurement is the exponent,
// and the remaining bits are the mantissa.  Effectively, a practically infinite amount of addressable precision.
type Float Phrase

// Float32 represents a phrase where the first measurement is a sign, the next eight are the exponent,
// and the remaining twenty-three are the mantissa. See IEEE 754.
type Float32 Float

// Float64 represents a phrase where the first measurement is a sign, the next eleven are the exponent,
// and the remaining fifty-two are the mantissa. See IEEE 754.
type Float64 Float

// Float128 represents a phrase where the first measurement is a sign, the next fifteen are the exponent,
// and the remaining one-hundred-and-twelve are the mantissa. See IEEE 754.
type Float128 Float

// Float256 represents a phrase where the first measurement is a sign, the next nineteen are the exponent,
// and the remaining two-hundred-and-thirty-six is the mantissa. See IEEE 754.
type Float256 Float
