package tiny

import (
	"github.com/ignite-laboratories/core"
	"regexp"
)

var usedNames = make(map[string]*core.GivenName)
var nameRollover = 1 << 14

// NameFilter is a standard function for returning a name which satisfies tiny's requirements for implicit naming.
//
// Names are explicitly filtered to ONLY the standard 26 letters of the English alphabet, with no special or spacing
// characters - meaning only single-word names.  This is specifically to ensure that operands can be represented as
// a single uninterrupted run of identifiable characters - a variable name, if you will.
//
// NOTE: This guarantees up to 2¹⁴ unique names before it begins recycling names.
func NameFilter(name core.GivenName) bool {
	var nonAlphaRegex = regexp.MustCompile(`^[a-zA-Z]+$`)

	if len(usedNames) >= nameRollover {
		usedNames = make(map[string]*core.GivenName)
	}

	if nonAlphaRegex.MatchString(name.Name) && usedNames[name.Name] == nil {
		usedNames[name.Name] = &name
		return true
	}
	return false
}
