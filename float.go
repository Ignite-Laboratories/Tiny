package tiny

type floatPhrase Phrase

// Float32Phrase represents a phrase where the first measurement is a sign, the next eight are the exponent,
// and the remaining twenty-three are the mantissa. See IEEE 754.
type Float32Phrase floatPhrase

// Float64Phrase represents a phrase where the first measurement is a sign, the next eleven are the exponent,
// and the remaining fifty-two are the mantissa. See IEEE 754.
type Float64Phrase floatPhrase

// Float128Phrase represents a phrase where the first measurement is a sign, the next fifteen are the exponent,
// and the remaining one-hundred-and-twelve are the mantissa. See IEEE 754.
type Float128Phrase floatPhrase

// Float256Phrase represents a phrase where the first measurement is a sign, the next nineteen are the exponent,
// and the remaining two-hundred-and-thirty-six is the mantissa. See IEEE 754.
type Float256Phrase floatPhrase

// FloatBigPhrase represents a phrase where the first measurement is a sign, the next measurement is the exponent,
// and the remaining bits are the mantissa.  Effectively - infinite precision.
type FloatBigPhrase floatPhrase
