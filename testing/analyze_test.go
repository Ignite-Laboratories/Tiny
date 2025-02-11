package testing

import (
	"github.com/ignite-laboratories/tiny"
	"testing"
)

/**
Bit/Data/Shade
*/

func Test_RemainderShade(t *testing.T) {
	light := tiny.NewRemainder([]byte{0, 0, 0, 0}, tiny.Create.Zeros(8)...)
	shadeTester(tiny.Analyze.RemainderShade(light), tiny.Light, false, t)

	dark := tiny.NewRemainder([]byte{255, 255, 255, 255}, tiny.Create.Ones(8)...)
	shadeTester(tiny.Analyze.RemainderShade(dark), tiny.Dark, true, t)

	jumbled := tiny.NewRemainder([]byte{22, 222, 111, 144}, []tiny.Bit{0, 1, 0, 0, 1, 0, 0, 1, 1}...)
	shadeTester(tiny.Analyze.RemainderShade(jumbled), tiny.Grey, true, t)

	lessThanHalfGrey := tiny.NewRemainder([]byte{7, 7, 7, 7}, []tiny.Bit{1, 1, 1, 0, 0, 0, 0, 0}...)
	shadeTester(tiny.Analyze.RemainderShade(lessThanHalfGrey), tiny.Grey, false, t)

	halfGrey := tiny.NewRemainder([]byte{15, 15, 15, 15}, []tiny.Bit{1, 1, 1, 1, 0, 0, 0, 0}...)
	shadeTester(tiny.Analyze.RemainderShade(halfGrey), tiny.Grey, false, t)

	moreThanHalfGrey := tiny.NewRemainder([]byte{31, 31, 31, 31}, []tiny.Bit{1, 1, 1, 1, 1, 0, 0, 0}...)
	shadeTester(tiny.Analyze.RemainderShade(moreThanHalfGrey), tiny.Grey, true, t)
}

func Test_ByteShade(t *testing.T) {
	light := []byte{0, 0, 0, 0}
	shadeTester(tiny.Analyze.ByteShade(light...), tiny.Light, false, t)

	dark := []byte{255, 255, 255, 255}
	shadeTester(tiny.Analyze.ByteShade(dark...), tiny.Dark, true, t)

	jumbled := []byte{22, 222, 111, 144}
	shadeTester(tiny.Analyze.ByteShade(jumbled...), tiny.Grey, true, t)

	lessThanHalfGrey := []byte{7, 7, 7, 7}
	shadeTester(tiny.Analyze.ByteShade(lessThanHalfGrey...), tiny.Grey, false, t)

	halfGrey := []byte{15, 15, 15, 15}
	shadeTester(tiny.Analyze.ByteShade(halfGrey...), tiny.Grey, false, t)

	moreThanHalfGrey := []byte{31, 31, 31, 31}
	shadeTester(tiny.Analyze.ByteShade(moreThanHalfGrey...), tiny.Grey, true, t)

}

func Test_BitShade(t *testing.T) {
	light := tiny.Create.Zeros(8)
	shadeTester(tiny.Analyze.BitShade(light...), tiny.Light, false, t)

	dark := tiny.Create.Ones(8)
	shadeTester(tiny.Analyze.BitShade(dark...), tiny.Dark, true, t)

	jumbledGrey := []tiny.Bit{0, 1, 0, 0, 1, 0, 0, 1, 1}
	shadeTester(tiny.Analyze.BitShade(jumbledGrey...), tiny.Grey, false, t)

	lessThanHalfGrey := []tiny.Bit{1, 1, 1, 0, 0, 0, 0, 0}
	shadeTester(tiny.Analyze.BitShade(lessThanHalfGrey...), tiny.Grey, false, t)

	halfGrey := []tiny.Bit{1, 1, 1, 1, 0, 0, 0, 0}
	shadeTester(tiny.Analyze.BitShade(halfGrey...), tiny.Grey, false, t)

	MoreThanHalfGrey := []tiny.Bit{1, 1, 1, 1, 1, 0, 0, 0}
	shadeTester(tiny.Analyze.BitShade(MoreThanHalfGrey...), tiny.Grey, true, t)
}

func shadeTester(analysis tiny.BinaryCount, shade tiny.Shade, predominantlyDark bool, t *testing.T) {
	if analysis.Shade != shade {
		t.Errorf("Expected %v, got %v", shade, analysis.Shade)
	}
	if analysis.PredominantlyDark != predominantlyDark {
		if predominantlyDark {
			t.Errorf("Data is predominantly dark but the analysis said it wasn't")
		} else {
			t.Errorf("Data is not predominantly dark but the analysis said it was")
		}
	}
}

/**
HasPrefix
*/

func Test_HasPrefix_Failure(t *testing.T) {
	data := []tiny.Bit{0, 1, 0, 0, 1, 0, 0, 1, 1}
	prefix := []tiny.Bit{1, 0, 1, 1}
	if tiny.Analyze.HasPrefix(data, prefix...) {
		t.Errorf("%v is not the prefix", tiny.To.String(prefix...))
	}
}

func Test_HasPrefix_Static(t *testing.T) {
	data := []tiny.Bit{0, 1, 0, 0, 1, 0, 0, 1, 1}
	prefix := []tiny.Bit{0, 1, 0, 0}
	if !tiny.Analyze.HasPrefix(data, prefix...) {
		t.Errorf("Data did not have the prefix %v", tiny.To.String(prefix...))
	}
}

func Test_HasPrefix_Random(t *testing.T) {
	data := tiny.Create.Random(8)
	prefix := data[:5]
	if !tiny.Analyze.HasPrefix(data, prefix...) {
		t.Errorf("Data did not have the prefix %v", tiny.To.String(prefix...))
	}
}

func Test_HasPrefix_Synthesized(t *testing.T) {
	// This walks all the values of a Note [0-7] and then checks if that prefix exists
	// in a synthesized pattern of that data.
	var prefixes [][]tiny.Bit
	for i := 0; i <= tiny.MaxNote; i++ {
		prefixes = append(prefixes, tiny.From.Number(i, 3))
	}

	for _, prefix := range prefixes {
		data := tiny.Create.Repeating(4, prefix...)
		if !tiny.Analyze.HasPrefix(data, prefix...) {
			t.Errorf("Data did not have the prefix %v", tiny.To.String(prefix...))
		}
	}
}

/**
OneDistribution
*/

func Test_OneDistribution_FullSpectrum(t *testing.T) {
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

func Test_OneDistribution_Light(t *testing.T) {
	bytes := []byte{0, 0, 0, 0, 0}
	ones := tiny.Analyze.OneDistribution(bytes...)
	for _, index := range ones {
		if index != 0 {
			t.Errorf("Should not have any ones")
		}
	}
}

func Test_OneDistribution_Dark(t *testing.T) {
	bytes := []byte{255, 255, 255, 255, 255}
	ones := tiny.Analyze.OneDistribution(bytes...)
	for _, index := range ones {
		if index != 5 {
			t.Errorf("Should have exactly 5 ones")
		}
	}
}

func Test_OneDistribution_Grey(t *testing.T) {
	bytes := []byte{42, 22, 88, 222, 133}
	expected := []int{2, 2, 1, 3, 3, 3, 3, 1}
	ones := tiny.Analyze.OneDistribution(bytes...)
	for i := 0; i < len(expected); i++ {
		if ones[i] != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], ones[i])
		}
	}
}
