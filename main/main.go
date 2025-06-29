package main

import (
	"github.com/ignite-laboratories/tiny"
)

var approxDepth = 3

func main() {
	dataLength := 32
	data := tiny.Synthesize.RandomPhrase(dataLength)
	a := tiny.Synthesize.Approximation(data, approxDepth)

	a.Modulate(tiny.Modulate.Toggle(4, false, 0, 1, 1))
}
