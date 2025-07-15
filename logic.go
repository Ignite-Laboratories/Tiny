package tiny

// TODO: Tile operations

// _logic represents a factory for logic gate binary functions.
type _logic int

// Logic provides access to logic gate functions.
var Logic _logic

// NOT applies the below truth table against the input Bit to produce an output Bit.
//
//	"The NOT Truth Table"
//
//	        𝑎 | 𝑜𝑢𝑡
//	        0 | 1
//	        1 | 0
func (_ _logic) NOT(i int, b Bit) Bit {
	return b ^ 1
}

// AND pairwise applies the below truth table against the input bits to produce an output Bit.
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
func (_ _logic) AND(i int, bits ...Bit) (Bit, any) {
	if len(bits) == 0 {
		return Zero, nil
	}
	if len(bits) == 1 {
		return bits[0], nil
	}

	bit := bits[0]
	for _, b := range bits[1:] {
		bit = bit & b
	}
	return bit, nil
}

// OR pairwise applies the below truth table against the input bits to produce an output Bit.
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
func (_ _logic) OR(i int, bits ...Bit) (Bit, any) {
	if len(bits) == 0 {
		return Zero, nil
	}
	if len(bits) == 1 {
		return bits[0], nil
	}

	bit := bits[0]
	for _, b := range bits[1:] {
		bit = bit | b
	}
	return bit, nil
}

// XOR pairwise applies the below truth table against the input bits to produce an output Bit.
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
func (_ _logic) XOR(i int, bits ...Bit) (Bit, any) {
	if len(bits) == 0 {
		return Zero, nil
	}
	if len(bits) == 1 {
		return bits[0], nil
	}

	bit := bits[0]
	for _, b := range bits[1:] {
		bit = bit ^ b
	}
	return bit, nil
}

// NAND pairwise applies the below truth table against the input bits to produce an output Bit.
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
func (_ _logic) NAND(i int, bits ...Bit) (Bit, any) {
	if len(bits) == 0 {
		return Zero, nil
	}
	if len(bits) == 1 {
		return bits[0], nil
	}

	bit := bits[0]
	for _, b := range bits[1:] {
		bit = 1 ^ (bit & b)
	}
	return bit, nil
}

// NOR pairwise applies the below truth table against the input bits to produce an output Bit.
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
func (_ _logic) NOR(i int, bits ...Bit) (Bit, any) {
	if len(bits) == 0 {
		return Zero, nil
	}
	if len(bits) == 1 {
		return bits[0], nil
	}

	bit := bits[0]
	for _, b := range bits[1:] {
		bit = 1 ^ (bit | b)
	}
	return bit, nil
}

// XNOR pairwise applies the below truth table against the input bits to produce an output Bit.
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
func (_ _logic) XNOR(i int, bits ...Bit) (Bit, any) {
	if len(bits) == 0 {
		return Zero, nil
	}
	if len(bits) == 1 {
		return bits[0], nil
	}

	bit := bits[0]
	for _, b := range bits[1:] {
		bit = 1 ^ (bit ^ b)
	}
	return bit, nil
}
