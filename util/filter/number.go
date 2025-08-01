package filter

import (
	"github.com/adamcolton/luce/math/ints"
)

func InRange[N ints.Number](start, end N) Filter[N] {
	return func(n N) bool {
		return n >= start && n < end
	}
}
