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
		"elem-string-true": {
			expected: true,
			f:        filter.IsKind(reflect.String).Elem(),
			v:        []string{},
		},
		"cannot-elem-false": {
			expected: false,
			f:        filter.IsKind(reflect.String).Elem(),
			v:        123,
		},
		"numIn-true": {
			expected: true,
			f:        filter.NumIn(filter.EQ(3)),
			v:        func(a, b, c int) {},
		},
		"numIn-false": {
			expected: false,
			f:        filter.NumIn(filter.EQ(4)),
			v:        func(a, b, c int) {},
		},
		"in-true": {
			expected: true,
			f:        filter.IsType(ltype.Int).In(1),
			v:        func(a, b, c int) {},
		},
		"in-negative-true": {
			expected: true,
			f:        filter.IsType(ltype.String).In(-1),
			v:        func(a, b int, c string) {},
		},
		"in-false": {
			expected: false,
			f:        filter.IsType(ltype.String).In(1),
			v:        func(a, b, c int) {},
		},
		"numOut-true": {
			expected: true,
			f:        filter.NumOut(filter.EQ(3)),
			v:        func() (a, b, c int) { return 1, 2, 3 },
		},
		"numOut-false": {
			expected: false,
			f:        filter.NumOut(filter.EQ(4)),
			v:        func() (a, b, c int) { return 1, 2, 3 },
		},
		"out-true": {
			expected: true,
			f:        filter.IsType(ltype.Int).Out(1),
			v:        func() (a, b, c int) { return 1, 2, 3 },
		},
		"out-negative-true": {
			expected: true,
			f:        filter.IsType(ltype.String).Out(-1),
			v:        func() (a, b int, c string) { return 1, 2, "3" },
		},
		"out-false": {
			expected: false,
			f:        filter.IsType(ltype.String).Out(1),
			v:        func() (a, b, c int) { return 1, 2, 3 },
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.f.OnInterface(tc.v))
		})
	}
}
