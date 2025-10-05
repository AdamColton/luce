package ljson

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// nameType is what Export reads the fields from. It is unexported and Export
// does not exist yet here, so the implementations are tested directly.
func TestNameType(t *testing.T) {
	type person struct {
		Name string
		Age  int
	}
	pt := reflect.TypeOf(person{})
	strT, intT := reflect.TypeOf(""), reflect.TypeOf(0)

	n, typ := FieldKey{Type: pt, Name: "Age"}.nameType()
	assert.Equal(t, "Age", n)
	assert.Equal(t, intT, typ)

	n, typ = FieldKey{Type: pt, Name: "Missing"}.nameType()
	assert.Empty(t, n)
	assert.Nil(t, typ)

	n, typ = deferGetFieldMarshal[string]{f: pt.Field(0), t: pt}.nameType()
	assert.Equal(t, "Name", n)
	assert.Equal(t, strT, typ)

	n, typ = marshalValToField[string]{name: "Nick", t: strT}.nameType()
	assert.Equal(t, "Nick", n)
	assert.Equal(t, strT, typ)

	n, typ = marshalFieldGenerator[person, string, string]{name: "Foo"}.nameType()
	assert.Equal(t, "Foo", n)
	assert.Equal(t, strT, typ)
}
