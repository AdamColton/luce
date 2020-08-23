package gothicgo

type coreType interface {
	PrefixWriterTo
	ImportsRegistrar
	PackageRef() PackageRef
	Kind() Kind
}

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
	Ptr() PointerType
	Slice() SliceType
	Array(size int) ArrayType
	AsMapElem(key Type) MapType
	AsMapKey(elem Type) MapType
}

type typeWrapper struct{ coreType }

func newTypeWrapper(t coreType) typeWrapper {
	if htw, ok := t.(typeWrapper); ok {
		return newTypeWrapper(htw.coreType)
	}
	return typeWrapper{t}
}

func newType(ct coreType) Type {
	if t, ok := ct.(Type); ok {
		return t
	}
	return typeWrapper{ct}
}

func (t typeWrapper) Named(name string) NameType {
	return NameType{name, t}
}

func (t typeWrapper) Unnamed() NameType {
	return NameType{"", t}
}

func (t typeWrapper) Ptr() PointerType {
	return PointerTo(t)
}

func (t typeWrapper) Slice() SliceType {
	return SliceOf(t)
}

func (t typeWrapper) Array(size int) ArrayType {
	return ArrayOf(t, size)
}

func (t typeWrapper) AsMapElem(key Type) MapType {
	return MapOf(key, t)
}

func (t typeWrapper) AsMapKey(elem Type) MapType {
	return MapOf(t, elem)
}
