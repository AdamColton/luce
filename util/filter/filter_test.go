package filter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/stretchr/testify/assert"
)

func TestSliceInPlace(t *testing.T) {
	type Foo struct {
		A, B int
	}
	f := filter.Filter[Foo](func(foo Foo) bool { return foo.A*foo.B%2 == 0 })

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
			ts, fs := f.SliceInPlace(foos)
			assert.Equal(t, len(foos), len(ts)+len(fs))
			i := 0
			for _, foo := range ts {
				assert.True(t, f(foo))
				assert.Equal(t, foos[i], foo)
				i++
			}
			for _, foo := range fs {
				assert.False(t, f(foo))
				assert.Equal(t, foos[i], foo)
				i++
			}
		})
	}
}
