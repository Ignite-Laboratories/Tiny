package tiny

import (
	"fmt"
	"reflect"
	"unsafe"
)

// GetBitWidth returns the bit width of the provided binary operand.
func GetBitWidth[T binary](operands ...T) uint {
	width := uint(0)
	for _, raw := range operands {
		switch operand := any(raw).(type) {
		case Phrase:
			width += operand.BitWidth()
		case Measurement:
			width += operand.BitWidth()
		case []byte:
			width += uint(len(operand) * 8)
		case []Bit:
			width += uint(len(operand))
		case byte:
			width += 8
		case Bit:
			width += 1
		default:
			panic("invalid binary type: " + reflect.TypeOf(operand).String())
		}
	}
	return width
}

// BleedEnd returns the ending bits of the operands and the operands missing those bits.
//
// All bleed operations are always returned in their original most→to→least significant order.
func BleedEnd[T binary](width uint, operands ...T) ([][]Bit, []T) {
	bits := make([][]Bit, 0, len(operands))

	for x := 0; x < int(width); x++ {
		for i, raw := range operands {
			if GetBitWidth(raw) == 0 {
				continue
			}

			switch operand := any(raw).(type) {
			case Phrase:
			case Measurement:
				var bit Bit
				bit, operand = operand.BleedLastBit()
				bits[i] = append([]Bit{bit}, bits[i]...)
				operands[i] = any(operand).(T)
			case []byte:
				panic("cannot bleed bits from static width elements")
			case []Bit:
				bits[i] = append([]Bit{operand[len(operand)-1]}, bits[i]...)
				operands[i] = any(operand[:len(operand)-1]).(T)
			case byte:
				panic("cannot bleed bits from static width elements")
			case Bit:
				panic("cannot bleed bits from static width elements")
			default:
				panic("invalid binary type: " + reflect.TypeOf(operand).String())
			}
		}
	}
	return bits, operands
}

// BleedStart returns the first bit of the operands and the operands missing that bit.
//
// All bleed operations are always returned in their original most→to→least significant order.
func BleedStart[T binary](width uint, operands ...T) ([][]Bit, []T) {
	bits := make([][]Bit, 0, len(operands))

	for x := 0; x < int(width); x++ {
		for i, raw := range operands {
			if GetBitWidth(raw) == 0 {
				continue
			}

			switch operand := any(raw).(type) {
			case Phrase:
			case Measurement:
				var bit Bit
				bit, operand = operand.BleedFirstBit()
				bits[i] = append([]Bit{bit}, bits[i]...)
				operands[i] = any(operand).(T)
			case []byte:
				panic("cannot bleed bits from static width elements")
			case []Bit:
				bits[i] = append([]Bit{operand[0]}, bits[i]...)
				operands[i] = any(operand[1:]).(T)
			case byte:
				panic("cannot bleed bits from static width elements")
			case Bit:
				panic("cannot bleed bits from static width elements")
			default:
				panic("invalid binary type: " + reflect.TypeOf(operand).String())
			}
		}
	}
	return bits, operands
}

// GetWidestOperand returns the widest bit width of the provided operands.
func GetWidestOperand[T binary](operands ...T) uint {
	var widest uint
	for _, o := range operands {
		width := GetBitWidth(o)
		if width > widest {
			widest = width
		}
	}
	return widest
}

// AlignOperand applies the provided Alignment scheme against the operand in order to place the measured binary information relative to the provided bit width.
func AlignOperand[T binary](operand T, width uint, scheme Alignment) T {
	switch scheme {
	case PadLeftSideWithZeros:
		return any(padLeftSideWithZeros(width, operand)[0]).(T)
	case PadLeftSideWithOnes:
		return any(padLeftSideWithOnes(width, operand)[0]).(T)
	case PadRightSideWithZeros:
		return any(padRightSideWithZeros(width, operand)[0]).(T)
	case PadRightSideWithOnes:
		return any(padRightSideWithOnes(width, operand)[0]).(T)
	case PadToMiddleWithZeros:
		return any(padToMiddleWithZeros(width, operand)[0]).(T)
	case PadToMiddleWithOnes:
		return any(padToMiddleWithOnes(width, operand)[0]).(T)
	default:
		panic("invalid alignment scheme")
	}
}

func ReverseOperands[T binary](operands ...T) []T {
	// Put your thing down, flip it, and reverse it
	reversed := make([]T, len(operands))
	limit := len(operands) - 1

	for i, raw := range operands {
		switch operand := any(raw).(type) {
		case Measurement:
			reversed[limit-i] = any(operand.Reverse()).(T)
		case Phrase:
			reversed[limit-i] = any(operand.Reverse()).(T)
		case []byte:
			r := make([]byte, len(operand))
			for ii := len(operand) - 1; ii >= 0; ii-- {
				r[limit-ii] = ReverseByte(operand[ii])
			}
			reversed[limit-i] = any(r).(T)
		case []Bit:
			r := make([]Bit, len(operand))
			for ii := len(operand) - 1; ii >= 0; ii-- {
				r[limit-ii] = operand[ii]
			}
			reversed[limit-i] = any(r).(T)
		case byte:
			reversed[limit-i] = any(ReverseByte(operand)).(T)
		default:
			panic(fmt.Errorf("invalid binary type: %T", operand))
		}
	}

	return reversed
}

// ReverseByte reverses all the bits of a byte.
func ReverseByte(b byte) byte {
	b = (b&0xF0)>>4 | (b&0x0F)<<4
	b = (b&0xCC)>>2 | (b&0x33)<<2
	return (b&0xAA)>>1 | (b&0x55)<<1
}

// SanityCheck ensures the provided bits are all either Zero, One, or Nil - as Bit is a byte underneath.  In the land of
// binary, that can break all logic without you ever knowing - thus, this intentionally hard panics with ErrorNotABit.
func SanityCheck(bits ...Bit) {
	for _, b := range bits {
		if b != Zero && b != One && b != Nil {
			panic(ErrorNotABit)
		}
	}
}

// Measure takes a named Measurement of the bits in any sized object at runtime and returns them as a Logical Phrase.  This
// automatically will determine the host architecture's endianness and reverse the bytes if they are found to be BigEndian.
// This ensures all tiny operations happen in LittleEndian byte order, regardless of the underlying hardware.
func Measure[T any](name string, value T, endian ...Endianness) Phrase {
	targetEndian := GetEndianness()
	if len(endian) > 0 {
		targetEndian = endian[0]
	}

	valueType := reflect.TypeOf(value)
	var size uintptr
	var dataPtr unsafe.Pointer

	// Handle slices differently from other types
	if valueType.Kind() == reflect.Slice {
		// Get slice width and element size
		sliceVal := reflect.ValueOf(value)
		elemSize := valueType.Elem().Size()
		width := sliceVal.Len()
		size = uintptr(width) * elemSize

		// Get pointer to first element
		if width > 0 {
			dataPtr = sliceVal.UnsafePointer()
		}
	} else {
		dataPtr = unsafe.Pointer(&value)
		size = valueType.Size()
	}

	if size == 0 {
		return NewPhrase(name, Logical)
	}

	bytes := unsafe.Slice((*byte)(dataPtr), size)

	if targetEndian == BigEndian {
		for i := 0; i < len(bytes); i++ {
			bytes[i] = bytes[len(bytes)-1-i]
		}
	}

	phrase := NewPhrase(name, Logical)
	phrase.Data = []Measurement{NewMeasurementOfBytes(bytes...)}
	return phrase
}

// ToType converts a Phrase of binary information into the specified type T, respecting the architecture's Endianness.
func ToType[T any](p Phrase, endian ...Endianness) T {
	// TODO: Entirely re-write this to utilize Emit and read operations, that way we aren't actually expanding ALL the bits into a full byte in the process.
	bits := p.GetAllBits()
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
