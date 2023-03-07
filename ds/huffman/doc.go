// Package huffman builds variable length codes for compression. A Huffman Tree
// is constructed from frequency data, so that the more frequent a symbol is, the
// shorter its code. The Tree can restore bits back to the original symbols and
// a Lookup made from it encodes symbols as bits.
//
// For instance, frequency data on letters in English text can be used to build
// a Huffman tree.
package huffman
