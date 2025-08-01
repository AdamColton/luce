package flow

// NilCheck returns t if it is not nil. If it is nil, it returns the result of
// invoking constructor. It does not change t.
func NilCheck[T any](t *T, constructor func() *T) *T {
	if t == nil {
		return constructor()
	}
	return t
}

// Tern implements the ternary operator. It returns a if cond is true and b if it
// is not. Both a and b are evaluated before Tern is called.
func Tern[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
