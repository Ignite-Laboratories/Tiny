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
	var bits []tiny.Bit
	var err error

	p := tiny.MeasureMany[byte](77, 22, 44, 88)
	fmt.Printf("\n#0 - Taking a direct binary measurement of several bytes into a phrase named %v:\n", p.Name)
	fmt.Printf("%v ← %v(byte{77}, byte{22}, byte{44}, byte{88})\n\n", p.StringPretty(), p.Name)

	random := rand.Int64()
	width := 17
	p = tiny.Measure[int64](random).AsPhrase(width)
	fmt.Printf("#1 - Measuring a random 64 bit number into a phrase aligned at %d bits-per-measurement named %v:\n", width, p.Name)
	fmt.Printf("%v ← %v(%v)\n\n", p.StringPretty(), p.Name, random)

	fmt.Printf("#2 - Emitting from the end of %v until a condition has been met:\n", p.Name)
	width = 11
	continueFn := func(i uint, data []tiny.Bit) bool {
		if len(data) < width {
			return true
		}
		return false
	}

	bits, err = p.EmitUntil(continueFn, travel.Westbound)
	fmt.Printf("%v ← %v %v while ( len(found) < %d )\n", bits, p.Name, travel.Westbound.StringFull(true), width)

	bits, err = p.EmitUntil(continueFn, travel.Eastbound)
	fmt.Printf("%v ← %v %v while ( len(found) < %d )\n\n", bits, p.Name, travel.Eastbound.StringFull(true), width)

	fmt.Printf("#3 - Emitting specific bits of %v:\n", p.Name)
	bits, _ = p.EmitBetween(11, 44)
	fmt.Printf("%v ← %v[11:44]\n\n", bits, p.Name)

	fmt.Printf("#4 - Gracefully emitting beyond the bounds of %v:\n", p.Name)
	bits, err = p.EmitBetween(55, 88)
	fmt.Printf("%v ← %v[55:88] - Error: %v\n\n", bits, p.Name, err)

	fmt.Printf("#5 -  Emitting the NOT of the last emitted bits from %v:\n", p.Name)
	notBits, _ := emit.NOT(bits...)
	fmt.Printf("%v ← !%v\n\n", notBits, p.Name)

	fmt.Println("#6 - Measuring an object in memory:")

	m := tiny.Measure[direction.Direction](direction.Future)
	fmt.Printf("%v ← Measurement of [%v]\n\n", m.StringPretty(), direction.Future)

	fmt.Println("#7 - Recreating the original object from the measurement:")
	fmt.Printf("%v ← Reconstructed Object\n\n", tiny.ToType[direction.Direction](m))

	fmt.Println("#8 - Pattern emission:")

	pattern := []tiny.Bit{1, 0, 0, 1, 1}
	fmt.Printf("%v ← %v %v\n", tiny.NewMeasurementOfPattern(22, travel.Westbound, pattern...).AsPhrase(-1).StringPretty(), pattern, travel.Westbound.StringFull())
	fmt.Printf("%v ← %v %v\n", tiny.NewMeasurementOfPattern(22, travel.Eastbound, pattern...).AsPhrase(-1).StringPretty(), pattern, travel.Eastbound.StringFull())
	fmt.Printf("%v ← %v %v\n", tiny.NewMeasurementOfPattern(22, travel.Inbound, pattern...).AsPhrase(-1).StringPretty(), pattern, travel.Inbound.StringFull())
	fmt.Printf("%v ← %v %v\n", tiny.NewMeasurementOfPattern(22, travel.Outbound, pattern...).AsPhrase(-1).StringPretty(), pattern, travel.Outbound.StringFull())
	width = 11
	fmt.Printf("%v ← %d repeating `0`s\n", tiny.NewMeasurementOfBit(width, 0).AsPhrase(-1).StringPretty(), width)
	fmt.Printf("%v ← %d repeating `1`s\n\n", tiny.NewMeasurementOfBit(width, 1).AsPhrase(-1).StringPretty(), width)
}
