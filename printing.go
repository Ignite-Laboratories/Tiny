package tiny

import (
	"fmt"
	"strconv"
	"strings"
)

type _print int

// IndexWidth is a convenience method to pretty print the width of an index between pipes.
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
func (_ _print) IndexWidth(width int, showWidth ...bool) string {
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
		return "| |"
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

// DeltaCharacter is used to build waveform visualizations in a vertically emitted text plot.
// It calculates if the values have changed and then prints the appropriate pipe character to visualize it.
//
//	if old == new -> '|'
//	if old < new -> '\'
//	if old > new -> '/'
func (_ _print) DeltaCharacter(old int, new int) string {
	switch {
	case old < new:
		return "\\"
	case old > new:
		return "/"
	default:
		return "|"
	}
}

// BitDrop prints a waveform showing the relative bit drop from the provided index width
//
// For instance, let's walk a note index and show the outputs for each point from this function:
//
//	|←3→|
//	||  | (
//
//	| 11|  (7)  ||  |
//	| 10|  (6)  ||  |
//	|  1|  (5)  |-| |
//	|  0|  (4)  |-| |
//	|  1|  (3)  |-| |
//	| 10|  (2)  ||  |
//	| 11|  (1)  ||  |
//	|100|  (0)  |   |
func (_ _print) BitDrop(point Phrase, index int) string {
	return ""
}
