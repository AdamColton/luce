package thresher

import (
	"reflect"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/serial/rye/compact"
	"github.com/adamcolton/luce/util/reflector"
)

func Graph(i any) *grapher {
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

	var c *codec
	var v reflect.Value
	switch t.Kind() {
	case reflect.Pointer:
		c = getCodec(t.Elem())
		v = ro.v.Elem()
	case reflect.Slice:
		c = makeSliceCodec(t)
		v = ro.v
	}

	for _, r := range c.roots(v) {
		_, found := g.ptrs[r.addr]
		if !found {
			g.walk(r)
		}
	}
}

func (g *grapher) enc() {
	for _, ro := range g.ptrs {
		var c *codec
		var i any
		var t reflect.Type
		switch ro.v.Kind() {
		case reflect.Pointer:
			e := ro.v.Elem()
			t = e.Type()
			c = getCodec(t)
			i = e.Interface()
		case reflect.Slice:
			c = makeSliceCodec(t)
			t = ro.v.Type()
			i = ro.v.Interface()
		}

		size := c.size(i)
		s := compact.MakeSerializer(int(size))
		c.enc(i, s)
		store[string(ro.id)] = storeRecord{
			t:    t,
			data: s.Data,
		}
	}
}
