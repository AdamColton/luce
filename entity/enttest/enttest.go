package enttest

import (
	"time"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/serial/wrap/gob"
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
	return entity.GetSerializer().SerializeType(f, buf)
}

func (f *Foo) TypeID32() uint32 {
	return 3918102473
}

func (f *Foo) EntLoad(k entity.Key, data []byte) error {
	d := entity.GetDeserializer()
	check := serial.DeserializeToTypeCheck[*Foo](d)
	return check(f, data)
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
	for i := 0; i < 200; i++ {
		time.Sleep(time.Millisecond)
		if entity.Store.Get(key).Found {
			return nil
		}
	}
	return ErrTimeout
}

func Setup() type32.TypeMap {
	entity.Store = lerr.Must(ephemeral.Factory(bytebtree.New, 10).FlatStore([]byte("testing")))
	m32 := type32.NewTypeMap()
	entity.SetSerializer(m32.Serializer(gob.Serializer{}))
	entity.SetDeserializer(m32.Deserializer(gob.Deserializer{}))
	lerr.Panic(serial.RegisterPtr[Foo](m32))
	return m32
}
