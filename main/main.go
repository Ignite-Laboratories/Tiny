package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
	"math/big"
)

var maxLength = maxbytes * 8
var maxbytes = 8

func main() {
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
