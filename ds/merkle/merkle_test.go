package merkle_test

import (
	"crypto/rand"
	"crypto/sha256"
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
