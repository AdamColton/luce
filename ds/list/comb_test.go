package list_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/math/ints"
	"github.com/adamcolton/luce/math/ints/comb"
	"github.com/stretchr/testify/assert"
)

func TestCombinator(t *testing.T) {
	a := list.Generator[float64]{
		Fn: func(i int) float64 {
			f := float64(i)
			return f * f
		},
		Length: 3,
	}
	b := list.Generator[string]{
		Fn: func(i int) string {
			return strconv.Itoa(i)
		},
		Length: 5,
	}
	c := list.Combinator(a, b, comb.Cross)
	assert.Equal(t, 15, c.Len())

	expected := []struct {
		A float64
		B string
	}{
		{A: 0, B: "0"}, {A: 1, B: "0"}, {A: 4, B: "0"},
		{A: 0, B: "1"}, {A: 1, B: "1"}, {A: 4, B: "1"},
		{A: 0, B: "2"}, {A: 1, B: "2"}, {A: 4, B: "2"},
		{A: 0, B: "3"}, {A: 1, B: "3"}, {A: 4, B: "3"},
		{A: 0, B: "4"}, {A: 1, B: "4"}, {A: 4, B: "4"},
	}

	for i, e := range expected {
		assert.Equal(t, e, c.AtIdx(i), i)
	}

	gf := list.GeneratorFactory(ints.Int)
	grid := list.Combinator(gf(3), gf(5), comb.Cross).Iter().Channel(0)

	for y := 0; y < 5; y++ {
		for x := 0; x < 3; x++ {
			pt := <-grid
			assert.Equal(t, x, pt.A)
			assert.Equal(t, y, pt.B)
		}
	}
}
