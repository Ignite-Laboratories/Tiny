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
	fmt.Printf("%v ← Phrase of [byte{77}, byte{22}, byte{44}, byte{88}]\n\n", tiny.Measure[byte]("bytes", 77, 22, 44, 88).StringPretty())

	fmt.Println("#1 - Measuring a random number:")
	random := rand.Int64()
	logical := tiny.Measure[int64]("random", random).AsLogical()
	fmt.Printf("%v ← %v\n\n", logical.StringPretty(), random)

	fmt.Println("#2 - Emitting specific bits from the logical phrase:")
	bits := emit.Between(11, 44, logical)
	fmt.Printf("%v ← Phrase[11:44]\n\n", bits)

	fmt.Println("#3 -  Emitting the NOT of the emitted bits:")
	fmt.Printf("%v ← !Phrase[11:44]\n\n", emit.NOT(bits...))

	fmt.Println("#4 - Measuring an object in memory into a phrase:")

	p := tiny.Measure[direction.Direction]("forward progress", direction.Future)
	fmt.Printf("%v ← Phrase of [%v]\n\n", p.Align().StringPretty(), direction.Future)

	fmt.Println("#5 - Recreating the original object from the phrase:")
	fmt.Printf("%v ← Reconstructed Object\n\n", tiny.ToType[direction.Direction](p))

	fmt.Println("#6 - Pattern emission:")

	fmt.Printf("%v ← Westbound `1, 0, 0, 1, 1`\n", tiny.NewMeasurementOfPattern(22, travel.Westbound, 1, 0, 0, 1, 1).StringPretty())
	fmt.Printf("%v ← Eastbound `1, 0, 0, 1, 1`\n", tiny.NewMeasurementOfPattern(22, travel.Eastbound, 1, 0, 0, 1, 1).StringPretty())
	fmt.Printf("%v ← Inward `1, 0, 0, 1, 1`\n", tiny.NewMeasurementOfPattern(22, travel.Inbound, 1, 0, 0, 1, 1).StringPretty())
	fmt.Printf("%v ← Outward `1, 0, 0, 1, 1`\n", tiny.NewMeasurementOfPattern(22, travel.Outbound, 1, 0, 0, 1, 1).StringPretty())
	fmt.Printf("%v ← Digit `0`\n", tiny.NewMeasurementOfBit(11, 0).StringPretty())
	fmt.Printf("%v ← Digit `1`\n", tiny.NewMeasurementOfBit(11, 1).StringPretty())

	pempty := tiny.NewPhrase("test").Append()
	pempty.Data = append(pempty.Data, tiny.NewMeasurement())

	pempty.Data[0] = pempty.Data[0].Append(1)
	pempty.Data[1] = pempty.Data[1].Append(0)
	fmt.Println(pempty.RollUp().StringPretty())
}
