package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
)

func main() {
	dataLength := 32
	data := tiny.Synthesize.RandomPhrase(dataLength)
	a := tiny.Synthesize.Approximation(data, 3)
	fmt.Println(a.Lower)
	fmt.Println(a.Upper)
	fmt.Println(data.StringBinary())
	fmt.Println(a.GetClosestValue().Text(2))
	fmt.Println(a.Delta.Text(2))

	l := len(a.Delta.Text(2))

	for i := 0; i < 64; i++ {
		fmt.Println()

		a.Refine()
		fmt.Println(a.Lower)
		fmt.Println(a.Upper)
		fmt.Println(data.StringBinary())
		fmt.Println(a.GetClosestValue().Text(2))
		fmt.Println(a.Delta.Text(2))
		nl := len(a.Delta.Text(2))
		fmt.Println(l - nl)
		l = nl
	}
}
