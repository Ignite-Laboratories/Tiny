package testing

import (
	"ignitelabs.net/git/tiny"
	"testing"
)

func Test_Measure_GetAllBits(t *testing.T) {
	bits := []tiny.Bit{0, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0}
	expected := append([]tiny.Bit{1, 0, 1, 0, 1, 0, 1, 0}, bits...)
	measure := tiny.NewMeasure([]byte{170}, bits...)
	result := measure.GetAllBits()
	CompareBitSlices(result, expected, t)
}

func Test_Measure_BitLength(t *testing.T) {
	bits := []tiny.Bit{0, 1, 0}
	bytes := []byte{170, 85}

	bitMeasure := tiny.NewMeasure([]byte{}, bits...)
	CompareValues(bitMeasure.BitLength(), 3, t)

	byteMeasure := tiny.NewMeasure(bytes)
	CompareValues(byteMeasure.BitLength(), 16, t)

	bothMeasure := tiny.NewMeasure(bytes, bits...)
	CompareValues(bothMeasure.BitLength(), 19, t)
}

func Test_Measure_ByteLength(t *testing.T) {
	bytes := []byte{170, 85, 38, 75}

	for i := 0; i < len(bytes); i++ {
		toTest := bytes[:len(bytes)-i]
		m := tiny.NewMeasure(toTest)
		CompareValues(m.ByteBitLength(), len(toTest)*8, t)
	}
}

func Test_Measure_Value(t *testing.T) {
	m := tiny.NewMeasure([]byte{170, 85}, 1, 0, 0, 1)
	CompareValues(m.Value(), 697689, t)
}

func Test_Measure_Toggle(t *testing.T) {
	bytes := []byte{255, 0, 128, 127, 77}
	inverseBytes := []byte{0, 255, 127, 128, 178}
	bits := tiny.From.Bits(0, 1, 0, 1)
	inverseBits := tiny.From.Bits(1, 0, 1, 0)
	m := tiny.NewMeasure(bytes, bits...)
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

func Test_Measure_ForEachBit(t *testing.T) {
	random := tiny.Synthesize.Random(22)
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

func Test_Measure_ReadFromBytes(t *testing.T) {
	measure := tiny.NewMeasure([]byte{170, 85})
	result := measure.Read(6, 10)
	CompareBitSlices(result, tiny.From.Bits(1, 0, 0, 1), t)
}

func Test_Measure_ReadFromBits(t *testing.T) {
	measure := tiny.NewMeasure([]byte{255}, 0, 1, 1, 0, 0)
	result := measure.Read(8, 12)
	CompareBitSlices(result, tiny.From.Bits(0, 1, 1, 0), t)
}

func Test_Measure_ReadAcrossBytesAndBits(t *testing.T) {
	measure := tiny.NewMeasure([]byte{170, 85}, 0, 1, 0, 1, 1, 0)
	expected := tiny.From.Bits(1, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1)
	result := measure.Read(6, 21)
	CompareBitSlices(result, expected, t)
}

/**
Append
*/

func Test_Measure_ApppendBytesLengthLimit(t *testing.T) {
	defer ShouldPanic(t)
	measure := tiny.NewMeasure([]byte{170, 22, 88}, 0, 1, 1)
	measure.AppendBytes(255)
}

func Test_Measure_ApppendBitsLengthLimit(t *testing.T) {
	defer ShouldPanic(t)
	measure := tiny.NewMeasure([]byte{170, 22, 88}, 0, 1, 1)
	measure.AppendBits(0, 1, 0, 1, 0, 1, 0, 1)
}

func Test_Measure_AppendBits(t *testing.T) {
	bits170 := tiny.From.Byte(170)
	measure := tiny.NewMeasure([]byte{170}, 0, 1, 1)
	measure.AppendBits(0, 0, 1)

	expected := append(bits170, tiny.From.Bits(0, 1, 1)...)
	expected = append(expected, tiny.From.Bits(0, 0, 1)...)
	CompareBitSlices(measure.GetAllBits(), expected, t)
}

func Test_Measure_AppendBytes(t *testing.T) {
	measure := tiny.NewMeasure([]byte{85}, 0, 1, 1)
	measure.AppendBytes(255, 170)

	bits85 := tiny.From.Byte(85)
	bits255 := tiny.From.Byte(255)
	bits170 := tiny.From.Byte(170)

	expected := append(bits85, 0, 1, 1)
	expected = append(expected, bits255...)
	expected = append(expected, bits170...)

	CompareBitSlices(measure.GetAllBits(), expected, t)
}

func Test_Measure_Append(t *testing.T) {
	measure := tiny.NewMeasure([]byte{85}, 0, 1, 1)
	measure.Append(tiny.NewMeasure([]byte{170}, 1, 0, 0))
	expected := append(tiny.From.Byte(85), tiny.From.Bits(0, 1, 1)...)
	expected = append(expected, tiny.From.Byte(170)...)
	expected = append(expected, tiny.From.Bits(1, 0, 0)...)
	CompareBitSlices(measure.GetAllBits(), expected, t)
}

/**
Prepend
*/

func Test_Measure_PrependBytesLengthLimit(t *testing.T) {
	defer ShouldPanic(t)
	measure := tiny.NewMeasure([]byte{170, 22, 88}, 0, 1, 1)
	measure.PrependBytes(255)
}

func Test_Measure_PrependBitsLengthLimit(t *testing.T) {
	defer ShouldPanic(t)
	measure := tiny.NewMeasure([]byte{170, 22, 88}, 0, 1, 1)
	measure.PrependBits(0, 1, 0, 1, 0, 1, 0, 1)
}

func Test_Measure_PrependBits(t *testing.T) {
	bits170 := tiny.From.Byte(170)
	measure := tiny.NewMeasure([]byte{170}, 0, 1, 1)
	measure.PrependBits(0, 0, 1)

	expected := append(tiny.From.Bits(0, 0, 1), bits170...)
	expected = append(expected, tiny.From.Bits(0, 1, 1)...)
	CompareBitSlices(measure.GetAllBits(), expected, t)
}

func Test_Measure_PrependBytes(t *testing.T) {
	measure := tiny.NewMeasure([]byte{85}, 0, 1, 1)
	measure.PrependBytes(255, 170)
	expected := append(tiny.From.Bytes(255, 170, 85), tiny.From.Bits(0, 1, 1)...)
	CompareBitSlices(measure.GetAllBits(), expected, t)
}

func Test_Measure_Prepend(t *testing.T) {
	measure := tiny.NewMeasure([]byte{85}, 0, 1, 1)
	measure.Prepend(tiny.NewMeasure([]byte{170}, 1, 0, 0))
	expected := append(tiny.From.Byte(170), tiny.From.Bits(1, 0, 0)...)
	expected = append(expected, tiny.From.Byte(85)...)
	expected = append(expected, tiny.From.Bits(0, 1, 1)...)
	CompareBitSlices(measure.GetAllBits(), expected, t)
}
