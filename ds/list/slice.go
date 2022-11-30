package list

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
	out := AllocBuf(buf, ln)
	for i := 0; i < ln; i++ {
		out = append(out, l.AtIdx(i))
	}
	return out
}

func AllocBuf[T any](buf []T, ln int) []T {
	if cap(buf) >= ln {
		return buf[:0]
	}
	return make([]T, 0, ln)
}
