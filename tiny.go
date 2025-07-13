package tiny

import (
	"reflect"
	"unsafe"
)

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
