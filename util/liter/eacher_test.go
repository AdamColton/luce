package liter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

func TestEacher(t *testing.T) {
	values := []string{"a", "b", "c"}
	e := liter.Eacher[string]{
		Func: func(fn liter.EachFn[string]) {
			done := false
			for i, v := range values {
				fn(i, v, &done)
				if done {
					return
				}
			}
		},
		L: len(values),
	}
	assert.Equal(t, 3, e.Len())

	var got []string
	e.Each(func(idx int, s string, done *bool) {
		assert.Equal(t, values[idx], s)
		got = append(got, s)
	})
	assert.Equal(t, values, got)

	// Setting done stops the iteration.
	got = nil
	e.Each(func(idx int, s string, done *bool) {
		got = append(got, s)
		*done = idx == 1
	})
	assert.Equal(t, []string{"a", "b"}, got)
}
