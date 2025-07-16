package main

import (
	"fmt"
	"math/rand/v2"
	"tiny"
)

func main() {
	fmt.Println("#0 - Taking a direct binary measurement -")
	m := tiny.NewMeasurementOfBytes(77, 22)
	fmt.Printf("%v ← Measurement of [byte{77}, byte{22}]\n\n", m.StringPretty())

	fmt.Println("#1 - Measuring a random number")
	random := rand.Int64()
	p := tiny.Measure("random", random)
	fmt.Printf("%v ← %v\n\n", p.StringPretty(), random)

	fmt.Println("#2 - Emitting specific bits from the phrase")
	bits := tiny.Emit(tiny.Bits.Between(11, 44), tiny.Unlimited, p)
	fmt.Printf("%v ← Phrase[11:44]\n\n", bits)

	fmt.Println("#3 -  Emitting the NOT of the emitted bits")
	NOTbits := tiny.Emit(tiny.Bits.Gate(tiny.Logic.NOT), tiny.Unlimited, bits...)
	fmt.Printf("%v ← !Phrase[11:44]\n\n", NOTbits)

	fmt.Println("#4 - Measuring an object in memory into a phrase -")
	future := tiny.Future
	p = tiny.Measure("data", future)
	fmt.Printf("%v ← Phrase of [%v]\n\n", p.Align().StringPretty(), future)

	fmt.Println("#5 - Recreating the original object from the phrase -")
	result := tiny.ToType[tiny.Direction](p)
	fmt.Printf("%v ← Reconstructed Object\n", result)
}
