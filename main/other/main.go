package main

import (
	"fmt"
	"github.com/ignite-laboratories/core"
	"tiny"
)

func main() {
	i := 0
	for _, n := range core.NewFilterableSlice[core.GivenName](core.Names...).Where(func(i int, entry core.GivenName) bool {
		return tiny.NameFilter(entry)
	}) {
		n.Name = ""
		i++
	}
	fmt.Println(i)
}
