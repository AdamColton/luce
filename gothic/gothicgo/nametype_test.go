package gothicgo

import (
	"testing"

	"github.com/testify/assert"
)

func TestNameType(t *testing.T) {
	nt := NameType{"Foo", IntType}
	assert.Equal(t, "Foo", nt.Name())
	assert.Equal(t, IntType, nt.Type())

	assert.Equal(t, "Foo int", PrefixWriteToString(nt, DefaultPrefixer))
}

func TestUnnamed(t *testing.T) {
	rs := Unnamed(IntType, StringType)
	assert.Len(t, rs, 2)

	expected := []NameType{
		{"", IntType},
		{"", StringType},
	}
	assert.Equal(t, expected, rs)
}
