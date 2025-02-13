package testing

import (
	"github.com/ignite-laboratories/tiny"
	"testing"
)

/**
Bit/Measure/Shade
*/

func Test_Analyze_Shade(t *testing.T) {
	light := tiny.NewMeasure([]byte{0, 0, 0, 0}, 0, 0, 0)
	lightExpected := tiny.BinaryShade{
		Zeros:             35,
		Ones:              0,
		Total:             35,
		Shade:             tiny.Light,
		PredominantlyDark: false,
		Distribution:      [8]int{0, 0, 0, 0, 0, 0, 0, 0},
	}
	shadeTester(tiny.Analyze.Shade(light), lightExpected, t)

	dark := tiny.NewMeasure([]byte{255, 255, 255, 255}, 1, 1, 1)
	darkExpected := tiny.BinaryShade{
		Zeros:             0,
		Ones:              35,
		Total:             35,
		Shade:             tiny.Dark,
		PredominantlyDark: true,
		Distribution:      [8]int{5, 5, 5, 4, 4, 4, 4, 4},
	}
	shadeTester(tiny.Analyze.Shade(dark), darkExpected, t)

	jumbled := tiny.NewMeasure([]byte{22, 222, 111, 144}, 0, 1, 0)
	jumbledExpected := tiny.BinaryShade{
		Zeros:             17,
		Ones:              18,
		Total:             35,
		Shade:             tiny.Grey,
		PredominantlyDark: true,
		Distribution:      [8]int{2, 3, 1, 3, 2, 3, 3, 1},
	}
	shadeTester(tiny.Analyze.Shade(jumbled), jumbledExpected, t)

	lessThanHalfGrey := tiny.NewMeasure([]byte{7, 7, 7, 7}, 0, 1, 0)
	lessThanHalfGreyExpected := tiny.BinaryShade{
		Zeros:             22,
		Ones:              13,
		Total:             35,
		Shade:             tiny.Grey,
		PredominantlyDark: false,
		Distribution:      [8]int{0, 1, 0, 0, 0, 4, 4, 4},
	}
	shadeTester(tiny.Analyze.Shade(lessThanHalfGrey), lessThanHalfGreyExpected, t)

	halfGrey := tiny.NewMeasure([]byte{15, 15, 15, 15}, 0, 1, 0)
	halfGreyExpected := tiny.BinaryShade{
		Zeros:             18,
		Ones:              17,
		Total:             35,
		Shade:             tiny.Grey,
		PredominantlyDark: false,
		Distribution:      [8]int{0, 1, 0, 0, 4, 4, 4, 4},
	}
	shadeTester(tiny.Analyze.Shade(halfGrey), halfGreyExpected, t)

	moreThanHalfGrey := tiny.NewMeasure([]byte{31, 31, 31, 31}, 0, 1, 0)
	moreThanHalfGreyExpected := tiny.BinaryShade{
		Zeros:             14,
		Ones:              21,
		Total:             35,
		Shade:             tiny.Grey,
		PredominantlyDark: true,
		Distribution:      [8]int{0, 1, 0, 4, 4, 4, 4, 4},
	}
	shadeTester(tiny.Analyze.Shade(moreThanHalfGrey), moreThanHalfGreyExpected, t)
}

func Test_Analyze_ByteShade(t *testing.T) {
	light := []byte{0, 0, 0, 0}
	lightExpected := tiny.BinaryShade{
		Zeros:             32,
		Ones:              0,
		Total:             32,
		Shade:             tiny.Light,
		PredominantlyDark: false,
		Distribution:      [8]int{0, 0, 0, 0, 0, 0, 0, 0},
	}
	shadeTester(tiny.Analyze.ByteShade(light...), lightExpected, t)

	dark := []byte{255, 255, 255, 255}
	darkExpected := tiny.BinaryShade{
		Zeros:             0,
		Ones:              32,
		Total:             32,
		Shade:             tiny.Dark,
		PredominantlyDark: true,
		Distribution:      [8]int{4, 4, 4, 4, 4, 4, 4, 4},
	}
	shadeTester(tiny.Analyze.ByteShade(dark...), darkExpected, t)

	jumbled := []byte{22, 222, 111, 144}
	jumbledExpected := tiny.BinaryShade{
		Zeros:             15,
		Ones:              17,
		Total:             32,
		Shade:             tiny.Grey,
		PredominantlyDark: true,
		Distribution:      [8]int{2, 2, 1, 3, 2, 3, 3, 1},
	}
	shadeTester(tiny.Analyze.ByteShade(jumbled...), jumbledExpected, t)

	lessThanHalfGrey := []byte{7, 7, 7, 7}
	lessThanHalfGreyExpected := tiny.BinaryShade{
		Zeros:             20,
		Ones:              12,
		Total:             32,
		Shade:             tiny.Grey,
		PredominantlyDark: false,
		Distribution:      [8]int{0, 0, 0, 0, 0, 4, 4, 4},
	}
	shadeTester(tiny.Analyze.ByteShade(lessThanHalfGrey...), lessThanHalfGreyExpected, t)

	halfGrey := []byte{15, 15, 15, 15}
	halfGreyExpected := tiny.BinaryShade{
		Zeros:             16,
		Ones:              16,
		Total:             32,
		Shade:             tiny.Grey,
		PredominantlyDark: false,
		Distribution:      [8]int{0, 0, 0, 0, 4, 4, 4, 4},
	}
	shadeTester(tiny.Analyze.ByteShade(halfGrey...), halfGreyExpected, t)

	moreThanHalfGrey := []byte{31, 31, 31, 31}
	moreThanHalfGreyExpected := tiny.BinaryShade{
		Zeros:             12,
		Ones:              20,
		Total:             32,
		Shade:             tiny.Grey,
		PredominantlyDark: true,
		Distribution:      [8]int{0, 0, 0, 4, 4, 4, 4, 4},
	}
	shadeTester(tiny.Analyze.ByteShade(moreThanHalfGrey...), moreThanHalfGreyExpected, t)
}

func Test_Analyze_BitShade(t *testing.T) {
	light := tiny.From.Byte(0)
	lightExpected := tiny.BinaryShade{
		Zeros:             8,
		Ones:              0,
		Total:             8,
		Shade:             tiny.Light,
		PredominantlyDark: false,
		Distribution:      [8]int{0, 0, 0, 0, 0, 0, 0, 0},
	}
	shadeTester(tiny.Analyze.BitShade(light...), lightExpected, t)

	dark := tiny.From.Byte(255)
	darkExpected := tiny.BinaryShade{
		Zeros:             0,
		Ones:              8,
		Total:             8,
		Shade:             tiny.Dark,
		PredominantlyDark: true,
		Distribution:      [8]int{1, 1, 1, 1, 1, 1, 1, 1},
	}
	shadeTester(tiny.Analyze.BitShade(dark...), darkExpected, t)

	unevenLengthOfBits := []tiny.Bit{0, 1, 0, 0, 1, 0, 0, 1, 1, 1}
	jumbledGreyExpected := tiny.BinaryShade{
		Zeros:             5,
		Ones:              5,
		Total:             10,
		Shade:             tiny.Grey,
		PredominantlyDark: false,
		Distribution:      [8]int{1, 2, 0, 0, 1, 0, 0, 1},
	}
	shadeTester(tiny.Analyze.BitShade(unevenLengthOfBits...), jumbledGreyExpected, t)

	lessThanHalfGrey := []tiny.Bit{1, 1, 1, 0, 0, 0, 0, 0}
	lessThanHalfGreyExpected := tiny.BinaryShade{
		Zeros:             5,
		Ones:              3,
		Total:             8,
		Shade:             tiny.Grey,
		PredominantlyDark: false,
		Distribution:      [8]int{1, 1, 1, 0, 0, 0, 0, 0},
	}
	shadeTester(tiny.Analyze.BitShade(lessThanHalfGrey...), lessThanHalfGreyExpected, t)

	halfGrey := []tiny.Bit{1, 1, 1, 1, 0, 0, 0, 0}
	halfGreyExpected := tiny.BinaryShade{
		Zeros:             4,
		Ones:              4,
		Total:             8,
		Shade:             tiny.Grey,
		PredominantlyDark: false,
		Distribution:      [8]int{1, 1, 1, 1, 0, 0, 0, 0},
	}
	shadeTester(tiny.Analyze.BitShade(halfGrey...), halfGreyExpected, t)

	moreThanHalfGrey := []tiny.Bit{1, 1, 1, 1, 1, 0, 0, 0}
	moreThanHalfGreyExpected := tiny.BinaryShade{
		Zeros:             3,
		Ones:              5,
		Total:             8,
		Shade:             tiny.Grey,
		PredominantlyDark: true,
		Distribution:      [8]int{1, 1, 1, 1, 1, 0, 0, 0},
	}
	shadeTester(tiny.Analyze.BitShade(moreThanHalfGrey...), moreThanHalfGreyExpected, t)
}

func shadeTester(analysis tiny.BinaryShade, expected tiny.BinaryShade, t *testing.T) {
	if analysis.Shade != expected.Shade {
		t.Errorf("Expected %v, got %v", expected.Shade, analysis.Shade)
	}
	if analysis.PredominantlyDark != expected.PredominantlyDark {
		if expected.PredominantlyDark {
			t.Errorf("Data is predominantly dark but the analysis said it wasn't")
		} else {
			t.Errorf("Data is not predominantly dark but the analysis said it was")
		}
	}
	if analysis.Ones != expected.Ones {
		t.Errorf("Expected %d ones, got %d", expected.Ones, analysis.Ones)
	}
	if analysis.Zeros != expected.Zeros {
		t.Errorf("Expected %d zeros, got %d", expected.Zeros, analysis.Zeros)

	}
	if analysis.Total != expected.Total {
		t.Errorf("Expected a total of %d, got %d", expected.Total, analysis.Total)
	}
	for i := 0; i < len(analysis.Distribution); i++ {
		if analysis.Distribution[i] != expected.Distribution[i] {
			t.Errorf("Expected %d at distribution index %d, got %d", expected.Distribution[i], i, analysis.Distribution[i])
		}
	}
}

/**
HasPrefix
*/

func Test_Analyze_HasPrefix_Failure(t *testing.T) {
	data := []tiny.Bit{0, 1, 0, 0, 1, 0, 0, 1, 1}
	prefix := []tiny.Bit{1, 0, 1, 1}
	if tiny.Analyze.HasPrefix(data, prefix...) {
		t.Errorf("%v is not the prefix", tiny.To.String(prefix...))
	}
}

func Test_Analyze_HasPrefix_Static(t *testing.T) {
	data := []tiny.Bit{0, 1, 0, 0, 1, 0, 0, 1, 1}
	prefix := []tiny.Bit{0, 1, 0, 0}
	if !tiny.Analyze.HasPrefix(data, prefix...) {
		t.Errorf("Data did not have the prefix %v", tiny.To.String(prefix...))
	}
}

func Test_Analyze_HasPrefix_Random(t *testing.T) {
	random := tiny.Synthesize.Random(8)
	bits := random.GetAllBits()
	prefix := bits[:5]
	if !tiny.Analyze.HasPrefix(bits, prefix...) {
		t.Errorf("Data did not have the prefix %v", tiny.To.String(prefix...))
	}
}

func Test_Analyze_HasPrefix_Synthesized(t *testing.T) {
	// This walks all the values of a Note [0-7] and then checks if that prefix exists
	// in a synthesized pattern of that data.
	var prefixes [][]tiny.Bit
	for i := 0; i <= tiny.MaxNote; i++ {
		prefixes = append(prefixes, tiny.From.Number(i, 3))
	}

	for _, prefix := range prefixes {
		repeating := tiny.Synthesize.Repeating(4, prefix...)
		bits := repeating.GetAllBits()
		if !tiny.Analyze.HasPrefix(bits, prefix...) {
			t.Errorf("Data did not have the prefix %v", tiny.To.String(prefix...))
		}
	}
}

/**
OneDistribution
*/

func Test_Analyze_OneDistribution_FullSpectrum(t *testing.T) {
	// With this arrangement each byte has one more one than the last in each index
	// This means the count of 1s per index should be i+1
	bytes := []byte{0, 1, 3, 7, 15, 31, 63, 127, 255}
	ones := tiny.Analyze.OneDistribution(bytes...)
	for i, index := range ones {
		if index != i+1 {
			t.Errorf("Invalid one count")
		}
	}
}

func Test_Analyze_OneDistribution_Light(t *testing.T) {
	bytes := []byte{0, 0, 0, 0, 0}
	ones := tiny.Analyze.OneDistribution(bytes...)
	for _, index := range ones {
		if index != 0 {
			t.Errorf("Should not have any ones")
		}
	}
}

func Test_Analyze_OneDistribution_Dark(t *testing.T) {
	bytes := []byte{255, 255, 255, 255, 255}
	ones := tiny.Analyze.OneDistribution(bytes...)
	for _, index := range ones {
		if index != 5 {
			t.Errorf("Should have exactly 5 ones")
		}
	}
}

func Test_Analyze_OneDistribution_Grey(t *testing.T) {
	bytes := []byte{42, 22, 88, 222, 133}
	expected := []int{2, 2, 1, 3, 3, 3, 3, 1}
	ones := tiny.Analyze.OneDistribution(bytes...)
	for i := 0; i < len(expected); i++ {
		if ones[i] != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], ones[i])
		}
	}
}
