package corpus

import (
	"github.com/adamcolton/luce/ds/document"
	"github.com/adamcolton/luce/ds/prefix"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/serial"
)

type CorpusRef = entity.Ref[Corpus, *Corpus]

var (
	corpusDeserializer   func(data []byte) (*CorpusBale, error)
	rootDeserializer     func(data []byte) (*RootBale, error)
	documentDeserializer func(data []byte) (*DocumentBale, error)
	serializer           serial.PrefixSerializer
)

func init() {
	entity.AddDeserializerListener(func() {
		d := entity.GetDeserializer()
		corpusDeserializer = serial.DeserializeTypeCheck[*CorpusBale](d)
		rootDeserializer = serial.DeserializeTypeCheck[*RootBale](d)
		documentDeserializer = serial.DeserializeTypeCheck[*DocumentBale](d)
		r, ok := d.Detyper.(serial.TypeRegistrar)
		if ok {
			r.RegisterType((*CorpusBale)(nil))
			r.RegisterType((*RootBale)(nil))
			r.RegisterType((*DocumentBale)(nil))
			//TODO: move these to their own packages
			r.RegisterType((*DocBaleType)(nil))
			r.RegisterType((*prefix.PrefixBale)(nil))
		}
	})
	entity.AddSerializerListener(func() {
		serializer = entity.GetSerializer()
	})
}

func (c *Corpus) EntKey() entity.Key {
	return c.key
}

func (c *Corpus) EntVal(buf []byte) ([]byte, error) {
	return serializer.SerializeType(c.Bale(), buf)
}

func (c *Corpus) EntInit() {
	c.SetDefaults()
}

func (c *Corpus) EntLoad(k entity.Key, data []byte) error {
	bale, err := corpusDeserializer(data)
	if err != nil {
		return err
	}

	bale.UnbaleTo(c)
	c.key = k
	c.ref = entity.Put(c)
	return nil
}

func (c *Corpus) Save() (*CorpusRef, error) {
	if !c.save {
		c.prefix.Save(nil)
		c.docs.Each(func(key document.ID, d *docRef, done *bool) {
			d.Save(nil)
		})
		c.save = true
	}
	c.ref = entity.Put(c)
	return c.ref, c.ref.Save(c)
}

func (c *Corpus) entSave() {
	entity.Save(c)
}

type saver interface {
	entSave()
}

func (c *Corpus) saveIf(toSave ...saver) {
	if len(toSave) == 0 {
		panic("bad saveIf")
	}
	if c.save {
		for _, s := range toSave {
			s.entSave()
		}
	}
}

func (d *Document) EntKey() entity.Key {
	return d.DocType.Key
}

func (d *Document) EntVal(buf []byte) ([]byte, error) {
	return serializer.SerializeType(d.Bale(), buf)
}

func (d *Document) EntLoad(k entity.Key, data []byte) error {
	bale, err := documentDeserializer(data)
	if err != nil {
		return err
	}

	bale.UnbaleTo(d)
	d.Key = k
	return nil
}
