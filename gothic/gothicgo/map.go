package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// MapType extends Type with map specific information
type MapType struct {
	Key, Val Type
}

// MapOf returns a MapType around with the given key and element types.
func MapOf(key, val Type) *MapType {
	return &MapType{
		Key: key,
		Val: val,
	}
}

// PrefixWriteTo fulfills Type.
func (m *MapType) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString("map[")
	m.Key.PrefixWriteTo(sw, p)
	sw.WriteRune(']')
	m.Val.PrefixWriteTo(sw, p)
	sw.Err = lerr.Wrap(sw.Err, "While writing map type")
	return sw.Rets()
}

// RegisterImports for the Key and Val types.
func (m *MapType) RegisterImports(i *Imports) {
	m.Val.RegisterImports(i)
	m.Key.RegisterImports(i)
}

// PackageRef fulfills Type, Returns PkgBuiltin.
func (*MapType) PackageRef() PackageRef { return pkgBuiltin }

// Elem returns the Val type.
func (m *MapType) Elem() Type { return m.Val }
