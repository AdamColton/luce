// Package rye is a helper for writing serializing and deserializing logic.
//
// There are two related compact formats. One for Uint64 and one for []byte. The
// []byte format is PrefixSlice. If the slice is one byte and that byte is less
// than 129, it is written directly. So all one byte values upto 128 are encoded
// in a single byte. The byte 129 is reseverved for a nil slice. For any slice
// up to 121 bytes, the length is encoded in a single slice as len+129. For
// values longer than 121 bytes, a meta length is encoded as metalength+250.
// Then the data length is encoded with that many bytes little-endian.
//
// Compact Uint64s are encoded similarly. Any value up to and including 128 is
// encoded as a single byte. Values above 128 have their length encoded as
// 121+len and then the value is encoded little-endian. It would be slightly
// more efficient to have length encoding start at 247, allowing for more single
// byte values. But this strategy was chosen for consistency. Converting a
// Uint64 to a little-endian byte slice and encoding that will produce the same
// output as encoding the Uint64 directly.
//
// For signed ints, simply casting them will not work well. Twos compliment
// means that for any negative value the most significant bit will be 1 and so
// all 8 bytes would always be written. Instead, the absolute value is bit
// shifted left one and the sign is placed in the least significant bit.
package rye
