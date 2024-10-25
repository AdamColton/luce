package testsuite

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

func TestMap(fn func(map[int]string) lmap.Wrapper[int, string], t *testing.T) {
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
	assert.Equal(t, m, fn(got))

	c := 0
	m.Each(func(key int, val string, done *bool) {
		c++
		*done = true
	})
	assert.Equal(t, 1, c)

	_, ok := m.New().(lmap.Map[int, string])
	assert.True(t, ok)
}
