package gothicgo

import (
	"bytes"
	"io"
)

// MemoryContext provides a Context that will not write to disk. Primarily
// intended for testing.
type MemoryContext struct {
	*BaseContext
	Files map[string]*bytes.Buffer
	Last  *bytes.Buffer
}

// NewMemoryContext sets up a MemoryContext
func NewMemoryContext() *MemoryContext {
	ctx := &MemoryContext{
		BaseContext: CtxFactory{}.New(),
		Files:       make(map[string]*bytes.Buffer),
	}

	ctx.CreateFile = func(path string) (io.Writer, error) {
		buf := ctx.Files[path]
		if buf == nil {
			buf = bytes.NewBuffer(nil)
			ctx.Files[path] = buf
		}
		ctx.Last = buf
		return buf, nil
	}

	ctx.MkdirAll = func(path string) error { return nil }
	ctx.Abs = func(path string) (string, error) { return path, nil }

	return ctx
}
