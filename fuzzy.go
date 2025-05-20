package tiny

type _fuzzy int

// SixtyFour expects up to a crumb key Measurement and returns a bit range up to six bits, yielding 64 unique values.
//
//	00 - 0
//	01 - 2
//	10 - 4
//	11 - 6
func (_ _fuzzy) SixtyFour(key Measurement) int {
	switch key.Value() {
	case 0:
		return 0
	case 1:
		return 2
	case 2:
		return 4
	case 3:
		return 6
	default:
		panic("fuzzy.SixtyFour: key measurement is greater than 3")
	}
}
