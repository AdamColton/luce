package bimap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/bimap"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/stretchr/testify/assert"
)

func TestM2M(t *testing.T) {
	m2m := bimap.NewM2M[string, int]()

	m2m.Add("A", 1)
	m2m.Add("A", 2)
	m2m.Add("A", 3)

	m2m.Add("B", 2)
	m2m.Add("B", 3)
	m2m.Add("B", 4)

	m2m.Add("C", 3)
	m2m.Add("C", 4)
	m2m.Add("C", 5)

	b := m2m.A("A")
	assert.True(t, b.Contains(1))
	assert.True(t, b.Contains(2))
	assert.True(t, b.Contains(3))

	b = m2m.A("B")
	assert.True(t, b.Contains(2))
	assert.True(t, b.Contains(3))
	assert.True(t, b.Contains(4))

	b = m2m.A("C")
	assert.True(t, b.Contains(3))
	assert.True(t, b.Contains(4))
	assert.True(t, b.Contains(5))

	a := m2m.B(1)
	assert.True(t, a.Contains("A"))

	a = m2m.B(2)
	assert.True(t, a.Contains("A"))
	assert.True(t, a.Contains("B"))

	a = m2m.B(3)
	assert.True(t, a.Contains("A"))
	assert.True(t, a.Contains("B"))
	assert.True(t, a.Contains("C"))

	a = m2m.B(4)
	assert.True(t, a.Contains("B"))
	assert.True(t, a.Contains("C"))

	a = m2m.B(5)
	assert.True(t, a.Contains("C"))

	assert.Equal(t, 3, m2m.LenA())
	assert.Equal(t, 5, m2m.LenB())

	expectA := slice.Slice[string]{"A", "B", "C"}
	gotA := m2m.SliceA(nil).Sort(slice.LT[string]())
	assert.Equal(t, expectA, gotA)

	expectB := slice.Slice[int]{1, 2, 3, 4, 5}
	gotB := m2m.SliceB(nil).Sort(slice.LT[int]())
	assert.Equal(t, expectB, gotB)

	aLens := map[string]int{
		"A": 3,
		"B": 3,
		"C": 3,
	}
	m2m.EachA(func(a string, b lset.Reader[int], done *bool) {
		assert.Equal(t, aLens[a], b.Len())
	})

	bLens := map[int]int{
		1: 1,
		2: 2,
		3: 3,
		4: 2,
		5: 1,
	}
	m2m.EachB(func(b int, a lset.Reader[string], done *bool) {
		assert.Equal(t, bLens[b], a.Len())
	})

	m2m.Remove("C", 5)
	a = m2m.B(5)
	assert.Nil(t, a)

	b = m2m.A("C")
	assert.True(t, b.Contains(3))
	assert.True(t, b.Contains(4))

	m2m.Remove("C", 4)
	m2m.Remove("C", 3)
	b = m2m.A("C")
	assert.Nil(t, b)

	type key struct {
		a string
		b int
	}

	c := 0
	m2m.Each(func(a string, b int, done *bool) {
		*done = true
		c++
	})
	assert.Equal(t, 1, c)

	expect := map[key]bool{
		{a: "A", b: 1}: false,
		{a: "A", b: 2}: false,
		{a: "A", b: 3}: false,
		{a: "B", b: 2}: false,
		{a: "B", b: 3}: false,
		{a: "B", b: 4}: false,
	}

	m2m.Each(func(a string, b int, done *bool) {
		expect[key{a: a, b: b}] = true
	})

	if assert.Len(t, expect, 6) {
		for k, b := range expect {
			assert.True(t, b, k)
		}
	}

	m2m.RemoveA("B")
	a = m2m.B(3)
	assert.Equal(t, 1, a.Len())
	assert.True(t, a.Contains("A"))
	assert.False(t, a.Contains("B"))

	m2m.RemoveB(3)
	b = m2m.A("A")
	assert.Equal(t, 2, b.Len())
	m2m.RemoveB(2)
	m2m.RemoveB(1)
	b = m2m.A("A")
	assert.Nil(t, b)
}
