package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// MapType extends Type with map specific information
type MapType interface {
	Type
	MapElem() Type
	Elem() Type
	MapKey() Type
}

// MapOf returns a MapType around with the given key and element types.
func MapOf(key, elem Type) MapType {
	return mapT{
		typeWrapper: typeWrapper{
			mapCT{
				key:  key,
				elem: elem,
			},
		},
	}
}

// mapCT is the inner part of the map nesting structure, it implements coreType.
// This nesting structure allows using the typeWrapper on the inner mapCT, while
// the outer wrapper extends it with Elem, MapElem and MapKey.

type mapCT struct {
	key, elem Type
}

func (m mapCT) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString("map[")
	m.key.PrefixWriteTo(sw, p)
	sw.WriteRune(']')
	m.elem.PrefixWriteTo(sw, p)
	sw.Err = lerr.Wrap(sw.Err, "While writing map type")
	return sw.Rets()
}
func (mapCT) Kind() Kind { return MapKind }

func (m mapCT) RegisterImports(i *Imports) {
	m.elem.RegisterImports(i)
	m.key.RegisterImports(i)
}

func (mapCT) PackageRef() PackageRef { return pkgBuiltin }

type mapT struct {
	typeWrapper
	size int
}

func (m mapT) MapKey() Type  { return m.coreType.(mapCT).key }
func (m mapT) Elem() Type    { return m.coreType.(mapCT).elem }
func (m mapT) MapElem() Type { return m.Elem() }
