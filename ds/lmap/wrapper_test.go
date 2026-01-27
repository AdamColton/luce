package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestWrapped(t *testing.T) {
	m := lmap.New(map[int]string{1: "1", 2: "2", 3: "3"})
	_, ok := m.Mapper.(lmap.Map[int, string])
	assert.True(t, ok)

	w := lmap.Wrap(m)
	_, ok = w.Mapper.(lmap.Map[int, string])
	assert.True(t, ok)

	assert.Equal(t, w.Mapper, w.Wrapped())
}

func TestWrapperEach(t *testing.T) {
	var m lmap.Wrapper[string, int]
	type kv struct {
		k string
		v int
	}
	var vals slice.Slice[kv]

	fn := func(k string, v int, done *bool) {
		vals = append(vals, kv{k: k, v: v})
	}
	m.Each(fn)
	assert.Len(t, vals, 0)

	m = lmap.New(map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	})
	m.Each(fn)
	vals.Sort(func(i, j kv) bool {
		return i.k < j.k
	})

	expected := slice.Slice[kv]{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	}
	assert.Equal(t, expected, vals)
}
