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
}

// HelpfulType fulfils Type. Currently a placeholder.
type HelpfulType interface {
	Type
	Named(string) NameType
	Unnamed() NameType
	Ptr() PointerType
	Slice() SliceType
	Array(size int) ArrayType
}

// HelpfulTypeWrapper turns any Type into a HelpfulType.
type HelpfulTypeWrapper struct{ Type }

// NewHelpfulTypeWrapper ensures that we're not re-wrapping a
// wHelpfulTypeWrapper.
func NewHelpfulTypeWrapper(t Type) HelpfulTypeWrapper {
	if htw, ok := t.(HelpfulTypeWrapper); ok {
		return NewHelpfulTypeWrapper(htw.Type)
	}
	return HelpfulTypeWrapper{t}
}

// NewHelpfulType checks if the type is already a HelpfulType and if not, wraps
// it in a HelpfulTypeWrapper.
func NewHelpfulType(t Type) HelpfulType {
	if ht, ok := t.(HelpfulType); ok {
		return ht
	}
	return HelpfulTypeWrapper{t}
}

// Named returns a NameType
func (h HelpfulTypeWrapper) Named(name string) NameType {
	return NameType{name, h}
}

// Unnamed returns a NameType with an empty string as the name.
func (h HelpfulTypeWrapper) Unnamed() NameType {
	return NameType{"", h}
}

// Ptr returns a Pointer to the type.
func (h HelpfulTypeWrapper) Ptr() PointerType {
	return PointerTo(h.Type)
}

// Slice of the underlying type.
func (h HelpfulTypeWrapper) Slice() SliceType {
	return SliceOf(h.Type)
}

// Array of the underlying type.
func (h HelpfulTypeWrapper) Array(size int) ArrayType {
	return ArrayOf(h.Type, size)
}
