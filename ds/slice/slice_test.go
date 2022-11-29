package slice_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/iter"
	"github.com/stretchr/testify/assert"
)

func TestClone(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	cp := slice.Clone(data)
	assert.Equal(t, data, cp)
	data[0] = 0
	assert.Equal(t, 3, cp[0])
}

func TestSwap(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	slice.Swap(data, 0, 1)
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
	slice.Less[int](func(i, j int) bool {
		return i < j
	}).Sort(got)
	expected := []int{1, 2, 3, 4, 5, 6}
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
	slice.Less[string](func(i, j string) bool {
		return i < j
	}).Sort(got)
	expected := []string{"1", "2", "3", "4", "5", "6"}
	assert.Equal(t, expected, got)
}

func TestLess(t *testing.T) {
	i := []int{6, 7, 9, 2, 3, 4, 1, 5, 8}
	slice.Less[int](func(i, j int) bool {
		return i < j
	}).Sort(i)
	expected := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	assert.Equal(t, expected, i)
}

func TestUnique(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	got := slice.Unique(data)
	expected := []int{3, 1, 4, 5, 9}
	l := slice.Less[int](func(i, j int) bool {
		return i < j
	})
	l.Sort(got)
	l.Sort(expected)
	assert.Equal(t, expected, got)

}

func TestIter(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5}
	it := slice.NewIter(s)
	c := 0
	doFn := func(i int) bool {
		assert.Equal(t, c, i)
		c++
		return false
	}
	iter.Do[int](it, doFn)
	assert.Len(t, s, c)

	c = 0
	slice.IterFactory(s).Do(doFn)
	assert.Len(t, s, c)

	it.I = 0
	c = 0
	iter.Do[int](it, func(i int) bool {
		assert.Equal(t, c, it.Idx())
		c++
		assert.True(t, i < 4)
		return i == 3
	})

	s[0] = 100
	i, done := it.Start()
	assert.Equal(t, 100, i)
	assert.False(t, done)

}

func TestForAll(t *testing.T) {
	s := []int{0, 1, 2, 3, 4, 5}
	c := 0
	fn := func(idx int, i int) {
		assert.Equal(t, idx, i)
		c++
	}
	slice.ForAll(s, fn).Wait()
	assert.Len(t, s, c)
}

func TestAppendNotDefault(t *testing.T) {
	got := slice.AppendNotDefault([]string{"Start"}, "", "Foo", "", "Bar", "Baz", "")
	expected := []string{"Start", "Foo", "Bar", "Baz"}
	assert.Equal(t, expected, got)
}

func TestRemove(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	data = slice.Remove(data, 5, 1, 3)
	expected := []int{3, 5, 4}
	assert.Equal(t, expected, data)

	data = []int{3, 1, 4, 1, 5, 9}
	data = slice.Remove(data, 0, 0, -1, 100)
	expected = []int{9, 1, 4, 1, 5}
	assert.Equal(t, expected, data)
}

func TestPop(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	i, got := slice.Pop(data)
	assert.Equal(t, 9, i)
	assert.Equal(t, data[:5], got)

	data = nil
	i, got = slice.Pop(data)
	assert.Equal(t, 0, i)
	assert.Nil(t, got)
}

func TestShift(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	i, got := slice.Shift(data)
	assert.Equal(t, 3, i)
	assert.Equal(t, data[1:6], got)

	data = nil
	i, got = slice.Shift(data)
	assert.Equal(t, 0, i)
	assert.Nil(t, got)
}
