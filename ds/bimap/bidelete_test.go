package bimap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/bimap"
	"github.com/stretchr/testify/assert"
)

func TestAddBidelete(t *testing.T) {
	bi := bimap.New[string, int](0)
	assert.Equal(t, bimap.Bidelete[string, int]{}, bi.Add("a", 1))

	// Adding an a that is already paired removes its old b.
	got := bi.Add("a", 2)
	assert.Equal(t, bimap.Bidelete[string, int]{B: bimap.Deleted[int]{Value: 1, Deleted: true}}, got)
	_, found := bi.B(1)
	assert.False(t, found)

	// Adding a b that is already paired removes its old a.
	got = bi.Add("z", 2)
	assert.Equal(t, bimap.Bidelete[string, int]{A: bimap.Deleted[string]{Value: "a", Deleted: true}}, got)
	_, found = bi.A("a")
	assert.False(t, found)

	// Adding a pair that is already there reports the b as deleted.
	got = bi.Add("z", 2)
	assert.True(t, got.B.Deleted)
	assert.Equal(t, 2, got.B.Value)
	v, found := bi.A("z")
	assert.True(t, found)
	assert.Equal(t, 2, v)
}
