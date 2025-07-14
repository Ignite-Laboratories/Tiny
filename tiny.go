package tiny

import (
	"reflect"
	"unsafe"
)

// Bits provides access to slice index accessor expressions.
//
// UnaryExpression.Position - yourSlice[pos] - Reads the provided index position of your slice.
//
// UnaryExpression.All - yourSlice[:] - Reads the entirety of your slice.
//
// UnaryExpression.From - yourSlice[low:] - Reads from the provided index to the end of your slice.
//
// UnaryExpression.To - yourSlice[:high] - Reads to the provided index from the start of your slice.
//
// UnaryExpression.Between - yourSlice[low:high] - Reads between the provided indexes of your slice.
//
// UnaryExpression.Between - yourSlice[low:high:mid] - Reads between the provided indexes of your slice up to the provided maximum.
//
// UnaryExpression.LogicGate - Performs a logical operation for every bit of your slice.
var Bits UnaryExpression

// GetBitLength returns the bit length of the provided binary operand.
func GetBitLength[T binary](operand T) int {
	switch concrete := any(operand).(type) {
	case Phrase:
		return concrete.BitLength()
	case Measurement:
		return concrete.BitLength()
	case []byte:
		return len(concrete) * 8
	case []Bit:
		return len(concrete)
	case byte:
		return 8
	case Bit:
		return 1
	default:
		panic("invalid binary type: " + reflect.TypeOf(concrete).String())
	}
}

// GetWidestOperand returns the bit length widest of the provided operands.
func GetWidestOperand[T binary](operands ...T) int {
	var largest int
	for _, o := range operands {
		length := GetBitLength(o)
		if length > largest {
			largest = length
		}
	}
	return largest
}

// SanityCheck ensures the provided bits are all either one or zero, as Bit is a byte underneath.  In the land of binary,
// that can break all logic without you ever knowing - thus, this intentionally hard panics with ErrorNotABit.
func SanityCheck(bits ...Bit) {
	for _, b := range bits {
		if b > 1 {
			panic(ErrorNotABit)
		}
	}
}

// Measure extracts bits from any sized object at runtime.  This automatically will determine
// the host architecture's endianness, but you may override that if desired.
func Measure[T any](value T, endian ...Endianness) Phrase {
	targetEndian := GetEndianness()
	if len(endian) > 0 {
		targetEndian = endian[0]
	}

	valueType := reflect.TypeOf(value)
	var size uintptr
	var dataPtr unsafe.Pointer

	// Handle slices differently from other types
	if valueType.Kind() == reflect.Slice {
		// Get slice length and element size
		sliceVal := reflect.ValueOf(value)
		elemSize := valueType.Elem().Size()
		length := sliceVal.Len()
		size = uintptr(length) * elemSize

		// Get pointer to first element
		if length > 0 {
			dataPtr = unsafe.Pointer(sliceVal.UnsafePointer())
		}
	} else {
		dataPtr = unsafe.Pointer(&value)
		size = valueType.Size()
	}

	if size == 0 {
		return NewLogicalPhrase()
	}

	bytes := unsafe.Slice((*byte)(dataPtr), size)
	result := make([]Bit, size*8)

	for byteIdx := 0; byteIdx < len(bytes); byteIdx++ {
		var currentByte byte
		if targetEndian == BigEndian {
			currentByte = bytes[len(bytes)-1-byteIdx]
		} else {
			currentByte = bytes[byteIdx]
		}

		for bitIdx := 0; bitIdx < 8; bitIdx++ {
			resultIdx := (byteIdx * 8) + bitIdx
			bit := (currentByte >> (7 - bitIdx)) & 1
			result[resultIdx] = Bit(bit)
		}
	}

	phrase := NewLogicalPhrase()
	phrase.Data = make([]Measurement, len(result))
	for i, b := range result {
		phrase.Data[i] = NewMeasurement(b)
	}

	return phrase
}

// ToType converts a slice of bits into the specified type T, respecting endianness
func ToType[T any](p Phrase, endian ...Endianness) T {
	bits := p.GetBits()
	var zero T
	typeOf := reflect.TypeOf(zero)

	targetEndian := GetEndianness()
	if len(endian) > 0 {
		targetEndian = endian[0]
	}

	// Handle slices
	if typeOf.Kind() == reflect.Slice {
		elemType := typeOf.Elem()
		elemSize := elemType.Size()

		numElements := len(bits) / (8 * int(elemSize))
		if numElements == 0 {
			return zero
		}

		sliceVal := reflect.MakeSlice(typeOf, numElements, numElements)
		slicePtr := unsafe.Pointer(sliceVal.UnsafePointer())
		resultBytes := unsafe.Slice((*byte)(slicePtr), numElements*int(elemSize))

		for byteIdx := 0; byteIdx < len(bits)/8; byteIdx++ {
			var currentByte byte
			for bitIdx := 0; bitIdx < 8; bitIdx++ {
				if bits[byteIdx*8+bitIdx] == 1 {
					currentByte |= 1 << (7 - bitIdx)
				}
			}

			if targetEndian == BigEndian {
				elementIdx := byteIdx / int(elemSize)
				byteOffset := byteIdx % int(elemSize)
				resultBytes[elementIdx*int(elemSize)+(int(elemSize)-1-byteOffset)] = currentByte
			} else {
				resultBytes[byteIdx] = currentByte
			}
		}

		return sliceVal.Interface().(T)
	}

	// Handle non-slices
	size := typeOf.Size()
	if len(bits) > int(size)*8 {
		panic("bit slice too large for target type")
	}

	result := zero
	resultPtr := unsafe.Pointer(&result)
	resultBytes := unsafe.Slice((*byte)(resultPtr), size)

	for byteIdx := 0; byteIdx < len(bits)/8; byteIdx++ {
		var currentByte byte
		for bitIdx := 0; bitIdx < 8; bitIdx++ {
			if bits[byteIdx*8+bitIdx] == 1 {
				currentByte |= 1 << (7 - bitIdx)
			}
		}

		if targetEndian == BigEndian {
			resultBytes[len(resultBytes)-1-byteIdx] = currentByte
		} else {
			resultBytes[byteIdx] = currentByte
		}
	}

	return result

}
