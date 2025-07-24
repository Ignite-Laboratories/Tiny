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
	fmt.Println("\n#0 - Taking a direct binary measurement:")
	fmt.Printf("%v ← Phrase of [byte{77}, byte{22}, byte{44}, byte{88}]\n\n", tiny.MeasureMany[byte](77, 22, 44, 88).StringPretty())

	fmt.Println("#1 - Measuring a random number:")
	random := rand.Int64()
	p := tiny.Measure[int64](random).AsPhrase(31)
	fmt.Printf("%v ← %v\n\n", p.StringPretty(), random)

	fmt.Println("#2 - Emitting specific bits from a logical phrase:")
	bits, _ := emit.Between(11, 44, p)
	fmt.Printf("%v ← Phrase[11:44]\n\n", bits)

	fmt.Println("#3 - Gracefully emitting beyond the phrase bounds:")
	bits2, err := emit.Between(55, 88, p)
	fmt.Printf("%v ← Phrase[55:88] - Error: %v\n\n", bits2, err)

	fmt.Println("#4 -  Emitting the NOT of the emitted bits:")
	notBits, _ := emit.NOT(bits...)
	fmt.Printf("%v ← !Phrase[11:44]\n\n", notBits)

	fmt.Println("#5 - Measuring an object in memory:")

	m := tiny.Measure[direction.Direction](direction.Future)
	fmt.Printf("%v ← Measurement of [%v]\n\n", m.StringPretty(), direction.Future)

	fmt.Println("#6 - Recreating the original object from the phrase:")
	fmt.Printf("%v ← Reconstructed Object\n\n", tiny.ToType[direction.Direction](m))

	fmt.Println("#7 - Pattern emission:")

	fmt.Printf("%v ← Westbound `1, 0, 0, 1, 1`\n", tiny.NewMeasurementOfPattern(22, travel.Westbound, 1, 0, 0, 1, 1).AsPhrase(-1).StringPretty())
	fmt.Printf("%v ← Eastbound `1, 0, 0, 1, 1`\n", tiny.NewMeasurementOfPattern(22, travel.Eastbound, 1, 0, 0, 1, 1).AsPhrase(-1).StringPretty())
	fmt.Printf("%v ← Inward `1, 0, 0, 1, 1`\n", tiny.NewMeasurementOfPattern(22, travel.Inbound, 1, 0, 0, 1, 1).AsPhrase(-1).StringPretty())
	fmt.Printf("%v ← Outward `1, 0, 0, 1, 1`\n", tiny.NewMeasurementOfPattern(22, travel.Outbound, 1, 0, 0, 1, 1).AsPhrase(-1).StringPretty())
	fmt.Printf("%v ← Digit `0`\n", tiny.NewMeasurementOfBit(11, 0).AsPhrase(-1).StringPretty())
	fmt.Printf("%v ← Digit `1`\n", tiny.NewMeasurementOfBit(11, 1).AsPhrase(-1).StringPretty())
}
