package slice_test

import (
	"fmt"
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/iter"
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

func TestKeys(t *testing.T) {
	data := map[int]string{
		1: "1",
		2: "2",
		3: "3",
		4: "4",
		5: "5",
		6: "6",
	}
	got := slice.Keys(data)
	slice.LT[int]().Sort(got)
	expected := slice.Slice[int]{1, 2, 3, 4, 5, 6}
	assert.Equal(t, expected, got)
}

func TestVals(t *testing.T) {
	data := map[int]string{
		1: "1",
		2: "2",
		3: "3",
		4: "4",
		5: "5",
		6: "6",
	}
	got := slice.Vals(data)
	slice.LT[string]().Sort(got)
	expected := slice.Slice[string]{"1", "2", "3", "4", "5", "6"}
	assert.Equal(t, expected, got)
}

func TestLess(t *testing.T) {
	i := []int{6, 7, 9, 2, 3, 4, 1, 5, 8}
	slice.GT[int]().Sort(i)
	expected := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	assert.Equal(t, expected, i)
}

func TestUnique(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	got := slice.Unique(data)
	expected := slice.Slice[int]{3, 1, 4, 5, 9}
	l := slice.LT[int]()
	l.Sort(got)
	l.Sort(expected)
	assert.Equal(t, expected, got)

}

func TestIter(t *testing.T) {
	s := slice.Slice[int]{3, 1, 4, 1, 5, 9, 2, 6, 5, 3}
	it := s.Iter()
	forFn := func(i, idx int) {
		assert.Equal(t, s[idx], i)
	}
	c := iter.ForIdx[int](it, forFn)
	assert.Len(t, s, c)

	st, _ := upgrade.To[iter.Starter[int]](it)
	st.Start()
	c = slice.IterFactory(s).ForIdx(forFn)
	assert.Len(t, s, c)
	c = iter.Factory[int](s.IterFactory).ForIdx(forFn)
	assert.Len(t, s, c)

	st.Start()
	iter.Seek[int](it, func(i int) bool {
		assert.True(t, i < 4)
		return i == 3
	})

	s[0] = 100
	i, done := st.Start()
	assert.Equal(t, 100, i)
	assert.False(t, done)
}

func TestForAll(t *testing.T) {
	s := slice.Slice[int]{0, 1, 2, 3, 4, 5}
	c := 0
	fn := func(idx int, i int) {
		assert.Equal(t, idx, i)
		c++
	}
	s.ForAll(fn).Wait()
	assert.Len(t, s, c)
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
