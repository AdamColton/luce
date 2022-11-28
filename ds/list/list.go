package list

import "github.com/adamcolton/luce/ds/slice"

type Generator[T any] struct {
	Fn     func(int) T
	Length int
}

func (g Generator[T]) AtIdx(idx int) T {
	return g.Fn(idx)
}

func (g Generator[T]) Len() int {
	return g.Length
}

type Reverse[T any] struct {
	List[T]
}

func (r Reverse[T]) AtIdx(idx int) T {
	return r.List.AtIdx(r.Len() - 1 - idx)
}

func (r Reverse[T]) Slice(buf []T) []T {
	ln := r.Len()
	out := slice.BufferEmpty(ln, buf)
	ln--
	for i := 0; i <= ln; i++ {
		out = append(out, r.List.AtIdx(ln-i))
	}
	return out
}
