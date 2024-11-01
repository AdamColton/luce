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
