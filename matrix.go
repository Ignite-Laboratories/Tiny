package tiny

// Table provides access to bit expression from 2D binary types.
//
// Matrix.Position - yourSlice[pos] - Reads the provided index position of your slice.
//
// Matrix.All - yourSlice[:] - Reads the entirety of your slice.
//
// Matrix.From - yourSlice[low:] - Reads from the provided index to the end of your slice.
//
// Matrix.To - yourSlice[:high] - Reads to the provided index from the start of your slice.
//
// Matrix.Between - yourSlice[low:high] - Reads between the provided indexes of your slice.
//
// Matrix.Between - yourSlice[low:high:mid] - Reads between the provided indexes of your slice up to the provided maximum.
//
// Matrix.Gate - Performs a logical operation for every bit of your slice.
var Table Matrix

// Matrix is a type of Expression that indicates to emit that its variadic input is a collection of binary rows, rather than linear information.
type Matrix Expression

var matrixRead = func(i int, column ...Bit) ([]Bit, int) {
	return column, 0
}

// First - yourSlice[0] - Reads the first index position of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Matrix) First(reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_matrix:      &True,
		_first:       &True,
		_reverse:     &isReverse,
		_matrixLogic: &matrixRead,
	}
}

// Last - yourSlice[𝑛 - 1] - Reads the last index position of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Matrix) Last(reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_matrix:      &True,
		_last:        &True,
		_reverse:     &isReverse,
		_matrixLogic: &matrixRead,
	}
}

// Position - yourSlice[pos] - Reads the provided index position of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Matrix) Position(pos uint, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_matrix:      &True,
		_pos:         &pos,
		_reverse:     &isReverse,
		_matrixLogic: &matrixRead,
	}
}

// From - yourSlice[low:] - Reads from the provided index to the end of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Matrix) From(low uint, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_matrix:      &True,
		_low:         &low,
		_reverse:     &isReverse,
		_matrixLogic: &matrixRead,
	}
}

// To - yourSlice[:high] - Reads to the provided index from the start of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Matrix) To(high uint, reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_matrix:      &True,
		_high:        &high,
		_reverse:     &isReverse,
		_matrixLogic: &matrixRead,
	}
}

// Between - yourSlice[low:high] - Reads between the provided indexes of your slice in most→to→least significant order.
//
// Between - yourSlice[low:high:mid] - Reads between the provided indexes of your slice up to the provided maximum in most→to→least significant order.
func (_ Matrix) Between(low uint, high uint, max ...uint) Expression {
	var m *uint
	if len(max) > 0 {
		m = &max[0]
	}

	return Expression{
		_matrix:      &True,
		_low:         &low,
		_high:        &high,
		_max:         m,
		_matrixLogic: &matrixRead,
	}
}

// BetweenReverse - yourSlice[low:high] - Reads between the provided indexes of your slice in least←to←most significant order.
//
// BetweenReverse - yourSlice[low:high:mid] - Reads between the provided indexes of your slice up to the provided maximum in least←to←most significant order.
func (_ Matrix) BetweenReverse(low uint, high uint, max ...uint) Expression {
	var m *uint
	if len(max) > 0 {
		m = &max[0]
	}

	return Expression{
		_matrix:      &True,
		_low:         &low,
		_high:        &high,
		_max:         m,
		_reverse:     &True,
		_matrixLogic: &matrixRead,
	}
}

// All - yourSlice[:] - Reads the entirety of your slice.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Matrix) All(reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_matrix:      &True,
		_reverse:     &isReverse,
		_matrixLogic: &matrixRead,
	}
}

// Gate - Pads all operands to equal length and calls the provided logic gate function for every column of bits to produce an output column.
//
// Expression operations happen in most→to→least significant order - if you would like least←to←most order, please indicate "reverse".
func (_ Matrix) Gate(logic func(int, ...Bit) ([]Bit, int), reverse ...bool) Expression {
	isReverse := len(reverse) > 0 && reverse[0]
	return Expression{
		_matrix:      &True,
		_matrixLogic: &logic,
		_reverse:     &isReverse,
	}
}
