package filter

// Compare builds simple filters.
type Compare byte

// Compare Constants
const (
	CmprEQ Compare = 1 << iota
	CmprGT
	not

	CmprNEQ = not | CmprEQ
	CmprLTE = not | CmprGT
	CmprLT  = CmprNEQ | CmprGT
	CmprGTE = CmprEQ | CmprGT
)

// EQ returns a filter that will check if a given value is equal to 'a'.
func EQ[T Comparable](a T) Filter[T] {
	return func(b T) bool { return a == b }
}

// EQ returns a filter that will check if a given value is greater than to 'a'.
func GT[T Comparable](a T) Filter[T] {
	return func(b T) bool { return a < b }
}

// NEQ returns a filter that will check if a given value is not equal to 'a'.
func NEQ[T Comparable](a T) Filter[T] {
	return func(b T) bool { return a != b }
}

// LTE returns a filter that will check if a given value is less than or equal
// to 'a'.
func LTE[T Comparable](a T) Filter[T] {
	return func(b T) bool { return a >= b }
}

// LT returns a filter that will check if a given value is not less than to 'a'.
func LT[T Comparable](a T) Filter[T] {
	return func(b T) bool { return a > b }
}

// GTE returns a filter that will check if a given value is greater than or
// equal to 'a'.
func GTE[T Comparable](a T) Filter[T] {
	return func(b T) bool { return a <= b }
}

// Comparable are the types that support ==, != and >.
type Comparable interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		string
}

// CompareFilter creates a Filter from a compare type.
func CompareFilter[T Comparable](c Compare, a T) Filter[T] {
	return func(b T) bool {
		return ((c&CmprEQ == CmprEQ && a == b) ||
			(c&CmprGT == CmprGT && b > a)) !=
			(c&not == not)
	}
}

// String return the string representing the Compare; fulfills Stringer.
func (c Compare) String() string {
	switch c {
	case CmprLT:
		return "<"
	case CmprLTE:
		return "<="
	case CmprEQ:
		return "=="
	case CmprGT:
		return ">"
	case CmprGTE:
		return ">="
	case CmprNEQ:
		return "!="
	}
	return "??"
}
