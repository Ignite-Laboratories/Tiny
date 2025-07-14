package tiny

// Bits provides access to bit expression from binary types.
//
// Expression.Position - yourSlice[pos] - Reads the provided index position of your slice.
//
// Expression.All - yourSlice[:] - Reads the entirety of your slice.
//
// Expression.From - yourSlice[low:] - Reads from the provided index to the end of your slice.
//
// Expression.To - yourSlice[:high] - Reads to the provided index from the start of your slice.
//
// Expression.Between - yourSlice[low:high] - Reads between the provided indexes of your slice.
//
// Expression.Between - yourSlice[low:high:mid] - Reads between the provided indexes of your slice up to the provided maximum.
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
	_matrix  *bool
	_logic   *LogicFunc
}

// LogicFunc takes in a single bit, and its index, and returns an output bit.
type LogicFunc func(int, Bit) Bit

// First - yourSlice[0] - Reads the first index position of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) First(reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_first:   &True,
		_reverse: &isReverse,
	}
}

// Last - yourSlice[𝑛 - 1] - Reads the last index position of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) Last(reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_last:    &True,
		_reverse: &isReverse,
	}
}

// Position - yourSlice[pos] - Reads the provided index position of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) Position(pos uint, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_pos:     &pos,
		_reverse: &isReverse,
	}
}

// From - yourSlice[low:] - Reads from the provided index to the end of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) From(low uint, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_low:     &low,
		_reverse: &isReverse,
	}
}

// To - yourSlice[:high] - Reads to the provided index from the start of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Expression) To(high uint, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_high:    &high,
		_reverse: &isReverse,
	}
}

// Between - yourSlice[low:high] - Reads between the provided indexes of your slice in most→to→least significant order.
//
// Between - yourSlice[low:high:mid] - Reads between the provided indexes of your slice up to the provided maximum in most→to→least significant order.
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

// BetweenReverse - yourSlice[low:high] - Reads between the provided indexes of your slice in least←to←most significant order.
//
// BetweenReverse - yourSlice[low:high:mid] - Reads between the provided indexes of your slice up to the provided maximum in least←to←most significant order.
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

// All - yourSlice[:] - Reads the entirety of your slice.
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
func (_ Expression) Gate(logic func(int, Bit) Bit, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_logic:   &logic,
		_reverse: &isReverse,
	}
}
