package tiny

// Composition provides all of the metadata necessary to reconstitute a Movement.
type Composition struct {
	Movement
}

// Catalyze takes the raw binary data and reconstitutes the original Movement.
func Catalyze(data ...byte) Movement {
	return nil
}
