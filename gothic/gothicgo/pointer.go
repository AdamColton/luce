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

// PointerTo returns a PointerType to the underlying type.
func PointerTo(t Type) PointerType {
	return pointerT{
		typeWrapper{
			pointerCT{t},
		},
	}
}

// pointerCT is the inner part of the pointer nesting structure, it inheirits
// ImportsRegistrar from the underlying type but overrides Kind, PackageRef and
// PrefixWriteTo. This nesting structure allows using the typeWrapper on the
// inner pointerCT, while the outer wrapper extends it with Elem and
// PointerElem.

type pointerCT struct {
	Type
}

func (p pointerCT) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteRune('*')
	p.Type.PrefixWriteTo(sw, pre)
	sw.Err = lerr.Wrap(sw.Err, "While writing pointer type")
	return sw.Rets()
}
func (pointerCT) Kind() Kind             { return PointerKind }
func (pointerCT) PackageRef() PackageRef { return pkgBuiltin }

type pointerT struct {
	typeWrapper
}

func (p pointerT) Elem() Type        { return p.coreType.(pointerCT).Type }
func (p pointerT) PointerElem() Type { return p.Elem() }
