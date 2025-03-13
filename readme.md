# Tiny

### A very simplistic helper package for working at the bit level in Go.

----

The goal of Tiny is to allow easy bit level manipulation of odd-length binary information.  Synthesizing random
strings of binary data, for example, is not a common operation - but one that the language is quite good at.

The central pillar of Tiny is the `Measurement` object.  This represents a container of bytes and bits that can
easily be manipulated through the `Tiny` toolkit.

I took the liberty of fleshing out the binary ranges from 1-8 with fitting terms, filling in the gaps with
as apt of terminology as I could imagine:

| Width | Range |  Name  |
|:-----:|:-----:|:------:|
|   1   |  0-1  |  Bit   |
|   2   |  0-3  | Crumb  |
|   3   |  0-7  |  Note  |
|   4   | 0-15  | Nibble |
|   5   | 0-31  | Flake  |
|   6   | 0-63  | Morsel |
|   7   | 0-127 | Shred  |
|   8   | 0-255 |  Byte  |
