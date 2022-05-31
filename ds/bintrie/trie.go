package bintrie

import "github.com/adamcolton/luce/serial/rye"

type Trie interface {
	Insert(uint32)
	Has(uint32) bool
	Delete(uint32)
	All() []*rye.Bits
	Copy() Trie
	Size() int
	InsertTrie(t Trie)
	private()
}

func New() Trie {
	return &node{}
}
