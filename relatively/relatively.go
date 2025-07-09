package relatively

// Relativity represents the abstract logical relationship of two entities, 𝑎 and 𝑏.
//
// Rather than imbuing 'size', 'value', or 'position', this aims to describe that '𝑎' has
// a logical relationship with '𝑏' that's understood contextually by the caller.  Whether
// in an ordered list, comparing physical dimensions, or relational timing - this provides
// a common language for describing the relationship between both entities.
//
// See Before, Same, After
type Relativity int

const (
	// Before indicates that 𝑎 logically comes before 𝑏.
	Before Relativity = -1
	// Same indicates that 𝑎 and 𝑏 are logically the same.
	Same = 0
	// After indicates that 𝑎 logically comes after 𝑏.
	After Relativity = 1
)
