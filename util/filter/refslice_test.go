package filter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/stretchr/testify/assert"
)

func TestRefSliceInPlace(t *testing.T) {
	type Foo struct {
		A, B int
	}
	f := filter.Filter[*Foo](func(foo *Foo) bool {
		return foo.A*foo.B%2 == 0
	})

	tt := map[string][]Foo{
		"Simple": {
			{A: 1, B: 3},
			{A: 2, B: 6},
			{A: 3, B: 9},
			{A: 4, B: 12},
			{A: 5, B: 15},
			{A: 6, B: 18},
			{A: 7, B: 21},
			{A: 8, B: 24},
		},
		"All-True": {
			{A: 2, B: 6},
			{A: 4, B: 12},
			{A: 6, B: 18},
			{A: 8, B: 24},
		},
		"All-False": {
			{A: 1, B: 3},
			{A: 3, B: 9},
			{A: 5, B: 15},
			{A: 7, B: 21},
		},
		"Empty": {},
	}

	for n, foos := range tt {
		t.Run(n, func(t *testing.T) {
			ln := len(foos)
			ts, fs := filter.RefSliceInPlace(f, foos)
			assert.Equal(t, ln, len(ts)+len(fs))
			i := 0
			for _, foo := range ts {
				assert.True(t, f(&foo))
				assert.Equal(t, foos[i], foo)
				i++
			}
			for _, foo := range fs {
				assert.False(t, f(&foo))
				assert.Equal(t, foos[i], foo)
				i++
			}
		})
	}
}

func TestRefSliceInPlacePointers(t *testing.T) {
	// The filter receives pointers into vals, not copies.
	vals := []int{3, 1, 4, 1, 5}
	seen := map[*int]bool{}
	f := filter.Filter[*int](func(i *int) bool {
		seen[i] = true
		return *i > 2
	})
	ts, fs := filter.RefSliceInPlace(f, vals)
	assert.ElementsMatch(t, []int{3, 4, 5}, ts)
	assert.ElementsMatch(t, []int{1, 1}, fs)
	addrs := map[*int]bool{}
	for i := range vals {
		addrs[&vals[i]] = true
	}
	assert.NotEmpty(t, seen)
	for p := range seen {
		assert.True(t, addrs[p])
	}
}
