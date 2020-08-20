package gothicgo

// Kind represents the different kinds of types. Two different structs will have
// different types, but the same kind (StructKind)
type Kind uint8

// Defined Kinds
const (
	NoneKind = Kind(iota)
	SliceKind
	PointerKind
	MapKind
	FuncKind
	StructKind
	UnknownKind
	BuiltinKind
	InterfaceKind
	TypeDefKind
	InterfaceTypeDefKind
	ArrayKind
	BoolKind
	ByteKind
	IntKind
	Int8Kind
	Int16Kind
	Int32Kind
	Int64Kind
	Complex128Kind
	Complex64Kind
	Float32Kind
	Float64Kind
	RuneKind
	StringKind
	UintKind
	Uint8Kind
	Uint16Kind
	Uint32Kind
	Uint64Kind
)

var kindStrs = map[Kind]string{
	SliceKind:            "SliceKind",
	PointerKind:          "PointerKind",
	MapKind:              "MapKind",
	FuncKind:             "FuncKind",
	StructKind:           "StructKind",
	UnknownKind:          "UnknownKind",
	BuiltinKind:          "BuiltinKind",
	InterfaceKind:        "InterfaceKind",
	TypeDefKind:          "TypeDefKind",
	InterfaceTypeDefKind: "InterfaceTypeDefKind",
	ArrayKind:            "ArrayKind",
	BoolKind:             "BoolKind",
	ByteKind:             "ByteKind",
	IntKind:              "IntKind",
	Int8Kind:             "Int8Kind",
	Int16Kind:            "Int16Kind",
	Int32Kind:            "Int32Kind",
	Int64Kind:            "Int64Kind",
	Complex128Kind:       "Complex128Kind",
	Complex64Kind:        "Complex64Kind",
	Float32Kind:          "Float32Kind",
	Float64Kind:          "Float64Kind",
	RuneKind:             "RuneKind",
	StringKind:           "StringKind",
	UintKind:             "UintKind",
	Uint8Kind:            "Uint8Kind",
	Uint16Kind:           "Uint16Kind",
	Uint32Kind:           "Uint32Kind",
	Uint64Kind:           "Uint64Kind",
}

func (k Kind) String() string {
	s, ok := kindStrs[k]
	if ok {
		return s
	}
	return "UndefinedKind"
}
