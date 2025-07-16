package main

import (
	"fmt"
	"tiny"
)

func main() {
	m := tiny.NewMeasurementOfBytes(77, 44)
	fmt.Println(m.StringPretty())

	bits := tiny.Emit(tiny.Bits.Between(2, 10), tiny.Unlimited, m)
	fmt.Println(bits)

	NOTbits := tiny.Emit(tiny.Bits.Gate(tiny.Logic.NOT), tiny.Unlimited, bits...)
	fmt.Println(NOTbits)
}
