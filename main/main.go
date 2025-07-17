package main

import (
	"fmt"
	"math/rand/v2"
	"tiny"
	"tiny/emit"
)

func main() {
	fmt.Println("#0 - Taking a direct binary measurement -")

	m := tiny.Measure[byte]("bytes", 77, 22)
	fmt.Printf("%v ← Measurement of [byte{77}, byte{22}]\n\n", m.StringPretty())

	fmt.Println("#1 - Measuring a random number")

	random := rand.Int64()
	p := tiny.Measure[int64]("random", random)
	fmt.Printf("%v ← %v\n\n", p.StringPretty(), random)

	fmt.Println("#2 - Emitting specific bits from the phrase")

	bits := emit.Between(11, 44).FromPhrase(p)
	fmt.Printf("%v ← Phrase[11:44]\n\n", bits)

	fmt.Println("#3 -  Emitting the NOT of the emitted bits")

	NOTbits := emit.NOT().FromBits(bits...)
	fmt.Printf("%v ← !Phrase[11:44]\n\n", NOTbits)

	fmt.Println("#4 - Measuring an object in memory into a phrase -")

	p = tiny.Measure[tiny.Direction]("forward progress", tiny.Future)
	fmt.Printf("%v ← Phrase of [%v]\n\n", p.Align().StringPretty(), tiny.Future)

	fmt.Println("#5 - Recreating the original object from the phrase -")

	result := tiny.ToType[tiny.Direction](p)
	fmt.Printf("%v ← Reconstructed Object\n", result)
}
