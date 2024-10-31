package list

type Reducer[T any] func(T, T) T

func Reduce[T any](l List[T], r Reducer[T]) (t T) {
	if ln := l.Len(); ln > 0 {
		t = l.AtIdx(0)
		for i := 1; i < ln; i++ {
			t = r(t, l.AtIdx(i))
		}
	}
	return
}
