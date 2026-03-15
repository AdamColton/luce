package morph_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/ds/morph"
	"github.com/adamcolton/luce/util/liter"
	"github.com/stretchr/testify/assert"
)

// intIter is a minimal liter.Iter over a slice of ints.
type intIter struct {
	vals []int
	idx  int
}

func newIter(vals ...int) *intIter { return &intIter{vals: vals} }

func (i *intIter) Next() (int, bool) { i.idx++; return i.Cur() }
func (i *intIter) Done() bool        { return i.idx >= len(i.vals) }
func (i *intIter) Idx() int          { return i.idx }
func (i *intIter) Cur() (int, bool) {
	if i.Done() {
		return 0, true
	}
	return i.vals[i.idx], false
}

// even includes the even numbers, as strings.
var even = morph.NewVal(func(i int) (string, bool) {
	return "n" + strconv.Itoa(i), i%2 == 0
})

func collect(it liter.Iter[string]) (vals []string, idxs []int) {
	for v, done := it.Cur(); !done; v, done = it.Next() {
		vals = append(vals, v)
		idxs = append(idxs, it.Idx())
	}
	return
}

func TestValIter(t *testing.T) {
	tt := map[string]struct {
		in   []int
		vals []string
	}{
		"mixed":         {[]int{1, 2, 3, 4, 5}, []string{"n2", "n4"}},
		"all-included":  {[]int{2, 4}, []string{"n2", "n4"}},
		"none-included": {[]int{1, 3}, nil},
		"empty":         {nil, nil},
		"one-included":  {[]int{2}, []string{"n2"}},
		"one-excluded":  {[]int{1}, nil},
	}
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			it := even.Iter(newIter(tc.in...))
			vals, idxs := collect(it)
			assert.Equal(t, tc.vals, vals)

			// The index counts the included values, not the underlying values.
			for i := range vals {
				assert.Equal(t, i, idxs[i])
			}

			// Once done, Cur is the zero value.
			v, done := it.Cur()
			assert.Equal(t, "", v)
			assert.True(t, done)
			assert.True(t, it.Done())
		})
	}
}

func TestValIterStartsAtCur(t *testing.T) {
	src := newIter(2, 4, 6)
	src.Next()
	vals, _ := collect(even.Iter(src))
	assert.Equal(t, []string{"n4", "n6"}, vals)
}

func TestValAllIter(t *testing.T) {
	double := morph.NewValAll(func(i int) int { return i * 2 })
	it := double.Iter(newIter(1, 2, 3))
	var got []int
	for v, done := it.Cur(); !done; v, done = it.Next() {
		got = append(got, v)
	}
	assert.Equal(t, []int{2, 4, 6}, got)
}

func TestValFactory(t *testing.T) {
	f := even.Factory(func() (liter.Iter[int], int, bool) {
		it := newIter(1, 2, 3, 4)
		v, done := it.Cur()
		return it, v, done
	})
	it, v, done := f()
	assert.Equal(t, "n2", v)
	assert.False(t, done)
	vals, _ := collect(it)
	assert.Equal(t, []string{"n2", "n4"}, vals)

	// Each call to the Factory starts a new iteration.
	it, _, _ = f()
	vals, _ = collect(it)
	assert.Equal(t, []string{"n2", "n4"}, vals)
}
