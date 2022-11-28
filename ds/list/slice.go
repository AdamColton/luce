package list

import "github.com/adamcolton/luce/ds/slice"

type SliceList[T any] []T

func (sl SliceList[T]) AtIdx(idx int) T {
	return sl[idx]
}

func (sl SliceList[T]) Len() int {
	return len(sl)
}

func (sl SliceList[T]) Slice(buf []T) []T {
	return sl
}

func ToSlice[T any](l List[T], buf []T) []T {
	if s, ok := l.(Slicer[T]); ok {
		return s.Slice(buf)
	}
	ln := l.Len()
	out := slice.BufferEmpty(ln, buf)
	for i := 0; i < ln; i++ {
		out = append(out, l.AtIdx(i))
	}
	return out
}
