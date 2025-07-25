package tiny

// Analyze is a way to glean information about existing binary information.
var Analyze _analyze

// From is a way to analyze binary slices from existing data.
//
// It's read left-to-right:
//
//	tiny.From [byte/number/etc...]
var From _from

// Fuzzy provides access to fuzzy ZLE operations.
var Fuzzy _fuzzy

// Modify is a way to alter existing binary information.
var Modify _modify

// Synthesize is a way to create binary slices from known parameters.
var Synthesize _synthesize

// To is a way to convert binary slices to other forms.
//
// It's read left-to-right:
//
//	tiny.To [byte/number/etc...]
var To _to

// Print contains miscellaneous functions for printing out textual information.
var Print _print
