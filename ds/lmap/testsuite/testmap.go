// Package testsuite provides a shared test for implementations of lmap.Wrapper.
package testsuite

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

// TestMap runs the shared tests for a Wrapper over a map[int]string. The fn
// must return a Wrapper that uses the map it is given directly, not a copy of
// it: the test checks the Wrapper against that map. Call it from a test of the
// implementation, for example testsuite.TestMap(lmap.NewSafe[int, string], t).
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
	assert.Equal(t, base, got)

	c := 0
	m.Each(func(key int, val string, done *bool) {
		c++
		*done = true
	})
	assert.Equal(t, 1, c)

	assert.Equal(t, reflect.TypeOf(m.Mapper), reflect.TypeOf(m.New()))
}
