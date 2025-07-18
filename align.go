package tiny

import (
	"fmt"
	"tiny/direction"
	"tiny/travel"
)

func middlePadOperands[T Binary](width uint, direction direction.Direction, travel travel.Travel, digits []Bit, operands ...T) []T {
	// TODO: Implement north/south padding
	out := make([]T, len(operands))

	for i, o := range operands {
		l := GetBitWidth(o)
		toPad := width - l
		left := toPad / 2
		right := toPad - left

		if travel == travel.Outward {
			out[i] = padOperands(left, direction, travel.Westbound, digits, o)[0]
			out[i] = padOperands(right, direction, travel.Eastbound, digits, out[i])[0]
		} else if travel == travel.Inward {
			out[i] = padOperands(left, direction, travel.Eastbound, digits, o)[0]
			out[i] = padOperands(right, direction, travel.Westbound, digits, out[i])[0]
		} else {
			out[i] = padOperands(left, direction, travel, digits, o)[0]
			out[i] = padOperands(right, direction, travel, digits, out[i])[0]
		}
	}

	return out
}

func padOperands[T Binary](width uint, direction direction.Direction, travel travel.Travel, digits []Bit, operands ...T) []T {
	out := make([]T, len(operands))
	for i, raw := range operands {
		paddingWidth := width - GetBitWidth(raw)
		if paddingWidth == 0 {
			out[i] = raw
			continue
		}

		padding := NewMeasurementOfDigit(int(paddingWidth), digits).GetAllBits()

		switch operand := any(raw).(type) {
		case Phrase:
		case Measurement:
			switch direction {
			case direction.West:
				out[i] = any(operand.Prepend(padding...)).(T)
			case direction.East:
				out[i] = any(operand.Append(padding...)).(T)
			case direction.North:
			case direction.South:
				// TODO: Implement north/south padding
			default:
				panic(fmt.Sprintf("cannot pad to '%v' direction", direction))
			}
		case []Bit:
			switch direction {
			case direction.West:
				out[i] = any(append(padding, operand...)).(T)
			case direction.East:
				out[i] = any(append(operand, padding...)).(T)
			case direction.North:
			case direction.South:
				// TODO: Implement north/south padding
			default:
				panic(fmt.Sprintf("cannot pad to '%v' direction", direction))
			}
		case []byte:
		case byte:
		case Bit:
			panic("cannot pad static width elements")
		default:
			panic("unknown operand type")
		}
	}
	return out
}
