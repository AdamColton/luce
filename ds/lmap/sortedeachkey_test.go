package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestSortedEachKey(t *testing.T) {
	m := lmap.New(map[string]int{"b": 2, "c": 3, "a": 1})
	less := func(i, j string) bool { return i < j }

	var got []string
	keys := m.SortedEachKey(less, nil, func(k string, v int, done *bool) {
		got = append(got, k)
	})
	assert.Equal(t, []string{"a", "b", "c"}, got)
	assert.Equal(t, []string{"a", "b", "c"}, []string(keys))
}
