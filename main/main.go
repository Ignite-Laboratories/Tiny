package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
)

// TODO: Test this library

func main() {
	PrintUpToValue[tiny.Note](tiny.MaxCrumb)
	PrintUpToValue[tiny.Flake](tiny.MaxCrumb)
	PrintUpToValue[tiny.Shred](tiny.MaxCrumb)
}

func PrintUpToValue[T tiny.SubByte](maxValue int) {
	for i := 0; i < maxValue; i++ {
		fmt.Println(T(i).String())
	}
}
