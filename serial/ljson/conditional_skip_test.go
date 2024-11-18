package ljson_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/stretchr/testify/assert"
)

func TestConditionalFieldsSkipsUnknown(t *testing.T) {
	ctx := ljson.NewMarshalContext(false)
	ctx.Sort = true
	keys := ljson.GetFieldKeys[fmPerson]()
	keys["Bogus"] = ljson.FieldKey{Type: reflect.TypeOf(fmPerson{}), Name: "Bogus"}
	never := func(ctx *ljson.MarshalContext[bool]) bool { return false }

	// "Nope" is not in the keys and "Bogus" is not a field of the type.
	ctx.TypesContext.ConditionalFields(never, keys, "Nope", "Bogus", "Age")
	str, err := ljson.Stringify(fmPerson{Name: "Adam", Age: 40}, ctx)
	assert.NoError(t, err)
	assert.Equal(t, `{"Name":"Adam"}`, str)
}

func TestConditionalFieldsKeepsFieldMarshaler(t *testing.T) {
	ctx := ljson.NewMarshalContext(true)
	ctx.Sort = true
	keys := ljson.GetFieldKeys[fmPerson]()
	fm := func(name string, v int, ctx *ljson.MarshalContext[bool]) (string, ljson.WriteNode, error) {
		wn, err := ljson.Marshal(v*2, ctx)
		return "Double", wn, err
	}
	ljson.AddFieldMarshal[int](keys["Age"], fm, ctx.TypesContext)
	always := func(ctx *ljson.MarshalContext[bool]) bool { return ctx.Context }
	ctx.TypesContext.ConditionalFields(always, keys, "Age")

	str, err := ljson.Stringify(fmPerson{Name: "Adam", Age: 40}, ctx)
	assert.NoError(t, err)
	assert.Equal(t, `{"Double":80,"Name":"Adam"}`, str)
}
