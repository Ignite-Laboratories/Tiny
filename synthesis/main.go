package main

import (
	"fmt"
	"math"
)

//var threshold = uint(10)
//
//func main2() {
//	random := int(rand.N(1 << threshold))
//	fmt.Println(random)
//	fmt.Println()
//
//	n := threshold
//	i := 0
//	for random > 0 {
//		i++
//		n -= 1
//		random = int(math.Abs(float64(1<<n - random))) // | t - 2^(n-1) |
//		fmt.Println(random)
//	}
//
//	fmt.Println()
//	fmt.Println(i)
//	fmt.Println(i - int(threshold))
//}

func main() {
	avg := 0
	threshold := 10
	cycles := 1 << threshold

	maximum := 0
	maxi := 0
	minimum := threshold

	offset := 0
	for x := 0; x < cycles; x++ {
		random := int(x)

		if random&1 == 0 {
			random += offset
			offset += 1
			if offset >= 4 {
				offset = 0
			}

			if random >= cycles {
				random = cycles
			}
		}

		n := threshold
		i := 0
		for random > 1 {
			i++
			n -= 1
			random = int(math.Abs(float64(random - 1<<n))) // | t - 2^(n-1) |
		}
		d := threshold - i
		avg += d

		fmt.Println(d)

		if d > maximum {
			maximum = d
			maxi = i
		}
		if d < minimum {
			minimum = d
		}
	}
	avg /= cycles

	fmt.Println()
	fmt.Println(avg)
	fmt.Println(maximum)
	fmt.Println(maxi)
	fmt.Println(minimum)
}
