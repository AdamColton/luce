package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// SliceType extends Type with slice specific information
type SliceType struct {
	Type
}

// SliceOf returns a SliceType around t.
func SliceOf(t Type) *SliceType {
	return &SliceType{t}
}

func (s *SliceType) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString("[]")
	s.Type.PrefixWriteTo(sw, p)
	sw.Err = lerr.Wrap(sw.Err, "While writing slice")
	return sw.Rets()
}
func (*SliceType) PackageRef() PackageRef { return pkgBuiltin }

func (s *SliceType) Elem() Type { return s.Type }
