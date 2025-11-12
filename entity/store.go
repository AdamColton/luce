package entity

import (
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/util/reflector"
)

type Typer interface {
	serial.Detyper
	serial.InterfaceTypePrefixer
}

var (
	entstore     store.FlatStore
	serializer   serial.Serializer
	deserializer serial.Deserializer

	typer Typer

	listeners []func()

	ErrMustBeRefser = lerr.Str("must fulfill type entity.Refser")
)

type Setup struct {
	Store        store.FlatStore
	Serializer   serial.Serializer
	Deserializer serial.Deserializer
	Typer        Typer
}

func (s Setup) Init() {
	//TODO: validate none nil
	//validate only call once
	entstore = s.Store
	serializer = s.Serializer
	deserializer = s.Deserializer
	typer = s.Typer
	for _, fn := range listeners {
		fn()
	}
}

func RegisterListener(fn func()) {
	listeners = append(listeners, fn)
}

func GetDeserializer() serial.Deserializer {
	return deserializer
}

func GetSerializer() serial.Serializer {
	return serializer
}

func GetTyper() Typer {
	return typer
}

const (
	ErrNilStore  = lerr.Str("entity.Store is nil")
	ErrNilEntRef = lerr.Str("entity pointer in EntRef is nil")
	ErrNoRecord  = lerr.Str("no record exists in the store for the given key")
)

type Referer interface {
	save(now bool) error
	EntKey() Key
	Clear(cacheRm bool)
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
	if entstore == nil {
		return ErrNilStore
	}
	e, ok := er.Get()
	if !ok {
		return ErrNilEntRef
	}
	if now {
		data := lerr.Must(typer.PrefixInterfaceType(e, nil))
		data = append(data, lerr.Must(e.EntVal(nil))...)
		lerr.Panic(entstore.Put(er.key, data))
	} else {
		DeferStrategy.DeferSave(er, er.saveNow)
	}

	return nil
}

func (er *Ref[T, E]) Delete() error {
	if entstore == nil {
		return ErrNilStore
	}
	err := entstore.Delete(er.key)
	if err != nil {
		return err
	}
	er.ent.p = nil
	er.allRefsRm()
	er.key = nil
	return nil
}

func (er *Ref[T, E]) load() error {
	if entstore == nil {
		return ErrNilStore
	}

	r := entstore.Get(er.key)
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
	_, data, _ := typer.GetType(r.Value)
	err := e.EntLoad(er.key, data)
	if err == nil {
		Put(e)
	}
	return err
}

func Save[T any, E EntPtr[T]](ent E) (*Ref[T, E], error) {
	er := Put(ent)
	return er, er.save(false)
}

func Load(k Key) (Entity, error) {
	ent, found := GetEnt(k)
	if found {
		return ent, nil
	}

	if entstore == nil {
		return nil, ErrNilStore
	}

	r := entstore.Get(k)
	if !r.Found {
		return nil, ErrNoRecord
	}

	t, data, err := typer.GetType(r.Value)
	if err != nil {
		return nil, err
	}
	i := reflector.Make(t).Interface().(Entity)
	//TODO: call addToAllRefs??
	if ei, ok := i.(EntIniter); ok {
		ei.EntInit()
	}

	err = i.EntLoad(k, data)
	if err != nil {
		return nil, err
	}
	return i, nil
	// if err == nil {
	// 	Put(i)
	// }
}
