package main

import (
	"fmt"
	"tiny"
)

func main() {
	m := tiny.NewMeasurementOfBytes(77, 44)
	fmt.Println(m)

	bits := tiny.Emit(tiny.Bits.Between(6, 10), tiny.Unlimited, m)
	fmt.Println(bits)
}
