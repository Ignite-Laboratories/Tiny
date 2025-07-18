package tiny

// Integer represents a phrase encoded as two measurements - a sign bit, and an arbitrary bit-width value.
//
// NOTE: The entire goal of tiny is to break away from the boundaries of overflow logic - if you explicitly
// require working with index-based overflow logic, please use an Index phrase.
type Integer Phrase
