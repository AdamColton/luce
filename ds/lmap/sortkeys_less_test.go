package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestSortKeysLess(t *testing.T) {
	m := lmap.New(map[string]int{"b": 2, "c": 3, "a": 1})
	got := m.SortKeys(func(i, j string) bool { return i > j }, nil)
	assert.Equal(t, []string{"c", "b", "a"}, []string(got))

	// The buffer is used if it has the capacity.
	buf := make([]string, 0, 3)
	got = m.SortKeys(func(i, j string) bool { return i < j }, buf)
	assert.Equal(t, []string{"a", "b", "c"}, []string(got))
	assert.Equal(t, []string{"a", "b", "c"}, buf[:3])
}
