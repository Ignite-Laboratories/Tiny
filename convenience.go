package tiny

import (
	"fmt"
	"tiny/travel"
)

// shouldReverseLongitudinally indicates if the direction of travel is westerly and the emission should be reversed, otherwise it panics.
//
// NOTE: This is entirely a convenience function for emission passthrough
func shouldReverseLongitudinally(traveling ...travel.Travel) bool {
	reverse := false
	if len(traveling) > 0 {
		t := traveling[0]
		switch t {
		case travel.Westbound:
			reverse = true
		case travel.Eastbound:
			reverse = false
		case travel.Inbound, travel.Outbound:
			panic(fmt.Sprintf("cannot emit in multiple directions [%v]", t.StringFull()))
		case travel.Northbound, travel.Southbound:
			panic(fmt.Sprintf("cannot emit latitudinally from a linear binary measurement [%v]", t.StringFull()))
		default:
			panic(fmt.Sprintf("unknown direction of travel [%v]", t))
		}
	}
	return reverse
}
