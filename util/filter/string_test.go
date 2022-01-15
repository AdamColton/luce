package filter_test

import (
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/stretchr/testify/assert"
)

func TestPrefix(t *testing.T) {
	tt := map[string]bool{
		"test":    true,
		"testing": true,
		"atest":   false,
		"abc":     false,
	}

	p := filter.Prefix("test")
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc, p(n))
		})
	}
}

func TestSuffix(t *testing.T) {
	tt := map[string]bool{
		"test":    true,
		"testing": false,
		"atest":   true,
		"abc":     false,
	}

	s := filter.Suffix("test")
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc, s(n))
		})
	}
}

func TestContains(t *testing.T) {
	tt := map[string]bool{
		"test":     true,
		"testing":  true,
		"atesting": true,
		"atsting":  false,
		"atest":    true,
		"abc":      false,
	}

	c := filter.Contains("test")
	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc, c(n))
		})
	}
}

func TestRegex(t *testing.T) {
	tt := map[string]map[string]bool{
		"ca*t": {
			"cat":         true,
			"ct":          true,
			"cot":         false,
			"acat":        true,
			"dogcatmouse": true,
		},
		"^ca*t$": {
			"cat":         true,
			"ct":          true,
			"cot":         false,
			"acat":        false,
			"dogcatmouse": false,
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			r := filter.MustRegex(n)
			for s, expected := range tc {
				assert.Equal(t, expected, r(s))
			}
			var err error
			r, err = filter.Regex(n)
			assert.NoError(t, err)
			for s, expected := range tc {
				assert.Equal(t, expected, r(s))
			}
		})
	}

	_, err := filter.Regex("bad [ regex")
	assert.Error(t, err)
}
