package filestore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"

	"github.com/adamcolton/luce/store"
)

type Encoder func([]byte) string
type Decoder func(string) []byte

type dir struct {
	path                string
	encoder, bktEncoder Encoder
	decoder, bktDecoder Decoder
}

func Store(path string, encoder, bktEncoder Encoder, decoder, bktDecoder Decoder) (store.Store, error) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return nil, err
	}
	if encoder == nil {
		encoder = EncoderCast
	}
	if bktEncoder == nil {
		bktEncoder = encoder
	}
	if decoder == nil {
		decoder = DecoderCast
	}
	if bktDecoder == nil {
		bktDecoder = decoder
	}
	return &dir{
		path:       path,
		encoder:    encoder,
		bktEncoder: bktEncoder,
		decoder:    decoder,
		bktDecoder: bktDecoder,
	}, nil
}

func Factory(path string, encoder, bktEncoder Encoder, decoder, bktDecoder Decoder) (store.Store, error) {
	return Store(path, encoder, bktEncoder, decoder, bktDecoder)
}

func (d *dir) encode(key []byte) string {
	if d.encoder != nil {
		return d.encoder(key)
	}
	return string(key)

}

func (d *dir) bkt(path string) *dir {
	return &dir{
		path:       path,
		encoder:    d.encoder,
		bktEncoder: d.bktEncoder,
		decoder:    d.decoder,
		bktDecoder: d.bktDecoder,
	}
}

func (d *dir) Store(bkt []byte) (store.Store, error) {
	n := path.Join(d.path, d.bktEncoder(bkt))
	s, _ := os.Stat(n)
	if s != nil && !s.IsDir() {
		return nil, fmt.Errorf("Value already exists at that key")
	}
	os.MkdirAll(n, 0777)
	return d.bkt(n), nil
}

func (d *dir) Put(key, value []byte) error {
	n := path.Join(d.path, d.encode(key))
	s, _ := os.Stat(n)
	if s != nil && s.IsDir() {
		return fmt.Errorf("Bucket already exists at that key")
	}
	f, err := os.Create(n)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(value)
	return err
}

func (d *dir) Get(key []byte) store.Record {
	n := path.Join(d.path, d.encode(key))
	r := store.Record{}
	s, err := os.Stat(n)
	if err != nil {
		return r
	}
	if s.IsDir() {
		r.Store = d.bkt(n)
		r.Found = true
		return r
	}

	f, err := os.Open(n)
	if err != nil {
		return r
	}
	defer f.Close()
	r.Value, _ = ioutil.ReadAll(f)
	r.Found = true
	return r
}

func (d *dir) Next(key []byte) []byte {
	n := d.encode(key)
	files, err := ioutil.ReadDir(d.path)
	if err != nil {
		return nil
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
	i := sort.Search(len(files), func(i int) bool {
		return files[i].Name() > n
	})
	if i < len(files) {
		n := files[i].Name()
		info, _ := os.Stat(path.Join(d.path, n))
		if info.IsDir() {
			return d.bktDecoder(n)
		}
		return d.decoder(n)
	}
	return nil
}

func (d *dir) Delete(key []byte) error {
	n := path.Join(d.path, d.encode(key))
	os.RemoveAll(n)
	return nil
}
