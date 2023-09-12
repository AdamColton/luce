package merkle

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataLeaf(t *testing.T) {
	var maxLeafSize uint32 = 200
	var ln uint32 = 5010
	data := make([]byte, ln)
	rand.Read(data)

	n := NewBuilder(maxLeafSize, sha256.New).Build(data).(*tree).node
	idx := 0
	var fn func(node) int
	fn = func(n node) int {
		if dl, ok := n.(*dataLeaf); ok {
			dld := dl.Data()
			ln := len(dld)
			assert.Equal(t, data[idx:idx+ln], dld)
			idx += ln
			assert.Equal(t, ln, dl.Len())
			return ln
		}
		b := n.(*branch)
		ln := fn(b.children[0])
		ln += fn(b.children[1])
		assert.Equal(t, ln, b.Len())
		return ln
	}
	fn(n)
	assert.Equal(t, int(ln), idx)
}
