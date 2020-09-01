package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/util/luceio"
)

// FuncCall represents a call to a function. It implements PrefixWriterTo and
// ImportsRegistrar.
type FuncCall struct {
	Args   []string
	Caller FuncCaller
}

// PrefixWriteTo fulfills PrefixWriterTo. It writes the function call.
func (fc *FuncCall) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString(fc.Caller.Call(pre, fc.Args...))
	return sw.Rets()
}

// RegisterImports fulfills ImportsRegistrar. It calls RegisterImports on the
// Type and the Value if it implements ImportsRegistrar.
func (fc *FuncCall) RegisterImports(i *Imports) {
	if r, ok := fc.Caller.(ImportsRegistrar); ok {
		r.RegisterImports(i)
	}
}
