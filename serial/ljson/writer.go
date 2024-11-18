package ljson

import (
	"bytes"
	"io"

	"github.com/adamcolton/luce/util/luceio"
)

// WriteContext is passed into a WriteNode
type WriteContext struct {
	EscapeHtml bool
	*luceio.SumWriter
}

// WriteNode writes a node of the json document
type WriteNode func(ctx *WriteContext)

// String invokes the WriteNode and returns the data written as a string.
func (wn WriteNode) String() string {
	buf := bytes.NewBuffer(nil)
	wn.WriteTo(buf)
	return buf.String()
}

// WriteTo fulfills io.WriterTo and writes the WriteNode to the Writer.
func (wn WriteNode) WriteTo(w io.Writer) (int64, error) {
	wctx := &WriteContext{
		SumWriter: luceio.NewSumWriter(w),
	}
	wn(wctx)
	return wctx.Rets()
}

// Stringify marshals the value given and returns a json string.
func Stringify[T any](v T, ctx *MarshalContext) (string, error) {
	wn, err := Marshal(v, ctx)
	if err != nil {
		return "", err
	}
	return wn.String(), nil
}
