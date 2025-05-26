package main

import (
	"fmt"
	"github.com/ignite-laboratories/tiny"
	"math/big"
)

func main() {
	source := tiny.Synthesize.RandomPhrase(1024)
	bitLength := source.BitLengthAsBigInt()

	target := source.AsBigInt()
	fmt.Println(target.Text(2))

	upper := new(big.Int).Exp(big.NewInt(2), bitLength, nil)

	divisor := upper.Div(upper, big.NewInt(8))
	multiplier := new(big.Int).Div(target, divisor)
	fmt.Println(multiplier.Text(2))
	difference := new(big.Int).Sub(target, multiplier.Mul(multiplier, divisor))

	fmt.Println(difference.Text(2))
}

func Shrink(source tiny.Phrase, timeline tiny.Phrase) (tiny.Phrase, tiny.Phrase) {
	bitLength := source.BitLengthAsBigInt()

	target := source.AsBigInt()

	upper := new(big.Int)
	upper.Exp(big.NewInt(2), bitLength, nil)

	divisor := upper.Div(upper, big.NewInt(8))
	multiplier := new(big.Int).Div(target, divisor)
	difference := new(big.Int).Sub(target, multiplier.Mul(multiplier, divisor))
	timeline = timeline.AppendBigInt(multiplier)
	fmt.Println(difference.Text(2))
	return nil, nil
}
