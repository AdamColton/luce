package enttest

import (
	"strconv"
	"time"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/serial/wrap/gob"
	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/ephemeral"
)

type Foo struct {
	ID   []byte
	Name string
	entity.Refs
}

func (f *Foo) EntKey() entity.Key {
	return f.ID
}

func (f *Foo) EntVal(buf []byte) ([]byte, error) {
	return entity.GetSerializer().Serialize(f, buf)
}

func (f *Foo) TypeID32() uint32 {
	return 3918102473
}

func (f *Foo) EntLoad(k entity.Key, data []byte) error {
	d := entity.GetDeserializer()
	return d.Deserialize(f, data)

	// TODO: bring this back
	//check := serial.DeserializeToTypeCheck[*Foo](d)
	//return check(f, data)
}

const ErrTimeout = lerr.Str("timeout")

func SaveAndWait[T any, E entity.EntPtr[T]](ref *entity.Ref[T, E], now bool) error {
	err := ref.Save(nil)
	if err != nil {
		return err
	}
	return Wait(ref.EntKey())
}

func Wait(key entity.Key) error {
	for range 200 {
		time.Sleep(time.Millisecond)
		if EntStore.Get(key).Found {
			return nil
		}
	}
	return ErrTimeout
}

var EntStore store.FlatStore

func Setup() type32.TypeMap {
	m32 := type32.NewTypeMap()
	EntStore = lerr.Must(ephemeral.Factory(bytebtree.New, 10).FlatStore([]byte("testing")))
	entity.Setup{
		Store:        EntStore,
		Typer:        m32,
		Serializer:   gob.Serializer{},
		Deserializer: gob.Deserializer{},
	}.Init()
	lerr.Panic(serial.RegisterPtr[Foo](m32))
	lerr.Panic(serial.RegisterPtr[String](m32))
	lerr.Panic(serial.RegisterPtr[Int](m32))
	return m32
}

type Stringer interface {
	String() string
	entity.Entity
}

type String struct {
	entity.Key32
	Str string
}

func (s *String) TypeID32() uint32 {
	return 3350641450
}

func (s *String) String() string {
	return s.Str
}

func (s *String) EntVal(buf []byte) ([]byte, error) {
	return entity.GetSerializer().Serialize(s, buf)
}

func (s *String) EntLoad(k entity.Key, data []byte) error {
	d := entity.GetDeserializer()

	//TODO: bring this back
	//check := serial.DeserializeToTypeCheck[*String](d)
	return d.Deserialize(s, data)
}

type Int struct {
	entity.Key32
	I int
}

func (i *Int) TypeID32() uint32 {
	return 2065445151
}

func (i *Int) String() string {
	return strconv.Itoa(i.I)
}

func (i *Int) EntVal(buf []byte) ([]byte, error) {
	return entity.GetSerializer().Serialize(i, buf)
}

func (i *Int) EntLoad(k entity.Key, data []byte) error {
	d := entity.GetDeserializer()
	return d.Deserialize(i, data)

	// Todo: make this work
	//check := serial.DeserializeToTypeCheck[*Int](d)
	//return check(i, data)
}
