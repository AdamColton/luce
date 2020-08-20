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
	BoolType       = HelpfulType(builtin{"bool", BoolKind})
	ByteType       = HelpfulType(builtin{"byte", ByteKind})
	IntType        = HelpfulType(builtin{"int", IntKind})
	Int8Type       = HelpfulType(builtin{"int8", Int8Kind})
	Int16Type      = HelpfulType(builtin{"int16", Int16Kind})
	Int32Type      = HelpfulType(builtin{"int32", Int32Kind})
	Int64Type      = HelpfulType(builtin{"int64", Int64Kind})
	Complex128Type = HelpfulType(builtin{"complex128", Complex128Kind})
	Complex64Type  = HelpfulType(builtin{"complex64", Complex64Kind})
	Float32Type    = HelpfulType(builtin{"float32", Float32Kind})
	Float64Type    = HelpfulType(builtin{"float64", Float64Kind})
	RuneType       = HelpfulType(builtin{"rune", RuneKind})
	StringType     = HelpfulType(builtin{"string", StringKind})
	UintType       = HelpfulType(builtin{"uint", UintKind})
	Uint8Type      = HelpfulType(builtin{"uint8", Uint8Kind})
	Uint16Type     = HelpfulType(builtin{"uint16", Uint16Kind})
	Uint32Type     = HelpfulType(builtin{"uint32", Uint32Kind})
	Uint64Type     = HelpfulType(builtin{"uint64", Uint64Kind})
	UintptrType    = HelpfulType(builtin{"uintptr", BuiltinKind})
)
