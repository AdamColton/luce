package slice_test

import (
	"fmt"
	"testing"

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
