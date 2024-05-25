package thresher

import (
	"reflect"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/adamcolton/luce/util/reflector"
)

/*
A graph is going to hold memory locations linked by pointers.
*/

func Graph(i any) *grapher {
	v := reflector.EnsurePointer(reflect.ValueOf(i))
	ro := rootObjByV(v)
	g := newGrapher()
	g.walk(ro)
	return g
}

func newGrapher() *grapher {
	return &grapher{
		types: lset.New[reflect.Type](),
		ptrs:  lmap.Map[uintptr, *rootObj]{},
	}
}

type grapher struct {
	types *lset.Set[reflect.Type]
	ptrs  lmap.Map[uintptr, *rootObj]
}

func (g *grapher) walk(ro *rootObj) {
	g.ptrs[ro.addr] = ro
	t := ro.v.Type()
	g.types.Add(t)
	c := getCodec(t.Elem())
	for _, r := range c.roots(ro.v.Elem()) {
		_, found := g.ptrs[r.addr]
		if !found {
			g.walk(r)
		}
	}
}

func (g *grapher) enc() []byte {
	size := uint64(len(g.ptrs)) * 8
	for _, ro := range g.ptrs {
		e := ro.v.Elem()
		c := getCodec(e.Type())
		size += c.size(e.Interface())
	}
	s := compact.MakeSerializer(int(size))
	for ptr, ro := range g.ptrs {
		s.Uint64(uint64(ptr))
		e := ro.v.Elem()
		c := getCodec(e.Type())
		c.enc(e.Interface(), s)
	}
	return s.Data
}
