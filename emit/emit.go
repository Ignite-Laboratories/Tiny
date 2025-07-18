// Package emit provides access to bit expression from binary types. This process walks a cursor across the binary information
// and selectively yields bits according to the rules defined by logical expressions. Expressions follow Go slice index accessor
// rules, meaning the low side is inclusive and the high side is exclusive.
//
// NOTE: You must also provide a maximum number of bits to be emitted with the expression - this may be Unlimited.
//
// Positions[𝑛₀,𝑛₁,𝑛₂,𝑛₃] - reads the provided index positions of your binary information in most→to→least significant order - regardless of the provided variadic order.
//
// PositionsFromEnd[𝑛₀,𝑛₁,𝑛₂,𝑛₃] - reads the provided index positions of your binary information in least←to←most significant order - regardless of the provided variadic order.
//
// All[:] - Reads the entirety of your binary information.
//
// Low[low:] - Reads from the provided index to the end of your binary information.
//
// High[:high] - Reads to the provided index from the start of your binary information.
//
// Between[low:high] - Reads between the provided indexes of your binary information.
//
// Gate - Performs a logical operation for every bit of your binary information.
//
// Pattern - XORs the provided pattern against the target bits in most→to→least significant order.
//
// PatternFromEnd - XORs the provided pattern against the target bits in least←to←most significant order.
package emit

import "tiny"

// Positions [𝑛₀,𝑛₁,𝑛₂...] creates a tiny.Expression which will read the provided index positions of your binary information in most→to→least significant order - regardless of the provided variadic order.
func Positions(positions ...uint) tiny.Expression {
	return tiny.Expression{
		Positions: &positions,
	}
}

// PositionsFromEnd [𝑛₀,𝑛₁,𝑛₂...] creates a tiny.Expression which will read the provided index positions of your binary information in least←to←most significant order - regardless of the provided variadic order.
func PositionsFromEnd(positions ...uint) tiny.Expression {
	return tiny.Expression{
		Positions: &positions,
		Reverse:   &tiny.True,
	}
}

// Width [𝑛] creates a tiny.Expression which will read the provided bit width in most→to→least significant order.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func Width[T tiny.Binary](width uint, operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{
		Low:  &tiny.Start,
		High: &width,
	}, operands...)
}

// WidthFromEnd [𝑛] creates a tiny.Expression which will read the provided bit width in least←to←most significant order.
func WidthFromEnd[T tiny.Binary](width uint, operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{
		Low:     &tiny.Start,
		High:    &width,
		Reverse: &tiny.True,
	}, operands...)
}

// First [0] creates a tiny.Expression which will read the first index position of your binary information.
func First[T tiny.Binary](operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{
		Positions: &tiny.Initial,
	}, operands...)
}

// Last [𝑛 - 1] creates a tiny.Expression which will read the last index position of your binary information.
func Last[T tiny.Binary](operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{
		Last: &tiny.True,
	}, operands...)
}

// Low [low:] creates a tiny.Expression which will read from the provided index to the end of your binary information.
func Low[T tiny.Binary](low uint, operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{
		Low: &low,
	}, operands...)
}

// LowFromEnd [low:] creates a tiny.Expression which will read from the provided index to the end of your binary information.
func LowFromEnd[T tiny.Binary](low uint, operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{
		Low:     &low,
		Reverse: &tiny.True,
	}, operands...)
}

// High [:high] creates a tiny.Expression which will read to the provided index from the start of your binary information.
func High[T tiny.Binary](high uint, operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{
		High: &high,
	}, operands...)
}

// HighFromEnd [:high] creates a tiny.Expression which will read to the provided index from the start of your binary information.
func HighFromEnd[T tiny.Binary](high uint, operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{
		High:    &high,
		Reverse: &tiny.True,
	}, operands...)
}

// Between [low:high:*] creates a tiny.Expression which will read between the provided indexes of your binary information up to the provided maximum in most→to→least significant order.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func Between[T tiny.Binary](low uint, high uint, operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{
		Low:  &low,
		High: &high,
	}, operands...)
}

// BetweenFromEnd [low:high:*] creates a tiny.Expression which will read between the provided indexes of your binary information up to the provided maximum in most→to→least significant order.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func BetweenFromEnd[T tiny.Binary](low uint, high uint, operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{
		Low:     &low,
		High:    &high,
		Reverse: &tiny.True,
	}, operands...)
}

// All [:] creates a tiny.Expression which will read the entirety of your binary information.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func All[T tiny.Binary](operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{}, operands...)
}

// Reversed [:] creates a tiny.Expression which will read the entirety of your binary information.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func Reversed[T tiny.Binary](operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{Reverse: &tiny.True}, operands...)
}

/**
Logic Gates
*/

// Gate creates a tiny.Expression which will apply the provided logic gate against every input bit.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func Gate[T tiny.Binary](logic tiny.BitLogicFunc, operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{
		BitLogic: &logic,
	}, operands...)
}

// GateFromEnd creates a tiny.Expression which will apply the provided logic gate against every input bit.
func GateFromEnd[T tiny.Binary](logic tiny.BitLogicFunc, operands ...T) []tiny.Bit {
	return tiny.Emit(tiny.Expression{
		BitLogic: &logic,
		Reverse:  &tiny.True,
	}, operands...)
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
func NOT[T tiny.Binary](operands ...T) []tiny.Bit {
	return Gate(func(i uint, bits ...tiny.Bit) ([]tiny.Bit, *tiny.Phrase) {
		if len(bits) == 0 {
			return tiny.SingleZero, nil
		}
		for _, b := range bits {
			bits[0] = b ^ 1
		}
		return bits, nil
	}, operands...)
}

/**
Pattern Emission
*/

// Pattern creates a tiny.Expression which will XOR the provided pattern against the input bits in most→to→least significant order.
func Pattern[T tiny.Binary](pattern []tiny.Bit, operands ...T) []tiny.Bit {
	return Gate(patternLogic(pattern...), operands...)
}

// PatternFromEnd creates a tiny.Expression which will XOR the provided pattern against the input bits in least←to←most significant order.
func PatternFromEnd[T tiny.Binary](pattern []tiny.Bit, operands ...T) []tiny.Bit {
	return GateFromEnd(patternLogic(pattern...), operands...)
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
