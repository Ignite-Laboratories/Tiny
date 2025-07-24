package main

import (
	"fmt"
	"math/rand/v2"
	"tiny"
)

func main() {
	m := tiny.Measure[uint64](rand.Uint64()).AsPhrase(11)
	fmt.Println(m.StringPretty())

	m = m.Reverse()
	fmt.Println(m.StringPretty())
}
