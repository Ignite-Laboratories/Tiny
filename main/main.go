package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
)

// TODO: Test this library

func main() {
	for i := 0; i < 10; i++ {
		data := tiny.Synthesize.Random(8)
		bits := data.GetAllBits()
		fmt.Println(bits)
	}
}
