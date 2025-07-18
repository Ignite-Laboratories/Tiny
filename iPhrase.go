package tiny

// IPhrase is an intermediate type used to allow the higher-order phrases to be grouped with a phrase in case statements and still access their composed phrase methods.
type IPhrase interface {
	GetData() []Measurement
	GetAllBits() []Bit
	BitWidth() uint
	BleedLastBit() (Bit, any)
	BleedFirstBit() (Bit, any)
	RollUp() any
	Reverse() any
	Append(bits ...Bit) any
	AppendBytes(bytes ...byte) any
	AppendMeasurement(m ...Measurement) any
	Prepend(bits ...Bit) any
	PrependBytes(bytes ...byte) any
	PrependMeasurement(m ...Measurement) any
	Align(width ...uint) any
}
