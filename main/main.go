package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
)

var maxLength = maxbytes * 8
var maxbytes = 8

func main() {
	data := tiny.Synthesize.RandomPhrase(maxbytes)
	//data := tiny.NewPhraseFromBits(0, 0, 1, 1, 1, 0, 1, 0, 1, 1, 0, 0, 0, 1, 0, 1)
	fmt.Println(data)

	fmt.Println()
	movement := tiny.Synthesize.Movement(data, 3)
	fmt.Println()
	fmt.Println(movement)
	fmt.Println()

	result := movement.Perform()
	fmt.Println()
	fmt.Println(result)
}
