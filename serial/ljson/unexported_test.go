package ljson_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/stretchr/testify/assert"
)

type mixed struct {
	A string
	b string
}

func TestUnexportedFieldsSkipped(t *testing.T) {
	ctx := ljson.NewMarshalContext(false)
	str, err := ljson.Stringify(mixed{A: "x", b: "y"}, ctx)
	assert.NoError(t, err)
	assert.Equal(t, `{"A":"x"}`, str)
}

func TestFieldKeyField(t *testing.T) {
	fk := ljson.FieldKey{Type: reflect.TypeOf(mixed{}), Name: "A"}
	f, found := fk.Field()
	assert.True(t, found)
	assert.Equal(t, "A", f.Name)

	fk.Name = "Missing"
	_, found = fk.Field()
	assert.False(t, found)
}
