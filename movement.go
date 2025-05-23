package tiny

// Movement is a slice of phrases which provides the ability to distill and catalyze binary information.
type Movement []Phrase

// NewMovementFromBytes creates a single-phrase movement from a slice of bytes.
func NewMovementFromBytes(data ...byte) Movement {
	return Movement{NewPhrase(data...)}
}

// NewMovement creates a multi-phrase movement from the provided phrases.
func NewMovement(phrases ...Phrase) Movement {
	return phrases
}

// Distill takes in a phrase and reduces it into its most simple form.
//
// The resulting Movement contains two phrases, the signature and the reduction.
func (m Movement) Distill() Movement {
	return nil
}

// Reconstitute takes the distilled binary movement and reconstitutes it.
func (m Movement) Reconstitute() Phrase {
	return nil
}
