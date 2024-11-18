package ljson

import "github.com/adamcolton/luce/util/luceio"

// WriteContext is passed into a WriteNode. EscapeHtml makes strings escape <, >
// and &, and the SumWriter is where the json is written.
type WriteContext struct {
	EscapeHtml bool
	*luceio.SumWriter
}

// WriteNode writes a node of the json document. A WriteNode reports an error by
// setting the Err of its SumWriter.
type WriteNode func(ctx *WriteContext)
