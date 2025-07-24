// Package travel provides access to the Travel enumeration.
package travel

// Travel represents a longitudinal or latitudinal direction.Direction of Travel.
//
// These directly relate to the cardinal directions of calculation.
//
// See Westbound, Eastbound, Northbound, Southbound, Outbound, and Inbound.
type Travel byte

const (
	// Westbound represents a westerly direction of travel.
	//
	// See Eastbound, Northbound, Southbound, Outbound, and Inbound.
	Westbound Travel = iota

	// Eastbound represents an easterly direction of travel.
	//
	// See Westbound, Northbound, Southbound, Outbound, and Inbound.
	Eastbound

	// Northbound represents a northerly direction of travel.
	//
	// See Westbound, Eastbound, Southbound, Outbound, and Inbound.
	Northbound

	// Southbound represents a southerly direction of travel.
	//
	// See Westbound, Eastbound, Northbound, Outbound, and Inbound.
	Southbound

	// Outbound represents an outward direction of travel.
	//
	// See Westbound, Eastbound, Northbound, Southbound, and Inbound.
	Outbound

	// Inbound represents an inward direction of travel.
	//
	// See Westbound, Eastbound, Northbound, Southbound, and Outbound.
	Inbound
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
	case Outbound:
		return "OUT"
	case Inbound:
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
	case Outbound:
		if lower {
			return "outward"
		}
		return "Outward"
	case Inbound:
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
