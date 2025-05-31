package tiny

import (
	"fmt"
	"github.com/ALTree/bigfloat"
	"math/big"
)

// Composition represents a working structure for synthesizing binary information.
//
// A composition consists of Movements which can be referenced and updated while distilling binary information.
//
// Passage - A passage is a collection of related phrases.
// This allows for looping phrases to be clustered together.
//
// Phrase - A phrase is a way of grouping together longer runs of bits.
// Since a measurement is limited in width, this allows any arbitrary length of bits to be grouped together.
//
// Measurement - A measurement is a variable width slice of bits up to MaxMeasurementBitLength.
type Composition struct {
	Movements   map[string][]Passage
	Target      *big.Int
	TargetFloat *big.Float
	Fuzzy       *big.Int

	pathway []Bit
}

func Distill(phrase Phrase) *Composition {
	c := &Composition{
		Movements: make(map[string][]Passage),
	}
	c.Movements[MovementStart] = make([]Passage, 0)
	c.Movements[MovementPathway] = make([]Passage, 0)
	c.Movements[MovementSeed] = make([]Passage, 0)
	c.pathway = make([]Bit, 0)
	c.Target = phrase.AsBigInt()
	c.TargetFloat = new(big.Float).SetInt(c.Target)
	c.synthesizeFuzzy()
	c.distill()
	return c
}

// AddPassageToStart appends the provided passage to the end of the start movement.
func (c *Composition) AddPassageToStart(passage Passage) {
	c.Movements[MovementStart] = append(c.Movements[MovementStart], passage)
}

// AddPassageToPathway appends the provided passage to the end of the pathway movement.
func (c *Composition) AddPassageToPathway(passage Passage) {
	c.Movements[MovementPathway] = append(c.Movements[MovementPathway], passage)
}

// AddPassageToSeed appends the provided passage to the end of the seed movement.
func (c *Composition) AddPassageToSeed(passage Passage) {
	c.Movements[MovementSeed] = append(c.Movements[MovementSeed], passage)
}

func (c *Composition) synthesizeFuzzy() {
	// TODO: Create a better synthetic fuzzy approximation
	c.Fuzzy = new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(c.Target.BitLen()-1)), nil)
	c.AddPassageToPathway(NewPassage())
}

func (c *Composition) distill(depth ...int) {
	//var d int
	//if len(depth) > 0 {
	//	d = depth[0]
	//}
	//
	//if c.Fuzzy.Cmp(c.Target) < 0 {
	//	// Multiply by > 1
	//	c.pathway = append(c.pathway, One)
	//	str := "1."
	//	char = "0"
	//} else {
	//	// Multiply by < 1
	//	c.pathway = append(c.pathway, Zero)
	//	str := "0."
	//	char = "9"
	//}
	//
	//// Threshold condition
	//if c.Target.BitLen()-c.Fuzzy.BitLen() > len(c.pathway) {
	//	return
	//}
}

func (c *Composition) distill2(target *big.Int) {
	targetFloat := new(big.Float).SetInt(target)

	//ptrn := Synthesize.Pattern(target.BitLen(), 1, 0)
	//mid, _ := new(big.Int).SetString(ptrn.StringBinary(), 2)
	//midFloat := new(big.Float).SetInt(mid)

	lower := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(target.BitLen()-1)), nil)
	lowerFloat := new(big.Float).SetInt(lower)

	correction := new(big.Float).Quo(lowerFloat, targetFloat)
	fmt.Println(target.Text(10))
	fmt.Println(correction)
	correction32, _ := correction.Float32()
	fmt.Println(correction32)
	correction = correction.SetFloat64(float64(correction32))
	corrected, _ := new(big.Float).Mul(targetFloat, correction).Int(nil)

	fmt.Println(target.Text(2))
	fmt.Println(corrected.Text(2))

	//c.distill(corrected)
}

func (c *Composition) distill1(target *big.Int) {
	targetFloat := new(big.Float).SetInt(target)

	upper := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(target.BitLen())), nil)
	upperFloat := new(big.Float).SetInt(upper)
	lower := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(target.BitLen()-1)), nil)
	lowerFloat := new(big.Float).SetInt(lower)

	upperDelta := new(big.Int).Sub(upper, target)
	lowerDelta := new(big.Int).Sub(target, lower)
	boundary := new(big.Float)

	if upperDelta.Cmp(lowerDelta) > 0 {
		fmt.Println("Going high")
		boundary = upperFloat
	} else {
		fmt.Println("Going low")
		boundary = lowerFloat
	}
	correction := new(big.Float).Quo(boundary, targetFloat)
	fmt.Println(correction)
	correction32, _ := correction.Float32()
	correction = correction.SetFloat64(float64(correction32))
	corrected, _ := new(big.Float).Mul(targetFloat, correction).Int(nil)

	fmt.Println(target.Text(2))
	fmt.Println(corrected.Text(2))

	//c.distill(corrected)
}

type approximation struct {
	Base       int
	Exponent   int
	Target     *big.Int
	Fuzzy      *big.Int
	Remainder  *big.Int
	Overshoot  Bit
	Correction float32
}

func (approx *approximation) Correct() {
	fuzzyFloat := new(big.Float).SetInt(approx.Fuzzy)
	targetFloat := new(big.Float).SetInt(approx.Target)

	correction := new(big.Float).Quo(targetFloat, fuzzyFloat)
	correction32, _ := correction.Float64()
	correction = correction.SetFloat64(correction32)
	corrected, _ := new(big.Float).Mul(fuzzyFloat, correction).Int(nil)
	approx.Fuzzy = corrected
	approx.Correction, _ = correction.Float32()

	if approx.Fuzzy.Cmp(approx.Target) > 0 {
		approx.Overshoot = One
		approx.Remainder = new(big.Int).Sub(approx.Fuzzy, approx.Target)
	} else {
		approx.Overshoot = Zero
		approx.Remainder = new(big.Int).Sub(approx.Target, approx.Fuzzy)
	}
}

/**
Utility Functions
*/

// roundBigFloat rounds the input big.Float to the nearest integer value and indicates whether it went up or down.
func roundBigFloat(x *big.Float) (*big.Int, bool) {
	withoutHalf, _ := x.Int(nil)
	withHalf, _ := new(big.Float).Add(x, new(big.Float).SetFloat64(0.5)).Int(nil)

	if withHalf.Cmp(withoutHalf) > 0 {
		return withHalf, true
	}
	return withoutHalf, false
}

// logBaseN performs a log_base(x) operation using big floats.
func logBaseN(x *big.Int, base int) *big.Float {
	// log_n(x) = ln(x) / ln(n)
	xFloat := new(big.Float).SetInt(x)
	baseFloat := new(big.Float).SetInt(big.NewInt(int64(base)))

	lnX := bigfloat.Log(xFloat)
	lnBase := bigfloat.Log(baseFloat)
	return new(big.Float).Quo(lnX, lnBase)
}

func findSmallest(input ...approximation) approximation {
	if len(input) == 0 {
		panic("need at least one value")
	}

	smallest := input[0]

	for _, v := range input[1:] {
		if v.Remainder.Cmp(smallest.Remainder) < 0 {
			smallest = v
		}
	}

	return smallest
}
