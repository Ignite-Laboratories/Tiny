package tiny

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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

// PrintIndexWidth is a convenience method to pretty print the width of an index between pipes.
// The output will always ensure the distance between the two pipes is the provided width of characters,
// with the width printed in the center and arrows pointing outwards as such:
//
//	|1|
//	|←→|
//	|←3→|
//	|← 4→|
//	|← 5 →|
//
// ...
//
//	|←      16      →|
//
// NOTE: If you'd like to simply print the pipes and arrows without the width, provide false to the optional showWidth.
func PrintIndexWidth(width int, showWidth ...bool) string {
	show := true
	if len(showWidth) > 0 {
		show = showWidth[0]
	}
	if width < 0 {
		width = 0
	}

	digits := len(strconv.Itoa(width))
	widthStr := fmt.Sprintf("%*v", digits, width)
	if !show {
		widthStr = strings.Repeat(" ", digits)
	}

	switch {
	case width == 0:
		return "||"
	case width == 1:
		return fmt.Sprintf("|%s|", widthStr)
	case width == 2:
		return "|←→|"
	case width == 3:
		return fmt.Sprintf("|←%s→|", widthStr)
	default:
		totalPadding := width - digits - 2
		left := totalPadding / 2
		right := totalPadding - left

		return fmt.Sprintf("|←%*v%s%*v→|", left, "", widthStr, right, "")
	}
}
