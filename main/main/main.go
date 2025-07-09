package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
)

func main() {
	b := tiny.NewPhraseFromBits(1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 0)
	a := tiny.NewPhraseFromBits(1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 1, 0)
	c, sign := a.Subtract(b)
	fmt.Println(sign, c.StringBinary())
}
