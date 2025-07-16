package main

import (
	"fmt"
	"math/rand/v2"
	"tiny"
)

func main() {
	fmt.Println("#0 - Creating a binary measurement directly -")
	m := tiny.NewMeasurementOfBytes(77, 22)
	fmt.Printf("%v ← Measurement of [byte{77}, byte{22}]\n\n", m.StringPretty())

	fmt.Println("#1 - Measuring an object directly out of memory into a phrase -")
	random := rand.Int64()
	p := tiny.Measure("data", random)
	fmt.Printf("%v ← Phrase of [%v]\n\n", p.Align().StringPretty(), random)

	fmt.Println("#2 - Emitting specific bits from the phrase")
	bits := tiny.Emit(tiny.Bits.Between(11, 44), tiny.Unlimited, p)
	fmt.Printf("%v ← Phrase[11:44]\n\n", bits)

	fmt.Println("#3 -  Emitting the NOT of the emitted bits")
	NOTbits := tiny.Emit(tiny.Bits.Gate(tiny.Logic.NOT), tiny.Unlimited, bits...)
	fmt.Printf("%v ← !Phrase[11:44]\n\n", NOTbits)

	fmt.Println("#4- Reconstructing the original phrase")
	start := tiny.Emit(tiny.Bits.To(11), tiny.Unlimited, p)                         // Get the start range
	NOTbits = tiny.Emit(tiny.Bits.Gate(tiny.Logic.NOT), tiny.Unlimited, NOTbits...) // NOT the NOT bits again
	end := tiny.Emit(tiny.Bits.From(44), tiny.Unlimited, p)                         // Get the end range

	reconstructed := append(start, NOTbits...)
	reconstructed = append(reconstructed, end...)
	p = tiny.NewPhraseFromBits("Reconstructed", tiny.Raw, p.Endianness, reconstructed...)
	fmt.Printf("%v ← Reconstructed Phrase\n\n", p.Align().StringPretty())

	// TODO: ToType and handle slices
}
