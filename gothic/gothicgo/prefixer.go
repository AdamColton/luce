package gothicgo

import "io"

// Prefixer takes a PackageRef and returns the correct prefix for it. If the
// reference is to the same pacakge we are in, it will return an empty string.
// If it's a package imported normally, it will return the package name followed
// by a period. If it is an aliased package, it will return the alias followed
// by a period.
type Prefixer interface {
	Prefix(ref PackageRef) string
}

// PrefixWriterTo is the base type for Go code generation. It uses the Prefixer
// to correctly prefix values in the generated code that is written to the
// Writer.
type PrefixWriterTo interface {
	PrefixWriteTo(io.Writer, Prefixer) (int64, error)
}

// DefaultPrefixer always returns the package prefix.
var DefaultPrefixer = defaultPrefixer{}

type defaultPrefixer struct{}

func (defaultPrefixer) Prefix(ref PackageRef) string {
	if ref.Name() == "" {
		return ""
	}
	return ref.Name() + "."
}
