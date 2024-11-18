package ljson_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/stretchr/testify/assert"
)

func TestMarshalString(t *testing.T) {
	ctx := ljson.NewMarshalContext("hello")

	str, err := ljson.Stringify("this is a test", ctx)
	assert.NoError(t, err)
	assert.Equal(t, `"this is a test"`, str)
}

type ctxStr string

func contextString(cs ctxStr, ctx *ljson.MarshalContext[bool]) (ljson.WriteNode, error) {
	str := ":foo"
	if ctx.Context {
		str = ":bar"
	}
	str = string(cs) + str
	return func(ctx *ljson.WriteContext) {
		ctx.WriteString(str)
	}, nil
}

func TestMarshalContext(t *testing.T) {
	ctx := ljson.NewMarshalContext(false)
	ljson.AddMarshaler(contextString, ctx.TypesContext)

	str, err := ljson.Stringify(ctxStr("test"), ctx)
	assert.NoError(t, err)
	assert.Equal(t, "test:foo", str)

	ctx.Context = true
	str, err = ljson.Stringify(ctxStr("baz"), ctx)
	assert.NoError(t, err)
	assert.Equal(t, "baz:bar", str)
}

func TestStruct(t *testing.T) {
	type Person struct {
		Name string
		Role string
	}
	p := Person{
		Name: "Adam",
		Role: "admin",
	}
	ctx := ljson.NewMarshalContext(false)
	ctx.Sort = true
	str, err := ljson.Stringify(p, ctx)
	assert.NoError(t, err)
	assert.Equal(t, `{"Name":"Adam","Role":"admin"}`, str)
}

func check[T any](t *testing.T, ctx *ljson.MarshalContext[bool], v T, expected string) {
	str, err := ljson.Stringify(v, ctx)
	assert.NoError(t, err)
	assert.Equal(t, expected, str)
}
func TestBuiltin(t *testing.T) {
	ctx := ljson.NewMarshalContext(false)
	check(t, ctx, int(123), "123")
	check(t, ctx, int8(123), "123")
	check(t, ctx, int16(123), "123")
	check(t, ctx, int32(123), "123")
	check(t, ctx, int64(123), "123")
	check(t, ctx, uint(123), "123")
	check(t, ctx, uint8(123), "123")
	check(t, ctx, uint16(123), "123")
	check(t, ctx, uint32(123), "123")
	check(t, ctx, uint64(123), "123")
	check(t, ctx, float32(123), "123")
	check(t, ctx, float64(123), "123")
	check(t, ctx, true, "true")
}

func TestStructPtr(t *testing.T) {
	type Person struct {
		ID   uint
		Name string
		Age  int
		Role string
	}
	p := &Person{
		ID:   123,
		Name: "Adam",
		Age:  40,
		Role: "admin",
	}
	ctx := ljson.NewMarshalContext(false)
	ctx.Sort = true
	str, err := ljson.Stringify(p, ctx)
	assert.NoError(t, err)
	assert.Equal(t, `{"Age":40,"ID":123,"Name":"Adam","Role":"admin"}`, str)
}

func TestSlice(t *testing.T) {
	ctx := ljson.NewMarshalContext(false)
	str, err := ljson.Stringify([]string{"a", "b", "c"}, ctx)
	assert.NoError(t, err)
	assert.Equal(t, `["a","b","c"]`, str)
}

func TestMarshalMap(t *testing.T) {
	m := map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
	}
	ctx := ljson.NewMarshalContext(false)
	ctx.Sort = true
	got, err := ljson.Stringify(m, ctx)
	assert.NoError(t, err)
	assert.Equal(t, `{"1":1,"2":2,"3":3}`, got)
}

func TestMarshalMapOfPtr(t *testing.T) {
	type Person struct {
		Name string
		Role string
	}
	m := map[int]*Person{
		1: {
			Name: "Adam",
			Role: "admin",
		},
		2: {
			Name: "Fletcher",
			Role: "user",
		},
	}
	ctx := ljson.NewMarshalContext(false)
	ctx.Sort = true
	str, err := ljson.Stringify(m, ctx)
	assert.NoError(t, err)
	assert.Equal(t, `{1:{"Name":"Adam","Role":"admin"},2:{"Name":"Fletcher","Role":"user"}}`, str)
}

type A struct {
	Name string
	B    *B
}
type B struct {
	Name string
	C    *C
}
type C struct {
	Name string
	A    *A
}

func TestCircularRefErr(t *testing.T) {
	a := A{Name: "A"}
	b := B{Name: "B"}
	c := C{Name: "C", A: &a}
	b.C = &c
	a.B = &b
	ctx := ljson.NewMarshalContext(false)
	_, err := ljson.Stringify(a, ctx)
	assert.Error(t, err)
}

func TestFieldOptions(t *testing.T) {
	type Person struct {
		ID   int
		Name string
		Age  int
		Role string
	}
	p := &Person{
		Name: "Adam",
		Age:  40,
		Role: "admin",
	}
	ctx := ljson.NewMarshalContext(false)
	ctx.Sort = true
	keys := ljson.GetFieldKeys[Person]()
	fm := func(name string, v int, ctx *ljson.MarshalContext[bool]) (string, ljson.WriteNode, error) {
		if !ctx.Context {
			return "", nil, nil
		}
		wn, err := ljson.Marshal(v, ctx)
		return name, wn, err
	}
	ljson.AddFieldMarshal[int](keys["Age"], fm, ctx.TypesContext)
	ctx.TypesContext.OmitFields(keys, "ID")
	str, err := ljson.Stringify(p, ctx)
	assert.NoError(t, err)
	assert.Equal(t, `{"Name":"Adam","Role":"admin"}`, str)
	ctx.Context = true
	str, err = ljson.Stringify(p, ctx)
	assert.NoError(t, err)
	assert.Equal(t, `{"Age":40,"Name":"Adam","Role":"admin"}`, str)
}
