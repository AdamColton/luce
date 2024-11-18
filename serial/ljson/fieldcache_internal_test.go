package ljson

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type cachedFields struct {
	A string
	B string
}

// The field marshalers are normally only set by the options that come later
// (AddFieldMarshal, OmitFields), so these tests set them directly.
func TestFieldMarshalersSetDirectly(t *testing.T) {
	tctx := NewTypesContext[bool]()
	typ := reflect.TypeOf(cachedFields{})

	// A field marshaler that is already set is used.
	tctx.fieldMarshalers[FieldKey{Type: typ, Name: "A"}] = func(name string, v reflect.Value, ctx *MarshalContext[bool]) (string, WriteNode) {
		wn, _ := MarshalString("X", ctx)
		return "Renamed", wn
	}
	// A nil field marshaler omits the field.
	tctx.fieldMarshalers[FieldKey{Type: typ, Name: "B"}] = nil

	ctx := tctx.NewMarshalContext(false)
	str, err := Stringify(cachedFields{A: "a", B: "b"}, ctx)
	assert.NoError(t, err)
	assert.Equal(t, `{"Renamed":"X"}`, str)
}
