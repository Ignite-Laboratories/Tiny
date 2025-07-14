package tiny

import (
	"fmt"
	"math"
)

// VariadicExpression represents the standard slice index accessor pattern, and expressions can be generated from the global Read variable.
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
//
// LogicGate - Performs a logical operation for every bit of your slice.
type VariadicExpression struct {
	pos   *uint
	low   *uint
	high  *uint
	max   *uint
	first *bool
	last  *bool
	logic *func(int, ...Bit) []Bit
}

// ManyBits provides access to slice index accessor expressions.
//
// VariadicExpression.Position - yourSlice[pos] - Reads the provided index position of your slice.
//
// VariadicExpression.All - yourSlice[:] - Reads the entirety of your slice.
//
// VariadicExpression.From - yourSlice[low:] - Reads from the provided index to the end of your slice.
//
// VariadicExpression.To - yourSlice[:high] - Reads to the provided index from the start of your slice.
//
// VariadicExpression.Between - yourSlice[low:high] - Reads between the provided indexes of your slice.
//
// VariadicExpression.Between - yourSlice[low:high:mid] - Reads between the provided indexes of your slice up to the provided maximum.
//
// VariadicExpression.LogicGate - Performs a logical operation for every bit of your slice.
var ManyBits VariadicExpression

// VariadicEmit expresses the bits of binary information according to logical rules.  A variadic expression will
// operate against as many operands as you provide, and the number of outputs can fluctuate based upon the operation.
func VariadicEmit[T binary](expr VariadicExpression, data ...[]T) [][]Bit {
	if expr.max != nil && (expr.low == nil || expr.high == nil) {
		panic("invalid slice expression: max requires both low and high to be set")
	}

	if len(data) == 0 {
		return make([][]Bit, 0)
	}

	expected := -1
	if expr.high != nil {
		expected = int(*expr.high)

		if expr.low != nil {
			expected = int(*expr.high) - int(*expr.low)
		}

		if expr.max != nil {
			maximum := int(*expr.max)
			if maximum < expected {
				expected = maximum
			}
		}
	}

	slice := make([]Bit, 0, int(math.Min(1<<10, float64(expected))))
	for _, element := range data {
		switch concrete := any(element).(type) {
		case Measurement:
			byteBits := VariadicEmit(expr, concrete.Bytes...)
			slice = append(slice, byteBits...)

			if expr.high != nil {
				high := *expr.high - uint(len(byteBits))
				expr.high = &high
			}
			if expr.low != nil {
				low := *expr.low - uint(len(byteBits))
				if low < 0 {
					low = 0
				}
				expr.low = &low
			}
			slice = append(slice, VariadicEmit(expr, concrete.Bits...)...)
		case Phrase:
			slice = append(slice, VariadicEmit(expr, concrete.Data...)...)
		case []byte:
			slice = append(slice, VariadicEmit(expr, concrete...)...)
		case []Bit:
			capacity := len(concrete)

			switch {
			case expr.logic != nil:
				out := make([]Bit, 0)
				for i, b := range concrete {
					out = append(out, (*expr.logic)(i, b))
				}
				return out
			case expr.first != nil:
				return []Bit{concrete[0]}
			case expr.last != nil:
				return []Bit{concrete[len(slice)-1]}
			case expr.pos != nil:
				return []Bit{concrete[*expr.pos]}
			case expr.low == nil && expr.high == nil:
				return concrete[:]
			case expr.low == nil && expr.high != nil:
				high := int(math.Min(float64(capacity), float64(*expr.high)))
				return concrete[:high]
			case expr.low != nil && expr.high == nil:
				low := int(math.Min(float64(capacity), float64(*expr.low)))
				return concrete[low:]
			case expr.low != nil && expr.high != nil:
				high := int(math.Min(float64(capacity), float64(*expr.high)))
				low := int(math.Min(float64(capacity), float64(*expr.low)))
				if expr.max == nil {
					return concrete[low:high]
				}
				return concrete[low:high:*expr.max]
			default:
				panic("invalid slice expression")
			}
		case byte:
			return VariadicEmit(expr, NewMeasurementFromBytes(concrete).GetAllBits()...)
		default:
			panic(fmt.Errorf("invalid bit expression type: %T", concrete))
		}

		if len(slice) >= expected {
			// Bailout when the requested range has accrued
			return slice[:expected]
		}
	}

	return slice
}

// First - yourSlice[0] - Reads the first index position of your slice.
func (_ VariadicExpression) First(pos uint) VariadicExpression {
	return VariadicExpression{
		first: &True,
	}
}

// Last - yourSlice[𝑛 - 1] - Reads the last index position of your slice.
func (_ VariadicExpression) Last(pos uint) VariadicExpression {
	return VariadicExpression{
		last: &True,
	}
}

// Position - yourSlice[pos] - Reads the provided index position of your slice.
func (_ VariadicExpression) Position(pos uint) VariadicExpression {
	return VariadicExpression{
		pos: &pos,
	}
}

// From - yourSlice[low:] - Reads from the provided index to the end of your slice.
func (_ VariadicExpression) From(low uint) VariadicExpression {
	return VariadicExpression{
		low: &low,
	}
}

// To - yourSlice[:high] - Reads to the provided index from the start of your slice.
func (_ VariadicExpression) To(high uint) VariadicExpression {
	return VariadicExpression{
		high: &high,
	}
}

// Between - yourSlice[low:high] - Reads between the provided indexes of your slice.
//
// Between - yourSlice[low:high:mid] - Reads between the provided indexes of your slice up to the provided maximum.
func (_ VariadicExpression) Between(low uint, high uint, max ...uint) VariadicExpression {
	var m *uint
	if len(max) > 0 {
		m = &max[0]
	}

	return VariadicExpression{
		low:  &low,
		high: &high,
		max:  m,
	}
}

// All - yourSlice[:] - Reads the entirety of your slice.
func (_ VariadicExpression) All() VariadicExpression {
	return VariadicExpression{}
}

// LogicGate - Reads every bit and calls the provided logic gate function to manipulate the binary output.
func (_ VariadicExpression) LogicGate(logic func(int, Bit) Bit) VariadicExpression {
	return VariadicExpression{
		logic: &logic,
	}
}
