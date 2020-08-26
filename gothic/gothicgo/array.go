package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// ArrayOf returns a ArrayType around t.
func ArrayOf(t Type, size int) ArrayType {
	if size < 0 {
		size = 0
	}
	return ArrayType{
		Type: t,
		Size: size,
	}
}

// ArrayType implements a Go array.
type ArrayType struct {
	Type
	Size int
}

// PrefixWriteTo fulfills PrefixWriterTo. Writes the array signature.
func (a ArrayType) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString("[")
	if a.Size > 0 {
		sw.WriteInt(a.Size)
	} else {
		sw.WriteString("...")
	}
	sw.WriteString("]")
	a.Type.PrefixWriteTo(sw, p)
	sw.Err = lerr.Wrap(sw.Err, "While writing array")
	return sw.Rets()
}

// Kind fulfills Type. Returns ArrayKind.
func (ArrayType) Kind() Kind { return ArrayKind }

// PackageRef fulfills Type. Returns pkgBuiltin.
func (ArrayType) PackageRef() PackageRef { return pkgBuiltin }

// Elem returns the underlying Type.
func (a ArrayType) Elem() Type { return a.Type }
