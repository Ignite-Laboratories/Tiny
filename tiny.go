package tiny

import (
	"reflect"
	"unsafe"
)

// SanityCheck ensures the provided bits are all either one or zero, as Bit is a byte underneath.  In the land of binary,
// that can break all logic without you ever knowing - thus, this intentionally hard panics with ErrorNotABit.
func SanityCheck(bits ...Bit) {
	for _, b := range bits {
		if b != 0 && b != 1 {
			panic(ErrorNotABit)
		}
	}
}

// GetEndianness returns the Endianness of the currently executing hardware.
func GetEndianness() Endianness {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf[0] {
	case 0xCD:
		return LittleEndian
	case 0xAB:
		return BigEndian
	default:
		panic("could not determine native endianness")
	}
}

// GetBits extracts bits from any sized object at runtime.  This automatically will determine
// the host architecture's endianness, but you may override that if desired.
func GetBits[T any](value T, endian ...Endianness) []Bit {
	// Default to little endian if not specified
	targetEndian := GetEndianness()
	if len(endian) > 0 {
		targetEndian = endian[0]
	}

	ptr := unsafe.Pointer(&value)
	size := reflect.TypeOf(value).Size()
	bytes := unsafe.Slice((*byte)(ptr), size)
	result := make([]Bit, size*8)

	for byteIdx := 0; byteIdx < len(bytes); byteIdx++ {
		// For big endian, read bytes from the end
		var currentByte byte
		if targetEndian == BigEndian {
			currentByte = bytes[len(bytes)-1-byteIdx]
		} else {
			currentByte = bytes[byteIdx]
		}

		// Extract bits from each byte, MSB to LSB
		for bitIdx := 0; bitIdx < 8; bitIdx++ {
			resultIdx := (byteIdx * 8) + bitIdx
			bit := (currentByte >> (7 - bitIdx)) & 1
			result[resultIdx] = Bit(bit)
		}
	}

	return result
}

// BitsToType converts a slice of bits into the specified type T, respecting endianness
func BitsToType[T any](bits []Bit, endian ...Endianness) T {
	var zero T
	size := reflect.TypeOf(zero).Size()
	if len(bits) > int(size)*8 {
		panic("bit slice too large for target type")
	}

	targetEndian := GetEndianness()
	if len(endian) > 0 {
		targetEndian = endian[0]
	}

	result := zero
	resultPtr := unsafe.Pointer(&result)
	resultBytes := unsafe.Slice(resultPtr, size)

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
