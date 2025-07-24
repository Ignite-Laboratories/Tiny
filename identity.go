package tiny

import (
	"github.com/ignite-laboratories/core"
	"regexp"
)

var usedNames = make(map[string]*core.GivenName)
var nameRollover = 1 << 14

// NameFilter is a standard function for returning a name which satisfies tiny's requirements for implicit naming.
// Currently, these are our explicit filters -
//
//   - Only the standard 26 letters of the English alphabet (case-insensitive)
//   - No whitespace or special characters (meaning only single word names)
//   - At least three characters in length
//
// These filters will never be reduced - if any changes are made, they will only be augmented.
//
// NOTE: This guarantees up to 2¹⁴ unique names before it begins recycling names.
func NameFilter(name core.GivenName) bool {
	var nonAlphaRegex = regexp.MustCompile(`^[a-zA-Z]+$`)

	if len(usedNames) >= nameRollover {
		usedNames = make(map[string]*core.GivenName)
	}

	if nonAlphaRegex.MatchString(name.Name) && usedNames[name.Name] == nil && len(name.Name) > 2 {
		usedNames[name.Name] = &name
		return true
	}
	return false
}
