package filter

import (
	"github.com/adamcolton/luce/math/ints"
)

// InRange returns a Filter that is true for values from start up to but not
// including end.
func InRange[N ints.Number](start, end N) Filter[N] {
	return func(n N) bool {
		return n >= start && n < end
	}
}
