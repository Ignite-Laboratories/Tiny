package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
)

func main() {
	a := tiny.NewPhrase(255)
	b := tiny.NewPhrase(255)
	c := a.Add(b)
	fmt.Println(c)
}
