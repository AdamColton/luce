package filter_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	tt := map[string]struct {
		expected bool
		f        filter.Type
		v        any
	}{
		"string-true": {
			expected: true,
			f:        filter.IsType(ltype.String),
			v:        "foo",
		},
		"string-false": {
			expected: false,
			f:        filter.IsType(ltype.String),
			v:        123,
		},
		"string-type-true": {
			expected: true,
			f:        filter.IsType(ltype.String),
			v:        "foo",
		},
		"string-kind-true": {
			expected: true,
			f:        filter.IsKind(reflect.String),
			v:        "foo",
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.f.OnInterface(tc.v))
		})
	}
}
