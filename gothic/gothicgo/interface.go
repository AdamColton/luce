package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// InterfaceEmbeddable allows one interface to be embedded in another
type InterfaceEmbeddable interface {
	InterfaceEmbed(w io.Writer, pre Prefixer) (int64, error)
}

// InterfaceType is used to generate an interface
type InterfaceType struct {
	embedded []InterfaceEmbeddable
}

// NewInterfaceType adds a new interface to an existing file
func NewInterfaceType(embed ...InterfaceEmbeddable) *InterfaceType {
	return &InterfaceType{
		embedded: embed,
	}
}

// Embed one interface in another.
func (i *InterfaceType) Embed(embed ...InterfaceEmbeddable) *InterfaceType {
	i.embedded = append(i.embedded, embed...)
	return i
}

// PrefixWriteTo writes the interface name and package if necessary.
func (i *InterfaceType) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	if len(i.embedded) == 0 {
		n, err := w.Write([]byte("interface{}"))
		return int64(n), err
	}
	sw := luceio.NewSumWriter(w)
	sw.WriteString("interface {")
	for _, e := range i.embedded {
		sw.WriteString("\n\t")
		e.InterfaceEmbed(sw, pre)
	}
	sw.WriteString("\n}")
	sw.Err = lerr.Wrap(sw.Err, "While writing interface:")
	return sw.Rets()
}

// RegisterImports on the Interface.
func (i *InterfaceType) RegisterImports(im *Imports) {
	for _, e := range i.embedded {
		if r, ok := e.(ImportsRegistrar); ok {
			r.RegisterImports(im)
		}
	}
}

// PackageRef for the package Interface is in, fulfills Type interface.
func (i *InterfaceType) PackageRef() PackageRef { return pkgBuiltin }
