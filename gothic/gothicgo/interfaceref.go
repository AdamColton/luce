package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// InterfaceRef references an interface type definition in a package.
type InterfaceRef struct {
	Name      string
	Interface *InterfaceType
	Pkg       PackageRef
}

// NewInterfaceRef creates an InterfaceRef.
func NewInterfaceRef(p PackageRef, name string) *InterfaceRef {
	return &InterfaceRef{
		Name: name,
		Pkg:  p,
	}
}

// PrefixWriteTo fulfills PrefixWriterTo.
func (i *InterfaceRef) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString(p.Prefix(i.Pkg))
	sw.WriteString(i.Name)
	sw.Err = lerr.Wrap(sw.Err, "While writing InterfaceRef %s", i.Name)
	return sw.Rets()
}

// PackageRef returns the Package in which the interface type is definined.
func (i *InterfaceRef) PackageRef() PackageRef { return i.Pkg }

// RegisterImports adds the Package as an import.
func (i *InterfaceRef) RegisterImports(im *Imports) {
	im.Add(i.Pkg)
}

// Elem returns the underlying Interface as a Type.
func (i *InterfaceRef) Elem() Type {
	return i.Interface
}

// InterfaceEmbed allows one interface to be embedded in another.
func (i *InterfaceRef) InterfaceEmbed(w io.Writer, pre Prefixer) (int64, error) {
	return i.PrefixWriteTo(w, pre)
}
