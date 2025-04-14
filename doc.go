// Package tiny provides a toolkit for interfacing with variable-width ranges of binary information.
// It is presented in a fluent style, leveraging being read from left-to-right.
// Thus, to get a tiny type from another uses 'tiny.From' - while going from a tiny
// type to another uses 'tiny.To'.
// The core component of tiny is the Measurement.  This is a standard container for up
// to 32 bits of information.
package tiny

import (
	"github.com/ignite-laboratories/core"
)

var ModuleName = "tiny"

func init() {
	core.ModuleReport(ModuleName)
}

func Report() {}
