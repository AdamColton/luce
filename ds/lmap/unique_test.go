package lmap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestUnique(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	expected := slice.Slice[int]{3, 1, 4, 5, 9}
	l := slice.LT[int]()
	l.Sort(expected)

	got := lmap.Unique(data, nil)
	l.Sort(got)
	assert.Equal(t, expected, got)

	got = lmap.Unique(data, data)
	l.Sort(got)
	assert.Equal(t, expected, got)
}
