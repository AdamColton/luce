package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/stretchr/testify/assert"
)

// mockMapper is a minimal Mapper so that Wrapper can be tested on its own.
type mockMapper struct {
	m map[int]string
}

func (mm *mockMapper) Get(k int) (string, bool) { v, ok := mm.m[k]; return v, ok }
func (mm *mockMapper) Len() int                 { return len(mm.m) }
func (mm *mockMapper) Map() map[int]string      { return mm.m }
func (mm *mockMapper) Set(k int, v string)      { mm.m[k] = v }
func (mm *mockMapper) Delete(k int)             { delete(mm.m, k) }
func (mm *mockMapper) New() lmap.Mapper[int, string] {
	return &mockMapper{m: make(map[int]string)}
}
func (mm *mockMapper) Each(fn lmap.EachFunc[int, string]) {
	done := false
	for k, v := range mm.m {
		fn(k, v, &done)
		if done {
			return
		}
	}
}

func TestWrapMapper(t *testing.T) {
	mm := &mockMapper{m: map[int]string{1: "one"}}
	w := lmap.Wrap[int, string](mm)
	assert.Equal(t, mm, w.Mapper)
	assert.Equal(t, mm, w.Wrapped())

	// Wrapping a Wrapper returns it unchanged.
	assert.Equal(t, w, lmap.Wrap[int, string](w))
}
