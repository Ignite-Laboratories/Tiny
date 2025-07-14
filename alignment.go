package tiny

// LeftPadWithZeros will pad the left side of the smaller operands to the length of the largest with zeros.
func LeftPadWithZeros[T binary](operands ...T) []T {
	return pad[T](GetWidestOperand[T](operands...), false, Zero, operands...)
}

// LeftPadWithOnes will pad the left side of the smaller operands to the length of the largest with ones.
func LeftPadWithOnes[T binary](operands ...T) []T {
	return pad[T](GetWidestOperand[T](operands...), false, One, operands...)
}

// RightPadWithZeros will pad the right side of the smaller operands to the length of the largest with zeros.
func RightPadWithZeros[T binary](operands ...T) []T {
	return pad[T](GetWidestOperand[T](operands...), true, Zero, operands...)
}

// RightPadWithOnes will pad the right side of the smaller operands to the length of the largest with ones.
func RightPadWithOnes[T binary](operands ...T) []T {
	return pad[T](GetWidestOperand[T](operands...), true, One, operands...)
}

// MiddlePadWithZeros will equally pad both sides of the smaller operands to the length of the largest with zeros, biased towards the left.
func MiddlePadWithZeros[T binary](operands ...T) []T {
	return middlePad(Zero, operands...)
}

// MiddlePadWithOnes will equally pad both sides of the smaller operands to the length of the largest with ones, biased towards the left.
func MiddlePadWithOnes[T binary](operands ...T) []T {
	return middlePad(One, operands...)
}

// TileLeftToRight will repeat the provided pattern across a new operand as long as the longest operand, starting from the most significant side and working towards the least.
func TileLeftToRight[T binary](pattern T, operands ...T) []T {

}

func middlePad[T binary](placeholder Bit, operands ...T) []T {
	longest := GetWidestOperand[T](operands...)
	out := make([]T, len(operands))

	for i, o := range operands {
		length := GetBitLength(o)
		toPad := longest - length
		left := toPad / 2
		right := toPad - left

		out[i] = pad(left, false, placeholder, o)[0]
		out[i] = pad(right, true, placeholder, out[i])[0]
	}

	return out
}

func pad[T binary](length int, right bool, bit Bit, operands ...T) []T {
	if len(operands) == 0 {
		return make([]T, 0)
	}

	out := make([]T, len(operands))
	for i, o := range operands {
		toPadLen := length - GetBitLength(o)
		toPad := make([]Bit, toPadLen)

		if toPadLen == 0 {
			out[i] = o
			continue
		}

		if bit == 1 {
			for ii, _ := range toPad {
				toPad[ii] = 1
			}
		}

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

//
//// TileLeftToRight will tile the smallest operand across the length of the largest starting from the left towards the right.
//TileLeftToRight
//
//// TileRightToLeft will tile the smallest operand across the length of the largest starting from the right towards the left.
//TileRightToLeft
