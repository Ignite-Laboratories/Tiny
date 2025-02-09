package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
)

// TODO: Test this library

func main() {
	val := tiny.From.Int(2345)
	fmt.Println(tiny.To.String(val...))

	test := tiny.To.Byte(1, 0, 1, 1)
	fmt.Println(test)
	fmt.Println(tiny.To.String(tiny.From.Byte(test)...))

	grey := tiny.Create.Grey(4, tiny.From.Bits(0, 1, 1, 0)...)
	fmt.Println(grey)
}
