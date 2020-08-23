package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// ArrayType extends Type with array specific information
type ArrayType interface {
	Type
	ArrayElem() Type
	Elem() Type
	// Size is the number of elements in the array
	Size() int
}

// ArrayOf returns a ArrayType around t.
func ArrayOf(t Type, size int) ArrayType {
	if size < 0 {
		size = 0
	}
	return arrayT{
		typeWrapper{
			arrayCT{
				Type: t,
				size: size,
			},
		},
	}
}

// arrayCT is the inner part of the array nesting structure, it inheirits
// ImportsRegistrar from the underlying type but overrides Kind, PackageRef and
// PrefixWriteTo. This nesting structure allows using the typeWrapper on the
// inner arrayCT, while the outer wrapper extends it with Elem, ArrayElem and
// Size.

type arrayCT struct {
	Type
	size int
}

func (a arrayCT) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString("[")
	if a.size > 0 {
		sw.WriteInt(a.size)
	} else {
		sw.WriteString("...")
	}
	sw.WriteString("]")
	a.Type.PrefixWriteTo(sw, p)
	sw.Err = lerr.Wrap(sw.Err, "While writing array")
	return sw.Rets()
}
func (arrayCT) Kind() Kind             { return ArrayKind }
func (arrayCT) PackageRef() PackageRef { return pkgBuiltin }

type arrayT struct {
	typeWrapper
}

func (a arrayT) Elem() Type      { return a.coreType.(arrayCT).Type }
func (a arrayT) ArrayElem() Type { return a.Elem() }
func (a arrayT) Size() int       { return a.coreType.(arrayCT).size }
