package testing

import (
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Measurement_GetAllBits(t *testing.T) {
	bits := []tiny.Bit{0, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0}
	expected := append([]tiny.Bit{1, 0, 1, 0, 1, 0, 1, 0}, bits...)
	measure := tiny.NewMeasurement([]byte{170}, bits...)
	result := measure.GetAllBits()
	CompareSlices(result, expected, t)
}

func Test_Measurement_BitLength(t *testing.T) {
	bits := []tiny.Bit{0, 1, 0}
	bytes := []byte{170, 85}

	bitMeasure := tiny.NewMeasurement([]byte{}, bits...)
	CompareValues(bitMeasure.BitLength(), 3, t)

	byteMeasure := tiny.NewMeasurement(bytes)
	CompareValues(byteMeasure.BitLength(), 16, t)

	bothMeasure := tiny.NewMeasurement(bytes, bits...)
	CompareValues(bothMeasure.BitLength(), 19, t)
}

func Test_Measurement_ByteLength(t *testing.T) {
	bytes := []byte{170, 85, 38, 75}

	for i := 0; i < len(bytes); i++ {
		toTest := bytes[:len(bytes)-i]
		m := tiny.NewMeasurement(toTest)
		CompareValues(m.ByteBitLength(), len(toTest)*8, t)
	}
}

func Test_Measurement_Value(t *testing.T) {
	m := tiny.NewMeasurement([]byte{170, 85}, 1, 0, 0, 1)
	CompareValues(m.Value(), 697689, t)
}

func Test_Measurement_Clear(t *testing.T) {
	m := tiny.NewMeasurement([]byte{170, 85}, 1, 0, 0, 1)
	m.Clear()
	CompareSlices(m.GetAllBits(), []tiny.Bit{}, t)
}

func Test_Measurement_Toggle(t *testing.T) {
	bytes := []byte{255, 0, 128, 127, 77}
	inverseBytes := []byte{0, 255, 127, 128, 178}
	bits := tiny.From.Bits(0, 1, 0, 1)
	inverseBits := tiny.From.Bits(1, 0, 1, 0)
	m := tiny.NewMeasurement(bytes, bits...)
	m.Toggle()

	for i, b := range m.Bytes {
		if b != inverseBytes[i] {
			t.Errorf("Expected %d, got %d", inverseBytes[i], b)
		}
	}
	for i, b := range m.Bits {
		if b != inverseBits[i] {
			t.Errorf("Expected %d, got %d", inverseBits[i], b)
		}
	}
}

func Test_Measurement_ForEachBit(t *testing.T) {
	random := tiny.Synthesize.RandomBits(22)[0]
	bits := random.GetAllBits()
	count := 0
	random.ForEachBit(func(i int, bit tiny.Bit) tiny.Bit {
		count++
		if bit != bits[i] {
			t.Errorf("Expected %d, got %d", bits[i], bit)
		}
		return bit
	})
	if count != len(bits) {
		t.Errorf("Expected %d, got %d", len(bits), count)
	}
}

/**
Read
*/

func Test_Measurement_ReadFromBytes(t *testing.T) {
	measure := tiny.NewMeasurement([]byte{170, 85})
	result := measure.Read(6, 10)
	CompareSlices(result, tiny.From.Bits(1, 0, 0, 1), t)
}

func Test_Measurement_ReadFromBits(t *testing.T) {
	measure := tiny.NewMeasurement([]byte{255}, 0, 1, 1, 0, 0)
	result := measure.Read(8, 12)
	CompareSlices(result, tiny.From.Bits(0, 1, 1, 0), t)
}

func Test_Measurement_ReadAcrossBytesAndBits(t *testing.T) {
	measure := tiny.NewMeasurement([]byte{170, 85}, 0, 1, 0, 1, 1, 0)
	expected := tiny.From.Bits(1, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1)
	result := measure.Read(6, 21)
	CompareSlices(result, expected, t)
}

/**
Append
*/

func Test_Measurement_ApppendBytesLengthLimit(t *testing.T) {
	defer ShouldPanic(t)
	measure := tiny.NewMeasurement([]byte{11, 33, 55, 99, 170, 22, 88}, 0, 1, 1)
	measure.AppendBytes(255)
}

func Test_Measurement_ApppendBitsLengthLimit(t *testing.T) {
	defer ShouldPanic(t)
	measure := tiny.NewMeasurement([]byte{11, 33, 55, 99, 170, 22, 88}, 0, 1, 1)
	measure.AppendBits(0, 1, 0, 1, 0, 1, 0, 1)
}

func Test_Measurement_AppendBits(t *testing.T) {
	bits170 := tiny.From.Byte(170)
	measure := tiny.NewMeasurement([]byte{170}, 0, 1, 1)
	measure.AppendBits(0, 0, 1)

	expected := append(bits170, tiny.From.Bits(0, 1, 1)...)
	expected = append(expected, tiny.From.Bits(0, 0, 1)...)
	CompareSlices(measure.GetAllBits(), expected, t)
}

func Test_Measurement_AppendBytes(t *testing.T) {
	measure := tiny.NewMeasurement([]byte{85}, 0, 1, 1)
	measure.AppendBytes(255, 170)

	bits85 := tiny.From.Byte(85)
	bits255 := tiny.From.Byte(255)
	bits170 := tiny.From.Byte(170)

	expected := append(bits85, 0, 1, 1)
	expected = append(expected, bits255...)
	expected = append(expected, bits170...)

	CompareSlices(measure.GetAllBits(), expected, t)
}

func Test_Measurement_Append(t *testing.T) {
	measure := tiny.NewMeasurement([]byte{85}, 0, 1, 1)
	measure.Append(tiny.NewMeasurement([]byte{170}, 1, 0, 0))
	expected := append(tiny.From.Byte(85), tiny.From.Bits(0, 1, 1)...)
	expected = append(expected, tiny.From.Byte(170)...)
	expected = append(expected, tiny.From.Bits(1, 0, 0)...)
	CompareSlices(measure.GetAllBits(), expected, t)
}

/**
Prepend
*/

func Test_Measurement_PrependBytesLengthLimit(t *testing.T) {
	defer ShouldPanic(t)
	measure := tiny.NewMeasurement([]byte{11, 33, 55, 99, 170, 22, 88}, 0, 1, 1)
	measure.PrependBytes(255)
}

func Test_Measurement_PrependBitsLengthLimit(t *testing.T) {
	defer ShouldPanic(t)
	measure := tiny.NewMeasurement([]byte{11, 33, 55, 99, 170, 22, 88}, 0, 1, 1)
	measure.PrependBits(0, 1, 0, 1, 0, 1, 0, 1)
}

func Test_Measurement_PrependBits(t *testing.T) {
	bits170 := tiny.From.Byte(170)
	measure := tiny.NewMeasurement([]byte{170}, 0, 1, 1)
	measure.PrependBits(0, 0, 1)

	expected := append(tiny.From.Bits(0, 0, 1), bits170...)
	expected = append(expected, tiny.From.Bits(0, 1, 1)...)
	CompareSlices(measure.GetAllBits(), expected, t)
}

func Test_Measurement_PrependBytes(t *testing.T) {
	measure := tiny.NewMeasurement([]byte{85}, 0, 1, 1)
	measure.PrependBytes(255, 170)
	expected := append(tiny.From.Bytes(255, 170, 85), tiny.From.Bits(0, 1, 1)...)
	CompareSlices(measure.GetAllBits(), expected, t)
}

func Test_Measurement_Prepend(t *testing.T) {
	measure := tiny.NewMeasurement([]byte{85}, 0, 1, 1)
	measure.Prepend(tiny.NewMeasurement([]byte{170}, 1, 0, 0))
	expected := append(tiny.From.Byte(170), tiny.From.Bits(1, 0, 0)...)
	expected = append(expected, tiny.From.Byte(85)...)
	expected = append(expected, tiny.From.Bits(0, 1, 1)...)
	CompareSlices(measure.GetAllBits(), expected, t)
}

func Test_Measurement_TrimStart(t *testing.T) {
	for i := 0; i < 256; i++ {
		for ii := 0; ii < 8; ii++ {
			measure := tiny.NewMeasurement([]byte{byte(i)})
			bits := measure.GetAllBits()
			measure.TrimStart(ii)
			expected := bits[ii:]
			CompareSlices(measure.GetAllBits(), expected, t)
		}
	}
}

func Test_Measurement_TrimEnd(t *testing.T) {
	for i := 0; i < 256; i++ {
		for ii := 0; ii < 8; ii++ {
			measure := tiny.NewMeasurement([]byte{byte(i)})
			bits := measure.GetAllBits()
			measure.TrimEnd(ii)
			end := len(bits) - ii - 1
			expected := bits[:end]
			CompareSlices(measure.GetAllBits(), expected, t)
		}
	}
}

/**
BreakApart
*/

func Test_Measurement_BreakApart(t *testing.T) {
	m := tiny.NewMeasurement([]byte{}, 0, 1, 1, 0, 1, 1, 0)
	left, right := m.BreakApart(3)
	CompareSlices(left.GetAllBits(), tiny.From.Bits(0, 1, 1), t)
	CompareSlices(right.GetAllBits(), tiny.From.Bits(0, 1, 1, 0), t)
}

func Test_Measurement_BreakApart_Zero(t *testing.T) {
	m := tiny.NewMeasurement([]byte{}, 0, 1, 1, 0, 1, 1, 0)
	left, right := m.BreakApart(0)
	CompareSlices(left.GetAllBits(), tiny.From.Bits(), t)
	CompareSlices(right.GetAllBits(), tiny.From.Bits(0, 1, 1, 0, 1, 1, 0), t)
}

func Test_Measurement_BreakApart_Negative(t *testing.T) {
	defer ShouldPanic(t)
	m := tiny.NewMeasurement([]byte{}, 0, 1, 1, 0, 1, 1, 0)
	m.BreakApart(-1)
}

func Test_Measurement_BreakApart_BeyondBounds(t *testing.T) {
	defer ShouldPanic(t)
	m := tiny.NewMeasurement([]byte{}, 0, 1, 1, 0, 1, 1, 0)
	m.BreakApart(42)
}

/**
Invert
*/

func Test_Measurement_Invert(t *testing.T) {
	expected := tiny.NewMeasurement([]byte{178, 233, 222}, 0, 1, 1, 0)
	m := tiny.NewMeasurement([]byte{77, 22, 33}, 1, 0, 0, 1)
	// |        77       |         22      |        33       |    9    | ← Input Values
	// | 0 1 0 0 1 1 0 1 | 0 0 0 1 0 1 1 0 | 0 0 1 0 0 0 0 1 | 1 0 0 1 | ← Input
	// | 1 0 1 1 0 0 1 0 | 1 1 1 0 1 0 0 1 | 1 1 0 1 1 1 1 0 | 0 1 1 0 | ← Inverted
	// |       178       |        233      |       222       |    6    | ← Inverted Values
	m.Invert()
	CompareMeasurements(m, expected, t)
}
