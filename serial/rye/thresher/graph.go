package thresher

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/lset"
	"github.com/adamcolton/luce/ds/slice"
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
	ptrs := g.ptrs.Keys(nil).Sort(slice.LT[uintptr]())
	ln := len(g.ptrs)
	sizes := make([]uint64, ln)
	ln64 := uint64(ln)
	size := ln64*4 + compact.SizeUint64(uint64(ln64))
	for i, ptr := range ptrs {
		ro := g.ptrs[ptr]
		e := ro.v.Elem()
		c := getCodec(e.Type())
		s := c.size(e.Interface())
		sizes[i] = s
		size += s + compact.SizeUint64(s) + compact.Size(ro.id)
	}
	s := compact.MakeSerializer(int(size))
	s.CompactUint64(ln64)
	for i, ptr := range ptrs {
		ro := g.ptrs[ptr]
		s.CompactSlice(ro.id)
		e := ro.v.Elem()
		t := e.Type()
		s.Uint32(type2id(t))
		s.CompactUint64(sizes[i])
		c := getCodec(t)
		c.enc(e.Interface(), s)
	}
	return s.Data
}

func (g *grapher) dec(data []byte) {
	d := compact.NewDeserializer(data)
	ln := d.CompactUint64()
	g.ptrs = make(lmap.Map[uintptr, *rootObj], ln)
	wg := sync.WaitGroup{}
	for !d.Done() {
		id := d.CompactSlice()
		t, _ := typeIDs.A(d.Uint32())
		ro := makeRootObj(t, id)
		g.ptrs[ro.addr] = ro
		c := getCodec(t)
		size := d.CompactUint64()
		sub := compact.NewDeserializer(d.Slice(int(size)))
		ch := c.dec(sub)
		wg.Add(1)
		go func() {
			var v reflect.Value
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
				wg.Done()
			}()
			v = reflect.ValueOf(<-ch)
			ro.v.Elem().Set(v)
		}()
	}
	wg.Wait()
}
