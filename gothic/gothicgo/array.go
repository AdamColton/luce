package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// ArrayType extends Type with array specific information
type ArrayType interface {
	HelpfulType
	ArrayElem() HelpfulType
	Elem() HelpfulType
	// Size is the number of elements in the array
	Size() int
}

// ArrayOf returns a ArrayType around t.
func ArrayOf(t Type, size int) ArrayType {
	if size < 0 {
		size = 0
	}
	return arrayHT{
		HelpfulTypeWrapper: HelpfulTypeWrapper{
			arrayT{
				HelpfulType: NewHelpfulType(t),
				size:        size,
			},
		},
	}
}

// arrayT is the inner part of the pointer array structure, it inheirits
// ImportsRegistrar from the underlying type but overrides Kind, PackageRef
// and PrefixWriteTo. It also ensures the underlying type is HelpfulType so that
// the outerlayer doesn't have to do any conversion to return a HelpfulType.
// This nesting structure allows using the HelpfulTypeWrapper on the inner
// array type, while the outer wrapper extends it with Elem and PointerElem.

type arrayT struct {
	HelpfulType
	size int
}

func (a arrayT) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString("[")
	if a.size > 0 {
		sw.WriteInt(a.size)
	} else {
		sw.WriteString("...")
	}
	sw.WriteString("]")
	a.HelpfulType.PrefixWriteTo(sw, p)
	sw.Err = lerr.Wrap(sw.Err, "While writing array")
	return sw.Rets()
}
func (arrayT) Kind() Kind             { return ArrayKind }
func (arrayT) PackageRef() PackageRef { return pkgBuiltin }

type arrayHT struct {
	HelpfulTypeWrapper
	size int
}

func (a arrayHT) Elem() HelpfulType      { return a.Type.(arrayT).HelpfulType }
func (a arrayHT) ArrayElem() HelpfulType { return a.Elem() }
func (a arrayHT) Size() int              { return a.Type.(arrayT).size }
