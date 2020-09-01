package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// PointerType extends the Type interface with pointer specific information
type PointerType interface {
	Type
	Elem() Type
	PointerElem() Type
}

type pointer struct {
	Type
}

// PointerTo returns a PointerType to the underlying type.
func PointerTo(t Type) PointerType {
	p := &pointer{t}
	if _, ok := t.(StructEmbeddable); ok {
		return embeddablePointerWrapper{p}
	}
	return p
}

// PrefixWriteTo fulfills Type.
func (p *pointer) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteRune('*')
	p.Type.PrefixWriteTo(sw, pre)
	sw.Err = lerr.Wrap(sw.Err, "While writing pointer type")
	return sw.Rets()
}

// PackageRef fulfills Type. Return PkgBuiltin.
func (*pointer) PackageRef() PackageRef { return pkgBuiltin }

// Elem returns the Type pointed to.
func (p *pointer) Elem() Type { return p.Type }

func (p *pointer) PointerElem() Type { return p.Type }

type embeddablePointerWrapper struct{ *pointer }

func (e embeddablePointerWrapper) StructEmbedName() string {
	return e.Type.(StructEmbeddable).StructEmbedName()
}
