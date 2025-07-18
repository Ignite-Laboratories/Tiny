package tiny

// Natural represents a phrase holding a value belonging to the set of natural numbers, including zero.
//
// To those who think zero shouldn't be included in the set of natural numbers, I present a counter-argument:
// Base 1 has only one identifier, meaning it can only "represent" zero by -not- holding a value in an observable
// location.  Subsequently, all bases are built upon determining the size of a value through "identification" - in
// binary, through a series of zeros or ones, in decimal through the identifiers 0-9.
//
// Now here's where it gets tricky: a value doesn't even EXIST until it is given a place to exist within, meaning its
// existence directly implies a void which has now been filled - an identifiable "zero" state.  In fact, the very first
// identifier of all higher order bases (zero) specifically identifies this state!  Counting, itself, comes from the act of observing
// the general relativistic -presence- of anything - fingers, digits, different length squiggles, feelings - meaning to exclude
// zero attempts to redefine the very fundamental definition of identification itself: it's PERFECTLY reasonable to -naturally-
// count -zero- hairs on a magnificently bald head!
//
//	tl;dr - to count naturally involves identification, which implies accepting -non-existence- as a countable state
//
// I should note this entire system hinges on one fundamental flaw - this container technically holds one additional value beyond
// the 'natural' number set: nil!  I call this the "programmatic set" of numbers, and I can't stop you from setting your natural
// phrase to it, but I can empower you with awareness =)
type Natural tiny.Phrase
