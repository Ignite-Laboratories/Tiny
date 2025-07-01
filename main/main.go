package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
)

var maxLength = maxbytes * 8
var maxbytes = 64

func main() {
	data := tiny.Synthesize.RandomPhrase(maxbytes)
	c := tiny.Synthesize.Movement(data, 3)
	fmt.Println(c)
}
