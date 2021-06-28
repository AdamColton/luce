package filter

// Compare builds simple filters.
type Compare byte

// Compare Constants
const (
	EQ Compare = 1 << iota
	GT
	not

	NEQ = not | EQ
	LTE = not | GT
	LT  = NEQ | GT
	GTE = EQ | GT
)

// Int creates a filter that will take a value and apply the comparison to the
// value given for a.
func (c Compare) Int(a int) Int {
	return func(b int) bool {
		return ((c&EQ == EQ && a == b) ||
			(c&GT == GT && b > a)) !=
			(c&not == not)
	}
}

// Float creates a filter that will take a value and apply the comparison to the
// value given for a.
func (c Compare) Float(a float64) Float {
	return func(b float64) bool {
		return ((c&EQ == EQ && a == b) ||
			(c&GT == GT && b > a)) !=
			(c&not == not)
	}
}

// String creates a filter that will take a value and apply the comparison to
// the value given for a.
func (c Compare) String(a string) String {
	return func(b string) bool {
		return ((c&EQ == EQ && a == b) ||
			(c&GT == GT && b > a)) !=
			(c&not == not)
	}
}

// Str return the string representing the Compare
func (c Compare) Str() string {
	switch c {
	case LT:
		return "<"
	case LTE:
		return "<="
	case EQ:
		return "=="
	case GT:
		return ">"
	case GTE:
		return ">="
	case NEQ:
		return "!="
	}
	return "??"
}
