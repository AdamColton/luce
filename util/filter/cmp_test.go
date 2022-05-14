package filter

import (
	"strconv"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	tt := map[string]struct {
		cmp Compare
		fn  func(a, b int) bool
		str string
	}{
		"lt": {
			cmp: CmprLT,
			fn:  func(a, b int) bool { return b < a },
			str: "<",
		},
		"lte": {
			cmp: CmprLTE,
			fn:  func(a, b int) bool { return b <= a },
			str: "<=",
		},
		"eq": {
			cmp: CmprEQ,
			fn:  func(a, b int) bool { return b == a },
			str: "==",
		},
		"gt": {
			cmp: CmprGT,
			fn:  func(a, b int) bool { return b > a },
			str: ">",
		},
		"gte": {
			cmp: CmprGTE,
			fn:  func(a, b int) bool { return b >= a },
			str: ">=",
		},
		"neq": {
			cmp: CmprNEQ,
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
					assert.Equal(t, tc.str, tc.cmp.String())
					str := bs + tc.cmp.String() + as

					assert.Equal(t, expect, CompareFilter(tc.cmp, a)(b), str)
					assert.Equal(t, expect, CompareFilter(tc.cmp, float64(a))(float64(b)), str)
					assert.Equal(t, expect, CompareFilter(tc.cmp, as)(bs), str)
				}
			}
		})
	}

	assert.Equal(t, "??", Compare(20).String())

	assert.True(t, EQ("x")("x"))
	assert.False(t, EQ("x")("y"))
	assert.False(t, NEQ("x")("x"))
	assert.True(t, NEQ("x")("y"))
	assert.False(t, LTE(5)(6))
	assert.True(t, LTE(5)(5))
	assert.True(t, LTE(5)(4))
}

func TestChecker(t *testing.T) {
	err := lerr.Str("Test Error")
	c := EQ("Test").Check(err)
	assert.NoError(t, c("Test"))
	assert.Equal(t, err, c("foo"))

	defer func() {
		assert.Equal(t, err, recover())
	}()
	c.Panic("foo")
}
