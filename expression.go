package tiny

// Bits provides access to bit expression from binary types. Bit Expression uses Go slice accessors under the hood,
// so you can treat the operations as the same.  All operations are performed from most→to→least significant order,
// but you may optionally reverse the order to use least←to←most.
//
// Expression.Position[𝑛] - Reads the provided index position of your slice.
//
// Expression.All[:] - Reads the entirety of your slice.
//
// Expression.From[low:] - Reads from the provided index to the end of your slice.
//
// Expression.To[:high] - Reads to the provided index from the start of your slice.
//
// Expression.Between[low:high] - Reads between the provided indexes of your slice.
//
// Expression.Between[low:high:max] - Reads between the provided indexes of your slice up to the provided maximum.
//
// Expression.BetweenReverse[low:high:*] - Reads the same as "between" (up to an optional maximum) but from least←to←most significant order.
//
// Expression.Gate - Performs a logical operation for every bit of your slice.
var Bits Expression

// Expression represents the standard slice index accessor pattern, and expressions can be generated from the global Read variable.
type Expression struct {
	_pos     *uint
	_low     *uint
	_high    *uint
	_max     *uint
	_first   *bool
	_last    *bool
	_reverse *bool
	_unary   *UnaryFunc
}

// UnaryFunc takes in a single bit, and its index, and returns an output bit.
type UnaryFunc func(int, Bit) Bit

// MatrixFunc takes in many bits and their collectively shared index and returns an output bit plus an artifact.
type MatrixFunc func(int, ...Bit) (Bit, any)

// First [0] reads the first index position of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) First(reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_first:   &True,
		_reverse: &isReverse,
	}
}

// Last [𝑛 - 1] reads the last index position of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) Last(reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_last:    &True,
		_reverse: &isReverse,
	}
}

// Position [𝑛] reads the provided index position of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) Position(pos uint, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_pos:     &pos,
		_reverse: &isReverse,
	}
}

// From [low:] reads from the provided index to the end of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) From(low uint, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_low:     &low,
		_reverse: &isReverse,
	}
}

// To [:high] reads to the provided index from the start of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) To(high uint, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_high:    &high,
		_reverse: &isReverse,
	}
}

// Between [low:high:*] reads between the provided indexes of your slice up to the provided maximum in most→to→least significant order.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) Between(low uint, high uint, max ...uint) Expression {
	var m *uint
	if len(max) > 0 {
		m = &max[0]
	}

	return Expression{
		_low:  &low,
		_high: &high,
		_max:  m,
	}
}

// BetweenReverse [low:high:*] reads between the provided indexes of your slice, up to an optional maximum, in least←to←most significant order.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) BetweenReverse(low uint, high uint, max ...uint) Expression {
	var m *uint
	if len(max) > 0 {
		m = &max[0]
	}

	return Expression{
		_low:     &low,
		_high:    &high,
		_max:     m,
		_reverse: &True,
	}
}

// All [:] reads the entirety of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) All(reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_reverse: &isReverse,
	}
}

// Gate - Reads every bit and calls the provided logic gate function to manipulate the output bit.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) Gate(logic UnaryFunc, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_unary:   &logic,
		_reverse: &isReverse,
	}
}
