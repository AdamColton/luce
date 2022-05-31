package bintrie

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUint32(t *testing.T) {
	bt := New()
	u := uint32(3141592)
	bt.Insert(u)
	assert.True(t, bt.Has(u))
	assert.False(t, bt.Has(u+1))

	all := bt.All()
	if assert.Len(t, all, 1) {
		assert.Equal(t, u, uint32(all[0].ReadUint(32)))
	}
}

func TestBools(t *testing.T) {
	var a Trie = &node{}
	ua := uint32(3141592)
	a.Insert(ua)
	var b Trie = &node{}
	ub := uint32(6535897)
	b.Insert(ub)

	both := uint32(793238)
	a.Insert(both)
	b.Insert(both)

	or := Or(a, b)
	assert.Equal(t, 3, or.Size())
	all := or.All()
	assert.Equal(t, ua, uint32(all[0].ReadUint(32)))
	assert.Equal(t, both, uint32(all[1].ReadUint(32)))
	assert.Equal(t, ub, uint32(all[2].ReadUint(32)))

	or = Or(b, a)
	assert.Equal(t, 3, or.Size())
	all = or.All()
	assert.Equal(t, ua, uint32(all[0].ReadUint(32)))
	assert.Equal(t, both, uint32(all[1].ReadUint(32)))
	assert.Equal(t, ub, uint32(all[2].ReadUint(32)))

	and := And(b, a)
	assert.Equal(t, 1, and.Size())
	all = and.All()
	assert.Equal(t, both, uint32(all[0].ReadUint(32)))

	nand := Nand(a, b).All()
	assert.Equal(t, ua, uint32(nand[0].ReadUint(32)))
	nand = Nand(b, a).All()
	assert.Equal(t, ub, uint32(nand[0].ReadUint(32)))
}

func TestDelete(t *testing.T) {
	var a Trie = &node{}
	ua := uint32(3141592)
	a.Insert(ua)
	assert.True(t, a.Has(ua))
	assert.Equal(t, 1, a.Size())
	a.Delete(ua)
	assert.False(t, a.Has(ua))
	assert.Equal(t, 0, a.Size())
	a.Insert(ua)
	assert.True(t, a.Has(ua))
	assert.Equal(t, 1, a.Size())
}

func TestInsertTrie(t *testing.T) {
	us := []uint32{
		31415, 92653, 58979,
	}
	a := New()
	a.Insert(us[0])

	b := New()
	b.Insert(us[1])
	b.Insert(us[2])

	a.InsertTrie(b)

	for _, u := range us {
		assert.True(t, a.Has(u))
	}

	a = New()
	a.Insert(us[0])
	b = New()
	b.Insert(us[1])
	b.Insert(us[2])

	b.InsertTrie(a)
	for _, u := range us {
		assert.True(t, b.Has(u))
	}
}

func TestUnion(t *testing.T) {
	var a Trie = &node{}
	ua := uint32(3141592)
	a.Insert(ua)
	var b Trie = &node{}
	ub := uint32(6535897)
	b.Insert(ub)

	both := uint32(793238)
	a.Insert(both)
	b.Insert(both)

	a.Union(b)
	assert.Equal(t, 1, a.Size())
	all := a.All()
	assert.Equal(t, both, uint32(all[0].ReadUint(32)))

}

func TestMulti(t *testing.T) {
	inAll := make([]uint32, 10)
	for i := range inAll {
		inAll[i] = rand.Uint32()
	}

	tries := make(Multi, 10)
	for i := range tries {
		t := New()
		for _, u := range inAll {
			t.Insert(u)
			t.Insert(uint32(i))
		}
		tries[i] = t
	}

	for _, u := range inAll {
		assert.True(t, tries.Has(u))
	}

	for i := range tries {
		assert.True(t, tries.Has(uint32(i)))
	}
	assert.False(t, tries.Has(uint32(len(tries))))
}
