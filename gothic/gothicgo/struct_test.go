package gothicgo

import (
	"testing"

	"github.com/adamcolton/luce/ds/bufpool"

	"github.com/testify/assert"
)

func TestStruct(t *testing.T) {
	s := MustStructType()
	f, found := s.Field("Foo")
	assert.Nil(t, f)
	assert.False(t, found)

	fs, err := s.AddFields(
		MustPackageRef("foo").NewTypeRef("Foo", nil),
		&Field{NameType: StringType.Named("Name")},
	)
	assert.NoError(t, err)
	f, found = s.Field("Foo")
	assert.True(t, found)

	fs[0].AddTag("tagTestKey", "tagTestVal")

	_, err = s.AddFields(Comment("not a valid struct field"))
	assert.Equal(t, ErrBadField, err)

	_, err = s.AddFields(&Field{NameType: StringType.Unnamed()})
	assert.Equal(t, ErrBadFieldName, err)

	_, err = s.AddFields(&Field{NameType: StringType.Named("Foo")})
	assert.Equal(t, `Field "Foo" already exists in struct`, err.Error())

	assert.Equal(t, []string{"Foo", "Name"}, s.Fields())
	assert.Equal(t, 2, s.FieldCount())

	im := NewImports(nil)
	s.RegisterImports(im)
	str := bufpool.MustWriterToString(im)
	assert.Contains(t, str, `"foo"`)

	str = PrefixWriteToString(s, DefaultPrefixer)
	assert.Equal(t, "struct {\n\tfoo.Foo `tagTestKey:\"tagTestVal\"`\n\tName string\n}", str)

	s = MustStructType()
	str = PrefixWriteToString(s, DefaultPrefixer)
	assert.Equal(t, "struct{}", str)
}

func TestField(t *testing.T) {
	f, err := NewField(NameType{"Foo", StringType}, "this", "is", "test")
	assert.NoError(t, err)

	str := PrefixWriteToString(f, DefaultPrefixer)
	assert.Equal(t, "Foo string `test this:\"is\"`", str)

	f.AddTag("this", "more")
	f.AddTag("foo", "bar")
	str = PrefixWriteToString(f, DefaultPrefixer)
	assert.Equal(t, "Foo string `foo:\"bar\" test this:\"is;more\"`", str)

	f, err = NewField(NameType{"Foo", IntType}, "A", "B", "A", "C")
	assert.NoError(t, err)

	str = PrefixWriteToString(f, DefaultPrefixer)
	assert.Equal(t, "Foo int `A:\"B;C\"`", str)
}

func TestStructEmbed(t *testing.T) {
	s := MustStructType()
	foo := MustPackageRef("foo").NewTypeRef("Foo", nil)
	f, err := s.AddField(foo)
	assert.NoError(t, err)
	str := PrefixWriteToString(s, DefaultPrefixer)
	assert.Equal(t, "struct {\n\tfoo.Foo\n}", str)
	assert.Equal(t, "Foo", f.Name())

	s = MustStructType()
	_, err = s.AddField(foo.Pointer())
	assert.NoError(t, err)
	str = PrefixWriteToString(s, DefaultPrefixer)
	assert.Equal(t, "struct {\n\t*foo.Foo\n}", str)
	assert.Equal(t, "Foo", f.Name())

	s = MustStructType()
	i := MustPackageRef("foo").NewInterfaceRef("Bar")
	f, err = s.AddField(i)
	assert.NoError(t, err)
	str = PrefixWriteToString(s, DefaultPrefixer)
	assert.Equal(t, "struct {\n\tfoo.Bar\n}", str)
	assert.Equal(t, "Bar", f.Name())

	s = MustStructType()
	_, err = s.AddField(StringType.Pointer())
	assert.Equal(t, ErrBadField, err)
}
