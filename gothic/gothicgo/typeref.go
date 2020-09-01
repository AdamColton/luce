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
func (t *TypeRef) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString(p.Prefix(t.Pkg))
	sw.WriteString(t.Name)
	sw.Err = lerr.Wrap(sw.Err, "While writing TypeRef %s", t.Name)
	return sw.Rets()
}

// PackageRef fulfills Type. Returns the ExternalPackageRef.
func (t *TypeRef) PackageRef() PackageRef { return t.Pkg }

// RegisterImports fulfills Type.
func (t *TypeRef) RegisterImports(i *Imports) {
	i.Add(t.Pkg)
}

// Elem returns the underlying Type. This may be nil.
func (t *TypeRef) Elem() Type {
	return t.T
}

// StructEmbedName fulfills StructEmbeddable allowin a TypeRef to be embedded in
// a Struct.
func (t *TypeRef) StructEmbedName() string {
	return t.Name
}
