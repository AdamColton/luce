package ljson_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/ljson"
	"github.com/adamcolton/luce/util/luceio"
	"github.com/stretchr/testify/assert"
)

func TestMarshalString(t *testing.T) {
	ctx := ljson.NewMarshalContext()
	wn, err := ljson.Marshal("this is a test", ctx)
	assert.NoError(t, err)
	buf, sw := luceio.BufferSumWriter()
	wctx := &ljson.WriteContext{
		SumWriter: sw,
	}
	wn(wctx)
	assert.Equal(t, `"this is a test"`, buf.String())
}
