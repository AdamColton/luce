package gothicgo

import (
	"io"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/luceio"
)

// MapType extends Type with map specific information
type MapType interface {
	HelpfulType
	MapElem() HelpfulType
	Elem() HelpfulType
	MapKey() HelpfulType
}

// MapOf returns a MapType around with the given key and element types.
func MapOf(key, elem Type) MapType {
	return mapHT{
		HelpfulTypeWrapper: HelpfulTypeWrapper{
			mapT{
				key:  NewHelpfulType(key),
				elem: NewHelpfulType(elem),
			},
		},
	}
}

// mapT is the inner part of the pointer map structure, it implements Type. It
// also ensures the underlying key and elem are HelpfulTypes so that the
// outerlayer doesn't have to do any conversion to return a HelpfulType. This
// nesting structure allows using the HelpfulTypeWrapper on the inner map type,
// while the outer wrapper extends it with Elem and PointerElem.

type mapT struct {
	key, elem HelpfulType
}

func (m mapT) PrefixWriteTo(w io.Writer, p Prefixer) (int64, error) {
	sw := luceio.NewSumWriter(w)
	sw.WriteString("map[")
	m.key.PrefixWriteTo(sw, p)
	sw.WriteRune(']')
	m.elem.PrefixWriteTo(sw, p)
	sw.Err = lerr.Wrap(sw.Err, "While writing map type")
	return sw.Rets()
}
func (mapT) Kind() Kind { return MapKind }

func (m mapT) RegisterImports(i *Imports) {
	m.elem.RegisterImports(i)
	m.key.RegisterImports(i)
}

func (mapT) PackageRef() PackageRef { return pkgBuiltin }

type mapHT struct {
	HelpfulTypeWrapper
	size int
}

func (m mapHT) MapKey() HelpfulType  { return m.Type.(mapT).key }
func (m mapHT) Elem() HelpfulType    { return m.Type.(mapT).elem }
func (m mapHT) MapElem() HelpfulType { return m.Elem() }
