package tiny

// Align represents a scheme of how to align operands relative to each other.
//
// PadLeftSideWithZeros will pad the left side of the smaller operands to the length of the largest with zeros.
//
// PadLeftSideWithOnes will pad the left side of the smaller operands to the length of the largest with ones.
//
// PadRightSideWithZeros will pad the right side of the smaller operands to the length of the largest with zeros.
//
// PadRightSideWithOnes will pad the right side of the smaller operands to the length of the largest with ones.
//
// PadToMiddleWithZeros will equally pad both sides of the smaller operands to the length of the largest with zeros, biased towards the left.
//
// PadToMiddleWithOnes will equally pad both sides of the smaller operands to the length of the largest with ones, biased towards the left.
type Align int

const (
	// PadLeftSideWithZeros will pad the left side of the smaller operands to the length of the largest with zeros.
	PadLeftSideWithZeros Align = iota

	// PadLeftSideWithOnes will pad the left side of the smaller operands to the length of the largest with ones.
	PadLeftSideWithOnes

	// PadRightSideWithZeros will pad the right side of the smaller operands to the length of the largest with zeros.
	PadRightSideWithZeros

	// PadRightSideWithOnes will pad the right side of the smaller operands to the length of the largest with ones.
	PadRightSideWithOnes

	// PadToMiddleWithZeros will equally pad both sides of the smaller operands to the length of the largest with zeros, biased towards the left.
	PadToMiddleWithZeros

	// PadToMiddleWithOnes will equally pad both sides of the smaller operands to the length of the largest with ones, biased towards the left.
	PadToMiddleWithOnes
)

func padLeftSideWithZeros[T binary](length int, operands ...T) []T {
	return pad[T](length, false, Zero, operands...)
}

func padLeftSideWithOnes[T binary](length int, operands ...T) []T {
	return pad[T](length, false, One, operands...)
}

func padRightSideWithZeros[T binary](length int, operands ...T) []T {
	return pad[T](length, true, Zero, operands...)
}

func padRightSideWithOnes[T binary](length int, operands ...T) []T {
	return pad[T](length, true, One, operands...)
}

func padToMiddleWithZeros[T binary](length int, operands ...T) []T {
	return middlePad(length, Zero, operands...)
}

func padToMiddleWithOnes[T binary](length int, operands ...T) []T {
	return middlePad(length, One, operands...)
}

func middlePad[T binary](length int, digit Bit, operands ...T) []T {
	out := make([]T, len(operands))

	for i, o := range operands {
		l := GetBitLength(o)
		toPad := length - l
		left := toPad / 2
		right := toPad - left

		out[i] = pad(left, false, digit, o)[0]
		out[i] = pad(right, true, digit, out[i])[0]
	}

	return out
}

func pad[T binary](length int, right bool, digit Bit, operands ...T) []T {
	if len(operands) == 0 {
		return make([]T, 0)
	}

	out := make([]T, len(operands))
	for i, o := range operands {
		toPadLen := length - GetBitLength(o)
		if toPadLen == 0 {
			out[i] = o
			continue
		}

		toPad := NewMeasurementOfDigit(toPadLen, digit).GetAllBits()

		switch concrete := any(o).(type) {
		case Phrase:
			if right {
				out[i] = any(concrete.Append(toPad...)).(T)
			} else {
				out[i] = any(concrete.Prepend(toPad...)).(T)
			}
		case Measurement:
			if right {
				out[i] = any(concrete.Append(toPad...)).(T)
			} else {
				out[i] = any(concrete.Prepend(toPad...)).(T)
			}
		case []byte:
			panic("cannot pad static width elements")
		case []Bit:
			if right {
				out[i] = any(append(concrete, toPad...)).(T)
			} else {
				out[i] = any(append(toPad, concrete...)).(T)
			}
		case byte:
			panic("cannot pad a static width element")
		case Bit:
			panic("cannot pad a static width element")
		default:
			panic("unknown operand type")
		}
	}
	return out
}
