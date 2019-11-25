// Package testsuite provides tests to validate that an implementation of
// byteid.Index is correct.
package testsuite

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid"
	"github.com/stretchr/testify/assert"
)

// TestAll runs every test in the suite against the Indexes that factory makes.
func TestAll(t *testing.T, factory byteid.IndexFactory) {
	TestBasicInsertGet(t, factory)
	TestDeleteRecycle(t, factory)
	TestNext(t, factory)
	TestNextEdges(t, factory)
}

// TestBasicInsertGet tests Insert and Get: inserting an ID that is already
// there returns its existing index, and an insert only reports that an append
// is needed once the index is past SliceLen.
func TestBasicInsertGet(t *testing.T, factory byteid.IndexFactory) {
	tr := factory(2)
	assert.Equal(t, 0, tr.Len())

	idx, app := tr.Insert([]byte{1, 2, 3})
	assert.Equal(t, 0, idx)
	assert.False(t, app)
	assert.Equal(t, 1, tr.Len())

	idx, app = tr.Insert([]byte{127, 4, 5})
	assert.Equal(t, 1, idx)
	assert.False(t, app)

	idx, app = tr.Insert([]byte{1, 2, 3})
	assert.Equal(t, 0, idx)
	assert.False(t, app)

	tr.SetSliceLen(3)
	assert.Equal(t, 3, tr.SliceLen())
	idx, app = tr.Insert([]byte{1, 2, 6})
	assert.Equal(t, 2, idx)
	assert.False(t, app)

	idx, app = tr.Insert([]byte{1, 2, 3})
	assert.Equal(t, 0, idx)
	assert.False(t, app)

	idx, app = tr.Insert([]byte{1, 2, 6})
	assert.Equal(t, 2, idx)
	assert.False(t, app)

	idx, found := tr.Get([]byte{10, 2, 3, 4, 5, 6, 7, 8, 9})
	assert.Equal(t, -1, idx)
	assert.False(t, found)

	idx, app = tr.Insert([]byte{10, 2, 3, 4, 5, 6, 7, 8, 9})
	assert.Equal(t, 3, idx)
	assert.True(t, app)

	idx, found = tr.Get([]byte{10, 2, 3, 4, 5, 6, 7, 8, 9})
	assert.Equal(t, 3, idx)
	assert.True(t, found)

	idx, app = tr.Insert([]byte{10, 2, 3, 4, 5, 6, 7, 8, 9})
	assert.Equal(t, 3, idx)
	assert.False(t, app)

	idx, app = tr.Insert([]byte{10, 2, 3, 4, 5, 6, 7, 8, 10})
	assert.Equal(t, 4, idx)
	assert.True(t, app)

	idx, app = tr.Insert([]byte{10, 2, 3, 4, 5, 6, 7, 8, 9})
	assert.Equal(t, 3, idx)
	assert.False(t, app)

	idx, app = tr.Insert([]byte{10, 2, 3, 4, 5, 6})
	assert.Equal(t, 5, idx)
	assert.True(t, app)

	idx, found = tr.Get([]byte{10, 2, 3, 4, 5, 6})
	assert.Equal(t, 5, idx)
	assert.True(t, found)

}

// TestDeleteRecycle tests Delete, and that the index of a deleted ID is reused
// by the next Insert.
func TestDeleteRecycle(t *testing.T, factory byteid.IndexFactory) {
	tr := factory(2)

	idx, app := tr.Insert([]byte{20, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	assert.Equal(t, 0, idx)
	assert.False(t, app)

	idx, app = tr.Insert([]byte{20, 2, 3, 4, 5, 6, 7, 8, 9, 11})
	assert.Equal(t, 1, idx)
	assert.False(t, app)

	idx, found := tr.Get([]byte{20, 2, 3, 4, 5, 6, 7, 8, 9, 11})
	assert.Equal(t, 1, idx)
	assert.True(t, found)

	idx, found = tr.Delete([]byte{20, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	assert.Equal(t, 0, idx)
	assert.True(t, found)

	idx, found = tr.Get([]byte{20, 2, 3, 4, 5, 6, 7, 8, 9, 11})
	assert.Equal(t, 1, idx)
	assert.True(t, found)

	idx, found = tr.Delete([]byte{20, 2, 3, 4, 5, 6, 7, 8, 9, 11})
	assert.Equal(t, 1, idx)
	assert.True(t, found)

	idx, found = tr.Get([]byte{20, 2, 3, 4, 5, 6, 7, 8, 9, 11})
	assert.Equal(t, -1, idx)
	assert.False(t, found)

	idx, found = tr.Delete([]byte{20, 2, 3, 4, 5, 6, 7, 8, 9, 11})
	assert.Equal(t, -1, idx)
	assert.False(t, found)

	idx, app = tr.Insert([]byte{1, 2, 3})
	assert.Equal(t, 1, idx)
	assert.False(t, app)

	idx, app = tr.Insert([]byte{4, 5, 6})
	assert.Equal(t, 0, idx)
	assert.False(t, app)
}

// TestNext tests that Next visits the IDs in order, starting from nil.
func TestNext(t *testing.T, factory byteid.IndexFactory) {
	tr := factory(2)
	tr.Insert([]byte{4, 5, 6})
	tr.Insert([]byte{1, 2, 3})
	tr.Insert([]byte{7, 8, 9})

	expIds := [][]byte{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	expIdx := []int{1, 0, 2}
	i := 0
	for id, idx := tr.Next(nil); id != nil; id, idx = tr.Next(id) {
		assert.Equal(t, expIdx[i], idx)
		assert.Equal(t, expIds[i], id)
		i++
	}
	assert.Equal(t, i, len(expIds))
}

// TestNextEdges tests Next on an empty Index, with an ID that is not in the
// Index and with the last ID.
func TestNextEdges(t *testing.T, factory byteid.IndexFactory) {
	tr := factory(2)
	id, idx := tr.Next(nil)
	assert.Nil(t, id)
	assert.Equal(t, -1, idx)

	tr.Insert([]byte{4, 5, 6})
	tr.Insert([]byte{1, 2, 3})
	tr.Insert([]byte{7, 8, 9})

	// The ID does not have to be in the Index.
	id, idx = tr.Next([]byte{5})
	assert.Equal(t, []byte{7, 8, 9}, id)
	assert.Equal(t, 2, idx)

	// There is nothing after the last ID.
	id, idx = tr.Next([]byte{7, 8, 9})
	assert.Nil(t, id)
	assert.Equal(t, -1, idx)
}
