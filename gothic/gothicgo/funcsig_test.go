package gothicgo

import (
	"testing"

	"github.com/adamcolton/luce/ds/bufpool"

	"github.com/testify/assert"
)

func TestNameTypeSliceToString(t *testing.T) {
	pkg := MustExternalPackageRef("foo")
	tt := map[string]struct {
		nts      []NameType
		variadic bool
		expected string
	}{
		"empty": {
			expected: "",
		},
		"named-basic": {
			nts: []NameType{
				{"foo", IntType},
				{"bar", StringType},
			},
			expected: "foo int, bar string",
		},
		"unnamed-basic": {
			nts: []NameType{
				IntType.Unnamed(),
				StringType.Unnamed(),
			},
			expected: "int, string",
		},
		"named-repeat-type": {
			nts: []NameType{
				{"foo", IntType},
				{"bar", IntType},
			},
			expected: "foo, bar int",
		},
		"unnamed-repeat-type": {
			nts: []NameType{
				IntType.Unnamed(),
				IntType.Unnamed(),
			},
			expected: "int, int",
		},
		"named-double-repeat-type": {
			nts: []NameType{
				{"a", IntType},
				{"b", IntType},
				{"c", StringType},
				{"d", StringType},
			},
			expected: "a, b int, c, d string",
		},
		"unnamed-double-repeat-type": {
			nts: []NameType{
				IntType.Unnamed(),
				IntType.Unnamed(),
				StringType.Unnamed(),
				StringType.Unnamed(),
			},
			expected: "int, int, string, string",
		},
		"named-double-repeat-external-type": {
			nts: []NameType{
				{"a", pkg.MustExternalType("Bar")},
				{"b", pkg.MustExternalType("Bar")},
				{"c", pkg.MustExternalType("Baz")},
				{"d", pkg.MustExternalType("Baz")},
			},
			expected: "a, b foo.Bar, c, d foo.Baz",
		},
		"unnamed-double-repeat-external-type": {
			nts: []NameType{
				pkg.MustExternalType("Bar").Unnamed(),
				pkg.MustExternalType("Bar").Unnamed(),
				pkg.MustExternalType("Baz").Unnamed(),
				pkg.MustExternalType("Baz").Unnamed(),
			},
			expected: "foo.Bar, foo.Bar, foo.Baz, foo.Baz",
		},
		"named-variadic": {
			nts: []NameType{
				{"foo", IntType},
				{"bar", StringType},
			},
			variadic: true,
			expected: "foo int, bar ...string",
		},
		"unnamed-variadic": {
			nts: []NameType{
				IntType.Unnamed(),
				StringType.Unnamed(),
			},
			variadic: true,
			expected: "int, ...string",
		},
		"named-variadic-repeat-type": {
			nts: []NameType{
				{"foo", IntType},
				{"bar", IntType},
				{"baz", IntType},
			},
			variadic: true,
			expected: "foo, bar int, baz ...int",
		},
		"unnamed-variadic-repeat-type": {
			nts: []NameType{
				IntType.Unnamed(),
				IntType.Unnamed(),
				IntType.Unnamed(),
			},
			variadic: true,
			expected: "int, int, ...int",
		},
		"named-double-repeat-type-variadic": {
			nts: []NameType{
				{"a", IntType},
				{"b", IntType},
				{"c", StringType},
				{"d", StringType},
				{"e", StringType},
			},
			variadic: true,
			expected: "a, b int, c, d string, e ...string",
		},
		"unnamed-double-repeat-type-variadic": {
			nts: []NameType{
				IntType.Unnamed(),
				IntType.Unnamed(),
				StringType.Unnamed(),
				StringType.Unnamed(),
				StringType.Unnamed(),
			},
			variadic: true,
			expected: "int, int, string, string, ...string",
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			str, err := nameTypeSliceToString(DefaultPrefixer, tc.nts, tc.variadic)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, str)
		})
	}

	_, err := nameTypeSliceToString(DefaultPrefixer, []NameType{IntType.Named("foo"), IntType.Unnamed()}, false)
	assert.Equal(t, ErrMixedParameters, err)
	_, err = nameTypeSliceToString(DefaultPrefixer, []NameType{IntType.Unnamed(), IntType.Named("foo")}, false)
	assert.Equal(t, ErrMixedParameters, err)
}

func TestFuncSig(t *testing.T) {
	args := []NameType{
		IntType.Named("a"),
		StringType.Named("b"),
	}
	fs := NewFuncSig("Foo", args...).
		UnnamedRets(StringType)
	str := PrefixWriteToString(fs, DefaultPrefixer)
	assert.Equal(t, "func Foo(a int, b string) string", str)

	fs = NewFuncSig("Foo", args...).
		Returns(StringType.Named("x"), StringType.Named("y"))
	fs.Variadic = true
	str = PrefixWriteToString(fs, DefaultPrefixer)
	assert.Equal(t, "func Foo(a int, b ...string) (x, y string)", str)

	fs = fs.AsType(false)
	str = PrefixWriteToString(fs, DefaultPrefixer)
	assert.Equal(t, "func Foo(int, ...string) (string, string)", str)

	pkg1 := MustExternalPackageRef("pkg1")
	pkg2 := MustExternalPackageRef("pkg2")
	fs = NewFuncSig("Foo2Foo", pkg1.MustExternalType("Foo").Unnamed()).
		UnnamedRets(pkg2.MustExternalType("Foo"))
	i := NewImports(nil)
	fs.RegisterImports(i)
	str = bufpool.MustWriterToString(i)
	assert.Contains(t, str, "pkg1")
	assert.Contains(t, str, "pkg2")

	str = PrefixWriteToString(fs, DefaultPrefixer)
	assert.Contains(t, str, "func Foo2Foo(pkg1.Foo) pkg2.Foo")

	args = []NameType{
		IntType.Named("a"),
		StringType.Named("b"),
		IntType.Named("c"),
		StringType.Named("d"),
	}
	fs = NewFuncSig("Foo", args...).
		UnnamedRets(StringType)
	str = PrefixWriteToString(fs, DefaultPrefixer)
	assert.Equal(t, "func Foo(a int, b string, c int, d string) string", str)

}
