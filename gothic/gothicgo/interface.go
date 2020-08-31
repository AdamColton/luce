package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// InterfaceEmbeddable allows one interface to be embedded in another
type InterfaceEmbeddable interface {
	Type
	InterfaceEmbedName() string
}

// InterfaceType is used to generate an interface
type InterfaceType struct {
	methods  []*FuncSig
	embedded []InterfaceEmbeddable
}

// NewInterfaceType adds a new interface to an existing file
func NewInterfaceType() *InterfaceType {
	return &InterfaceType{}
}

// AddMethod to the interface.
func (i *InterfaceType) AddMethod(funcSig *FuncSig) {
	i.methods = append(i.methods, funcSig)
}

// Embed one interface in another.
func (i *InterfaceType) Embed(embed InterfaceEmbeddable) {
	i.embedded = append(i.embedded, embed)
}

// PrefixWriteTo writes the interface name and package if necessary.
func (i *InterfaceType) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	if len(i.methods) == 0 && len(i.embedded) == 0 {
		n, err := w.Write([]byte("interface{}"))
		return int64(n), err
	}
	sw := luceio.NewSumWriter(w)
	sw.WriteString("interface {")
	for _, e := range i.embedded {
		sw.WriteString("\n\t")
		e.PrefixWriteTo(sw, pre)
	}
	for _, m := range i.methods {
		sw.WriteString("\n\t")
		sumPrefixWriteTo(sw, pre, m.AsType(false))
	}
	sw.WriteString("\n}")
	sw.Err = lerr.Wrap(sw.Err, "While writing interface:")
	return sw.Rets()
}

// RegisterImports on the Interface.
func (i *InterfaceType) RegisterImports(im *Imports) {
	for _, m := range i.methods {
		m.RegisterImports(im)
	}
	for _, e := range i.embedded {
		im.Add(e.PackageRef())
	}
}

// PackageRef for the package Interface is in, fulfills Type interface.
func (i *InterfaceType) PackageRef() PackageRef { return pkgBuiltin }
