package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// TypeRef informs a gothic program of types available in external packages. It
// does not generate these but can use them. For instance
// MustPackageRef("time").NewTypeRef("Time", nil) creates a reference to the
// Time type in the time package. Passing in nil for the Type is acceptible.
type TypeRef struct {
	Name string
	T    Type
	Pkg  PackageRef
}

// NewTypeRef constructs a TypeRef.
func NewTypeRef(p PackageRef, name string, t Type) *TypeRef {
	return &TypeRef{
		Name: name,
		T:    t,
		Pkg:  p,
	}
}

// PrefixWriteTo fulfills Type. Writes the TypeRef with prefixing.
func (e *TypeRef) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString(p.Prefix(e.Pkg))
	sw.WriteString(e.Name)
	sw.Err = lerr.Wrap(sw.Err, "While writing external type %s", e.Name)
	return sw.Rets()
}

// PackageRef fulfills Type. Returns the ExternalPackageRef.
func (e *TypeRef) PackageRef() PackageRef { return e.Pkg }

// RegisterImports fulfills Type.
func (e *TypeRef) RegisterImports(i *Imports) {
	i.Add(e.Pkg)
}

// Elem returns the underlying Type. This may be nil.
func (e *TypeRef) Elem() Type {
	return e.T
}
