package gothicgo

// The Type interface represents a type in Go. Name is the type without the
// package, String is the type with the package and PrefixString takes a package name
// and return the string representing the type with the package included.
//
// PackageName returns a string representing the package. Package will return
// a *gothicgo.Package if the Type is part of the Gothic generation.
type Type interface {
	PrefixWriterTo
	ImportsRegistrar
	PackageRef() PackageRef
	Kind() Kind
	Named(string) NameType
	Unnamed() NameType
	Pointer() PointerType
	Slice() *SliceType
	Array(size int) *ArrayType
	AsMapElem(key Type) *MapType
	AsMapKey(elem Type) *MapType
}
