package morph

// Deref returns the value t points to. It will panic if t is nil.
func Deref[T any](t *T) T {
	return *t
}

// NewDeref returns Deref as a ValAll transform. This is convenient for
// converting a slice or iterator of pointers to values.
func NewDeref[T any]() ValAll[*T, T] {
	return Deref[T]
}
