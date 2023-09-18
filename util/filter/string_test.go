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
