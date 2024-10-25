package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestLen(t *testing.T) {
	m := lmap.New(map[int]string{1: "1", 2: "2", 3: "3"})

	assert.Equal(t, 3, m.Len())

	v, found := m.Get(1)
	assert.Equal(t, "1", v)
	assert.True(t, found)

	v, found = m.Get(4)
	assert.Equal(t, "", v)
	assert.False(t, found)

	m.Set(4, "4")
	v, found = m.Get(4)
	assert.Equal(t, "4", v)
	assert.True(t, found)

	assert.Equal(t, 4, m.Len())

	m.Delete(1)
	v, found = m.Get(1)
	assert.Equal(t, "", v)
	assert.False(t, found)

	got := make(map[int]string)
	m.Each(func(key int, val string) (done bool) {
		got[key] = val
		return false
	})
	assert.Equal(t, m, lmap.New(got))
}
