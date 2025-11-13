package corpus

import (
	"github.com/adamcolton/luce/ds/document"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial"
)

type CorpusRef = entity.Ref[Corpus, *Corpus]

var (
	corpusDeserializer   func(data []byte) (*CorpusBale, error)
	rootDeserializer     func(data []byte) (*RootBale, error)
	documentDeserializer func(data []byte) (*DocumentBale, error)
	serializer           serial.Serializer
)

func init() {
	entity.RegisterListener(func() {
		d := entity.GetDeserializer()
		corpusDeserializer = func(data []byte) (*CorpusBale, error) {
			out := &CorpusBale{}
			err := d.Deserialize(out, data)
			return out, err
		}
		rootDeserializer = func(data []byte) (*RootBale, error) {
			out := &RootBale{}
			err := d.Deserialize(out, data)
			return out, err
		}
		documentDeserializer = func(data []byte) (*DocumentBale, error) {
			out := &DocumentBale{}
			err := d.Deserialize(out, data)
			return out, err
		}
		r, ok := entity.GetTyper().(serial.TypeRegistrar)
		if ok {
			lerr.Panic(r.RegisterType((*Corpus)(nil)))
			lerr.Panic(r.RegisterType((*root)(nil)))
			lerr.Panic(r.RegisterType((*Document)(nil)))
			//TODO: move these to their own packages
			//r.RegisterType((*DocBaleType)(nil))
			//r.RegisterType((*prefix.PrefixBale)(nil))
		}
		serializer = entity.GetSerializer()

	})
}

func (c *Corpus) EntKey() entity.Key {
	return c.key
}

func (c *Corpus) EntVal(buf []byte) ([]byte, error) {
	return serializer.Serialize(c.Bale(), buf)
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

func (*Corpus) EntRefs(data []byte) ([]entity.Key, error) {
	bale, err := corpusDeserializer(data)
	return bale.EntRefs(), err
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

func (*Document) EntRefs(data []byte) ([]entity.Key, error) {
	return nil, nil
}

func (d *Document) EntVal(buf []byte) ([]byte, error) {
	return serializer.Serialize(d.Bale(), buf)
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
