package ljson_test

import (
	"reflect"
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/stretchr/testify/assert"
)

// asMap lets the tests compare a result whether it is a plain map or a named
// map type.
func asMap(m map[string]reflect.Type) map[string]reflect.Type { return m }

type exportPlain struct {
	A int
	B string
}

func TestExportType(t *testing.T) {
	ctx := ljson.NewMarshalContext("")
	expected := map[string]reflect.Type{
		"A": reflector.Type[int](),
		"B": reflector.Type[string](),
	}

	got, err := ljson.ExportType(reflect.TypeOf(exportPlain{}), ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, asMap(got))

	got, err = ljson.ExportType(reflect.TypeOf(&exportPlain{}), ctx)
	assert.NoError(t, err, "a pointer to a struct is followed")
	assert.Equal(t, expected, asMap(got))

	_, err = ljson.ExportType(reflect.TypeOf(1), ctx)
	assert.Error(t, err, "only structs can be exported")
}

func TestExportConvert(t *testing.T) {
	ctx := ljson.NewMarshalContext("")
	_, err := ljson.Stringify(exportPlain{A: 1, B: "b"}, ctx)
	assert.NoError(t, err)

	// A struct that is converted is not written as a struct, but its fields can
	// still be exported, also after it was marshaled.
	ljson.Convert(func(p exportPlain, ctx *ljson.MarshalContext[string]) string {
		return p.B
	}, ctx.TypesContext)
	got, err := ljson.Export[exportPlain](ctx)
	assert.NoError(t, err)
	assert.Equal(t, map[string]reflect.Type{
		"A": reflector.Type[int](),
		"B": reflector.Type[string](),
	}, asMap(got))
}

type exportNode struct {
	Val  int
	Next *exportNode
}

func TestExportRecursive(t *testing.T) {
	got, err := ljson.Export[exportNode](ljson.NewMarshalContext(""))
	assert.NoError(t, err)
	assert.Equal(t, map[string]reflect.Type{
		"Val":  reflector.Type[int](),
		"Next": reflector.Type[*exportNode](),
	}, asMap(got))
}

type exportPet struct {
	Name string
}

type exportAddress struct {
	Street string
}

type exportPerson struct {
	Name    string
	Home    exportAddress
	Work    *exportAddress
	Pets    map[string]exportPet
	Friends []exportPerson
	Keys    [2]int
}

func TestExportAll(t *testing.T) {
	got, err := ljson.ExportAll[exportPerson](ljson.NewMarshalContext(""))
	assert.NoError(t, err)

	flat := make(map[reflect.Type]map[string]reflect.Type, len(got))
	for t, fields := range got {
		flat[t] = asMap(fields)
	}

	expected := map[reflect.Type]map[string]reflect.Type{
		reflector.Type[exportAddress](): {
			"Street": reflector.Type[string](),
		},
		reflector.Type[exportPet](): {
			"Name": reflector.Type[string](),
		},
		reflector.Type[exportPerson](): {
			"Name":    reflector.Type[string](),
			"Home":    reflector.Type[exportAddress](),
			"Work":    reflector.Type[*exportAddress](),
			"Pets":    reflector.Type[map[string]exportPet](),
			"Friends": reflector.Type[[]exportPerson](),
			"Keys":    reflector.Type[[2]int](),
		},
	}
	assert.Equal(t, expected, flat)
}
