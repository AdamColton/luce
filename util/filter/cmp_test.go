package filter

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	tt := map[string]struct {
		cmp Compare
		fn  func(a, b int) bool
		str string
	}{
		"lt": {
			cmp: LT,
			fn:  func(a, b int) bool { return b < a },
			str: "<",
		},
		"lte": {
			cmp: LTE,
			fn:  func(a, b int) bool { return b <= a },
			str: "<=",
		},
		"eq": {
			cmp: EQ,
			fn:  func(a, b int) bool { return b == a },
			str: "==",
		},
		"gt": {
			cmp: GT,
			fn:  func(a, b int) bool { return b > a },
			str: ">",
		},
		"gte": {
			cmp: GTE,
			fn:  func(a, b int) bool { return b >= a },
			str: ">=",
		},
		"neq": {
			cmp: NEQ,
			fn:  func(a, b int) bool { return b != a },
			str: "!=",
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			for a := 0; a < 3; a++ {
				as := strconv.Itoa(a)
				for b := 0; b < 3; b++ {
					expect := tc.fn(a, b)
					bs := strconv.Itoa(b)
					assert.Equal(t, tc.str, tc.cmp.Str())
					str := bs + tc.cmp.Str() + as
					assert.Equal(t, expect, tc.cmp.Int(a)(b), str)
					assert.Equal(t, expect, tc.cmp.Float(float64(a))(float64(b)), str)
					assert.Equal(t, expect, tc.cmp.String(as)(bs), str)
				}
			}
		})
	}

	assert.Equal(t, "??", Compare(20).Str())
}
