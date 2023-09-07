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
	var fn func(node)
	fn = func(n node) {
		if dl, ok := n.(*dataLeaf); ok {
			dld := dl.Data()
			ln := len(dld)
			assert.Equal(t, data[idx:idx+ln], dld)
			idx += ln
			return
		}
		b := n.(*branch)
		fn(b.children[0])
		fn(b.children[1])
	}
	fn(n)
	assert.Equal(t, int(ln), idx)
}
