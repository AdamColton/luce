package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestFromIterDuplicateKeys(t *testing.T) {
	kvs := slice.Slice[lmap.KeyVal[string, int]]{
		lmap.NewKV("a", 1),
		lmap.NewKV("b", 2),
		lmap.NewKV("a", 3),
	}
	m := lmap.FromIter(kvs.Iter())
	assert.Equal(t, 2, m.Len())
	v, _ := m.Get("a")
	assert.Equal(t, 1, v)
}
