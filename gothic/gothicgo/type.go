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
}

// HelpfulTypeWrapper turns any Type into a HelpfulType.
type HelpfulTypeWrapper struct{ Type }

// Named returns a NameType
func (h HelpfulTypeWrapper) Named(name string) NameType {
	return NameType{name, h}
}

// Unnamed returns a NameType with an empty string as the name.
func (h HelpfulTypeWrapper) Unnamed() NameType {
	return NameType{"", h}
}
