package list

import "github.com/adamcolton/luce/math/ints/comb"

type combinator[A, B any] struct {
	fn comb.Combinator[int]
	a  List[A]
	b  List[B]
	ln int
}

func Combinator[A, B any](a List[A], b List[B], factory comb.CombinatorFactory[int]) Wrapper[struct {
	A A
	B B
}] {
	c := &combinator[A, B]{
		a: a,
		b: b,
	}
	c.fn, c.ln = factory(a.Len(), b.Len())
	return Wrapper[struct {
		A A
		B B
	}]{c}
}

func (c *combinator[A, B]) AtIdx(idx int) (out struct {
	A A
	B B
}) {
	a, b := c.fn(idx)
	out.A, out.B = c.a.AtIdx(a), c.b.AtIdx(b)
	return
}

func (c *combinator[A, B]) Len() int {
	return c.ln
}
