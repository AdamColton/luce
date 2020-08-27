package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// PointerType extends the Type interface with pointer specific information
type PointerType struct {
	Type
}

// PointerTo returns a PointerType to the underlying type.
func PointerTo(t Type) *PointerType {
	return &PointerType{t}
}

// PrefixWriteTo fulfills Type.
func (p *PointerType) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteRune('*')
	p.Type.PrefixWriteTo(sw, pre)
	sw.Err = lerr.Wrap(sw.Err, "While writing pointer type")
	return sw.Rets()
}

// PackageRef fulfills Type. Return PkgBuiltin.
func (*PointerType) PackageRef() PackageRef { return pkgBuiltin }

// Elem returns the Type pointed to.
func (p *PointerType) Elem() Type { return p.Type }
