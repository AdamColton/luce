package ljson_test

import (
	"errors"
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/stretchr/testify/assert"
)

type fmPerson struct {
	Name string
	Age  int
	b    int
}

func TestGetFieldKeys(t *testing.T) {
	keys := ljson.GetFieldKeys[fmPerson]()
	assert.Len(t, keys, 2)
	assert.Contains(t, keys, "Name")
	assert.Contains(t, keys, "Age")

	// Only a struct has field keys.
	assert.Nil(t, ljson.GetFieldKeys[int]())
}

func TestAddFieldMarshalPanics(t *testing.T) {
	ctx := ljson.NewMarshalContext(false)
	keys := ljson.GetFieldKeys[fmPerson]()
	fm := func(name string, v int, ctx *ljson.MarshalContext[bool]) (string, ljson.WriteNode, error) {
		return name, nil, nil
	}

	// The key is not a field of its type.
	missing := keys["Age"]
	missing.Name = "Missing"
	assert.Panics(t, func() { ljson.AddFieldMarshal[int](missing, fm, ctx.TypesContext) })

	// The type of the field is a string, not an int.
	assert.Panics(t, func() { ljson.AddFieldMarshal[int](keys["Name"], fm, ctx.TypesContext) })
}

func TestAddFieldMarshalNilOmits(t *testing.T) {
	ctx := ljson.NewMarshalContext(false)
	keys := ljson.GetFieldKeys[fmPerson]()
	ljson.AddFieldMarshal[int, bool](keys["Age"], nil, ctx.TypesContext)

	str, err := ljson.Stringify(fmPerson{Name: "Adam", Age: 40}, ctx)
	assert.NoError(t, err)
	assert.Equal(t, `{"Name":"Adam"}`, str)
}

func TestAddFieldMarshalError(t *testing.T) {
	ctx := ljson.NewMarshalContext(false)
	keys := ljson.GetFieldKeys[fmPerson]()
	errBoom := errors.New("boom")
	fm := func(name string, v int, ctx *ljson.MarshalContext[bool]) (string, ljson.WriteNode, error) {
		return "", nil, errBoom
	}
	ljson.AddFieldMarshal[int](keys["Age"], fm, ctx.TypesContext)

	_, err := ljson.Stringify(fmPerson{Name: "Adam", Age: 40}, ctx)
	assert.Equal(t, errBoom, err)
}

func TestOmitFieldsUnknown(t *testing.T) {
	ctx := ljson.NewMarshalContext(false)
	keys := ljson.GetFieldKeys[fmPerson]()
	ctx.TypesContext.OmitFields(keys, "Nope", "Age")

	str, err := ljson.Stringify(fmPerson{Name: "Adam", Age: 40}, ctx)
	assert.NoError(t, err)
	assert.Equal(t, `{"Name":"Adam"}`, str)
}
