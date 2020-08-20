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
	BoolType       = HelpfulType(HelpfulTypeWrapper{builtin{"bool", BoolKind}})
	ByteType       = HelpfulType(HelpfulTypeWrapper{builtin{"byte", ByteKind}})
	IntType        = HelpfulType(HelpfulTypeWrapper{builtin{"int", IntKind}})
	Int8Type       = HelpfulType(HelpfulTypeWrapper{builtin{"int8", Int8Kind}})
	Int16Type      = HelpfulType(HelpfulTypeWrapper{builtin{"int16", Int16Kind}})
	Int32Type      = HelpfulType(HelpfulTypeWrapper{builtin{"int32", Int32Kind}})
	Int64Type      = HelpfulType(HelpfulTypeWrapper{builtin{"int64", Int64Kind}})
	Complex128Type = HelpfulType(HelpfulTypeWrapper{builtin{"complex128", Complex128Kind}})
	Complex64Type  = HelpfulType(HelpfulTypeWrapper{builtin{"complex64", Complex64Kind}})
	Float32Type    = HelpfulType(HelpfulTypeWrapper{builtin{"float32", Float32Kind}})
	Float64Type    = HelpfulType(HelpfulTypeWrapper{builtin{"float64", Float64Kind}})
	RuneType       = HelpfulType(HelpfulTypeWrapper{builtin{"rune", RuneKind}})
	StringType     = HelpfulType(HelpfulTypeWrapper{builtin{"string", StringKind}})
	UintType       = HelpfulType(HelpfulTypeWrapper{builtin{"uint", UintKind}})
	Uint8Type      = HelpfulType(HelpfulTypeWrapper{builtin{"uint8", Uint8Kind}})
	Uint16Type     = HelpfulType(HelpfulTypeWrapper{builtin{"uint16", Uint16Kind}})
	Uint32Type     = HelpfulType(HelpfulTypeWrapper{builtin{"uint32", Uint32Kind}})
	Uint64Type     = HelpfulType(HelpfulTypeWrapper{builtin{"uint64", Uint64Kind}})
	UintptrType    = HelpfulType(HelpfulTypeWrapper{builtin{"uintptr", BuiltinKind}})
)
