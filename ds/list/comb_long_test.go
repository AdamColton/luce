package list_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/list"
	"github.com/adamcolton/luce/math/ints"
	"github.com/stretchr/testify/assert"
)

// Combinator only combines two lists, so a third Setter in the chain is not
// applied.
func TestCombinatorLongChain(t *testing.T) {
	type out struct {
		A, B, C int
	}
	gen := func(ln int) list.List[int] {
		return list.NewGenerator(ln, func(i int) int { return i + 1 })
	}
	s := list.Chain(nil, gen(2), func(v int, o *out) { o.A = v })
	list.Chain(&s.Next, gen(3), func(v int, o *out) { o.B = v })
	next := s.Next.(*list.ListSetter[int, *out])
	list.Chain(&next.Next, gen(4), func(v int, o *out) { o.C = v })

	c := list.Combinator(s, ints.Cross)
	assert.Equal(t, 6, c.Len())
	assert.Equal(t, &out{A: 2, B: 2}, c.AtIdx(3))
}
