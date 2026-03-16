package lmap_test

import (
	"fmt"
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

// lenIter adds Len to an Iter.
type lenIter struct {
	liter.Iter[string]
	ln int
}

func (l lenIter) Len() int { return l.ln }

func TestEachKey(t *testing.T) {
	m := lmap.New(map[string]int{"a": 1, "b": 2, "c": 3})
	keys := slice.Slice[string]{"c", "x", "a"}

	// Keys come in the order of the iterator and missing keys are skipped.
	var got []string
	e := m.EachKey(keys.Iter())
	e.Each(func(k string, v int, done *bool) {
		got = append(got, fmt.Sprint(k, v))
	})
	assert.Equal(t, []string{"c3", "a1"}, got)
	assert.Equal(t, 0, e.Len())

	// done stops the iteration.
	got = nil
	m.EachKey(keys.Iter()).Each(func(k string, v int, done *bool) {
		got = append(got, k)
		*done = true
	})
	assert.Equal(t, []string{"c"}, got)

	// Len comes from an iterator that fulfills Lener.
	e = m.EachKey(lenIter{Iter: keys.Iter(), ln: 3})
	assert.Equal(t, 3, e.Len())
}
