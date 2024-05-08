package entity

import (
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/wrap/gob"
	"github.com/adamcolton/luce/serial/wrap/json"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/util/reflector"
)

type Builder struct {
	Factory store.Factory
	serial.Serializer
	serial.Deserializer
}

var (
	entsKey = []byte("ents")
	idxsKey = []byte("idxs")
)

func NewStore[E Entity](b Builder, name string, init func() E) (*EntStore[E], error) {
	s, err := b.Factory.Store([]byte(name))
	if err != nil {
		return nil, err
	}
	ents, err := s.Store(entsKey)
	if err != nil {
		return nil, err
	}
	idxs, err := s.Store(idxsKey)
	if err != nil {
		return nil, err
	}

	if init == nil {
		t := reflector.Type[E]()
		init = func() E {
			return reflector.Make(t).Interface().(E)
		}
	}
	es := &EntStore[E]{
		Store:        ents,
		Serializer:   b.Serializer,
		Deserializer: b.Deserializer,
		IdxStore:     idxs,
		Init:         init,
	}
	AddGetter(es)

	return es, nil
}

func NewJsonBuilder(f store.Factory) Builder {
	return Builder{
		Factory:      f,
		Serializer:   json.NewSerializer("", ""),
		Deserializer: json.Deserializer{},
	}
}

func NewGobBuilder(f store.Factory) Builder {
	return Builder{
		Factory:      f,
		Serializer:   gob.Serializer{},
		Deserializer: gob.Deserializer{},
	}
}
