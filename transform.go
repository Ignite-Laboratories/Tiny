package tiny

// Transform
func Transform[T any](expr Expression, data ...T) ([]T, uint) {

}

// TransformBitsFunc takes in a column of bits and their index, then returns an output bit and a flow int to incorporate
// as a carry (overflow) or borrow (underflow) on the forthcoming matrix bits.
type TransformBitsFunc func(int, ...Bit) (out Bit, flow int)

// TransformPhraseFunc takes in a column of phrases and their index, then returns an output phrase and an artifact to
// incorporate into forthcoming matrix phrases.
type TransformPhraseFunc func(int, ...Phrase) (out Phrase, artifact Phrase)
