// Package compact extends rye with variable-length encodings that reduce the
// size of serialized data.
package compact

// == projects.Code.luce.compact ==
// [ ] CompactUint64 sort order
//  Encoded values should sort with bytes.Compare in numeric order, for store
//  keys. The value bytes are little-endian, so 511 [131,0xFF,0x01] sorts after
//  512 [131,0x00,0x02]; write them big-endian. Decoding also accepts
//  non-canonical input, [130,5] reads as 5.
// [ ] CompactInt64 sort order
//  The sign in the least significant bit doesn't keep order (-1 is 3, 1 is 2).
//  Sorting signed keys needs an offset-binary mapping.
