package tiny

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"reflect"
	"tiny/endian"
	"tiny/pad"
	"tiny/travel"
	"unsafe"
)

// GetRandomName returns a randomly generated name which conforms to the NameFilter rules.
func GetRandomName() string {
	return core.RandomNameFiltered(NameFilter).Name
}

// GetBitWidth returns the bit width of the provided binary operand.
func GetBitWidth[T Operable](operands ...T) uint {
	width := uint(0)
	for _, raw := range operands {
		switch operand := any(raw).(type) {
		case Phrase:
			width += operand.BitWidth()
		case Index:
			width += operand.BitWidth()
		case Real:
			width += operand.BitWidth()
		case Complex:
			width += operand.Real.BitWidth()
			width += operand.Imaginary.BitWidth()
		case Natural:
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
func BleedEnd[T Operable](width uint, operands ...T) ([][]Bit, []T) {
	bits := make([][]Bit, 0, len(operands))

	for x := 0; x < int(width); x++ {
		for i, raw := range operands {
			if GetBitWidth(raw) == 0 {
				continue
			}

			switch operand := any(raw).(type) {
			case Phrase, Complex, Index, Real, Natural:
				// TODO: Implement this
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
func BleedStart[T Operable](width uint, operands ...T) ([][]Bit, []T) {
	bits := make([][]Bit, 0, len(operands))

	for x := 0; x < int(width); x++ {
		for i, raw := range operands {
			if GetBitWidth(raw) == 0 {
				continue
			}

			switch operand := any(raw).(type) {
			case Phrase, Complex, Index, Real, Natural:
				// TODO: Implement this
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
func GetWidestOperand[T Operable](operands ...T) uint {
	var widest uint
	for _, o := range operands {
		width := GetBitWidth(o)
		if width > widest {
			widest = width
		}
	}
	return widest
}

// AlignOperands applies the provided padding scheme against the operands to align the measured binary information relative to the provided bit width.
//
// You must provide at least one digit to pad the data with, but you may provide a pattern of digits.  The pattern is emitted across the operand starting
// from the West side and working towards the East.  If working latitudinally, the pattern bits are applied longitudinally across each operand in the same way.
//
// NOTE: If you wish for
func AlignOperands[T Operable](operands []T, width uint, scheme pad.Scheme, traveling travel.Travel, digits ...Bit) []T {
	// TODO: alignment
	return operands
}

// ReverseOperands reverses the provided input operands.  If they are an Operable type, the internal bits
// are reversed - otherwise, the operands themselves are returned in reverse order.
func ReverseOperands[T any](operands ...T) []T {
	// Put your thing down, flip it, and reverse it
	reversed := make([]T, len(operands))
	limit := len(operands) - 1

	for i, raw := range operands {
		switch operand := any(raw).(type) {
		case Real, Complex:
			panic(fmt.Errorf("cannot reverse real or complex numbers - please first convert to a phrase"))
		case Phrase:
			reversed[limit-i] = any(operand.Reverse()).(T)
		case Index:
			reversed[limit-i] = any(operand.Reverse()).(T)
		case Natural:
			reversed[limit-i] = any(operand.Reverse()).(T)
		case Measurement:
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
			reversed[limit-i] = any(operand).(T)
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

// SanityCheck ensures the provided bits are all either Zero or One - as Bit is a byte underneath.  In the land of
// binary, that can break all logic without you ever knowing - thus, this intentionally hard panics with ErrorNotABit.
//
// NOTE: This does not account for a 'nil' bit - for that, please use SanityCheckWithNil.
func SanityCheck(bits ...Bit) {
	for _, b := range bits {
		if b != Zero && b != One {
			panic(ErrorNotABit)
		}
	}
}

// SanityCheckWithNil ensures the provided bits are all either Zero, One, or Nil - as Bit is a byte underneath.  In the land of
// binary, that can break all logic without you ever knowing - thus, this intentionally hard panics with ErrorNotABit.
//
// NOTE: This accounts for a 'nil' bit - if you wish to work with "traditional" bits, please use SanityCheck.
func SanityCheckWithNil(bits ...Bit) {
	for _, b := range bits {
		if b != Zero && b != One && b != Nil {
			panic(ErrorNotABitWithNil)
		}
	}
}

// Measure takes a Measurement of any sized object at runtime.
func Measure[T any](value T) Measurement {
	m := NewMeasurementOfBytes(measure(value)...)
	m.Endianness = endian.GetArchitectureEndianness()
	return m
}

// MeasureMany takes measurements of many objects at runtime and returns the result as a single Phrase.
func MeasureMany[T any](values ...T) Phrase {
	p := NewPhrase()
	for _, v := range values {
		p = p.AppendMeasurement(Measure(v))
	}
	return p
}

func measure[T any](value T) []byte {
	var size uintptr
	switch any(value).(type) {
	case byte, int8, bool:
		size = 1
	case uint16, int16:
		size = 2
	case uint32, int32, float32:
		size = 4
	case uint64, int64, float64, uint, int:
		size = 8
	case complex64:
		size = 8
	case complex128:
		size = 16
	case string:
		return []byte(any(value).(string))
	default:
		// Handle other types including slices using reflection
		val := reflect.ValueOf(value)
		if val.Kind() == reflect.Slice {
			if val.Len() == 0 {
				return []byte{}
			}
			elemSize := val.Type().Elem().Size()
			totalSize := uintptr(val.Len()) * elemSize
			size = totalSize
		} else {
			size = reflect.TypeOf(value).Size()
		}
	}

	if size == 0 {
		return []byte{}
	}

	var dataPtr unsafe.Pointer
	if val := reflect.ValueOf(value); val.Kind() == reflect.Slice {
		dataPtr = unsafe.Pointer(val.UnsafePointer())
	} else {
		dataPtr = unsafe.Pointer(reflect.ValueOf(&value).Elem().UnsafeAddr())
	}

	bytes := make([]byte, size)
	copy(bytes, (*[1 << 30]byte)(dataPtr)[:size:size])
	return bytes
}

// ToType converts a Measurement of binary information into the specified type T.
func ToType[T any](m Measurement) T {
	bits := m.GetAllBits()
	var zero T
	typeOf := reflect.TypeOf(zero)

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

		byteI := (len(bits) / 8) - 1
		i := len(bits) - 1
		for i > 0 {
			var currentByte byte
			for ii := 0; ii < 8; ii++ {
				if bits[i] == 1 {
					currentByte |= 1 << ii
				}
				i--
			}

			resultBytes[byteI] = currentByte
			byteI--
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

	byteI := (len(bits) / 8) - 1
	i := len(bits) - 1
	for i > 0 {
		var currentByte byte
		for ii := 0; ii < 8; ii++ {
			if bits[i] == 1 {
				currentByte |= 1 << ii
			}
			i--
		}

		resultBytes[byteI] = currentByte
		byteI--
	}

	return result
}
