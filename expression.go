package tiny

// Bits provides access to Bit Expression from binary types. This process walks a cursor across the binary information and selectively yields bits according to one of the below logical expressions.
//
// NOTE: When calling Emit, you may also provide a maximum number of bits to be emitted with the expression.
//
// Expression.Positions[𝑛₀,𝑛₁,𝑛₂,𝑛₃] - reads the provided index positions of your binary information in most→to→least significant order - regardless of the provided variadic order.
//
// Expression.PositionsReverse[𝑛₀,𝑛₁,𝑛₂,𝑛₃] - reads the provided index positions of your binary information in least←to←most significant order - regardless of the provided variadic order.
//
// Expression.All[:] - Reads the entirety of your binary information.
//
// Expression.From[low:] - Reads from the provided index to the end of your binary information.
//
// Expression.To[:high] - Reads to the provided index from the start of your binary information.
//
// Expression.Between[low:high] - Reads between the provided indexes of your binary information.
//
// Expression.Gate - Performs a logical operation for every bit of your binary information.
//
// Expression.Pattern - XORs the provided pattern against the target bits in most→to→least significant order.
//
// Expression.PatternReverse - XORs the provided pattern against the target bits in least←to←most significant order.
var Bits Expression

// Expression represents the standard slice index accessor pattern, and expressions can be generated from the global Read variable.
type Expression struct {
	_positions *[]uint
	_low       *uint
	_high      *uint
	_max       *uint
	_last      *bool
	_reverse   *bool
	_bitLogic  *BitLogicFunc
	_artifact  *ArtifactFunc
	_limit     uint
}

// BitLogicFunc takes in many bits and their collectively shared index and returns an output bit plus a nilable artifact.
type BitLogicFunc func(int, ...Bit) ([]Bit, *Phrase)

// ArtifactFunc applies the artifact from a single round of calculation against the provided operand bits.
type ArtifactFunc func(i int, artifact Phrase, operands ...Phrase) []Phrase

// Positions [𝑛₀,𝑛₁,𝑛₂...] reads the provided index positions of your binary information in most→to→least significant order - regardless of the provided variadic order.
func (_ Expression) Positions(positions ...uint) Expression {
	return Expression{
		_positions: &positions,
	}
}

// PositionsReverse [𝑛₀,𝑛₁,𝑛₂...] reads the provided index positions of your binary information in least←to←most significant order - regardless of the provided variadic order.
func (_ Expression) PositionsReverse(positions ...uint) Expression {
	return Expression{
		_positions: &positions,
		_reverse:   &True,
	}
}

// First [0] reads the first index position of your binary information.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) First(reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	zero := []uint{0}
	return Expression{
		_positions: &zero,
		_reverse:   &isReverse,
	}
}

// Last [𝑛 - 1] reads the last index position of your binary information.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) Last(reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_last:    &True,
		_reverse: &isReverse,
	}
}

// From [low:] reads from the provided index to the end of your binary information.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) From(low uint, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_low:     &low,
		_reverse: &isReverse,
	}
}

// To [:high] reads to the provided index from the start of your binary information.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) To(high uint, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_high:    &high,
		_reverse: &isReverse,
	}
}

// Between [low:high:*] reads between the provided indexes of your binary information up to the provided maximum in most→to→least significant order.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) Between(low uint, high uint, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_low:     &low,
		_high:    &high,
		_reverse: &isReverse,
	}
}

// All [:] reads the entirety of your binary information.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) All(reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_reverse: &isReverse,
	}
}

// Gate - Reads every bit of your binary information and calls the provided logic gate function to manipulate the output bit.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) Gate(logic BitLogicFunc, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_bitLogic: &logic,
		_reverse:  &isReverse,
	}
}

// Pattern - XORs the provided pattern against the target bits of your binary information in most→to→least significant order.
func (_ Expression) Pattern(pattern ...Bit) Expression {
	logic := patternLogic(pattern...)
	return Expression{
		_bitLogic: &logic,
	}
}

// PatternReverse - XORs the provided pattern against the target bits of your binary information in least←to←most significant order.
func (_ Expression) PatternReverse(pattern ...Bit) Expression {
	logic := patternLogic(pattern...)
	return Expression{
		_bitLogic: &logic,
		_reverse:  &True,
	}
}

func patternLogic(pattern ...Bit) BitLogicFunc {
	limit := len(pattern)
	step := 0
	return func(i int, operands ...Bit) ([]Bit, *Phrase) {
		for _, b := range pattern {
			operands[i] = b ^ pattern[i]
		}
		step++
		if step >= limit {
			step = 0
		}
		return operands, nil
	}
}
