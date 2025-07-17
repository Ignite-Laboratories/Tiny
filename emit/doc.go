// Package emit provides access to bit expression from binary types. This process walks a cursor across the binary information
// and selectively yields bits according to the rules defined by logical expressions. Expressions follow Go slice index accessor
// rules, meaning the low side is inclusive and the high side is exclusive.
//
// NOTE: You must also provide a maximum number of bits to be emitted with the expression - this may be Unlimited.
//
// Positions[𝑛₀,𝑛₁,𝑛₂,𝑛₃] - reads the provided index positions of your binary information in most→to→least significant order - regardless of the provided variadic order.
//
// PositionsFromEnd[𝑛₀,𝑛₁,𝑛₂,𝑛₃] - reads the provided index positions of your binary information in least←to←most significant order - regardless of the provided variadic order.
//
// All[:] - Reads the entirety of your binary information.
//
// Low[low:] - Reads from the provided index to the end of your binary information.
//
// High[:high] - Reads to the provided index from the start of your binary information.
//
// Between[low:high] - Reads between the provided indexes of your binary information.
//
// Gate - Performs a logical operation for every bit of your binary information.
//
// Pattern - XORs the provided pattern against the target bits in most→to→least significant order.
//
// PatternFromEnd - XORs the provided pattern against the target bits in least←to←most significant order.
package emit
