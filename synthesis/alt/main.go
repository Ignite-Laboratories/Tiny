package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

//func main() {
//	width := 32
//	dark, _ := strconv.ParseUint(tiny.NewMeasurementOfPattern(width, travel.Eastbound, 0, 0).String(), 2, width)
//	oneZero, _ := strconv.ParseUint(tiny.NewMeasurementOfPattern(width, travel.Eastbound, 0, 1).String(), 2, width)
//	zeroOne, _ := strconv.ParseUint(tiny.NewMeasurementOfPattern(width, travel.Eastbound, 1, 0).String(), 2, width)
//	light, _ := strconv.ParseUint(tiny.NewMeasurementOfPattern(width, travel.Eastbound, 1, 1).String(), 2, width)
//
//	fmt.Println(dark)
//	fmt.Println(oneZero)
//	fmt.Println(zeroOne)
//	fmt.Println(light)
//
//	fmt.Println(strconv.FormatUint(dark, 2))
//	fmt.Println(strconv.FormatUint(oneZero, 2))
//	fmt.Println(strconv.FormatUint(zeroOne, 2))
//	fmt.Println(strconv.FormatUint(light, 2))
//}

var index = uint(64)

func main() {
	cycles := 1 << 10

	avg := 0
	for x := 0; x < cycles; x++ {
		y := int(rand.Uint64())

		for y < 0 {
			y = int(rand.Uint64())
		}

		tester := func(target int) int {
			n := index
			i := 0
			for target > 1 && n <= index {
				i++
				n -= 1
				target = int(math.Abs(float64(1<<n - target))) // | t - 2^(n-1) |
			}
			return i
		}

		result := 0
		for i := 0; i < 1<<16; i++ {
			yi := y + i // 4578816351899506843
			if yi >= (1<<(index-1) - 1) {
				yi = 1 << index
			}

			drop := int(index) - tester(y)
			if drop > result {
				result = drop
			}
		}
		avg += result
	}
	avg /= cycles
	fmt.Println(avg)

	//random := rand.N(1 << index)
	//fmt.Println(random)
	//fmt.Println()
	//

	//fmt.Println(tester(random ^ light))
	//fmt.Println(tester(random ^ zeroOne))
	//fmt.Println(tester(random ^ oneZero))
	//fmt.Println(tester(random ^ dark))

	//fmt.Println()
	//fmt.Println(i)
	//fmt.Println(i - int(threshold))
}
