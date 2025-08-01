package flow

// NilCheck returns t if it is not nil and invokes the constructor if it is.
func NilCheck[T any](t *T, constructor func() *T) *T {
	if t == nil {
		return constructor()
	}
	return t
}

// Tern implements ternary operator.
func Tern[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}
