package filter

// Comparable are the types that support ==, != and >.
type Comparable interface {
	// TODO: is this already defined in a standard package
	// should I move this to ltype?
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~string
}

// EQ returns a Filter that will check if a given value is equal to 'a'.
func EQ[T Comparable](a T) Filter[T] {
	return func(b T) bool { return a == b }
}

// EQ returns a Filter that will check if a given value is greater than to 'a'.
func GT[T Comparable](a T) Filter[T] {
	return func(b T) bool { return a < b }
}

// NEQ returns a Filter that will check if a given value is not equal to 'a'.
func NEQ[T Comparable](a T) Filter[T] {
	return func(b T) bool { return a != b }
}

// LTE returns a Filter that will check if a given value is less than or equal
// to 'a'.
func LTE[T Comparable](a T) Filter[T] {
	return func(b T) bool { return a >= b }
}

// LT returns a Filter that will check if a given value is not less than to 'a'.
func LT[T Comparable](a T) Filter[T] {
	return func(b T) bool { return a > b }
}

// GTE returns a Filter that will check if a given value is greater than or
// equal to 'a'.
func GTE[T Comparable](a T) Filter[T] {
	return func(b T) bool { return a <= b }
}
