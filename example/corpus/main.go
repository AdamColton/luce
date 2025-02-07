package main

import (
	"fmt"
	"io"
	"time"

	"github.com/adamcolton/luce/ds/document"
	"github.com/adamcolton/luce/ds/document/corpus"
	"github.com/adamcolton/luce/entity"
	"github.com/adamcolton/luce/entity/entdefertimer"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/serial/wrap/gob"
	"github.com/adamcolton/luce/store/badgerstore"
	"github.com/adamcolton/luce/util/lfile"
)

type Config struct {
	Root       string
	HiddenDirs bool
	Exts       []string
}

var defertoq = entdefertimer.NewToq(time.Millisecond * 10)

func main() {
	cfg := &Config{}
	err := lfile.JsonConfig("", "config.json", cfg)
	lerr.Panic(err)

	//f := bstore.Factory("store.bolt", 0777, nil)
	//entity.Store = lerr.Must(f.NestedStore([]byte("entstore")))

	f := badgerstore.Factory("badger")
	entity.Store = lerr.Must(f.FlatStore([]byte("badger.db")))

	entity.DeferStrategy = defertoq

	// Todo: forgetting to setup the serializers does not product a useful error
	m32 := type32.NewTypeMap()
	entity.SetSerializer(m32.Serializer(gob.Serializer{}))
	entity.SetDeserializer(m32.Deserializer(gob.Deserializer{}))
	m32.RegisterType((*ID2Path)(nil))

	ref := entity.KeyRef[corpus.Corpus](entity.Key("corpus-root"))
	c, found := ref.Get()
	var i2p *ID2Path
	if !found {
		c, i2p = populateCorpus(ref.EntKey(), cfg)
	} else {
		var ok bool
		i2p, ok = entity.KeyRef[ID2Path](i2pkey).Get()
		if !ok {
			panic("not found")
		}
	}

	ids := c.Find("scanner")
	for _, id := range ids.Slice(nil) {
		fmt.Println(i2p.Map[id])
	}

	total := 0
	count := 0
	for k := entity.Store.Next(nil); k != nil; k = entity.Store.Next(k) {
		r := entity.Store.Get(k)
		total += len(r.Value)
		count++
	}
	fmt.Println("Total Bytes Written:", total)
	fmt.Println("Records:", count)

	lerr.Panic(entity.Store.(io.Closer).Close())
}

type Syncer interface {
	Sync() error
}

func populateCorpus(key entity.Key, cfg *Config) (*corpus.Corpus, *ID2Path) {
	c := corpus.NewKey(key)
	total := 0
	i2p := &ID2Path{
		Map: make(map[document.ID]string),
	}
	hdlr := lfile.IterHandlerFn(func(i lfile.Iterator) {
		d := i.Data()
		total += len(d)
		doc := c.AddDoc(string(d))
		i2p.Map[doc.ID] = i.Path()
	})

	exts := lfile.Exts(cfg.HiddenDirs, cfg.Exts...)
	mtch := lfile.MustRegexMatch(exts, "", "")
	mr := mtch.Root(cfg.Root)
	err := lfile.RunHandlerSource(mr, hdlr)
	lerr.Panic(err)
	fmt.Println("Total Bytes in Directory:", total)

	c.Save()
	entity.Save(i2p)

	start := time.Now()
	ms := time.Millisecond
	for time.Sleep(ms); !defertoq.Done(); time.Sleep(ms) {
	}
	lerr.Panic(entity.Store.(Syncer).Sync())
	fmt.Println(int(time.Since(start)/time.Second), "seconds")

	return c, i2p
}

// ID2Path maps document IDs to the path of the document.
// This is also a good example of a basic entity
type ID2Path struct {
	Map map[document.ID]string
}

func (i2p *ID2Path) TypeID32() uint32 {
	return 2107647393
}

var i2pkey = entity.Key("id2path")

func (i2p *ID2Path) EntKey() entity.Key {
	return i2pkey
}

func (i2p *ID2Path) EntVal(buf []byte) ([]byte, error) {
	return entity.GetSerializer().SerializeType(i2p, buf)
}

func (i2p *ID2Path) EntRefs() []entity.Key {
	return nil
}

func (i2p *ID2Path) EntLoad(k entity.Key, data []byte) error {
	i, err := entity.GetDeserializer().DeserializeType(data)
	if err != nil {
		return err
	}
	*i2p = *(i.(*ID2Path))
	return nil
}
