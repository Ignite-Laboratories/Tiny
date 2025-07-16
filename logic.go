package tiny

// _logic represents a factory for logic gate binary functions.
type _logic int

// Logic provides access to logic gate functions.
var Logic _logic

// NOT applies the below truth table against every input Bit to produce a slice of output bits.
//
// NOTE: If no bits are provided, Zero is returned.
//
//	"The NOT Truth Table"
//
//	        𝑎 | 𝑜𝑢𝑡
//	        0 | 1
//	        1 | 0
func (_ _logic) NOT(i uint, operands ...Bit) ([]Bit, *Phrase) {
	if len(operands) == 0 {
		return SingleZero, nil
	}
	for _, b := range operands {
		operands[0] = b ^ 1
	}
	return operands, nil
}

// AND pairwise applies the below truth table against the input bits to produce a single output Bit.
//
// NOTE: If no bits are provided, Zero is returned.  If a single bit is provided, it is returned.
//
//	"The AND Truth Table"
//
//	     𝑎 | 𝑏 | 𝑜𝑢𝑡
//	     0 | 0 | 0
//	     0 | 1 | 0
//	     1 | 0 | 0
//	     1 | 1 | 1
func (_ _logic) AND(i uint, operands ...Bit) ([]Bit, *Phrase) {
	if len(operands) == 0 {
		return SingleZero, nil
	}
	if len(operands) == 1 {
		return operands, nil
	}

	bit := operands[0]
	for _, b := range operands[1:] {
		bit = bit & b
	}
	return []Bit{bit}, nil
}

// OR pairwise applies the below truth table against the input bits to produce a single output Bit.
//
// NOTE: If no bits are provided, Zero is returned.  If a single bit is provided, it is returned.
//
//	"The OR Truth Table"
//
//	     𝑎 | 𝑏 | 𝑜𝑢𝑡
//	     0 | 0 | 0
//	     0 | 1 | 1
//	     1 | 0 | 1
//	     1 | 1 | 1
func (_ _logic) OR(i uint, operands ...Bit) ([]Bit, *Phrase) {
	if len(operands) == 0 {
		return SingleZero, nil
	}
	if len(operands) == 1 {
		return operands, nil
	}

	bit := operands[0]
	for _, b := range operands[1:] {
		bit = bit | b
	}
	return []Bit{bit}, nil
}

// XOR pairwise applies the below truth table against the input bits to produce a single output Bit.
//
// NOTE: If no bits are provided, Zero is returned.  If a single bit is provided, it is returned.
//
//	"The XOR Truth Table"
//
//	     𝑎 | 𝑏 | 𝑜𝑢𝑡
//	     0 | 0 | 0
//	     0 | 1 | 1
//	     1 | 0 | 1
//	     1 | 1 | 0
func (_ _logic) XOR(i uint, operands ...Bit) ([]Bit, *Phrase) {
	if len(operands) == 0 {
		return SingleZero, nil
	}
	if len(operands) == 1 {
		return operands, nil
	}

	bit := operands[0]
	for _, b := range operands[1:] {
		bit = bit ^ b
	}
	return []Bit{bit}, nil
}

// NAND pairwise applies the below truth table against the input bits to produce a single output Bit.
//
// NOTE: If no bits are provided, Zero is returned.  If a single bit is provided, it is returned.
//
//	"The NAND Truth Table"
//
//	     𝑎 | 𝑏 | 𝑜𝑢𝑡
//	     0 | 0 | 1
//	     0 | 1 | 1
//	     1 | 0 | 1
//	     1 | 1 | 0
func (_ _logic) NAND(i uint, operands ...Bit) ([]Bit, *Phrase) {
	if len(operands) == 0 {
		return SingleZero, nil
	}
	if len(operands) == 1 {
		return operands, nil
	}

	bit := operands[0]
	for _, b := range operands[1:] {
		bit = 1 ^ (bit & b)
	}
	return []Bit{bit}, nil
}

// NOR pairwise applies the below truth table against the input bits to produce a single output Bit.
//
// NOTE: If no bits are provided, Zero is returned.  If a single bit is provided, it is returned.
//
//	"The NOR Truth Table"
//
//	     𝑎 | 𝑏 | 𝑐
//	     0 | 0 | 1
//	     0 | 1 | 0
//	     1 | 0 | 0
//	     1 | 1 | 0
func (_ _logic) NOR(i uint, operands ...Bit) ([]Bit, *Phrase) {
	if len(operands) == 0 {
		return SingleZero, nil
	}
	if len(operands) == 1 {
		return operands, nil
	}

	bit := operands[0]
	for _, b := range operands[1:] {
		bit = 1 ^ (bit | b)
	}
	return []Bit{bit}, nil
}

// XNOR pairwise applies the below truth table against the input bits to produce a single output Bit.
//
// NOTE: If no bits are provided, Zero is returned.  If a single bit is provided, it is returned.
//
//	"The XNOR Truth Table"
//
//	     𝑎 | 𝑏 | 𝑐
//	     0 | 0 | 1
//	     0 | 1 | 0
//	     1 | 0 | 0
//	     1 | 1 | 1
func (_ _logic) XNOR(i uint, operands ...Bit) ([]Bit, *Phrase) {
	if len(operands) == 0 {
		return SingleZero, nil
	}
	if len(operands) == 1 {
		return operands, nil
	}

	bit := operands[0]
	for _, b := range operands[1:] {
		bit = 1 ^ (bit ^ b)
	}
	return []Bit{bit}, nil
}
