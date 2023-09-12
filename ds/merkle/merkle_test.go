package merkle_test

import (
	"crypto/rand"
	"crypto/sha256"
	mrand "math/rand"
	"testing"

	"github.com/adamcolton/luce/ds/merkle"
	"github.com/stretchr/testify/assert"
)

func TestLeaf(t *testing.T) {
	var maxLeafSize uint32 = 200
	var ln uint32 = 5000
	data := make([]byte, ln)
	rand.Read(data)

	b := merkle.NewBuilder(maxLeafSize, sha256.New)
	m := b.Build(nil)
	assert.Nil(t, m)
	m = b.Build(data)
	assert.Equal(t, data, m.Data())

	start := 0
	h := sha256.New()
	buf := make([]byte, 32)
	leaves := m.Leaves()
	for i := 0; i < leaves; i++ {
		l := m.Leaf(i)
		d := l.Data
		end := start + len(d)
		assert.Equal(t, i, int(l.Index))
		assert.Equal(t, data[start:end], d)
		assert.Equal(t, m.Digest(), l.Digest(h, buf))
		start = end
	}
	assert.Nil(t, m.Leaf(leaves))
}

func TestAssembler(t *testing.T) {
	var maxLeafSize uint32 = 200
	var ln uint32 = 5000
	data := make([]byte, ln)
	rand.Read(data)

	m := merkle.NewBuilder(maxLeafSize, sha256.New).Build(data)
	assert.Equal(t, data, m.Data())

	leaves := make([]*merkle.Leaf, m.Leaves())
	for i := range leaves {
		leaves[i] = m.Leaf(i)
	}
	mrand.New(mrand.NewSource(31415)).Shuffle(len(leaves), func(i, j int) {
		leaves[i], leaves[j] = leaves[j], leaves[i]
	})

	a := m.Description().Assembler(sha256.New())

	l := m.Leaf(0)
	l.Index = uint32(m.Leaves()) + 1
	assert.False(t, a.Add(l))

	for _, l := range leaves {
		done, tr := a.Done()
		assert.False(t, done)
		assert.Nil(t, tr)
		assert.True(t, a.Add(l), l.Index)
		assert.False(t, a.Add(l))
	}
	done, tr := a.Done()
	assert.True(t, done)
	assert.Equal(t, data, tr.Data())
}
