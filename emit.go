package tiny

import (
	"fmt"
	"math"
)

// Emit expresses the underlying bits of the provided binary data according to logical rules.
func Emit[T binary](expr Expression, data ...T) []Bit {
	if expr._max != nil && (expr._low == nil || expr._high == nil) {
		panic("invalid slice expression: max requires both low and high to be set")
	}

	if len(data) == 0 {
		return make([]Bit, 0)
	}

	var reverse bool
	if expr._reverse != nil && *expr._reverse {
		reverse = true
	}

	expected := -1
	if expr._high != nil {
		expected = int(*expr._high)

		if expr._low != nil {
			expected = int(*expr._high) - int(*expr._low)
		}

		if expr._max != nil {
			maximum := int(*expr._max)
			if maximum < expected {
				expected = maximum
			}
		}
	}

	yield := make([]Bit, 0, int(math.Min(1<<10, float64(expected))))

	//if expr._matrix != nil && *expr._matrix {
	//	/**
	//	Matrix Logic
	//	*/
	//
	//	if expr._matrixLogic == nil {
	//		panic("matrix expressions require a logic function")
	//	}
	//
	//	calculate := *expr._matrixLogic
	//
	//	if expr._alignment == nil {
	//		align := PadLeftSideWithZeros
	//		expr._alignment = &align
	//	}
	//
	//	longest := GetLongestOperand(data...)
	//
	//	if longest <= 0 {
	//		return yield, 0
	//	}
	//
	//	subExpr := expr
	//	subExpr._matrix = &False
	//
	//	// The underlying table is ordered [Col][Row]Bit
	//	table := make([][]Bit, longest)
	//	for i, raw := range data {
	//		data[i] = AlignOperand(raw, longest, *expr._alignment)
	//		table[i], _ = Emit[T](subExpr, raw)
	//	}
	//
	//	// TODO: We can't walk using longest because longest will grow as we carry - instead we need to just walk until we are out of bits to walk and pass the walk count to the matrix func
	//
	//	for i := 0; i < longest; i++ {
	//		colId := i
	//		if reverse {
	//			colId = longest - i - 1
	//		}
	//
	//		column := make([]Bit, len(table))
	//		for rowId, row := range table {
	//			column[rowId] = row[colId]
	//		}
	//		calculated, overflow := calculate(colId, column...)
	//
	//		// TODO: Insert the overflow binary value BELOW the upcoming columns in the direction of calculation
	//
	//		if reverse {
	//			yield = append(yield, calculated)
	//		} else {
	//			yield = append([]Bit{b}, yield...)
	//		}
	//	}
	//
	//	linear := make([]Bit, 0, len(matrix)*longest)
	//	for _, element := range matrix {
	//		linear = append(linear, element...)
	//	}
	//
	//	yield = linear
	//	count = uint(longest) // TODO: Align all the operands and set this to the number of returned operands
	//} else {
	/**
	Linear Logic
	*/

	if reverse {
		out := make([]T, len(data))
		c := len(data) - 1

		for i := len(data) - 1; i >= 0; i-- {
			switch operand := any(data[i]).(type) {
			case Measurement:
				out[c-i] = any(operand.Reverse()).(T)
			case Phrase:
				out[c-i] = any(operand.Reverse()).(T)
			case []byte:
				for ii := len(operand) - 1; ii >= 0; ii-- {
					out[c-i] = any(ReverseByte(operand[ii])).(T)
				}
			case []Bit:
				for ii := len(operand) - 1; ii >= 0; ii-- {
					out[c-i] = any(operand[ii]).(T)
				}
			case byte:
				out[c-i] = any(ReverseByte(operand)).(T)
			}
		}

		data = out
	}

	for _, raw := range data {
		switch operand := any(raw).(type) {
		case Measurement:
			byteBits := Emit(expr, operand.Bytes...)
			yield = append(yield, byteBits...)

			if expr._high != nil {
				high := *expr._high - uint(len(byteBits))
				expr._high = &high
			}
			if expr._low != nil {
				low := *expr._low - uint(len(byteBits))
				if low < 0 {
					low = 0
				}
				expr._low = &low
			}
			yield = append(yield, Emit(expr, operand.Bits...)...)
		case Phrase:
			yield = append(yield, Emit(expr, operand.Data...)...)
		case []byte:
			yield = append(yield, Emit(expr, operand...)...)
		case []Bit:
			capacity := len(operand)

			switch {
			case expr._unary != nil:
				out := make([]Bit, 0)
				for i, b := range operand {
					out = append(out, (*expr._unary)(i, b))
				}
				return out
			case expr._first != nil:
				return []Bit{operand[0]}
			case expr._last != nil:
				return []Bit{operand[len(operand)-1]}
			case expr._pos != nil:
				return []Bit{operand[*expr._pos]}
			case expr._low == nil && expr._high == nil:
				return operand[:]
			case expr._low == nil && expr._high != nil:
				high := int(math.Min(float64(capacity), float64(*expr._high)))
				return operand[:high]
			case expr._low != nil && expr._high == nil:
				low := int(math.Min(float64(capacity), float64(*expr._low)))
				return operand[low:]
			case expr._low != nil && expr._high != nil:
				high := int(math.Min(float64(capacity), float64(*expr._high)))
				low := int(math.Min(float64(capacity), float64(*expr._low)))
				if expr._max == nil {
					return operand[low:high]
				}
				return operand[low:high:*expr._max]
			default:
				panic("invalid slice expression")
			}
		case byte:
			return Emit(expr, NewMeasurementFromBytes(operand).GetAllBits()...)
		default:
			panic(fmt.Errorf("invalid bit expression type: %T", operand))
		}

		if len(yield) >= expected {
			// Bailout when the requested range has accrued
			return yield[:expected]
		}
	}
	//}

	return yield
}
