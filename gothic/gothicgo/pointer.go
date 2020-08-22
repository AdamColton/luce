package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// PointerType extends the Type interface with pointer specific information
type PointerType interface {
	HelpfulType
	Elem() HelpfulType
	PointerElem() HelpfulType
}

// PointerTo returns a PointerType to the underlying type.
func PointerTo(t Type) PointerType {
	return pointerHT{
		HelpfulTypeWrapper{pointerT{NewHelpfulType(t)}},
	}
}

// pointerT is the inner part of the pointer nesting structure, it inheirits
// ImportsRegistrar and PackageRef from the underlying type but overrides Kind
// and PrefixWriteTo. It also ensures the underlying type is HelpfulType so that
// the outerlayer doesn't have to do any conversion to return a HelpfulType.
// This nesting structure allows using the HelpfulTypeWrapper on the inner
// pointer type, while the outer wrapper extends it with Elem and PointerElem.

type pointerT struct {
	HelpfulType
}

func (p pointerT) PrefixWriteTo(w io.Writer, pre Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteRune('*')
	p.HelpfulType.PrefixWriteTo(sw, pre)
	sw.Err = lerr.Wrap(sw.Err, "While writing pointer type")
	return sw.Rets()
}
func (pointerT) Kind() Kind { return PointerKind }

type pointerHT struct {
	HelpfulTypeWrapper
}

func (p pointerHT) Elem() HelpfulType        { return p.Type.(pointerT).HelpfulType }
func (p pointerHT) PointerElem() HelpfulType { return p.Elem() }
