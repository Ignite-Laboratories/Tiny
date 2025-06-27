package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
)

func main() {
	dataLength := 256
	data := tiny.Synthesize.RandomPhrase(dataLength)
	a := tiny.Synthesize.Approximation(data, 8)
	fmt.Println(a.Target.AsBigInt())
	fmt.Println(a.Value)
	fmt.Println(a.Delta.Text(2))
	fmt.Println(a.Signature)
	fmt.Println(a.CalculateBitDrop())

	m := -1
	mi := -1
	for i := 0; i < 64; i++ {
		a.Refine()
		drop := a.CalculateBitDrop()
		if drop > m {
			mi = i
			m = drop
		}
	}

	fmt.Println()
	fmt.Println(m)
	fmt.Println(mi)
}
