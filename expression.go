package tiny

// Expression represents the standard slice index accessor pattern, and expressions can be generated from the global Read variable.
//
// Position - yourSlice[pos] - Reads the provided index position of your slice.
//
// All - yourSlice[:] - Reads the entirety of your slice.
//
// From - yourSlice[low:] - Reads from the provided index to the end of your slice.
//
// To - yourSlice[:high] - Reads to the provided index from the start of your slice.
//
// Between - yourSlice[low:high] - Reads between the provided indexes of your slice.
//
// Between - yourSlice[low:high:mid] - Reads between the provided indexes of your slice up to the provided maximum.
type Expression struct {
	pos   *uint
	low   *uint
	high  *uint
	max   *uint
	first *bool
	last  *bool
}

// Read provides access to slice index accessor expressions.
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
var Read Expression

// Express expresses the expression against the provided slice.
func Express[T any](expr Expression, slice []T) []T {
	if expr.max != nil && (expr.low == nil || expr.high == nil) {
		panic("invalid slice expression: max requires both low and high to be set")
	}

	switch {
	case expr.first != nil:
		return []T{slice[0]}
	case expr.last != nil:
		return []T{slice[len(slice)-1]}
	case expr.pos != nil:
		return []T{slice[*expr.pos]}
	case expr.low == nil && expr.high == nil:
		return slice[:]
	case expr.low == nil && expr.high != nil:
		return slice[:*expr.high]
	case expr.low != nil && expr.high == nil:
		return slice[*expr.low:]
	case expr.low != nil && expr.high != nil:
		if expr.max == nil {
			return slice[*expr.low:*expr.high]
		}
		return slice[*expr.low:*expr.high:*expr.max]
	default:
		panic("invalid slice expression")
	}
}

// First - yourSlice[0] - Reads the first index position of your slice.
func (_ Expression) First(pos uint) Expression {
	return Expression{
		pos:   nil,
		low:   nil,
		high:  nil,
		max:   nil,
		first: &True,
		last:  nil,
	}
}

// Last - yourSlice[𝑛 - 1] - Reads the last index position of your slice.
func (_ Expression) Last(pos uint) Expression {
	return Expression{
		pos:   nil,
		low:   nil,
		high:  nil,
		max:   nil,
		first: nil,
		last:  &True,
	}
}

// Position - yourSlice[pos] - Reads the provided index position of your slice.
func (_ Expression) Position(pos uint) Expression {
	return Expression{
		pos:   &pos,
		low:   nil,
		high:  nil,
		max:   nil,
		first: nil,
		last:  nil,
	}
}

// From - yourSlice[low:] - Reads from the provided index to the end of your slice.
func (_ Expression) From(low uint) Expression {
	return Expression{
		pos:   nil,
		low:   &low,
		high:  nil,
		max:   nil,
		first: nil,
		last:  nil,
	}
}

// To - yourSlice[:high] - Reads to the provided index from the start of your slice.
func (_ Expression) To(high uint) Expression {
	return Expression{
		pos:   nil,
		low:   nil,
		high:  &high,
		max:   nil,
		first: nil,
		last:  nil,
	}
}

// Between - yourSlice[low:high] - Reads between the provided indexes of your slice.
//
// Between - yourSlice[low:high:mid] - Reads between the provided indexes of your slice up to the provided maximum.
func (_ Expression) Between(low uint, high uint, max ...uint) Expression {
	var m *uint
	if len(max) > 0 {
		m = &max[0]
	}

	return Expression{
		pos:   nil,
		low:   &low,
		high:  &high,
		max:   m,
		first: nil,
		last:  nil,
	}
}

// All - yourSlice[:] - Reads the entirety of your slice.
func (_ Expression) All() Expression {
	return Expression{
		pos:   nil,
		low:   nil,
		high:  nil,
		max:   nil,
		first: nil,
		last:  nil,
	}
}
