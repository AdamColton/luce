package slice_test

import (
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

	var st iter.Starter[int]
	upgrade.Upgrade(it, &st)
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
	s := []int{3, 1, 4, 1, 5, 9}
	it := slice.NewIter(s)
	got := slice.IterSlice[int](it, nil)
	assert.Equal(t, s, got)

	got = slice.IterSlice[int](it, make([]int, 0, len(s)))
	assert.Equal(t, s, got)
	s[0] = 4
	assert.NotEqual(t, s, got)

	got = slice.IterSlice[int](&genIter{}, nil)
	expected := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	assert.Equal(t, expected, got)

	got = slice.IterSlice[int](&genIterLen{&genIter{}}, nil)
	assert.Equal(t, expected, got)
}

func TestCheckCapacity(t *testing.T) {
	expected := slice.Slice[int]{3, 1, 4, 1, 5}
	data := slice.Slice[int](make([]int, 0, 10))
	data = append(data, expected...)

	data = data.CheckCapacity(7)
	assert.Equal(t, expected, data)
	assert.True(t, cap(data) >= 7)

	data = data.CheckCapacity(15)
	assert.Equal(t, expected, data)
	assert.True(t, cap(data) >= 15)
}
