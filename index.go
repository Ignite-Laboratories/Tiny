package tiny

// Index represents an implicitly fixed-width phrase of raw binary information.
//
// See Natural, Real, Complex, and Binary
type Index struct {
	Phrase
	Width uint
	Flow  int
}

// TODO: Make indexes support a signed "flow" value that can be reset on demand, indicating how many times it over or under-flowed
