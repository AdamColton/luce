package reflidx

import (
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
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
)

func New[E entity.Entity](r Rounder, f store.Factory) *entity.EntStore[E] {
	t := reflector.Type[E]()
	// TODO: something better
	name := t.Name()

	init := func() E {
		return reflector.Make(t).Interface().(E)
	}

	str := lerr.Must(f.Store(EntsBucket))
	str = lerr.Must(str.Store([]byte(name)))

	idx := lerr.Must(f.Store(IndexesBucket))
	idx = lerr.Must(idx.Store([]byte(name)))

	es := &entity.EntStore[E]{
		Init:         init,
		Store:        str,
		Serializer:   r,
		Deserializer: r,
		IdxStore:     idx,
	}

	return es
}
