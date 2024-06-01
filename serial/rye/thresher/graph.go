package thresher

import (
	"reflect"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/adamcolton/luce/util/reflector"
)

func Save(i any) (id []byte) {
	g := graph(i)
	g.enc()
	return rootObjByV(reflect.ValueOf(i)).id
}

func graph(i any) *grapher {
	v := reflect.ValueOf(i)
	if v.Kind() != reflect.Slice {
		v = reflector.EnsurePointer(v)
	}
	ro := rootObjByV(v)
	g := newGrapher()
	g.walk(ro)
	return g
}

func newGrapher() *grapher {
	return &grapher{
		ptrs: lmap.Map[uintptr, *rootObj]{},
	}
}

type grapher struct {
	ptrs lmap.Map[uintptr, *rootObj]
}

func (g *grapher) walk(ro *rootObj) {
	g.ptrs[ro.addr] = ro
	v := ro.baseValue()
	c := getBaseCodec(v.Type())

	for _, r := range c.roots(v) {
		_, found := g.ptrs[r.addr]
		if !found {
			g.walk(r)
		}
	}
}

func (g *grapher) enc() {
	for _, ro := range g.ptrs {
		v := ro.baseValue()

		i := v.Interface()
		t := v.Type()
		c := getBaseCodec(t)

		size := c.size(i)
		s := compact.MakeSerializer(int(size))
		c.enc(i, s)
		store[string(ro.id)] = s.Data
	}
}
