package tiny

// Expression represents the standard slice index accessor pattern, and expressions can be generated from the global Read variable.
type Expression struct {
	Positions *[]uint
	Low       *uint
	High      *uint
	Last      *bool
	Reverse   *bool
	BitLogic  *BitLogicFunc
	Artifact  *ArtifactFunc
	limit     uint
}

func (e Expression) FromPhrase(p ...Phrase) []Bit {
	return Emit(e, p...)
}

func (e Expression) FromMeasurement(m ...Measurement) []Bit {
	return Emit(e, m...)
}

func (e Expression) FromByte(b ...byte) []Bit {
	return Emit(e, b...)
}

func (e Expression) FromBits(b ...Bit) []Bit {
	return Emit(e, b...)
}

// BitLogicFunc takes in many bits and their collectively shared index and returns an output bit plus a nilable artifact.
type BitLogicFunc func(uint, ...Bit) ([]Bit, *Phrase)

// ArtifactFunc applies the artifact from a single round of calculation against the provided operand bits.
type ArtifactFunc func(i uint, artifact Phrase, operands ...Phrase) []Phrase
