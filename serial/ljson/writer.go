package ljson

import "github.com/adamcolton/luce/util/luceio"

// WriteContext is passed into a WriteNode
type WriteContext struct {
	EscapeHtml bool
	*luceio.SumWriter
}

// WriteNode writes a node of the json document
type WriteNode func(ctx *WriteContext)
