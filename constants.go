package tiny

// Zero is an implicit Bit{0}.
const Zero Bit = 0

// One is an implicit Bit{1}.
const One Bit = 1

// ZeroZero is an implicit Crumb{00}.
const ZeroZero Crumb = 0

// ZeroOne is an implicit Crumb{01}.
const ZeroOne Crumb = 1

// OneZero is an implicit Crumb{10}.
const OneZero Crumb = 2

// OneOne is an implicit Crumb{11}.
const OneOne Crumb = 3

// MaxCrumb is the maximum value a Crumb can hold.
const MaxCrumb = 3

// MaxNote is the maximum value a Note can hold.
const MaxNote = 7

// MaxNibble is the maximum value a Nibble can hold.
const MaxNibble = 15

// MaxFlake is the maximum value a Flake can hold.
const MaxFlake = 31

// MaxMorsel is the maximum value a Morsel can hold.
const MaxMorsel = 63

// MaxShred is the maximum value a Shred can hold.
const MaxShred = 127

// MaxByte is the maximum value a byte can hold.
const MaxByte = 255

// WidthBit is the number of binary positions a Bit represents
const WidthBit = 1

// WidthCrumb is the number of binary positions a Crumb represents
const WidthCrumb = 2

// WidthNote is the number of binary positions a Note represents
const WidthNote = 3

// WidthNibble is the number of binary positions a Nibble represents
const WidthNibble = 4

// WidthFlake is the number of binary positions a Flake represents
const WidthFlake = 5

// WidthMorsel is the number of binary positions a Morsel represents
const WidthMorsel = 6

// WidthShred is the number of binary positions a Shred represents
const WidthShred = 7

// WidthByte is the number of binary positions a Byte represents
const WidthByte = 8
