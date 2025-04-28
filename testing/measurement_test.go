package testing

import (
	"github.com/ignite-laboratories/core/test"
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Measurement_GetAllBits(t *testing.T) {
	bits := []tiny.Bit{0, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0}
	expected := append([]tiny.Bit{1, 0, 1, 0, 1, 0, 1, 0}, bits...)
	measure := tiny.NewMeasurement([]byte{170}, bits...)
	result := measure.GetAllBits()
	test.CompareSlices(result, expected, t)
}

func Test_Measurement_BitLength(t *testing.T) {
	bits := []tiny.Bit{0, 1, 0}
	bytes := []byte{170, 85}

	bitMeasure := tiny.NewMeasurement([]byte{}, bits...)
	test.CompareValues(bitMeasure.BitLength(), 3, t)

	byteMeasure := tiny.NewMeasurement(bytes)
	test.CompareValues(byteMeasure.BitLength(), 16, t)

	bothMeasure := tiny.NewMeasurement(bytes, bits...)
	test.CompareValues(bothMeasure.BitLength(), 19, t)
}

func Test_Measurement_ByteLength(t *testing.T) {
	bytes := []byte{170, 85, 38, 75}

	for i := 0; i < len(bytes); i++ {
		toTest := bytes[:len(bytes)-i]
		m := tiny.NewMeasurement(toTest)
		test.CompareValues(m.ByteBitLength(), len(toTest)*8, t)
	}
}

func Test_Measurement_Value(t *testing.T) {
	m := tiny.NewMeasurement([]byte{170, 85}, 1, 0, 0, 1)
	test.CompareValues(m.Value(), 697689, t)
}

func Test_Measurement_Clear(t *testing.T) {
	m := tiny.NewMeasurement([]byte{170, 85}, 1, 0, 0, 1)
	m.Clear()
	test.CompareSlices(m.GetAllBits(), []tiny.Bit{}, t)
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

func Test_Measurement_ReadFromBytes(t *testing.T) {
	measure := tiny.NewMeasurement([]byte{170, 85})
	result := measure.Read(6, 10)
	test.CompareSlices(result, tiny.From.Bits(1, 0, 0, 1), t)
}

func Test_Measurement_ReadFromBits(t *testing.T) {
	measure := tiny.NewMeasurement([]byte{255}, 0, 1, 1, 0, 0)
	result := measure.Read(8, 12)
	test.CompareSlices(result, tiny.From.Bits(0, 1, 1, 0), t)
}

func Test_Measurement_ReadAcrossBytesAndBits(t *testing.T) {
	measure := tiny.NewMeasurement([]byte{170, 85}, 0, 1, 0, 1, 1, 0)
	expected := tiny.From.Bits(1, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1)
	result := measure.Read(6, 21)
	test.CompareSlices(result, expected, t)
}

/**
Append
*/

func Test_Measurement_ApppendBytesLengthLimit(t *testing.T) {
	defer test.ShouldPanic(t)
	measure := tiny.NewMeasurement([]byte{170, 22, 88}, 0, 1, 1)
	measure.AppendBytes(255)
}

func Test_Measurement_ApppendBitsLengthLimit(t *testing.T) {
	defer test.ShouldPanic(t)
	measure := tiny.NewMeasurement([]byte{170, 22, 88}, 0, 1, 1)
	measure.AppendBits(0, 1, 0, 1, 0, 1, 0, 1)
}

func Test_Measurement_AppendBits(t *testing.T) {
	bits170 := tiny.From.Byte(170)
	measure := tiny.NewMeasurement([]byte{170}, 0, 1, 1)
	measure.AppendBits(0, 0, 1)

	expected := append(bits170, tiny.From.Bits(0, 1, 1)...)
	expected = append(expected, tiny.From.Bits(0, 0, 1)...)
	test.CompareSlices(measure.GetAllBits(), expected, t)
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

	test.CompareSlices(measure.GetAllBits(), expected, t)
}

func Test_Measurement_Append(t *testing.T) {
	measure := tiny.NewMeasurement([]byte{85}, 0, 1, 1)
	measure.Append(tiny.NewMeasurement([]byte{170}, 1, 0, 0))
	expected := append(tiny.From.Byte(85), tiny.From.Bits(0, 1, 1)...)
	expected = append(expected, tiny.From.Byte(170)...)
	expected = append(expected, tiny.From.Bits(1, 0, 0)...)
	test.CompareSlices(measure.GetAllBits(), expected, t)
}

/**
Prepend
*/

func Test_Measurement_PrependBytesLengthLimit(t *testing.T) {
	defer test.ShouldPanic(t)
	measure := tiny.NewMeasurement([]byte{170, 22, 88}, 0, 1, 1)
	measure.PrependBytes(255)
}

func Test_Measurement_PrependBitsLengthLimit(t *testing.T) {
	defer test.ShouldPanic(t)
	measure := tiny.NewMeasurement([]byte{170, 22, 88}, 0, 1, 1)
	measure.PrependBits(0, 1, 0, 1, 0, 1, 0, 1)
}

func Test_Measurement_PrependBits(t *testing.T) {
	bits170 := tiny.From.Byte(170)
	measure := tiny.NewMeasurement([]byte{170}, 0, 1, 1)
	measure.PrependBits(0, 0, 1)

	expected := append(tiny.From.Bits(0, 0, 1), bits170...)
	expected = append(expected, tiny.From.Bits(0, 1, 1)...)
	test.CompareSlices(measure.GetAllBits(), expected, t)
}

func Test_Measurement_PrependBytes(t *testing.T) {
	measure := tiny.NewMeasurement([]byte{85}, 0, 1, 1)
	measure.PrependBytes(255, 170)
	expected := append(tiny.From.Bytes(255, 170, 85), tiny.From.Bits(0, 1, 1)...)
	test.CompareSlices(measure.GetAllBits(), expected, t)
}

func Test_Measurement_Prepend(t *testing.T) {
	measure := tiny.NewMeasurement([]byte{85}, 0, 1, 1)
	measure.Prepend(tiny.NewMeasurement([]byte{170}, 1, 0, 0))
	expected := append(tiny.From.Byte(170), tiny.From.Bits(1, 0, 0)...)
	expected = append(expected, tiny.From.Byte(85)...)
	expected = append(expected, tiny.From.Bits(0, 1, 1)...)
	test.CompareSlices(measure.GetAllBits(), expected, t)
}

func Test_Measurement_UnQuarterSplit(t *testing.T) {
	for i := 0; i < 256; i++ {
		expected := []byte{byte(i)}
		data := tiny.NewMeasurement(expected)
		data.QuarterSplit()
		data.UnQuarterSplit()
		test.CompareSlices(expected, data.Bytes, t)
	}
}

func Test_Measurement_QuarterSplit(t *testing.T) {
	for i := 0; i < 256; i++ {
		measure := tiny.NewMeasurement([]byte{byte(i)})
		measure.QuarterSplit()

		if i < 64 {
			if measure.BitLength() != 7 {
				t.Errorf("Expected reduction for 0-63")
			}
			if measure.Read(0, 1)[0] != 0 {
				t.Errorf("Expected a code of '0' for 0-63")
			}
			valueBits := measure.Read(1, 7)
			value := tiny.To.Byte(valueBits...)
			if value != byte(i) {
				t.Errorf("Expected a value of %d, got %d", i, value)
			}
		} else {
			code := measure.Read(0, 2)
			if i < 128 {
				if measure.BitLength() != 8 {
					t.Errorf("Expected no reduction for 64-127")
				}
				if code[0] != 1 && code[1] != 0 {
					t.Errorf("Expected a code of '10' for 64-127")
				}
				valueBits := measure.Read(2, 8)
				value := tiny.To.Byte(valueBits...)
				value += 64
				if value != byte(i) {
					t.Errorf("Expected a value of %d, got %d", i, value)
				}
			} else {
				if measure.BitLength() != 9 {
					t.Errorf("Expected 1 growth for 128+")
				}
				if code[0] != 1 && code[1] != 1 {
					t.Errorf("Expected a code of '11' for 128+")
				}
				valueBits := measure.Read(2, 9)
				value := tiny.To.Byte(valueBits...)
				value += 128
				if value != byte(i) {
					t.Errorf("Expected a value of %d, got %d", i, value)
				}
			}
		}
	}
}

func Test_Measurement_TrimStart(t *testing.T) {
	for i := 0; i < 256; i++ {
		for ii := 0; ii < 8; ii++ {
			measure := tiny.NewMeasurement([]byte{byte(i)})
			bits := measure.GetAllBits()
			measure.TrimStart(ii)
			expected := bits[ii:]
			test.CompareSlices(measure.GetAllBits(), expected, t)
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
			test.CompareSlices(measure.GetAllBits(), expected, t)
		}
	}
}

/**
BreakApart
*/

func Test_Measurement_BreakApart(t *testing.T) {
	m := tiny.NewMeasurement([]byte{}, 0, 1, 1, 0, 1, 1, 0)
	left, right := m.BreakApart(3)
	test.CompareSlices(left.GetAllBits(), tiny.From.Bits(0, 1, 1), t)
	test.CompareSlices(right.GetAllBits(), tiny.From.Bits(0, 1, 1, 0), t)
}

func Test_Measurement_BreakApart_Zero(t *testing.T) {
	m := tiny.NewMeasurement([]byte{}, 0, 1, 1, 0, 1, 1, 0)
	left, right := m.BreakApart(0)
	test.CompareSlices(left.GetAllBits(), tiny.From.Bits(), t)
	test.CompareSlices(right.GetAllBits(), tiny.From.Bits(0, 1, 1, 0, 1, 1, 0), t)
}

func Test_Measurement_BreakApart_Negative(t *testing.T) {
	defer test.ShouldPanic(t)
	m := tiny.NewMeasurement([]byte{}, 0, 1, 1, 0, 1, 1, 0)
	m.BreakApart(-1)
}

func Test_Measurement_BreakApart_BeyondBounds(t *testing.T) {
	defer test.ShouldPanic(t)
	m := tiny.NewMeasurement([]byte{}, 0, 1, 1, 0, 1, 1, 0)
	m.BreakApart(42)
}
