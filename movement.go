package tiny

// Movement represents cyclical Passage information.
type Movement struct {
	Signature Phrase
	Delta     Phrase
	Cycles    int
	BitWidth  int
}

// Perform uses the current movement information to re-build the original information.
func (c Movement) Perform() Phrase {
	//for i := 0; i < c.Cycles; i++ {
	//	m := Passage{
	//		Signature: c.Signature,
	//		Delta:     c.Delta,
	//		DeltaWidth: c.BitWidth,
	//		InitialWidth: c.BitWidth,
	//	}
	//}
	return nil
}
