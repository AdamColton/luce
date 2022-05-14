package filter_test

import (
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	tt := map[string]struct {
		fn    func(a, b int) bool
		cmpFn func(a int) filter.Filter[int]
	}{
		"lt": {
			fn:    func(a, b int) bool { return b < a },
			cmpFn: filter.LT[int],
		},
		"lte": {
			fn:    func(a, b int) bool { return b <= a },
			cmpFn: filter.LTE[int],
		},
		"eq": {
			fn:    func(a, b int) bool { return b == a },
			cmpFn: filter.EQ[int],
		},
		"gt": {
			fn:    func(a, b int) bool { return b > a },
			cmpFn: filter.GT[int],
		},
		"gte": {
			fn:    func(a, b int) bool { return b >= a },
			cmpFn: filter.GTE[int],
		},
		"neq": {
			fn:    func(a, b int) bool { return b != a },
			cmpFn: filter.NEQ[int],
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			for a := 0; a < 3; a++ {
				fn := tc.cmpFn(a)
				for b := 0; b < 3; b++ {
					expect := tc.fn(a, b)
					got := fn(b)
					assert.Equal(t, expect, got)
				}
			}
		})
	}
}

func TestChecker(t *testing.T) {
	err := lerr.Str("Test Error")
	c := filter.EQ("Test").Check(func(s string) error {
		return err
	})
	assert.NoError(t, c("Test"))
	assert.Equal(t, err, c("foo"))

	defer func() {
		assert.Equal(t, err, recover())
	}()
	c.Panic("foo")
}
