package ljson_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/adamcolton/luce/util/luceio"
	"github.com/stretchr/testify/assert"
)

func TestStringWriter(t *testing.T) {
	strWriter := ljson.MarshalString("this is a test")
	buf, sw := luceio.BufferSumWriter()
	wctx := &ljson.WriteContext{
		SumWriter: sw,
	}
	strWriter(wctx)
	assert.Equal(t, `"this is a test"`, buf.String())
}
