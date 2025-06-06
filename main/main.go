package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
)

func main() {
	source := tiny.Synthesize.RandomPhrase(1024)
	composition := tiny.Distill(source)
	fmt.Println(len(composition.Movements[tiny.MovementPathway]))
}
