package filter_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/filter"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/reflector/ltype"
	"github.com/stretchr/testify/assert"
)

type fooer interface {
	foo()
}

type foo struct{}

func (foo) foo() {}

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
			f:        filter.NumInEq(3),
			v:        func(a, b, c int) {},
		},
		"numIn-false": {
			expected: false,
			f:        filter.NumInEq(4),
			v:        func(a, b, c int) {},
		},
		"in-true": {
			expected: true,
			f:        filter.InType(1, ltype.Int),
			v:        func(a, b, c int) {},
		},
		"in-negative-true": {
			expected: true,
			f:        filter.InType(-1, ltype.String),
			v:        func(a, b int, c string) {},
		},
		"in-false": {
			expected: false,
			f:        filter.InType(1, ltype.String),
			v:        func(a, b, c int) {},
		},
		"numOut-true": {
			expected: true,
			f:        filter.NumOutEq(3),
			v:        func() (a, b, c int) { return 1, 2, 3 },
		},
		"numOut-false": {
			expected: false,
			f:        filter.NumOutEq(4),
			v:        func() (a, b, c int) { return 1, 2, 3 },
		},
		"out-true": {
			expected: true,
			f:        filter.OutType(1, ltype.Int),
			v:        func() (a, b, c int) { return 1, 2, 3 },
		},
		"out-negative-true": {
			expected: true,
			f:        filter.OutType(-1, ltype.String),
			v:        func() (a, b int, c string) { return 1, 2, "3" },
		},
		"out-false": {
			expected: false,
			f:        filter.OutType(1, ltype.String),
			v:        func() (a, b, c int) { return 1, 2, 3 },
		},
		"and": {
			expected: true,
			f: filter.InType(0, ltype.Int).
				And(filter.InType(0, ltype.Int)),
			v: func(a, b, c int) (d, e, f int) { return 1, 2, 3 },
		},
		"or": {
			expected: true,
			f: filter.InType(0, ltype.String).
				Or(filter.InType(0, ltype.Int)),
			v: func(a, b, c int) (d, e, f int) { return 1, 2, 3 },
		},
		"not": {
			expected: true,
			f:        filter.InType(0, ltype.String).Not(),
			v:        func(a, b, c int) (d, e, f int) { return 1, 2, 3 },
		},
		"implements": {
			expected: true,
			f:        filter.Implements[fooer](),
			v:        foo{},
		},
		"func": {
			expected: true,
			f: filter.Func(
				[]any{ltype.Int, filter.Implements[fooer](), filter.IsType(ltype.String).Filter},
				[]any{ltype.Int, filter.IsType(ltype.Int), filter.IsType(ltype.Int).Filter},
			),
			v: func(a int, b foo, c string) (d, e, f int) { return 1, 2, 3 },
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.f.OnInterface(tc.v))
		})
	}
}

func TestTypeChecker(t *testing.T) {
	expectedErr := lerr.Str("expected string, got int")

	errFn := func(t reflect.Type) error {
		return lerr.Str("expected string, got " + t.String())
	}

	c := filter.IsType(ltype.String).Check(errFn)

	ct, err := c("test")
	assert.NoError(t, err)
	assert.Equal(t, ltype.String, ct)

	ct, err = c(123)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, ltype.Int, ct)

	defer func() {
		assert.Equal(t, expectedErr, recover())
	}()

	ct = lerr.Must(c("test"))
	assert.Equal(t, ltype.String, ct)

	lerr.Must(c(123))
}

func TestMethodName(t *testing.T) {
	f := filter.MethodName(filter.Prefix("Err"))
	assert.True(t, f(reflector.MethodOn(t, "Error")))
	assert.True(t, f(reflector.MethodOn(t, "Errorf")))
	assert.False(t, f(reflector.MethodOn(t, "Log")))
}

func TestMethodFirst(t *testing.T) {
	// TODO: use something other than t
	ms := reflector.MethodsOn(t)
	f, _ := filter.NumOut(filter.EQ(2)).Method().First(ms...)
	assert.Equal(t, "Deadline", f.Name)
}

func TestAnyType(t *testing.T) {
	a := filter.AnyType()
	assert.True(t, a.OnInterface("string"))
}

func TestTypeErr(t *testing.T) {
	fn := filter.TypeErr("Foo: %v")
	got := fn(ltype.String)
	assert.Equal(t, "Foo: string", got.Error())
}
