package main

import (
	_ "embed"
	"fmt"
	"github.com/ignite-laboratories/tiny"
	"math/big"
)

//go:embed messages.pdf
var bytes string

var maxLength = maxbytes * 8
var maxbytes = 32

func main() {
	cycles := 1 << 12
	var average int
	for i := 0; i < cycles; i++ {
		average += test()
	}
	average /= cycles

	fmt.Println(average)
}

func test() int {
	//data := tiny.NewPhraseFromBytesAndBits([]byte(bytes)[:maxbytes])

	data := tiny.Synthesize.RandomPhrase(maxbytes)
	//dataStr := data.StringBinary()
	//fmt.Println(dataStr)

	delta := data.AsBigInt()
	signature := tiny.NewPhrase()

	thresholdPoint := -1
	//hitThreshold := false

	for i := maxLength; i >= 0; i-- {
		mid := tiny.Synthesize.Midpoint(i)

		delta = new(big.Int).Sub(mid.AsBigInt(), delta)
		if delta.Sign() < 0 {
			signature = signature.AppendBits(1)
		} else {
			signature = signature.AppendBits(0)
		}
		signature = signature.Align(tiny.GetArchitectureBitWidth())
		delta = new(big.Int).Abs(delta)
		//deltaStr := delta.Text(2)
		//fmt.Println(deltaStr)

		//if delta.Cmp(big.NewInt(1)) <= 0 && !hitThreshold {
		if i == 3 {
			thresholdPoint = int(delta.Int64())
			//hitThreshold = true
		}

		//diff := len(dataStr) - len(deltaStr) - signature.BitLength()
	}
	return thresholdPoint
}
