# Tiny

### A very simplistic helper package for working at the bit level in Go.

----

The goal of Tiny is to allow easy bit level manipulation of odd-length binary information.  Synthesizing random
strings of binary data, for example, is not a common operation - but one that the language is quite good at.

The project has several recursive layers of structure - starting with the _measurement_

----
**Measurement** - First dimension: point value

A container of up to the host architecture's bit width worth of binary information

Allows a variable binary measurement that's _always_ convertable to a Go `int`

----
**Phrase** - Second dimension: ordered sets of point values

A slice of measurements 

Allows an arbitrarily long binary measurement to be taken

----
**Passage** - Third dimension: unordered sets of ordered sets of point values

A map of named phrases 

Allows clustering of arbitrary measurements

----
**Movement** - Fourth dimension: the timing of how to perform a singular transformation cycle

Holds the logical steps to perform a single cycle of binary transformation

----
**Composition** - Executes transformations

Holds the starting conditions to coordinate recursively performing movements

----
**Growth Scheme**

The _implied_ rules in which a recursive pathway should be walked, such as the standard passage bit width

----
This project entirely operates on the concept of _variable ranges of binary information._
While the information, itself, is a value - each measurement of binary information may or may not contain
multiple _sub-measurements_ which store the actual calculable values.

Because of this, these are _ranges_ of bits - not _values_ - meaning terms like 'int16' or 'int64' are misnomers.
When referring to a 32-bit range of data it wouldn't make sense to say 'read a 32-bit integer _value_.'
Rather, you'd 'read a _cadence_ of binary information.'
That being said, not all terms are misleading - for instance 'read a _byte_ of binary information' is a 
_perfectly_ reasonable statement.

I took the liberty of fleshing out the appropriate binary ranges for binary synthesis with fitting terms, filling 
in the gaps with as apt of terminology as I could imagine:

| Bit Width | Value Upper Limit |  Name   |
|:---------:|:-----------------:|:-------:|
|     1     |         2         |   Bit   |
|     2     |         4         |  Crumb  |
|     3     |         8         |  Note   |
|     4     |        16         | Nibble  |
|     5     |        32         |  Flake  |
|     6     |        64         | Morsel  |
|     7     |        128        |  Shred  |
|     8     |        256        |  Byte   |
|    12     |       4,096       |  Scale  |
|    16     |      65,536       |  Motif  |
|    24     |        2²⁴        |  Riff   |
|    32     |        2³²        | Cadence |
|    48     |        2⁴⁸        |  Hook   |
|    64     |        2⁶⁴        | Melody  |
|    128    |       2¹²⁸        |  Verse  |

There is genuine thought behind these terms, as the dichotomous nature of linguistics are a big part of what 
drove me to begin on this work in the first place.
If we all speak using terms that are inherently imbued with meaning and definition, we collectively reach a 
higher state of collaboration with one another - one where the next engineer doesn't have to _think_ as much 
as the last =)

For instance, 12 bits represents a _Scale_ because a piano has 12 unique keys that define its standard scale.
A _verse_ would contain _melodies_, a _hook_ might include a _cadence_, and all the above would be comprised of _notes_. 

Ultimately, these terms allow us to create richer method names - as `EncodeVerse([]Phrase)` is a lot more explicit 
than `EncodeInt128([]Phrase)`