package filter_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/math/numiter"
	"github.com/adamcolton/luce/util/filter"
	"github.com/stretchr/testify/assert"
)

func TestMapSliceAndMap(t *testing.T) {
	i := numiter.NewRange(1, 13, 1).Wrap().Iter()
	m := lmap.Iter(i, func(i, idx int) (string, bool) {
		return fmt.Sprintf("%02d", i), true
	})

	tt := map[string]struct {
		filter filter.Map[int, string]
		keys   slice.Slice[int]
	}{
		"int_mod_2": {
			filter: filter.NewMap[int, string](func(i int) bool { return i%2 == 0 }, nil),
			keys:   slice.Slice[int]{2, 4, 6, 8, 10, 12},
		},
		"prefix_1": {
			filter: filter.NewMap[int, string](nil, filter.Prefix("1")),
			keys:   slice.Slice[int]{10, 11, 12},
		},
		"mod_1_and_prefix_0": {
			filter: filter.NewMap(func(i int) bool { return i%2 == 1 }, filter.Prefix("0")),
			keys:   slice.Slice[int]{1, 3, 5, 7, 9},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			ks, vs := tc.filter.Slice(m, nil, nil, filter.ReturnKeys)
			sort.IntSlice(ks).Sort()
			assert.Equal(t, tc.keys, ks)
			assert.Nil(t, vs)

			expectedVals := slice.TransformSlice(tc.keys, func(k, idx int) (string, bool) {
				return m[k], true
			})

			ks, vs = tc.filter.Slice(m, nil, nil, filter.ReturnVals)
			sort.StringSlice(vs).Sort()
			assert.Equal(t, expectedVals, vs)
			assert.Nil(t, ks)

			ks, vs = tc.filter.Slice(m, nil, nil, filter.ReturnBoth)
			for i, k := range ks {
				assert.Equal(t, m[k], vs[i])
			}
			sort.IntSlice(ks).Sort()
			sort.StringSlice(vs).Sort()
			assert.Equal(t, tc.keys, ks)
			assert.Equal(t, expectedVals, vs)

			m2 := tc.filter.Map(m, nil)
			assert.Len(t, m2, len(tc.keys))
			for _, k := range tc.keys {
				v, found := m2[k]
				assert.True(t, found)
				assert.Equal(t, m[k], v)
			}
		})
	}

	f := filter.NewMap(func(i int) bool { return true }, func(s string) bool { return true })
	ks, vs := f.Slice(m, nil, nil, 0)
	assert.Nil(t, ks)
	assert.Nil(t, vs)
}
