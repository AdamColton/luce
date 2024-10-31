package list

// Reducer combines two values into one.
type Reducer[T any] func(T, T) T

// Reduce combines the values of the List from the first to the last, so the
// result for values a, b, c is r(r(a, b), c). An empty List returns the zero
// value and a List with one value returns that value without calling r.
func Reduce[T any](l List[T], r Reducer[T]) (t T) {
	if ln := l.Len(); ln > 0 {
		t = l.AtIdx(0)
		for i := 1; i < ln; i++ {
			t = r(t, l.AtIdx(i))
		}
	}
	return
}
