package list_test

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/math/ints"
	"github.com/stretchr/testify/assert"
)

func TestCombinator(t *testing.T) {
	a := list.NewGenerator(3, func(i int) float64 {
		f := float64(i)
		return f * f
	})
	b := list.NewGenerator(5, func(i int) string {
		return strconv.Itoa(i)
	})

	type out struct {
		Float  float64
		String string
	}
	s := list.Chain(nil, a, func(f float64, o *out) {
		o.Float = f
	})
	list.Chain(&s.Next, b, func(s string, o *out) {
		o.String = s
	})

	c := list.Combinator(s, ints.Cross)
	assert.Equal(t, 15, c.Len())

	expected := []*out{
		{Float: 0, String: "0"}, {Float: 1, String: "0"}, {Float: 4, String: "0"},
		{Float: 0, String: "1"}, {Float: 1, String: "1"}, {Float: 4, String: "1"},
		{Float: 0, String: "2"}, {Float: 1, String: "2"}, {Float: 4, String: "2"},
		{Float: 0, String: "3"}, {Float: 1, String: "3"}, {Float: 4, String: "3"},
		{Float: 0, String: "4"}, {Float: 1, String: "4"}, {Float: 4, String: "4"},
	}

	for i, e := range expected {
		assert.Equal(t, e, c.AtIdx(i), i)
	}
}
