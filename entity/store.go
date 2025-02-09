package entity

import (
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/store"
)

var (
	Store        store.FlatStore
	serializer   serial.PrefixSerializer
	deserializer serial.PrefixDeserializer

	deserializerListeners []func()
	serializerListeners   []func()

	ErrMustBeRefser = lerr.Str("must fulfill type entity.Refser")
)

func GetDeserializer() serial.PrefixDeserializer {
	return deserializer
}

func GetSerializer() serial.PrefixSerializer {
	return serializer
}

func AddDeserializerListener(fn func()) {
	deserializerListeners = append(deserializerListeners, fn)
}

func AddSerializerListener(fn func()) {
	serializerListeners = append(serializerListeners, fn)
}

func SetSerializer(s serial.PrefixSerializer) {
	serializer = s
	for _, fn := range serializerListeners {
		fn()
	}
}

func SetDeserializer(d serial.PrefixDeserializer) {
	deserializer = d
	for _, fn := range deserializerListeners {
		fn()
	}
}

const (
	ErrNilStore  = lerr.Str("entity.Store is nil")
	ErrNilEntRef = lerr.Str("entity pointer in EntRef is nil")
	ErrNoRecord  = lerr.Str("no record exists in the store for the given key")
)

type Referer interface {
	save(now bool) error
	EntKey() Key
}

const ErrNilPtr = lerr.Str("entity reference has nil pointer")

func (er *Ref[T, E]) Save(ent E) error {
	if er.ent.p == nil {
		if ent == nil {
			return ErrNilPtr
		}
		er.ent.p = ent
	}
	return er.save(false)
}

func (er *Ref[T, E]) saveNow() error {
	return er.save(true)
}

func (er *Ref[T, E]) save(now bool) error {
	if Store == nil {
		return ErrNilStore
	}
	e, ok := er.Get()
	if !ok {
		return ErrNilEntRef
	}
	if now {
		data := lerr.Must(e.EntVal(nil))
		lerr.Panic(Store.Put(er.key, data))
	} else {
		DeferStrategy.DeferSave(er, er.saveNow)
	}

	return nil
}

func (er *Ref[T, E]) Delete() error {
	if Store == nil {
		return ErrNilStore
	}
	err := Store.Delete(er.key)
	if err != nil {
		return err
	}
	er.ent.p = nil
	er.allRefsRm()
	er.key = nil
	return nil
}

func (er *Ref[T, E]) load() error {
	if Store == nil {
		return ErrNilStore
	}

	r := Store.Get(er.key)
	if !r.Found {
		return ErrNoRecord
	}

	if er.ent.p == nil {
		var t T
		er.ent.p = &t
		er.addToAllRefs()
		var a any = er.ent.p
		ei, ok := a.(EntIniter)
		if ok {
			ei.EntInit()
		}
	}

	e := E(er.ent.p)
	err := e.EntLoad(er.key, r.Value)
	if err == nil {
		Put(e)
	}
	return err
}

func Save[T any, E EntPtr[T]](ent E) (*Ref[T, E], error) {
	er := Put(ent)
	return er, er.save(false)
}
