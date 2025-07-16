package tiny

// Alignment represents a scheme of how to align operands relative to each other - alignment operations are applied with
// no respect to the reversal of bit emission and instead operate as "left" or "right" spatially.
//
// NOTE: You can only align variable width binary types - such as slices of bytes or bits, or measurements or phrases.
// These will panic if you attempt to "pad a byte," for instance, as it's a static-width element.
//
// PadLeftSideWithZeros will pad the left side of the smaller operands to the width of the largest with zeros.
//
// PadLeftSideWithOnes will pad the left side of the smaller operands to the width of the largest with ones.
//
// PadRightSideWithZeros will pad the right side of the smaller operands to the width of the largest with zeros.
//
// PadRightSideWithOnes will pad the right side of the smaller operands to the width of the largest with ones.
//
// PadToMiddleWithZeros will equally pad both sides of the smaller operands to the width of the largest with zeros, biased towards the left.
//
// PadToMiddleWithOnes will equally pad both sides of the smaller operands to the width of the largest with ones, biased towards the left.
type Alignment int

const (
	// PadLeftSideWithZeros will pad the left side of the smaller operands to the width of the largest with zeros.
	PadLeftSideWithZeros Alignment = iota

	// PadLeftSideWithOnes will pad the left side of the smaller operands to the width of the largest with ones.
	PadLeftSideWithOnes

	// PadRightSideWithZeros will pad the right side of the smaller operands to the width of the largest with zeros.
	PadRightSideWithZeros

	// PadRightSideWithOnes will pad the right side of the smaller operands to the width of the largest with ones.
	PadRightSideWithOnes

	// PadToMiddleWithZeros will equally pad both sides of the smaller operands to the width of the largest with zeros, biased towards the left.
	PadToMiddleWithZeros

	// PadToMiddleWithOnes will equally pad both sides of the smaller operands to the width of the largest with ones, biased towards the left.
	PadToMiddleWithOnes
)

func padLeftSideWithZeros[T binary](width uint, operands ...T) []T {
	return pad[T](width, false, Zero, operands...)
}

func padLeftSideWithOnes[T binary](width uint, operands ...T) []T {
	return pad[T](width, false, One, operands...)
}

func padRightSideWithZeros[T binary](width uint, operands ...T) []T {
	return pad[T](width, true, Zero, operands...)
}

func padRightSideWithOnes[T binary](width uint, operands ...T) []T {
	return pad[T](width, true, One, operands...)
}

func padToMiddleWithZeros[T binary](width uint, operands ...T) []T {
	return middlePad(width, Zero, operands...)
}

func padToMiddleWithOnes[T binary](width uint, operands ...T) []T {
	return middlePad(width, One, operands...)
}

func middlePad[T binary](width uint, digit Bit, operands ...T) []T {
	out := make([]T, len(operands))

	for i, o := range operands {
		l := GetBitWidth(o)
		toPad := width - l
		left := toPad / 2
		right := toPad - left

		out[i] = pad(left, false, digit, o)[0]
		out[i] = pad(right, true, digit, out[i])[0]
	}

	return out
}

func pad[T binary](width uint, right bool, digit Bit, operands ...T) []T {
	out := make([]T, len(operands))
	for i, raw := range operands {
		paddingWidth := width - GetBitWidth(raw)
		if paddingWidth == 0 {
			out[i] = raw
			continue
		}

		padding := NewMeasurementOfDigit(int(paddingWidth), digit).GetAllBits()

		switch operand := any(raw).(type) {
		case Phrase:
			if right {
				out[i] = any(operand.Append(padding...)).(T)
			} else {
				out[i] = any(operand.Prepend(padding...)).(T)
			}
		case Measurement:
			if right {
				out[i] = any(operand.Append(padding...)).(T)
			} else {
				out[i] = any(operand.Prepend(padding...)).(T)
			}
		case []byte:
			panic("cannot pad static width elements")
		case []Bit:
			if right {
				out[i] = any(append(operand, padding...)).(T)
			} else {
				out[i] = any(append(padding, operand...)).(T)
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
