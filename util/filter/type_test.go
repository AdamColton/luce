package filter

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

func TestType(t *testing.T) {
	tt := map[string]struct {
		expected bool
		f        Filter[reflect.Type]
		reflect.Type
	}{
		"string-true": {
			expected: true,
			f:        IsType(""),
			Type:     reflect.TypeOf("foo"),
		},
		"string-false": {
			expected: false,
			f:        IsType(""),
			Type:     reflect.TypeOf(123),
		},
		"string-type-true": {
			expected: true,
			f:        IsType(reflect.TypeOf("")),
			Type:     reflect.TypeOf("foo"),
		},
		"string-kind-true": {
			expected: true,
			f:        IsKind(reflect.String),
			Type:     reflect.TypeOf("foo"),
		},
		"elem-string-true": {
			expected: true,
			f:        Elem(IsKind(reflect.String)),
			Type:     reflect.TypeOf(([]string)(nil)),
		},
		"cannot-elem-false": {
			expected: false,
			f:        Elem(IsKind(reflect.String)),
			Type:     reflect.TypeOf(123),
		},
		"isKind-nil-false": {
			expected: false,
			f:        IsKind(reflect.String),
			Type:     nil,
		},
		"elem-nil-false": {
			expected: false,
			f:        Elem(IsKind(reflect.String)),
			Type:     nil,
		},
		"isType-nil-false": {
			expected: false,
			f:        IsType(""),
			Type:     nil,
		},
		"numIn-true": {
			expected: true,
			f:        NumIn(EQ(3)),
			Type:     reflect.TypeOf(func(a, b, c int) {}),
		},
		"numIn-false": {
			expected: false,
			f:        NumIn(EQ(4)),
			Type:     reflect.TypeOf(func(a, b, c int) {}),
		},
		"in-true": {
			expected: true,
			f:        In(1, IsNilRef((*int)(nil))),
			Type:     reflect.TypeOf(func(a, b, c int) {}),
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.f(tc.Type))
		})
	}
}

func TestTypeChecker(t *testing.T) {
	expectedErr := lerr.Str("expected string, got int")
	c := TypeCheck(IsType(""), func(t reflect.Type) error {
		return lerr.Str("expected string, got " + t.String())
	})

	ct, err := c("test")
	assert.NoError(t, err)
	assert.Equal(t, reflect.TypeOf(""), ct)

	ct, err = c(123)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, reflect.TypeOf(456), ct)

	defer func() {
		assert.Equal(t, expectedErr, recover())
	}()

	ct = c.Panic("test")
	assert.Equal(t, reflect.TypeOf(""), ct)

	c.Panic(123)
}
