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
	p := tiny.MeasureMany[byte](77, 22, 44, 88)
	fmt.Printf("\n#0 - Taking a direct binary measurement of several bytes into a phrase named %v:\n", p.Name)
	fmt.Printf("%v ← %v(byte{77}, byte{22}, byte{44}, byte{88})\n\n", p.StringPretty(), p.Name)

	random := rand.Int64()
	p = tiny.Measure[int64](random).AsPhrase(17)
	fmt.Printf("#1 - Measuring a random 64 bit number into a phrase aligned at 17 bits-per-measurement named %v:\n", p.Name)
	fmt.Printf("%v ← %v(%v)\n\n", p.StringPretty(), p.Name, random)

	fmt.Printf("#2 - Emitting specific bits of %v:\n", p.Name)
	bits, _ := emit.Between(11, 44, p)
	fmt.Printf("%v ← %v[11:44]\n\n", bits, p.Name)

	fmt.Printf("#3 - Gracefully emitting beyond the bounds of %v:\n", p.Name)
	var err error
	bits, err = emit.Between(55, 88, p)
	fmt.Printf("%v ← %v[55:88] - Error: %v\n\n", bits, p.Name, err)

	fmt.Printf("#4 -  Emitting the NOT of the last emitted bits from %v:\n", p.Name)
	notBits, _ := emit.NOT(bits...)
	fmt.Printf("%v ← !%v\n\n", notBits, p.Name)

	fmt.Println("#5 - Measuring an object in memory:")

	m := tiny.Measure[direction.Direction](direction.Future)
	fmt.Printf("%v ← Measurement of [%v]\n\n", m.StringPretty(), direction.Future)

	fmt.Println("#6 - Recreating the original object from the measurement:")
	fmt.Printf("%v ← Reconstructed Object\n\n", tiny.ToType[direction.Direction](m))

	fmt.Println("#7 - Pattern emission:")

	fmt.Printf("%v ← `1, 0, 0, 1, 1` Westbound\n", tiny.NewMeasurementOfPattern(22, travel.Westbound, 1, 0, 0, 1, 1).AsPhrase(-1).StringPretty())
	fmt.Printf("%v ← `1, 0, 0, 1, 1` Eastbound\n", tiny.NewMeasurementOfPattern(22, travel.Eastbound, 1, 0, 0, 1, 1).AsPhrase(-1).StringPretty())
	fmt.Printf("%v ← `1, 0, 0, 1, 1` Inbound\n", tiny.NewMeasurementOfPattern(22, travel.Inbound, 1, 0, 0, 1, 1).AsPhrase(-1).StringPretty())
	fmt.Printf("%v ← `1, 0, 0, 1, 1` Outbound\n", tiny.NewMeasurementOfPattern(22, travel.Outbound, 1, 0, 0, 1, 1).AsPhrase(-1).StringPretty())
	fmt.Printf("%v ← `0` Repeating\n", tiny.NewMeasurementOfBit(11, 0).AsPhrase(-1).StringPretty())
	fmt.Printf("%v ← `1` Repeating\n", tiny.NewMeasurementOfBit(11, 1).AsPhrase(-1).StringPretty())
}
