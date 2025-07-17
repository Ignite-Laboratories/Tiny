package emit

import "tiny"

// Positions [𝑛₀,𝑛₁,𝑛₂...] creates a tiny.Expression which will read the provided index positions of your binary information in most→to→least significant order - regardless of the provided variadic order.
func Positions(positions ...uint) tiny.Expression {
	return tiny.Expression{
		Positions: &positions,
	}
}

// PositionsReverse [𝑛₀,𝑛₁,𝑛₂...] creates a tiny.Expression which will read the provided index positions of your binary information in least←to←most significant order - regardless of the provided variadic order.
func PositionsReverse(positions ...uint) tiny.Expression {
	return tiny.Expression{
		Positions: &positions,
		Reverse:   &tiny.True,
	}
}

// Width [𝑛] creates a tiny.Expression which will read the provided bit width in most→to→least significant order.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func Width(width ...uint) tiny.Expression {
	w := tiny.Unlimited
	if len(width) >= 0 {
		w = width[0]
	}
	return tiny.Expression{
		Low:  &tiny.Start,
		High: &w,
	}
}

// WidthReverse [𝑛] creates a tiny.Expression which will read the provided bit width in least←to←most significant order.
func WidthReverse(width ...uint) tiny.Expression {
	w := tiny.Unlimited
	if len(width) >= 0 {
		w = width[0]
	}
	return tiny.Expression{
		Low:     &tiny.Start,
		High:    &w,
		Reverse: &tiny.True,
	}
}

// First [0] creates a tiny.Expression which will read the first index position of your binary information.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func First(reverse ...bool) tiny.Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return tiny.Expression{
		Positions: &tiny.Initial,
		Reverse:   &isReverse,
	}
}

// Last [𝑛 - 1] creates a tiny.Expression which will read the last index position of your binary information.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func Last(reverse ...bool) tiny.Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return tiny.Expression{
		Last:    &tiny.True,
		Reverse: &isReverse,
	}
}

// From [low:] creates a tiny.Expression which will read from the provided index to the end of your binary information.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func From(low uint, reverse ...bool) tiny.Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return tiny.Expression{
		Low:     &low,
		Reverse: &isReverse,
	}
}

// To [:high] creates a tiny.Expression which will read to the provided index from the start of your binary information.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func To(high uint, reverse ...bool) tiny.Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return tiny.Expression{
		High:    &high,
		Reverse: &isReverse,
	}
}

// Between [low:high:*] creates a tiny.Expression which will read between the provided indexes of your binary information up to the provided maximum in most→to→least significant order.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func Between(low uint, high uint, reverse ...bool) tiny.Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return tiny.Expression{
		Low:     &low,
		High:    &high,
		Reverse: &isReverse,
	}
}

// All [:] creates a tiny.Expression which will read the entirety of your binary information.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func All(reverse ...bool) tiny.Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return tiny.Expression{
		Reverse: &isReverse,
	}
}

// NOT creates a tiny.Expression which will apply the below truth table against every input bit.
//
// NOTE: If no bits are provided, Zero is returned.
//
//	"The NOT Truth Table"
//
//	        𝑎 | 𝑜𝑢𝑡
//	        0 | 1
//	        1 | 0
func NOT(reverse ...bool) tiny.Expression {
	return Gate(tiny.Logic.NOT, reverse...)
}

// Gate creates a tiny.Expression which will apply the provided logic gate against every input bit.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func Gate(logic tiny.BitLogicFunc, reverse ...bool) tiny.Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return tiny.Expression{
		BitLogic: &logic,
		Reverse:  &isReverse,
	}
}

// Pattern creates a tiny.Expression which will XOR the provided pattern against the input bits in most→to→least significant order.
func Pattern(pattern ...tiny.Bit) tiny.Expression {
	logic := patternLogic(pattern...)
	return tiny.Expression{
		BitLogic: &logic,
	}
}

// PatternReverse creates a tiny.Expression which will XOR the provided pattern against the input bits in least←to←most significant order.
func PatternReverse(pattern ...tiny.Bit) tiny.Expression {
	logic := patternLogic(pattern...)
	return tiny.Expression{
		BitLogic: &logic,
		Reverse:  &tiny.True,
	}
}

func patternLogic(pattern ...tiny.Bit) tiny.BitLogicFunc {
	limit := len(pattern)
	step := 0
	return func(i uint, operands ...tiny.Bit) ([]tiny.Bit, *tiny.Phrase) {
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
