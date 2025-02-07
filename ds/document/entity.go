package document

import (
	"reflect"

	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/util/reflector"
)

type docTypeKey struct {
	WordID, VariantID reflect.Type
}

var (
	deserializers map[docTypeKey]any
)

func init() {
	entity.AddDeserializerListener(func() {
		deserializers = make(map[docTypeKey]any)
	})
}

func getDeserializer[WordID, VariantID comparable]() func(data []byte) (*DocumentBale[WordID, VariantID], error) {
	dtk := docTypeKey{
		WordID:    reflector.Type[WordID](),
		VariantID: reflector.Type[VariantID](),
	}
	got, ok := deserializers[dtk]
	if ok {
		return got.(func(data []byte) (*DocumentBale[WordID, VariantID], error))
	}

	d := entity.GetDeserializer()
	docDeserializer := serial.DeserializeTypeCheck[*DocumentBale[WordID, VariantID]](d)
	r, ok := d.Detyper.(serial.TypeRegistrar)
	if ok {
		r.RegisterType((*DocumentBale[WordID, VariantID])(nil))
	}
	return docDeserializer
}

func (doc *Document[WordID, VariantID]) EntVal(buf []byte) ([]byte, error) {
	s := entity.GetSerializer()
	//TODO: this is redundant, only needs to happen once
	r, ok := s.InterfaceTypePrefixer.(serial.TypeRegistrar)
	if ok {
		r.RegisterType((*DocumentBale[WordID, VariantID])(nil))
	}
	return s.SerializeType(doc.Bale(), buf)
}

func (doc *Document[WordID, VariantID]) EntLoad(k entity.Key, data []byte) error {
	doc.Key = k
	doc.save = true
	d := getDeserializer[WordID, VariantID]()
	bale, err := d(data)
	if err != nil {
		return err
	}
	bale.UnbaleTo(doc)
	return nil
}
