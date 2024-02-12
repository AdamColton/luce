package slice_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/liter"
	"github.com/adamcolton/luce/util/upgrade"
	"github.com/stretchr/testify/assert"
)

func TestClone(t *testing.T) {
	data := slice.New([]int{3, 1, 4, 1, 5, 9})
	cp := data.Clone()
	assert.Equal(t, data, cp)
	data[0] = 0
	assert.Equal(t, 3, cp[0])
}

func TestSwap(t *testing.T) {
	data := slice.Slice[int]{3, 1, 4, 1, 5, 9}
	data.Swap(0, 1)
	assert.Equal(t, 1, data[0])
	assert.Equal(t, 3, data[1])

}

func TestLess(t *testing.T) {
	i := []int{6, 7, 9, 2, 3, 4, 1, 5, 8}
	slice.GT[int]().Sort(i)
	expected := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	assert.Equal(t, expected, i)

	slice.LT[int]().Sort(i)
	expected = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	assert.Equal(t, expected, i)
}

func TestSort(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}
	people := slice.Slice[person]{
		{"Adam", 39},
		{"Lauren", 37},
		{"Stephen", 38},
		{"Alex", 35},
		{"Fletcher", 5},
	}.Sort(func(i, j person) bool {
		return i.Age < j.Age
	})
	expect := slice.Slice[person]{
		{"Fletcher", 5},
		{"Alex", 35},
		{"Lauren", 37},
		{"Stephen", 38},
		{"Adam", 39},
	}
	assert.Equal(t, expect, people)
}

func TestIter(t *testing.T) {
	s := slice.Slice[int]{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	it := s.Iter()
	forFn := func(i, idx int) {
		assert.Equal(t, s[idx], i)
	}
	c := liter.ForIdx[int](it, forFn)
	assert.Len(t, s, c)

	st, _ := upgrade.To[liter.Starter[int]](it)
	st.Start()
	c = slice.IterFactory(s).ForIdx(forFn)
	assert.Len(t, s, c)
	c = liter.Factory[int](s.IterFactory).ForIdx(forFn)
	assert.Len(t, s, c)

	st.Start()
	liter.Seek[int](it, func(i int) bool {
		assert.True(t, i < 4)
		return i == 3
	})

	s[0] = 100
	i, done := st.Start()
	assert.Equal(t, 100, i)
	assert.False(t, done)
}

func TestAppendNotZero(t *testing.T) {
	got := slice.Slice[string]{"Start"}.AppendNotZero("", "Foo", "", "Bar", "Baz", "")
	expected := []string{"Start", "Foo", "Bar", "Baz"}
	assert.Equal(t, expected, got)

	gotAny := slice.Slice[any]{}.AppendNotZero(1, 0, 2.0, 0.0, "", "test")
	expectedAny := []any{1, 2.0, "test"}
	assert.Equal(t, expectedAny, gotAny)
}

func TestAppendNotZeroInterface(t *testing.T) {
	var s slice.Slice[fmt.Stringer]
	var n fmt.Stringer
	s = s.AppendNotZero(n)
	assert.Len(t, s, 0)
}

func TestRemove(t *testing.T) {
	data := slice.Slice[int]{3, 1, 4, 1, 5, 9}
	data = data.Remove(5, 1, 3)
	expected := slice.Slice[int]{3, 5, 4}
	assert.Equal(t, expected, data)

	data = slice.Slice[int]{3, 1, 4, 1, 5, 9}
	data = data.Remove(0, 0, -1, 100)
	expected = slice.Slice[int]{9, 1, 4, 1, 5}
	assert.Equal(t, expected, data)
}

func TestRemoveOrdered(t *testing.T) {
	tt := map[string]struct {
		start    slice.Slice[int]
		remove   []int
		expected slice.Slice[int]
	}{
		"basic": {
			start:    slice.Slice[int]{3, 1, 4, 1, 5, 9},
			remove:   []int{1},
			expected: slice.Slice[int]{3, 4, 1, 5, 9},
		},
		"first": {
			start:    slice.Slice[int]{3, 1, 4, 1, 5, 9},
			remove:   []int{0},
			expected: slice.Slice[int]{1, 4, 1, 5, 9},
		},
		"last": {
			start:    slice.Slice[int]{3, 1, 4, 1, 5, 9},
			remove:   []int{5},
			expected: slice.Slice[int]{3, 1, 4, 1, 5},
		},
		"two": {
			start:    slice.Slice[int]{3, 1, 4, 1, 5, 9},
			remove:   []int{3, 1},
			expected: slice.Slice[int]{3, 4, 5, 9},
		},
		"two-first": {
			start:    slice.Slice[int]{3, 1, 4, 1, 5, 9},
			remove:   []int{3, 0},
			expected: slice.Slice[int]{1, 4, 5, 9},
		},
		"two-last": {
			start:    slice.Slice[int]{3, 1, 4, 1, 5, 9},
			remove:   []int{5, 3},
			expected: slice.Slice[int]{3, 1, 4, 5},
		},
		"first-last": {
			start:    slice.Slice[int]{3, 1, 4, 1, 5, 9},
			remove:   []int{5, 0},
			expected: slice.Slice[int]{1, 4, 1, 5},
		},
		"repeat": {
			start:    slice.Slice[int]{3, 1, 4, 1, 5, 9},
			remove:   []int{3, 3},
			expected: slice.Slice[int]{3, 1, 4, 5, 9},
		},
		"repeat-first": {
			start:    slice.Slice[int]{3, 1, 4, 1, 5, 9},
			remove:   []int{0, 0},
			expected: slice.Slice[int]{1, 4, 1, 5, 9},
		},
		"repeat-last": {
			start:    slice.Slice[int]{3, 1, 4, 1, 5, 9},
			remove:   []int{5, 5},
			expected: slice.Slice[int]{3, 1, 4, 1, 5},
		},
		"none": {
			start:    slice.Slice[int]{3, 1, 4, 1, 5, 9},
			remove:   []int{},
			expected: slice.Slice[int]{3, 1, 4, 1, 5, 9},
		},
		"out-of-range": {
			start:    slice.Slice[int]{3, 1, 4, 1, 5, 9},
			remove:   []int{-1, 6, 3, -3},
			expected: slice.Slice[int]{3, 1, 4, 5, 9},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.start.RemoveOrdered(tc.remove...))
		})
	}
}

func TestPop(t *testing.T) {
	data := slice.Slice[int]{3, 1, 4, 1, 5, 9}
	i, got := data.Pop()
	assert.Equal(t, 9, i)
	assert.Equal(t, data[:5], got)

	data = nil
	i, got = data.Pop()
	assert.Equal(t, 0, i)
	assert.Nil(t, got)
}

func TestShift(t *testing.T) {
	data := slice.Slice[int]{3, 1, 4, 1, 5, 9}
	i, got := data.Shift()
	assert.Equal(t, 3, i)
	assert.Equal(t, data[1:6], got)

	data = nil
	i, got = data.Shift()
	assert.Equal(t, 0, i)
	assert.Nil(t, got)
}

type genIter struct {
	idx int
}

func (g *genIter) Next() (t int, done bool) {
	if !g.Done() {
		g.idx++
	}
	return g.idx, g.Done()
}
func (g *genIter) Cur() (t int, done bool) { return g.idx, g.Done() }
func (g *genIter) Done() bool              { return g.idx >= 10 }
func (g *genIter) Idx() int                { return g.idx }

type genIterLen struct {
	*genIter
}

func (g *genIterLen) Len() int { return 10 }

func TestIterSlice(t *testing.T) {
	s := slice.New([]int{3, 1, 4, 1, 5, 9})
	it := slice.NewIter(s)
	got := slice.FromIter[int](it, nil)
	assert.Equal(t, s, got)

	got = slice.FromIter[int](it, make([]int, 0, len(s)))
	assert.Equal(t, s, got)
	s[0] = 4
	assert.NotEqual(t, s, got)

	got = slice.FromIter[int](&genIter{}, nil)
	expected := slice.New([]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	assert.Equal(t, expected, got)

	got = slice.FromIter[int](&genIterLen{&genIter{}}, nil)
	assert.Equal(t, expected, got)

	fact := slice.IterFactory(expected)
	got = slice.FromIterFactory(fact, got)
	assert.Equal(t, expected, got)

	fact = slice.IterFactory([]int{})
	got = slice.FromIterFactory(fact, got)
	assert.Nil(t, got)
}

func TestCheckCapacity(t *testing.T) {
	expected := slice.Slice[int]{3, 1, 4, 1, 5}
	data := slice.Make[int](0, 10)
	data = append(data, expected...)

	data = data.CheckCapacity(7)
	assert.Equal(t, expected, data)
	assert.True(t, cap(data) >= 7)

	data = data.CheckCapacity(15)
	assert.Equal(t, expected, data)
	assert.True(t, cap(data) >= 15)
}

func TestMake(t *testing.T) {
	data := slice.Make[int](0, 10)
	assert.Equal(t, 10, cap(data))
	assert.Equal(t, 0, len(data))

	data = slice.Make[int](15, 0)
	assert.Equal(t, 15, cap(data))
	assert.Equal(t, 15, len(data))
}

func TestSearch(t *testing.T) {
	data := slice.Slice[int]{2, 3, 5, 7, 11, 13, 17, 19, 23}
	fn := func(i int) bool { return i >= 5 }
	idx := data.Search(fn)
	assert.Equal(t, 5, data[idx])

	fn = func(i int) bool { return i >= 10 }
	idx = data.Search(fn)
	assert.Equal(t, 11, data[idx])

	fn = func(i int) bool { return i >= 0 }
	idx = data.Search(fn)
	assert.Equal(t, 0, idx)

	fn = func(i int) bool { return i >= 24 }
	idx = data.Search(fn)
	assert.Equal(t, len(data), idx)
}

func TestIdxCheck(t *testing.T) {
	data := slice.Slice[int]{2, 3, 5, 7, 11, 13, 17, 19, 23}
	assert.False(t, data.IdxCheck(-3))
	assert.False(t, data.IdxCheck(-1))
	assert.True(t, data.IdxCheck(0))
	assert.True(t, data.IdxCheck(5))
	assert.True(t, data.IdxCheck(len(data)-1))
	assert.False(t, data.IdxCheck(len(data)))
	assert.False(t, data.IdxCheck(len(data)+1))
}

func TestLen(t *testing.T) {
	s := []int{2, 3, 5, 7, 11, 13, 17, 19, 23}
	assert.Equal(t, len(s), slice.Len(s))
}

func TestTransform(t *testing.T) {
	in := slice.New([]int{3, -8, -15, 1, -2, 4, 1, -55, -66, 5, -7, 9})
	fn := func(i, idx int) (string, bool) {
		assert.Equal(t, i, in[idx])
		if i < 0 {
			return "", false
		}
		return strconv.Itoa(i), true
	}
	got := slice.Transform(in.Iter(), nil, fn)
	expected := slice.Slice[string]{"3", "1", "4", "1", "5", "9"}
	assert.Equal(t, expected, got)
	assert.Equal(t, len(got), cap(got))

	got = slice.TransformSlice(in, nil, fn)
	assert.Equal(t, expected, got)
	assert.Equal(t, len(got), cap(got))

	expected = slice.Slice[string]{"3", "-8", "-15", "1", "-2", "4", "1", "-55", "-66", "5", "-7", "9"}
	tfn := slice.ForAll(strconv.Itoa)
	got = tfn.Slice(in, nil)
	assert.Equal(t, expected, got)
	got = tfn.Factory(in.IterFactory, nil)
	assert.Equal(t, expected, got)

	got = slice.Transform(in.Iter(), nil, func(i, idx int) (string, bool) {
		return "", false
	})
	assert.Nil(t, got)

	in = slice.New([]int{})
	got = slice.Transform(in.Iter(), nil, fn)
	assert.Nil(t, got)
}

func TestIdx(t *testing.T) {
	s := slice.Slice[int]{3, 1, 4, 1, 5}
	tt := map[int]struct {
		expected int
		ok       bool
	}{
		0: {
			expected: 0,
			ok:       true,
		},
		4: {
			expected: 4,
			ok:       true,
		},
		5: {
			expected: 5,
			ok:       false,
		},
		-1: {
			expected: 4,
			ok:       true,
		},
		-5: {
			expected: 0,
			ok:       true,
		},
		-6: {
			expected: -1,
			ok:       false,
		},
	}

	for n, tc := range tt {
		t.Run(strconv.Itoa(n), func(t *testing.T) {
			idx, ok := s.Idx(n)
			assert.Equal(t, tc.expected, idx)
			assert.Equal(t, tc.ok, ok)
		})
	}
}

func TestReverse(t *testing.T) {
	s := slice.Slice[int]{2, 3, 5, 7, 11, 13, 17, 19, 23}
	ln := len(s)
	expected := make(slice.Slice[int], ln)
	for i := 0; i < ln; i++ {
		expected[i] = s[ln-1-i]
	}

	s.Reverse()
	assert.Equal(t, expected, s)
}

func TestSliceList(t *testing.T) {
	var l list.List[string] = slice.New([]string{"apple", "banana", "cantaloupe", "date"})

	assert.Equal(t, 4, l.Len())
	assert.Equal(t, "banana", l.AtIdx(1))
}
