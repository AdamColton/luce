package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	m := lmap.Empty[int, string](10)
	assert.Equal(t, 0, m.Len())
	m = lmap.New[int, string](nil)
	assert.Equal(t, 0, m.Len())
}

func TestEmptySafe(t *testing.T) {
	m := lmap.EmptySafe[int, string](10)
	assert.Equal(t, 0, m.Len())
	m = lmap.NewSafe[int, string](nil)
	assert.Equal(t, 0, m.Len())
}

func testMap(fn func(map[int]string) lmap.Wrapper[int, string], t *testing.T) {
	base := map[int]string{1: "1", 2: "2", 3: "3"}
	m := fn(base)

	assert.Equal(t, 3, m.Len())
	assert.Equal(t, base, m.Map())

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
	m.Each(func(key int, val string, done *bool) {
		got[key] = val
	})
	assert.Equal(t, base, got)

	c := 0
	m.Each(func(key int, val string, done *bool) {
		c++
		*done = true
	})
	assert.Equal(t, 1, c)
}

func TestMap(t *testing.T) {
	testMap(lmap.New[int, string], t)
}

func TestSafe(t *testing.T) {
	testMap(lmap.NewSafe[int, string], t)
}
