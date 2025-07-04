package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
	"math/big"
)

var maxLength = maxbytes * 8
var maxbytes = 8

func main() {
	maximum := 32

	average := 0
	for i := 0; i < 1<<maximum; i++ {
		midpoint := tiny.Synthesize.Midpoint(maximum)
		delta := new(big.Int).Sub(midpoint.AsBigInt(), big.NewInt(int64(i)))

		fullSize := delta.Text(2)
		if delta.Sign() >= 0 {
			fullSize = "_" + fullSize
		}
		average += maximum - len(fullSize)
	}
	average /= 1 << maximum
	fmt.Println(average)
}

func main2() {
	average := 0
	count := 1 << 32
	for i := 0; i < count; i++ {
		data := tiny.Synthesize.RandomPhrase(maxbytes)
		delta := Midpoint(data, maxLength)
		average += data.BitLength() - delta.BitLength()
	}
	average /= count
	fmt.Println(average)
}

func Midpoint(target tiny.Phrase, width int) tiny.Phrase {
	signature := tiny.NewPhrase()
	delta := target.AsBigInt()

	midpoint := tiny.Synthesize.Midpoint(width)

	delta = new(big.Int).Sub(delta, midpoint.AsBigInt())
	if delta.Sign() < 0 {
		signature = signature.AppendBits(1)
	} else {
		signature = signature.AppendBits(0)
	}
	signature = signature.Align()
	return tiny.NewPhraseFromBigInt(delta)
}
