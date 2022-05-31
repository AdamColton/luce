package bintrie

import "github.com/adamcolton/luce/serial/rye"

type Uint interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type TrieReader[U Uint] interface {
	Has(U) bool
	All() []*rye.Bits
	Copy() Trie[U]
}

type Trie[U Uint] interface {
	TrieReader[U]
	Insert(U)
	Delete(U)
	Size() int
	InsertTrie(t Trie[U])
	Union(t Trie[U])

	private()
}

func New[U Uint]() Trie[U] {
	return &node[U]{}
}
