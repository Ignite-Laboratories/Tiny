package tiny

import (
	"math"
)

var _bitWidth = -1

// GetArchitectureBitWidth queries the maximum uint value at runtime to determine the architecture's bit width.
func GetArchitectureBitWidth() int {
	if _bitWidth > 0 {
		return _bitWidth
	}

	// NOTE: I know there's no reason to check this to 2¹² - but for future proofing, I'm gonna do it anyway.
	for i := 2; i <= MaxScale; i++ {
		if math.MaxUint == (uint(1)<<i)-1 {
			return i
		}
	}
	panic("what be you, beast!?  over 2¹² bits of computation!??")
}

// GetBitWidth returns the number of binary digits necessary to represent the provided number
func GetBitWidth(number int) int {
	if number == 0 {
		return 1
	} else {
		return int(math.Floor(math.Log2(float64(number)))) + 1
	}
}

// GetBase10MaxWidth gets the number of base-10 digits necessary to represent the limit of provided width index.
func GetBase10MaxWidth(width int) int {
	return int(math.Floor(math.Log10(float64(1<<width)))) + 1
}
