package testing

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
	"testing"
)

func Test_Synthesize_FuzzyApproximate(t *testing.T) {
	data := tiny.Synthesize.RandomPhrase(32)
	_, data, _ = data.ReadBit()
	data = data.PrependBits(1)

	approx := tiny.Fuzzy.Approximation(data.AsBigInt(), 3)

	fmt.Println(approx.Indices)
	fmt.Println(approx.Target.Text(2))
	fmt.Println(approx.Value.Text(2))
	fmt.Println(approx.Delta.Text(2))

	fmt.Println(approx.Target)
	fmt.Println(approx.Value)
	fmt.Println(approx.Delta)
	fmt.Println(approx.Relativity)
}

func Test_Synthesize_FuzzyApproximate3(t *testing.T) {
	counts := make([]int, 8)

	for i := 0; i < 256; i++ {
		results := approximate()

		for ii := 0; ii < 8; ii++ {
			counts[ii] += results[ii]
		}
	}

	for ii := 0; ii < 8; ii++ {
		counts[ii] /= 256
		fmt.Printf("[%d] %d\n", ii, counts[ii])
	}
}

func approximate() []int {
	out := make([]int, 8)
	data := tiny.Synthesize.RandomPhrase(32)
	_, data, _ = data.ReadBit()
	data = data.PrependBits(1)

	target := data.AsBigInt()
	bitLen := target.BitLen()

	approx := tiny.Fuzzy.Approximation(target, 1)
	out[0] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 2)
	out[1] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 3)
	out[2] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 4)
	out[3] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 5)
	out[4] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 6)
	out[5] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 7)
	out[6] = bitLen - approx.Delta.BitLen()
	approx = tiny.Fuzzy.Approximation(target, 8)
	out[7] = bitLen - approx.Delta.BitLen()
	return out
}
