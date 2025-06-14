package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
)

func main() {
	x := 0
	y := 4096

	for i := 0; i < y; i++ {
		d := test()
		x += d
	}
	fmt.Println(x / y)
}

func test() int {
	msbSize := 8
	patternSize := 5
	dataLength := 32

	data := tiny.Synthesize.RandomPhrase(dataLength)
	// Ensure the first bit is a 1 so we don't drop any zeros when converting to a big.Int
	_, data = data.Read(1)
	data = data.PrependBits(1)

	target := data.AsBigInt()
	targetStr := target.Text(2)
	msbs, data := data.Read(msbSize)
	smallest := target
	//pattern := big.NewInt(-1)

	for i := 0; i <= (1<<patternSize)-1; i++ {
		bits := tiny.From.Number(i, patternSize)
		p := tiny.Synthesize.Pattern((dataLength*8)-msbSize, bits...).Prepend(msbs)
		bigInt := p.AsBigInt()

		delta := bigInt.Sub(target, bigInt)
		if delta.CmpAbs(smallest) < 0 {
			//pattern = p.AsBigInt()
			smallest = delta
		}
	}

	//fmt.Println(pattern.Text(2))
	smallestStr := smallest.Text(2)
	bitDrop := len(targetStr) - len(smallestStr)
	keySize := msbSize + patternSize
	difference := bitDrop - keySize
	//fmt.Println(smallestStr)
	//fmt.Println(targetStr)
	//fmt.Println(difference)
	return difference
}
