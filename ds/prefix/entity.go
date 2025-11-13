package prefix

import (
	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/serial"
)

var (
	deserializer serial.Deserializer
	serializer   serial.Serializer
	initCh       = make(chan struct{})
)

func init() {
	entity.RegisterListener(func() {
		t := entity.GetTyper()
		r, ok := t.(serial.TypeRegistrar)
		if ok {
			r.RegisterType((*node)(nil))
			r.RegisterType((*Prefix)(nil))
		}
		serializer = entity.GetSerializer()
		deserializer = entity.GetDeserializer()
		close(initCh)
	})
	baleChildren = lmap.TransformVal[rune](lmap.ForAll((*node).bale))
	unbaleChildren = lmap.TransformVal[rune](lmap.ForAll((*NodeBale).unbale))
}

func Wait() {
	<-initCh
}

func (p *Prefix) EntKey() entity.Key {
	return p.key
}

func (p *Prefix) EntVal(buf []byte) ([]byte, error) {
	p.Purge()
	return serializer.Serialize(p.Bale(), buf)
}

func (p *Prefix) Save() (*entity.Ref[Prefix, *Prefix], error) {
	er, err := entity.Save(p)
	p.save = true
	return er, err
}

func (p *Prefix) EntRefs(data []byte) ([]entity.Key, error) {
	return nil, nil
}

func (p *Prefix) EntLoad(k entity.Key, data []byte) error {
	bale := &PrefixBale{}
	err := deserializer.Deserialize(bale, data)
	if err != nil {
		return err
	}
	bale.UnbaleTo(p)
	p.key = k

	return nil
}

func (*Prefix) TypeID32() uint32 {
	return 1399013115
}
