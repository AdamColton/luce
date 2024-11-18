package ljson_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/stretchr/testify/assert"
)

func TestOmitEmptyPanics(t *testing.T) {
	ctx := ljson.NewMarshalContext(false)

	// A name that is not in the StructKeys.
	keys := ljson.GetFieldKeys[fmPerson]()
	assert.Panics(t, func() { ctx.TypesContext.OmitEmpty(keys, "Nope") })

	// A FieldKey for a field the type does not have.
	bad := ljson.StructKeys{"Nope": {Type: reflect.TypeOf(fmPerson{}), Name: "Nope"}}
	assert.Panics(t, func() { ctx.TypesContext.OmitEmpty(bad, "Nope") })
}
