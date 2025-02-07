package prefix

import (
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/serial"
)

var (
	prefixDeserializer func(data []byte) (*PrefixBale, error)
	serializer         serial.TypeSerializer
)

func init() {
	entity.AddDeserializerListener(func() {
		d := entity.GetDeserializer()
		prefixDeserializer = serial.DeserializeTypeCheck[*PrefixBale](d)
		r, ok := d.Detyper.(serial.TypeRegistrar)
		if ok {
			r.RegisterType((*NodeBale)(nil))
			r.RegisterType((*PrefixBale)(nil))
		}
	})
	entity.AddSerializerListener(func() {
		serializer = entity.GetSerializer()
	})
	baleChildren = lmap.TransformVal[rune](lmap.ForAll((*node).bale))
	unbaleChildren = lmap.TransformVal[rune](lmap.ForAll((*NodeBale).unbale))
}

func (p *Prefix) EntKey() entity.Key {
	return p.key
}

func (p *Prefix) EntVal(buf []byte) ([]byte, error) {
	p.Purge()
	return serializer.SerializeType(p.Bale(), buf)
}

func (p *Prefix) Save() (*entity.Ref[Prefix, *Prefix], error) {
	er, err := entity.Save(p)
	p.save = true
	return er, err
}

func (p *Prefix) EntLoad(k entity.Key, data []byte) error {
	bale, err := prefixDeserializer(data)
	if err != nil {
		return err
	}
	bale.UnbaleTo(p)
	p.key = k

	return nil
}
