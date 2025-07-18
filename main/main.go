package main

import (
	"fmt"
	"math/rand/v2"
	"tiny"
	"tiny/direction"
	"tiny/emit"
	"tiny/travel"
)

func main() {
	fmt.Println("#0 - Taking a direct binary measurement:")

	m := tiny.Measure[byte]("bytes", 77, 22, 44, 88)
	fmt.Printf("%v ← Phrase of [byte{77}, byte{22}, byte{44}, byte{88}]\n\n", m.StringPretty())

	fmt.Println("#1 - Measuring a random number:")

	random := rand.Int64()
	logical := tiny.Measure[int64]("random", random).AsLogical()
	fmt.Printf("%v ← %v\n\n", logical.StringPretty(), random)

	fmt.Println("#2 - Emitting specific bits from the logical phrase:")

	bits := emit.Between(11, 44, logical)
	fmt.Printf("%v ← Phrase[11:44]\n\n", bits)

	fmt.Println("#3 -  Emitting the NOT of the emitted bits:")

	NOTbits := emit.NOT(bits...)
	fmt.Printf("%v ← !Phrase[11:44]\n\n", NOTbits)

	fmt.Println("#4 - Measuring an object in memory into a phrase:")

	p := tiny.Measure[direction.Direction]("forward progress", direction.Future)
	fmt.Printf("%v ← Phrase of [%v]\n\n", p.Align().StringPretty(), direction.Future)

	fmt.Println("#5 - Recreating the original object from the phrase:")

	result := tiny.ToType[direction.Direction](p)
	fmt.Printf("%v ← Reconstructed Object\n\n", result)

	fmt.Println("#6 - Pattern emission:")

	fmt.Printf("%v ← Westbound `1, 0, 0, 1, 0`\n", tiny.NewMeasurementOfPattern(22, travel.Westbound, 1, 0, 0, 1, 0))
	fmt.Printf("%v ← Eastbound `1, 0, 0, 1, 0`\n", tiny.NewMeasurementOfPattern(22, travel.Eastbound, 1, 0, 0, 1, 0))
	fmt.Printf("%v ← Inward `1, 0, 0, 1, 0`\n", tiny.NewMeasurementOfPattern(22, travel.Inward, 1, 0, 0, 1, 0))
	fmt.Printf("%v ← Outward `1, 0, 0, 1, 0`\n", tiny.NewMeasurementOfPattern(22, travel.Outward, 1, 0, 0, 1, 0))
}
