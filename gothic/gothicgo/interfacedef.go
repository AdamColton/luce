package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"

	"github.com/adamcolton/luce/util/luceio"
)

// InterfaceDef represents a Type Definition in a file of an interface.
type InterfaceDef struct {
	Interface *InterfaceType
	name      string
	file      *File
	Comment   string
}

// NewInterfaceDef adds an InterfaceDef to a file
func (f *File) NewInterfaceDef(name string, embed ...InterfaceEmbeddable) (*InterfaceDef, error) {
	i := &InterfaceDef{
		Interface: NewInterfaceType(embed...),
		name:      name,
		file:      f,
	}
	return i, f.AddGenerator(i)
}

// MustInterfaceDef invokes NewInterfaceDef on the File and panics if there is
// an error.
func (f *File) MustInterfaceDef(name string, embed ...InterfaceEmbeddable) *InterfaceDef {
	i, err := f.NewInterfaceDef(name, embed...)
	lerr.Panic(err)
	return i
}

// NewInterfaceDef creates a file with the same name as the interface and
// invokes InterfaceDef on the file.
func (p *Package) NewInterfaceDef(name string, embed ...InterfaceEmbeddable) (*InterfaceDef, error) {
	return p.File(name).NewInterfaceDef(name, embed...)
}

// MustInterfaceDef creates a file with the same name as the interface and
// invokes MustInterfaceDef on the file.
func (p *Package) MustInterfaceDef(name string, embed ...InterfaceEmbeddable) *InterfaceDef {
	return p.File(name).MustInterfaceDef(name, embed...)
}

// NewInterfaceDef creates a package and file with the same name as the
// interface and invokes NewInterfaceDef on the file.
func (c *BaseContext) NewInterfaceDef(name string, embed ...InterfaceEmbeddable) (*InterfaceDef, error) {
	pkg, err := c.Package(name)
	if err != nil {
		return nil, err
	}
	return pkg.NewInterfaceDef(name, embed...)
}

// MustInterfaceDef creates a package and file with the same name as the
// interface and invokes MustInterfaceDef on the file.
func (c *BaseContext) MustInterfaceDef(name string, embed ...InterfaceEmbeddable) *InterfaceDef {
	pkg, err := c.Package(name)
	lerr.Panic(err)
	return pkg.MustInterfaceDef(name, embed...)
}

// PrefixWriteTo fulfills PrefixWriterTo. Writes the interface type def to the
// writer.
func (i *InterfaceDef) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteStrings("type ", i.name, " ")
	sumPrefixWriteTo(sw, pre, i.Interface)
	return sw.Rets()
}

// Ref gets a reference to the interface type definition.
func (i *InterfaceDef) Ref() *InterfaceRef {
	return &InterfaceRef{
		Name:      i.name,
		Pkg:       i.file.Package(),
		Interface: i.Interface,
	}
}

// RegisterImports for all argument and return types spedified or the packages
// of any interfaces embedded.
func (i *InterfaceDef) RegisterImports(im *Imports) {
	i.Interface.RegisterImports(im)
}

// PackageRef of the Interface Definition.
func (i *InterfaceDef) PackageRef() PackageRef {
	return i.file.Package()
}

// File of the Interface Definition.
func (i *InterfaceDef) File() *File {
	return i.file
}

// ScopeName fulfills Namer.
func (i *InterfaceDef) ScopeName() string {
	return i.name
}

// Elem returns the underlying Interface as a Type.
func (i *InterfaceDef) Elem() Type {
	return i.Interface
}

// Embed values into the InterfaceDef. InterfaceDef is returned for chaining.
func (i *InterfaceDef) Embed(embed ...InterfaceEmbeddable) *InterfaceDef {
	i.Interface.Embed(embed...)
	return i
}
