package bintrie

import "github.com/adamcolton/luce/serial/rye"

type TrieReader interface {
	Has(uint32) bool
	All() []*rye.Bits
	Copy() Trie
}

type Trie interface {
	TrieReader
	Insert(uint32)
	Delete(uint32)
	Size() int
	InsertTrie(t Trie)
	Union(t Trie)

	private()
}

func New() Trie {
	return &node{}
}
