package tiny

import (
	"github.com/ignite-laboratories/core"
	"regexp"
)

// NameFilter is a standard function for returning a name which satisfies tiny's requirements for implicit naming.
//
// Names are explicitly filtered to ONLY the standard 26 letters of the English alphabet, with no special or spacing
// characters - meaning only single-word names.  This is specifically to ensure that variables can be captured between
// double-quotes and without any confusion.
func NameFilter(name core.GivenName) bool {
	var nonAlphaRegex = regexp.MustCompile(`^[a-zA-Z]+$`)
	return nonAlphaRegex.MatchString(name.Name)
}
