package lstr_test

import (
	"testing"

	"github.com/adamcolton/luce/util/lstr"
	"github.com/stretchr/testify/assert"
)

func TestMatcher(t *testing.T) {
	tt := map[string]struct {
		expected map[rune]bool
		lstr.Matcher
	}{
		"rune_a": {
			Matcher: lstr.Rune('a'),
			expected: map[rune]bool{
				'a': true,
				'b': false,
				'A': false,
			},
		},
		"range_d_j": {
			Matcher: lstr.Range{'d', 'j'},
			expected: map[rune]bool{
				'd': true,
				'f': true,
				'j': true,
				'c': false,
				'k': false,
				';': false,
			},
		},
		"not(range_d_j)": {
			Matcher: lstr.Not{lstr.Range{'d', 'j'}},
			expected: map[rune]bool{
				'd': false,
				'f': false,
				'j': false,
				'c': true,
				'k': true,
				';': true,
			},
		},
		"rune_a_OR_range_d_j": {
			Matcher: lstr.Or{lstr.Rune('a'), lstr.Range{'d', 'j'}},
			expected: map[rune]bool{
				'a': true,
				'd': true,
				'f': true,
				'j': true,
				'A': false,
				'c': false,
				'k': false,
				';': false,
			},
		},
		"range_e_k_AND_range_d_j": {
			Matcher: lstr.And{lstr.Range{'e', 'k'}, lstr.Range{'d', 'j'}},
			expected: map[rune]bool{
				'd': false,
				'e': true,
				'f': true,
				'g': true,
				'j': true,
				'k': false,
				'l': false,
			},
		},
		"odd_match_fn": {
			Matcher: lstr.MatcherFunc(func(r rune) bool { return int(r)%2 == 1 }),
			expected: map[rune]bool{
				'd': false,
				'e': true,
				'f': false,
				'g': true,
			},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			for r, expected := range tc.expected {
				assert.Equal(t, expected, tc.Matcher.Matches(r))
			}
		})
	}
}
