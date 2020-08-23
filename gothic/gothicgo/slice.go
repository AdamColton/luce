package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// SliceType extends Type with slice specific information
type SliceType interface {
	Type
	SliceElem() Type
	Elem() Type
}

// SliceOf returns a SliceType around t.
func SliceOf(t Type) SliceType {
	return sliceT{
		typeWrapper{sliceCT{newType(t)}},
	}
}

// sliceCT is the inner part of the slice nesting structure, it inheirits
// ImportsRegistrar from the underlying type but overrides Kind, PackageRef and
// PrefixWriteTo. This nesting structure allows using the typeWrapper on the
// inner sliceCT, while the outer wrapper extends it with Elem and SliceElem.

type sliceCT struct {
	Type
}

func (s sliceCT) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString("[]")
	s.Type.PrefixWriteTo(sw, p)
	sw.Err = lerr.Wrap(sw.Err, "While writing slice")
	return sw.Rets()
}
func (sliceCT) Kind() Kind             { return SliceKind }
func (sliceCT) PackageRef() PackageRef { return pkgBuiltin }

type sliceT struct {
	typeWrapper
}

func (s sliceT) Elem() Type      { return s.coreType.(sliceCT).Type }
func (s sliceT) SliceElem() Type { return s.Elem() }
