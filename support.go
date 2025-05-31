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
