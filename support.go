package tiny

import (
	"fmt"
	"math"
)

// GetArchitectureBitWidth queries the maximum uint value at runtime to determine the architecture's bit width.
//
// NOTE: I don't expect us to jump to other size architectures any time soon, but this errors to ensure we don't
// mask that condition in the future without realizing it.
func GetArchitectureBitWidth() (int, error) {
	switch {
	case math.MaxUint == (1<<64)-1:
		return 64, nil
	case math.MaxUint == (1<<32)-1:
		return 32, nil
	case math.MaxUint == (1<<16)-1:
		return 16, nil
	default:
		return -1, fmt.Errorf("unknown architecture width")
	}
}
