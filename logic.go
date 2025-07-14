package tiny

// TODO: Tile operations

// _logic represents a factory for logic gate binary functions.
type _logic int

var Logic _logic

// NOT applies the below truth table against the input Bit to produce an output Bit.
//
//	"The NOT Truth Table"
//
//	        𝑎 | 𝑏
//	        0 | 1
//	        1 | 0
func (_ _logic) NOT(i int, b Bit) Bit {
	return b ^ 1
}
