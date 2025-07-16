package main

import (
	"fmt"
	"math/rand/v2"
	"tiny"
)

func main() {
	// #0 - Creating a binary measurement directly
	m := tiny.NewMeasurementOfBytes(77, 22)
	fmt.Printf("%v ← Measurement of [byte{77}, byte{22}]\n", m.StringPretty())

	// #1 - Measuring an object directly out of memory into a phrase
	random := rand.Int64()
	p := tiny.Measure("data", random)
	fmt.Printf("%v ← Phrase of [%v]\n", p.Align().StringPretty(), random)

	// #2 - Emitting specific bits from the phrase
	bits := tiny.Emit(tiny.Bits.Between(11, 44), tiny.Unlimited, p)
	fmt.Printf("%v ← Phrase[11:44]\n", bits)

	// #3 -  Emitting the NOT of the emitted bits
	NOTbits := tiny.Emit(tiny.Bits.Gate(tiny.Logic.NOT), tiny.Unlimited, bits...)
	fmt.Printf("%v ← !Phrase[11:44]\n", NOTbits)

	// #4- Converting it back to its original type
	start := tiny.Emit(tiny.Bits.To(11), tiny.Unlimited, p)                         // Get the start range
	NOTbits = tiny.Emit(tiny.Bits.Gate(tiny.Logic.NOT), tiny.Unlimited, NOTbits...) // NOT the NOT bits again
	end := tiny.Emit(tiny.Bits.From(44), tiny.Unlimited, p)                         // Get the end range

	reconstructed := append(start, NOTbits...)
	reconstructed = append(reconstructed, end...)
	p = tiny.NewPhraseFromBits("Reconstructed", tiny.Raw, reconstructed...)
	fmt.Printf("%v ← Reconstructed Phrase\n", p.Align().StringPretty())
}
