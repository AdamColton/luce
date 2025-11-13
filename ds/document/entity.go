package document

import (
	"reflect"

	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/util/reflector"
)

type docTypeKey struct {
	WordID, VariantID reflect.Type
}

var (
	deserializers map[docTypeKey]any
	typeRegistrar serial.TypeRegistrar
)

func init() {
	entity.RegisterListener(func() {
		deserializers = make(map[docTypeKey]any)
		typer := entity.GetTyper()
		typeRegistrar, _ = typer.(serial.TypeRegistrar)
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
	if typeRegistrar != nil {
		lerr.Panic(typeRegistrar.RegisterType((*Document[WordID, VariantID])(nil)))
	}
	return func(data []byte) (*DocumentBale[WordID, VariantID], error) {
		out := &DocumentBale[WordID, VariantID]{}
		err := d.Deserialize(out, data)
		return out, err
	}
}

func (doc *Document[WordID, VariantID]) EntVal(buf []byte) ([]byte, error) {
	s := entity.GetSerializer()
	//TODO: this is redundant, only needs to happen once
	if typeRegistrar != nil {
		lerr.Panic(typeRegistrar.RegisterType((*Document[WordID, VariantID])(nil)))
	}
	return s.Serialize(doc.Bale(), buf)
}

func (doc *Document[WordID, VariantID]) EntRefs(data []byte) ([]entity.Key, error) {
	return nil, nil
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

func (doc *Document[WordID, VariantID]) TypeID32() uint32 {
	k := docTypeKey{
		WordID:    reflector.Type[WordID](),
		VariantID: reflector.Type[VariantID](),
	}
	id, _ := typeIDs.B(k)
	return id
}
