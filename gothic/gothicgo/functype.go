package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/util/luceio"
)

// FuncType represents a function as a type. It is the function signature
// prefixed with the "func" literal.
type FuncType struct {
	*FuncSig
}

// NewFuncType returns a FuncType.
func NewFuncType(name string, args ...NameType) *FuncType {
	return &FuncType{NewFuncSig(name, args...)}
}

// PrefixWriteTo fulfils PrefixWriterTo. Writes the function type.
func (f *FuncType) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString("func")
	if f.Name != "" {
		sw.WriteRune(' ')
	}
	sumPrefixWriteTo(sw, pre, f.FuncSig)
	return sw.Rets()
}

// Returns sets the returns on the FuncType.
func (f *FuncType) Returns(rets ...NameType) *FuncType {
	f.FuncSig.Returns(rets...)
	return f
}

// UnnamedRets sets the returns on the FuncType to unnamed values.
func (f *FuncType) UnnamedRets(rets ...Type) *FuncType {
	f.FuncSig.UnnamedRets(rets...)
	return f

}
