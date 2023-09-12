package merkle

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	var maxLeafSize uint32 = 20
	var ln uint32 = 255
	data := make([]byte, ln)
	for i := range data {
		data[i] = byte(i)
	}

	m := NewBuilder(maxLeafSize, sha256.New).Build(data)
	assert.Equal(t, data, m.Data())

	var recurse func(node, int) (int, int)
	leaves := 0
	recurse = func(n node, start int) (int, int) {
		if dl, ok := n.(*dataLeaf); ok {
			d := dl.Data()
			assert.LessOrEqual(t, len(d), int(ln))
			end := start + len(d)
			assert.Equal(t, data[start:end], dl.Data())
			leaves++
			return start, end
		}
		b := n.(*branch)
		start, mid := recurse(b.children[0], start)
		_, end := recurse(b.children[1], mid)
		assert.Equal(t, data[start:end], b.Data())
		return start, end
	}
	recurse(m.(*tree).node, 0)
	assert.Equal(t, leaves, m.Leaves())

	assert.Nil(t, NewBuilder(maxLeafSize, sha256.New).Build(nil))
}
