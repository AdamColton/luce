package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

type InterfaceRef struct {
	Name      string
	Interface *InterfaceType
	Pkg       PackageRef
}

func NewInterfaceRef(p PackageRef, name string) *InterfaceRef {
	return &InterfaceRef{
		Name: name,
		Pkg:  p,
	}
}

func (i *InterfaceRef) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString(p.Prefix(i.Pkg))
	sw.WriteString(i.Name)
	sw.Err = lerr.Wrap(sw.Err, "While writing InterfaceRef %s", i.Name)
	return sw.Rets()
}

func (i *InterfaceRef) PackageRef() PackageRef { return i.Pkg }

func (i *InterfaceRef) RegisterImports(im *Imports) {
	im.Add(i.Pkg)
}

func (i *InterfaceRef) Elem() Type {
	return i.Interface
}

func (i *InterfaceRef) InterfaceEmbed(w io.Writer, pre Prefixer) (int64, error) {
	return i.PrefixWriteTo(w, pre)
}
