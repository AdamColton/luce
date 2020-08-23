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
	BoolType       = Type(typeWrapper{builtin{"bool", BoolKind}})
	ByteType       = Type(typeWrapper{builtin{"byte", ByteKind}})
	IntType        = Type(typeWrapper{builtin{"int", IntKind}})
	Int8Type       = Type(typeWrapper{builtin{"int8", Int8Kind}})
	Int16Type      = Type(typeWrapper{builtin{"int16", Int16Kind}})
	Int32Type      = Type(typeWrapper{builtin{"int32", Int32Kind}})
	Int64Type      = Type(typeWrapper{builtin{"int64", Int64Kind}})
	Complex128Type = Type(typeWrapper{builtin{"complex128", Complex128Kind}})
	Complex64Type  = Type(typeWrapper{builtin{"complex64", Complex64Kind}})
	Float32Type    = Type(typeWrapper{builtin{"float32", Float32Kind}})
	Float64Type    = Type(typeWrapper{builtin{"float64", Float64Kind}})
	RuneType       = Type(typeWrapper{builtin{"rune", RuneKind}})
	StringType     = Type(typeWrapper{builtin{"string", StringKind}})
	UintType       = Type(typeWrapper{builtin{"uint", UintKind}})
	Uint8Type      = Type(typeWrapper{builtin{"uint8", Uint8Kind}})
	Uint16Type     = Type(typeWrapper{builtin{"uint16", Uint16Kind}})
	Uint32Type     = Type(typeWrapper{builtin{"uint32", Uint32Kind}})
	Uint64Type     = Type(typeWrapper{builtin{"uint64", Uint64Kind}})
	UintptrType    = Type(typeWrapper{builtin{"uintptr", BuiltinKind}})
)
