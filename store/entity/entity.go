package entity

import (
	"fmt"
	"reflect"

	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/store"
)

type Entity interface {
	EntKey() []byte
}

type EntStore struct {
	store.Store
	Pather
	serial.Serializer
	serial.Deserializer
}

type Pather interface {
	EntPath(Entity) ([][]byte, bool)
}

type ErrPathNotFound struct {
	e Entity
}

func (err ErrPathNotFound) Error() string {
	return fmt.Sprintf("Path not found for %s", reflect.TypeOf(err.e))
}

func (es *EntStore) Put(e Entity, buf []byte) ([]byte, error) {
	s, err := es.PathStore(e)
	if err != nil {
		return nil, err
	}

	v, err := es.Serialize(e, buf)
	if err != nil {
		return nil, err
	}
	return v, s.Put(e.EntKey(), v)
}

func (es *EntStore) PathStore(e Entity) (store.Store, error) {
	ps, ok := es.EntPath(e)
	if !ok {
		return nil, ErrPathNotFound{e}
	}
	var err error
	s := es.Store
	for _, p := range ps {
		s, err = s.Store(p)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func (es *EntStore) Load(e Entity) error {
	return es.Get(e.EntKey(), e)
}

func (es *EntStore) Get(key []byte, e Entity) error {
	s, err := es.PathStore(e)
	if err != nil {
		return err
	}

	r := s.Get(key)
	if r.Value == nil {
		return fmt.Errorf("key not found")
	}

	return es.Deserialize(e, r.Value)
}

type KeyFilter func([]byte) bool

type EntKeySetter interface {
	SetEntKey([]byte)
}

var ErrGetSliceType = fmt.Errorf("entities must be a pointer to a slice of elements that fulfil entity")

func (es *EntStore) GetSlice(fn KeyFilter, entities interface{}) error {
	ptr := reflect.ValueOf(entities)
	if ptr.Kind() != reflect.Ptr {
		return ErrGetSliceType
	}
	slc := ptr.Elem()
	if slc.Kind() != reflect.Slice {
		return ErrGetSliceType
	}

	el := slc.Type().Elem()
	elPtr := el.Kind() == reflect.Ptr
	if elPtr {
		el = el.Elem()
	}
	i := reflect.New(el).Interface()
	zero, ok := i.(Entity)
	if !ok {
		return ErrGetSliceType
	}
	s, err := es.PathStore(zero)
	if err != nil {
		return err
	}

	for key := s.Next(nil); key != nil; key = s.Next(key) {
		if fn != nil && !fn(key) {
			continue
		}
		r := s.Get(key)
		if r.Value != nil {
			ev := reflect.New(el)
			ei := ev.Interface()
			err = es.Deserialize(ei, r.Value)
			if err != nil {
				return err
			}
			if set, ok := ei.(EntKeySetter); ok {
				set.SetEntKey(key)
			}
			if elPtr {
				slc = reflect.Append(slc, ev)
			} else {
				slc = reflect.Append(slc, ev.Elem())
			}
		}
	}
	ptr.Elem().Set(slc)
	return nil
}
