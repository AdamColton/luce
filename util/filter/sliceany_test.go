package filter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

func TestSliceAnyInPlaceEdges(t *testing.T) {
	isString := filter.IsType(ltype.String)
	tt := map[string]struct {
		vals             []any
		passing, failing int
	}{
		"all pass":  {[]any{"a", "b", "c"}, 3, 0},
		"none pass": {[]any{1, 2, 3}, 0, 3},
		"one pass":  {[]any{"a"}, 1, 0},
		"one fail":  {[]any{1}, 0, 1},
		"empty":     {nil, 0, 0},
	}
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			passing, failing := isString.SliceAnyInPlace(tc.vals)
			assert.Len(t, passing, tc.passing)
			assert.Len(t, failing, tc.failing)
		})
	}
}
