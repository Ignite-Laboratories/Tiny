// Package tiny provides a toolkit for interfacing with variable-width ranges of binary information.
//
// "tiny" types are anything sub-byte in width - Bit (1), Crumb (2), Note (3), Nibble (4), Flake (5), Morsel (6), Shred (7)
//
// The package is presented in a fluent style, leveraging being read from left-to-right. Thus, to get a tiny type
// from another uses 'tiny.From' - while going from a tiny type to another uses 'tiny.To'.
//
// The core component of tiny is the Measurement.  This is a standard container for up to 32 bits of information.
//
// The next major component of tiny is the Phrase.  This is a standard way of working with slices of measurements.
package tiny
