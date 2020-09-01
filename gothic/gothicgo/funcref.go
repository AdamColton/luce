package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/util/luceio"
)

// FuncRef represents a function as a type.
type FuncRef struct {
	FuncType *FuncType
	Pkg      PackageRef
}

// NewFuncRef creates a FuncRef representing a Func as a Type.
func NewFuncRef(pkg PackageRef, name string, args ...NameType) *FuncRef {
	return &FuncRef{
		FuncType: NewFuncType(name, args...),
		Pkg:      pkg,
	}
}

// Call produces a invocation of the function and fulfills the FuncCaller
// interface
func (f *FuncRef) Call(pre Prefixer, args ...string) string {
	return funcCall(pre, f.FuncType.FuncSig.Name, args, f.Pkg)
}

// PrefixWriteTo fulfills PrefixWriterTo. Writes the package and func name.
func (f *FuncRef) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteStrings(pre.Prefix(f.Pkg), f.FuncType.FuncSig.Name)
	return sw.Rets()
}

// RegisterImports fulfills ImportsRegistrar. It calls RegisterImports on the
// Type and the Value if it implements ImportsRegistrar.
func (f *FuncRef) RegisterImports(i *Imports) {
	i.Add(f.Pkg)
}

// Caller returns a FuncCall to this function.
func (f *FuncRef) Caller(args ...string) *FuncCall {
	return &FuncCall{
		Args:   args,
		Caller: f,
	}
}
