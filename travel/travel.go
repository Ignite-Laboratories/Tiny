package travel

// Travel represents a longitudinal or latitudinal direction.Direction of Travel.
//
// See Westbound, Eastbound, Northbound, Southbound, Outward, and Inward.
type Travel byte

const (
	// Westbound represents a westerly direction of travel.
	//
	// See Eastbound, Northbound, Southbound, Outward, and Inward.
	Westbound Travel = iota

	// Eastbound represents an eastern direction of travel.
	//
	// See Westbound, Northbound, Southbound, Outward, and Inward.
	Eastbound

	// Northbound represents a northward direction of travel.
	//
	// See Westbound, Eastbound, Southbound, Outward, and Inward.
	Northbound

	// Southbound represents a southernly direction of travel.
	//
	// See Westbound, Eastbound, Northbound, Outward, and Inward.
	Southbound

	// Outward represents an outward direction of travel.
	//
	// See Westbound, Eastbound, Northbound, Southbound, and Inward.
	Outward

	// Inward represents an inward direction of travel.
	//
	// See Westbound, Eastbound, Northbound, Southbound, and Outward.
	Inward
)

// String prints a two (or three) character representation of the Travel direction -
//
//	 Westbound: WB
//	Northbound: NB
//	 Eastbound: EB
//	Southbound: SB
//	   Outward: OUT
//	    Inward: IN
func (t Travel) String() string {
	switch t {
	case Westbound:
		return "WB"
	case Northbound:
		return "NB"
	case Eastbound:
		return "EB"
	case Southbound:
		return "SB"
	case Outward:
		return "OUT"
	case Inward:
		return "IN"
	default:
		return "Unknown"
	}
}

// StringFull prints an uppercase full word representation of the Travel direction.
//
// You may optionally pass true for a lowercase representation.
func (t Travel) StringFull(lowercase ...bool) string {
	lower := len(lowercase) > 0 && lowercase[0]
	switch t {
	case Westbound:
		if lower {
			return "westbound"
		}
		return "Westbound"
	case Northbound:
		if lower {
			return "northbound"
		}
		return "Northbound"
	case Eastbound:
		if lower {
			return "eastbound"
		}
		return "Eastbound"
	case Southbound:
		if lower {
			return "southbound"
		}
		return "Southbound"
	case Outward:
		if lower {
			return "outward"
		}
		return "Outward"
	case Inward:
		if lower {
			return "inward"
		}
		return "Inward"
	default:
		if lower {
			return "unknown"
		}
		return "Unknown"
	}
}
