package comb_test

import (
	"testing"

	"github.com/adamcolton/luce/math/ints/comb"
	"github.com/stretchr/testify/assert"
)

func TestCombinator(t *testing.T) {
	intsTable := map[string]struct {
		expected     [][2]int
		comb         comb.CombinatorFactory[int]
		lnA, lnB, ln int
	}{
		"pair": {
			expected: [][2]int{
				{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4},
				{5, 0}, {6, 1}, {7, 2}, {8, 3}, {9, 4},
				{0, 0}, {1, 1}, {2, 2}, {3, 3}, {4, 4},
			},
			lnA:  10,
			lnB:  5,
			ln:   10,
			comb: comb.Pair[int],
		},
		"cross": {
			expected: [][2]int{
				{0, 0}, {1, 0}, {2, 0},
				{0, 1}, {1, 1}, {2, 1},
				{0, 2}, {1, 2}, {2, 2},
				{0, 3}, {1, 3}, {2, 3},
				{0, 4}, {1, 4}, {2, 4},
			},
			lnA:  3,
			lnB:  5,
			ln:   15,
			comb: comb.Cross[int],
		},
	}

	for n, tc := range intsTable {
		t.Run(n, func(t *testing.T) {
			c, ln := tc.comb(tc.lnA, tc.lnB)
			assert.Equal(t, tc.ln, ln)
			a0, b0 := c(0)
			aLn, bLn := c(ln)
			assert.Equal(t, a0, aLn)
			assert.Equal(t, b0, bLn)
			for i, ab := range tc.expected {
				a, b := c(i)
				assert.Equal(t, ab[0], a, i)
				assert.Equal(t, ab[1], b, i)
			}
		})
	}
}
