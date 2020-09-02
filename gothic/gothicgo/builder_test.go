package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestBuilder(t *testing.T) {
	i := NewImports(nil)
	c := Comment("this is a test")
	fn := MustPackageRef("foo").NewFuncRef("Foo", StringType.Named("x"))

	b, err := NewBuilder(i, c, fn, "\ntesting")
	assert.NoError(t, err)

	b.RegisterImports(i)

	str := PrefixWriteToString(b, DefaultPrefixer)
	assert.Equal(t, "import (\n\t\"foo\"\n)\n// this is a test\nfoo.Foo\ntesting", str)
}

func TestLayout(t *testing.T) {
	i := NewImports(nil)
	fn := MustPackageRef("foo").NewFuncRef("Foo", StringType.Named("x"))
	v := &PackageVar{
		NT:    fn.FuncType.Named("fn"),
		Value: fn,
	}

	l := NewLayout()
	_, err := l.Section("imports", Comment("imports section"), i)
	assert.NoError(t, err)
	_, err = l.Section("func", "// func section\n")
	assert.NoError(t, err)
	_, err = l.Section("func", v)
	assert.NoError(t, err)
	_, err = l.Section("imports", Comment("end of imports section"), i)
	assert.NoError(t, err)

	l.RegisterImports(i)

	str := PrefixWriteToString(l, DefaultPrefixer)
	assert.Equal(t, "// imports section\nimport (\n\t\"foo\"\n)\n// end of imports section\nimport (\n\t\"foo\"\n)\n// func section\nvar fn func(string) = foo.Foo", str)
}
