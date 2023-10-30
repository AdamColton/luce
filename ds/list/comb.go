package list

import (
	"github.com/adamcolton/luce/math/ints"
)

type Setter[T any] interface {
	set(idx []int, t T) T
	ln([]int)
	chainLen() int
}

func NewListSetter[From, To any](l List[From], fn func(From, To) To) *ListSetter[From, To] {
	return &ListSetter[From, To]{
		List: l,
		Set:  fn,
	}
}

func InitSetter[From, To any](l List[From], fn func(From, *To)) *ListSetter[From, *To] {
	sfn := func(f From, _ *To) *To {
		t := new(To)
		fn(f, t)
		return t
	}
	return NewListSetter(l, sfn)
}

func NewSetter[From, To any](l List[From], fn func(From, To)) *ListSetter[From, To] {
	sfn := func(f From, t To) To {
		fn(f, t)
		return t
	}
	return NewListSetter(l, sfn)
}

func Chain[From, To any](next *Setter[*To], l List[From], fn func(From, *To)) (out *ListSetter[From, *To]) {
	if next == nil {
		out = InitSetter(l, fn)
	} else {
		out = NewSetter(l, fn)
		*next = out
	}
	return
}

type ListSetter[From, To any] struct {
	List[From]
	Set  func(From, To) To
	Next Setter[To]
}

func (ls *ListSetter[From, To]) set(idx []int, to To) To {
	if len(idx) == 0 || idx[0] > ls.Len() {
		return to
	}
	to = ls.Set(ls.AtIdx(idx[0]), to)
	if ls.Next != nil {
		to = ls.Next.set(idx[1:], to)
	}
	return to
}

func (ls *ListSetter[From, To]) chainLen() int {
	if ls.Next == nil {
		return 1
	}
	return ls.Next.chainLen() + 1
}

func (ls *ListSetter[From, To]) ln(ln []int) {
	ln[0] = ls.Len()
	if ls.Next != nil {
		ls.Next.ln(ln[1:])
	}
}

type combinatorSetter[T any] struct {
	fn ints.Combinator[int]
	ln int
	s  Setter[T]
}

func (c *combinatorSetter[T]) AtIdx(idx int) (out T) {
	return c.s.set(c.fn(idx), out)
}

func (c *combinatorSetter[T]) Len() int {
	return c.ln
}

func Combinator[T any](setter Setter[T], factory ints.CombinatorFactory[int]) Wrapper[T] {
	lns := make([]int, setter.chainLen())
	setter.ln(lns)
	fn, ln := factory(lns[0], lns[1])
	return Wrapper[T]{&combinatorSetter[T]{
		fn: fn,
		ln: ln,
		s:  setter,
	}}
}

func SliceCombinator[T any](factory ints.CombinatorFactory[int], ls ...List[T]) Wrapper[[]T] {
	lns := make([]int, len(ls))
	for i, l := range ls {
		lns[i] = l.Len()
	}
	fn, ln := factory(lns...)
	return Wrapper[[]T]{&sliceCombinator[T]{
		ls: ls,
		fn: fn,
		ln: ln,
	}}
}

type sliceCombinator[T any] struct {
	ls []List[T]
	fn ints.Combinator[int]
	ln int
}

func (c *sliceCombinator[T]) AtIdx(idx int) []T {
	out := make([]T, len(c.ls))
	for i, idx := range c.fn(idx) {
		out[i] = c.ls[i].AtIdx(idx)
	}
	return out
}

func (c *sliceCombinator[T]) Len() int {
	return c.ln
}
