package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// SliceType extends Type with slice specific information
type SliceType interface {
	HelpfulType
	SliceElem() HelpfulType
	Elem() HelpfulType
}

// SliceOf returns a SliceType around t.
func SliceOf(t Type) SliceType {
	return sliceHT{
		HelpfulTypeWrapper{sliceT{NewHelpfulType(t)}},
	}
}

// sliceT is the inner part of the pointer slice structure, it inheirits
// ImportsRegistrar and PackageRef from the underlying type but overrides Kind
// and PrefixWriteTo. It also ensures the underlying type is HelpfulType so that
// the outerlayer doesn't have to do any conversion to return a HelpfulType.
// This nesting structure allows using the HelpfulTypeWrapper on the inner
// slice type, while the outer wrapper extends it with Elem and PointerElem.

type sliceT struct {
	HelpfulType
}

func (s sliceT) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString("[]")
	s.HelpfulType.PrefixWriteTo(sw, p)
	sw.Err = lerr.Wrap(sw.Err, "While writing slice")
	return sw.Rets()
}
func (sliceT) Kind() Kind { return SliceKind }

type sliceHT struct {
	HelpfulTypeWrapper
}

func (s sliceHT) Elem() HelpfulType      { return s.Type.(sliceT).HelpfulType }
func (s sliceHT) SliceElem() HelpfulType { return s.Elem() }
