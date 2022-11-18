package filter_test

import (
	"sort"
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

func TestMapKeyFilter(t *testing.T) {
	f := filter.MapKeyFilter[int, string](func(i int) bool { return i%2 == 0 })
	m := map[int]string{
		1: "1",
		2: "2",
		3: "3",
		4: "4",
		5: "5",
		6: "6",
		7: "7",
		8: "8",
	}

	ks := f.KeySlice(m)
	sort.IntSlice(ks).Sort()
	assert.Equal(t, []int{2, 4, 6, 8}, ks)

	vs := f.ValSlice(m)
	sort.StringSlice(vs).Sort()
	assert.Equal(t, []string{"2", "4", "6", "8"}, vs)

	got := f.Map(m)
	expected := map[int]string{
		2: "2",
		4: "4",
		6: "6",
		8: "8",
	}
	assert.Equal(t, expected, got)
}
