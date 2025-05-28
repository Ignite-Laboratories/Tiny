package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
)

func main() {
	test(4294967296, 32)
	test(65536, 16)
	test(256, 8)
	test(128, 7)
	test(64, 6)
	test(32, 5)
	test(16, 4)
	test(8, 3)
	test(4, 2)
	test(2, 1)
}

func examine(subdivisions int, width int, c *tiny.Composition) {
	highCount := 0
	lowCount := 0
	for _, passage := range c.Movements[tiny.MovementPathway] {
		if passage[0].Bits()[0] == 0 {
			lowCount++
		} else {
			highCount++
		}
	}

	fmt.Printf("[%d, %d] %v [Low: %d, High: %d]\n", width, subdivisions, len(c.Movements[tiny.MovementPathway]), lowCount, highCount)
}

func test(subdivisions int, width int) {
	source := tiny.Synthesize.RandomPhrase(1024)
	composition := tiny.Distill(source, subdivisions, width)
	examine(subdivisions, width, composition)
}
