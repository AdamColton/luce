package testsuite

import (
	"reflect"
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
	assert.Equal(t, "1", m.GetVal(1))
	assert.True(t, found)

	v, found = m.Get(4)
	assert.Equal(t, "", v)
	assert.Equal(t, "", m.GetVal(4))
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

	assert.Equal(t, reflect.TypeOf(m.Mapper), reflect.TypeOf(m.New()))

	m.Set(5, "5")
	m.Set(6, "6")
	m.Set(7, "7")
	assert.Equal(t, 6, m.Len())
	m.DeleteMany([]int{5, 6, 7})
	assert.Equal(t, 3, m.Len())
}
