// Package huffman is used to build variable length codes for compression.
// Frequency data is used to construct a Huffman Tree. A Lookup can be created
// from that Tree. The Lookup can be used to encode arbitrary symbols as bits
// and then the Tree can be used to restore those bits back to the original
// symbols.
//
// For instance, frequency data on letters in English text can be used to build
// an Huffman tree
package huffman
