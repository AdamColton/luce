package thresher

import (
	"hash/crc64"
	"reflect"

	"github.com/adamcolton/luce/serial/rye"
)

type Skip struct {
	Model   interface{}
	Exclude []string
}

type field struct {
	name string
	kind uint64
}

func (f field) serialize() []byte {
	s := rye.Serializer{
		Data: make([]byte, len(f.name)+8),
	}
	s.Uint64(f.kind)
	s.String(f.name)
	return s.Data
}

func deserializeField(data []byte) field {
	d := rye.Deserializer{
		Data: data,
	}
	return field{
		kind: d.Uint64(),
		name: d.String(len(data) - 8),
	}
}

var tab = crc64.MakeTable(crc64.ISO)

func (f field) id() uint64 {
	return crc64.Checksum(f.serialize(), tab)
}

func typeID(t reflect.Type) uint64 {
	switch k := t.Kind(); k {
	case reflect.Ptr:
		f := field{
			name: "*",
			kind: typeID(t.Elem()),
		}
		return f.id()
	default:
		return uint64(k)
	}
}
