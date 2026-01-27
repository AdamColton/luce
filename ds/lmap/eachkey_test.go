package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestEachKey(t *testing.T) {
	m := lmap.New(map[string]int{"a": 1, "b": 2, "c": 3})
	keys := slice.Slice[string]{"c", "x", "a"}

	// Keys come in the order of the iterator and missing keys are skipped.
	var got []string
	m.EachKey(keys.Iter(), func(k string, v int, done *bool) {
		got = append(got, k)
	})
	assert.Equal(t, []string{"c", "a"}, got)

	// done stops the iteration.
	got = nil
	m.EachKey(keys.Iter(), func(k string, v int, done *bool) {
		got = append(got, k)
		*done = true
	})
	assert.Equal(t, []string{"c"}, got)
}
