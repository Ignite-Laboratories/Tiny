package tiny

import (
	"fmt"
	"math"
)

// UnaryExpression represents the standard slice index accessor pattern, and expressions can be generated from the global Read variable.
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
type UnaryExpression struct {
	pos     *uint
	low     *uint
	high    *uint
	max     *uint
	first   *bool
	last    *bool
	reverse *bool
	logic   *func(int, Bit) Bit
}

// UnaryEmit expresses the bits of binary information according to logical rules.
func UnaryEmit[T binary](expr UnaryExpression, data ...T) []Bit {
	if expr.max != nil && (expr.low == nil || expr.high == nil) {
		panic("invalid slice expression: max requires both low and high to be set")
	}

	if len(data) == 0 {
		return make([]Bit, 0)
	}

	if expr.reverse != nil && *expr.reverse {
		out := make([]T, len(data))
		count := len(data) - 1

		for i := len(data) - 1; i >= 0; i-- {
			switch concrete := any(data[i]).(type) {
			case Measurement:
				out[count-i] = any(concrete.Reverse()).(T)
			case Phrase:
				out[count-i] = any(concrete.Reverse()).(T)
			case []byte:
				for ii := len(concrete) - 1; ii >= 0; ii-- {
					out[count-i] = any(ReverseByte(concrete[ii])).(T)
				}
			case []Bit:
				for ii := len(concrete) - 1; ii >= 0; ii-- {
					out[count-i] = any(concrete[ii]).(T)
				}
			case byte:
				out[count-i] = any(ReverseByte(concrete)).(T)
			}
		}

		data = out
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
			byteBits := UnaryEmit(expr, concrete.Bytes...)
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
			slice = append(slice, UnaryEmit(expr, concrete.Bits...)...)
		case Phrase:
			slice = append(slice, UnaryEmit(expr, concrete.Data...)...)
		case []byte:
			slice = append(slice, UnaryEmit(expr, concrete...)...)
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
				return []Bit{concrete[len(concrete)-1]}
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
			return UnaryEmit(expr, NewMeasurementFromBytes(concrete).GetAllBits()...)
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
func (_ UnaryExpression) First(reverse ...bool) UnaryExpression {
	isReverse := len(reverse) > 0 && reverse[0]
	return UnaryExpression{
		first:   &True,
		reverse: &isReverse,
	}
}

// Last - yourSlice[𝑛 - 1] - Reads the last index position of your slice.
func (_ UnaryExpression) Last(reverse ...bool) UnaryExpression {
	isReverse := len(reverse) > 0 && reverse[0]
	return UnaryExpression{
		last:    &True,
		reverse: &isReverse,
	}
}

// Position - yourSlice[pos] - Reads the provided index position of your slice.
func (_ UnaryExpression) Position(pos uint, reverse ...bool) UnaryExpression {
	isReverse := len(reverse) > 0 && reverse[0]
	return UnaryExpression{
		pos:     &pos,
		reverse: &isReverse,
	}
}

// From - yourSlice[low:] - Reads from the provided index to the end of your slice.
func (_ UnaryExpression) From(low uint, reverse ...bool) UnaryExpression {
	isReverse := len(reverse) > 0 && reverse[0]
	return UnaryExpression{
		low:     &low,
		reverse: &isReverse,
	}
}

// To - yourSlice[:high] - Reads to the provided index from the start of your slice.
func (_ UnaryExpression) To(high uint, reverse ...bool) UnaryExpression {
	isReverse := len(reverse) > 0 && reverse[0]
	return UnaryExpression{
		high:    &high,
		reverse: &isReverse,
	}
}

// Between - yourSlice[low:high] - Reads between the provided indexes of your slice in most→to→least significant order.
//
// Between - yourSlice[low:high:mid] - Reads between the provided indexes of your slice up to the provided maximum in most→to→least significant order.
func (_ UnaryExpression) Between(low uint, high uint, max ...uint) UnaryExpression {
	var m *uint
	if len(max) > 0 {
		m = &max[0]
	}

	return UnaryExpression{
		low:  &low,
		high: &high,
		max:  m,
	}
}

// BetweenReverse - yourSlice[low:high] - Reads between the provided indexes of your slice in least←to←most significant order.
//
// BetweenReverse - yourSlice[low:high:mid] - Reads between the provided indexes of your slice up to the provided maximum in least←to←most significant order.
func (_ UnaryExpression) BetweenReverse(low uint, high uint, max ...uint) UnaryExpression {
	var m *uint
	if len(max) > 0 {
		m = &max[0]
	}

	return UnaryExpression{
		low:     &low,
		high:    &high,
		max:     m,
		reverse: &True,
	}
}

// All - yourSlice[:] - Reads the entirety of your slice.
func (_ UnaryExpression) All(reverse ...bool) UnaryExpression {
	isReverse := len(reverse) > 0 && reverse[0]
	return UnaryExpression{
		reverse: &isReverse,
	}
}

// LogicGate - Reads every bit and calls the provided logic gate function to manipulate the binary output.
func (_ UnaryExpression) LogicGate(logic func(int, Bit) Bit, reverse ...bool) UnaryExpression {
	isReverse := len(reverse) > 0 && reverse[0]
	return UnaryExpression{
		logic:   &logic,
		reverse: &isReverse,
	}
}
