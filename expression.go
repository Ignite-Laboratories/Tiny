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
	pos      *uint
	low      *uint
	high     *uint
	max      *uint
	first    *bool
	last     *bool
	selector *func(int, any) bool
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
	case expr.selector != nil:
		out := make([]T, 0)
		for i, v := range slice {
			if (*expr.selector)(i, v) {
				out = append(out, v)
			}
		}
		return out
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

// Select - Reads every bit and calls your selector function, allowing you to select if the bit should be yielded in the results.
func (_ Expression) Select(selector func(int, any) bool) Expression {
	return Expression{
		selector: &selector,
	}
}

// First - yourSlice[0] - Reads the first index position of your slice.
func (_ Expression) First(pos uint) Expression {
	return Expression{
		first: &True,
	}
}

// Last - yourSlice[𝑛 - 1] - Reads the last index position of your slice.
func (_ Expression) Last(pos uint) Expression {
	return Expression{
		last: &True,
	}
}

// Position - yourSlice[pos] - Reads the provided index position of your slice.
func (_ Expression) Position(pos uint) Expression {
	return Expression{
		pos: &pos,
	}
}

// From - yourSlice[low:] - Reads from the provided index to the end of your slice.
func (_ Expression) From(low uint) Expression {
	return Expression{
		low: &low,
	}
}

// To - yourSlice[:high] - Reads to the provided index from the start of your slice.
func (_ Expression) To(high uint) Expression {
	return Expression{
		high: &high,
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
		low:  &low,
		high: &high,
		max:  m,
	}
}

// All - yourSlice[:] - Reads the entirety of your slice.
func (_ Expression) All() Expression {
	return Expression{}
}
