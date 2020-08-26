package gothicgo

import (
	"io"
)

type builtin struct {
	name string
	kind Kind
}

func (b builtin) PackageRef() PackageRef { return pkgBuiltin }
func (b builtin) Kind() Kind             { return b.kind }
func (b builtin) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	n, err := w.Write([]byte(b.name))
	return int64(n), err
}

func (b builtin) RegisterImports(i *Imports) {}

// Built in Go types
var (
	BoolType       = builtin{"bool", BoolKind}
	ByteType       = builtin{"byte", ByteKind}
	IntType        = builtin{"int", IntKind}
	Int8Type       = builtin{"int8", Int8Kind}
	Int16Type      = builtin{"int16", Int16Kind}
	Int32Type      = builtin{"int32", Int32Kind}
	Int64Type      = builtin{"int64", Int64Kind}
	Complex128Type = builtin{"complex128", Complex128Kind}
	Complex64Type  = builtin{"complex64", Complex64Kind}
	Float32Type    = builtin{"float32", Float32Kind}
	Float64Type    = builtin{"float64", Float64Kind}
	RuneType       = builtin{"rune", RuneKind}
	StringType     = builtin{"string", StringKind}
	UintType       = builtin{"uint", UintKind}
	Uint8Type      = builtin{"uint8", Uint8Kind}
	Uint16Type     = builtin{"uint16", Uint16Kind}
	Uint32Type     = builtin{"uint32", Uint32Kind}
	Uint64Type     = builtin{"uint64", Uint64Kind}
	UintptrType    = builtin{"uintptr", BuiltinKind}
)
