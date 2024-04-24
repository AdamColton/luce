package reflidx

import (
	"reflect"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/rye"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/entity"
	"github.com/adamcolton/luce/util/reflector"
)

type Namer interface {
	Name() string
}

// TODO: Move this to serial
type Rounder interface {
	serial.Serializer
	serial.Deserializer
}

type Round struct {
	serial.Serializer
	serial.Deserializer
}

var (
	EntsBucket    = []byte("ents")
	IndexesBucket = []byte("idxs")
	// TODO: unexpose, make an object, add registers
	Funcs lmap.Map[string, any]
)

type EntStore[E entity.Entity] struct {
	*entity.EntStore[E]
}

func (es EntStore[E]) AddMethod(name string, multi bool) {
	t := reflector.Type[E]()
	m, found := t.MethodByName(name)
	// Todo: lerr.panic on false
	if !found {
		panic("method not found")
	}
	fn := m.Func.Interface().(func(E) []byte)
	es.AddIndex(name, multi, fn)
}

func New[E entity.Entity](r Rounder, f store.Factory) EntStore[E] {
	t := reflector.Type[E]()
	// TODO: something better
	name := []byte(t.Name())

	init := func() E {
		return reflector.Make(t).Interface().(E)
	}

	str := lerr.Must(f.Store(EntsBucket))
	str = lerr.Must(str.Store(name))

	es := &entity.EntStore[E]{
		Init:         init,
		Store:        str,
		Serializer:   r,
		Deserializer: r,
	}

	addIndexes[E](es, t, name, f)

	return EntStore[E]{es}
}

func addIndexes[E entity.Entity](es *entity.EntStore[E], t reflect.Type, name []byte, f store.Factory) {
	var idxs slice.Slice[entity.BaseIndexer[E]]
	isPtr := t.Kind() == reflect.Pointer
	if isPtr {
		t = t.Elem()
	}
	fields := t.NumField()
	for i := 0; i < fields; i++ {
		f := t.Field(i)
		if _, ok := f.Tag.Lookup("entityIndexSingle"); ok {
			idxs = append(idxs, makeFieldIdx[E](f.Name, f.Index, isPtr))
		}
	}

	if len(idxs) == 0 {
		return
	}

	idx := lerr.Must(f.Store(IndexesBucket))
	es.IdxStore = lerr.Must(idx.Store(name))
	for _, i := range idxs {
		es.AddIndexer(i)
	}
}

func makeFieldIdx[E entity.Entity](fname string, fidx []int, isPtr bool) entity.BaseIndexer[E] {
	var fn func(E) []byte
	if isPtr {
		fn = func(e E) []byte {
			v := reflect.ValueOf(e).Elem().FieldByIndex(fidx)
			return rye.Serialize.Any(v.Interface(), nil)
		}
	} else {
		fn = func(e E) []byte {
			v := reflect.ValueOf(e).FieldByIndex(fidx)
			return rye.Serialize.Any(v.Interface(), nil)
		}
	}

	return entity.BaseIndexer[E]{
		IndexName: fname,
		Fn:        fn,
	}
}
