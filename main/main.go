package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
	"math/big"
)

func main() {
	msbSize := 128
	patternSize := 5
	dataLength := 32

	data := tiny.Synthesize.RandomPhrase(dataLength)
	_, data = data.Read(1)     // Drop one bit
	data = data.PrependBits(1) // Add '1' into the first bit position
	// ^^^ This ensures the conversion to big.Int stays the same bit width

	target := data.AsBigInt()
	targetStr := target.Text(2)

	msbs, data := data.Read(msbSize) // Read out the MSBs
	smallest := target
	pattern := big.NewInt(-1)

	// Walk the pattern index's address space
	for i := 0; i <= (1<<patternSize)-1; i++ {
		// Create the initial pattern bits
		bits := tiny.From.Number(i, patternSize)

		// Synthesize the full pattern
		p := tiny.Synthesize.Pattern((dataLength*8)-msbSize, bits...).Prepend(msbs).AsBigInt()

		// Get the delta value
		delta := new(big.Int).Sub(target, p)
		if delta.CmpAbs(smallest) < 0 {
			// Save off the best result
			pattern = p
			smallest = delta
		}
	}

	fmt.Println(pattern.Text(2))
	smallestStr := smallest.Text(2)
	bitDrop := len(targetStr) - len(smallestStr)
	keySize := msbSize + patternSize
	difference := bitDrop - keySize
	fmt.Println(smallestStr)
	fmt.Println(targetStr)
	fmt.Println(difference)
}
