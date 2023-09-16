package filestore

import (
	"io/ioutil"
	"os"
	"path"
	"sort"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store"
)

const (
	// ErrBktExists is returned when attempting to Put to key that defines a
	// bucket.
	ErrBktExists = lerr.Str("Bucket already exists at that key")

	// ErrValExists is returned when attempting to create a Store with a key
	// that defines a value.
	ErrValExists = lerr.Str("Value already exists at that key")
)

// Encoder is used to convert []byte to string for directory and file paths.
type Encoder func([]byte) string

// Decoder is used to convert a directory or file path to a key.
type Decoder func(string) []byte

type dir struct {
	path string
	*factory
}

type factory struct {
	encoder, bktEncoder Encoder
	decoder, bktDecoder Decoder
}

// Store fulfills store.Factory.
func (f *factory) Store(path []byte) (store.Store, error) {
	return f.dir(string(path))
}

func (f *factory) dir(path string) (*dir, error) {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return nil, err
	}
	return &dir{
		path:    path,
		factory: f,
	}, nil
}

// NewFactory creates a factory with the given encoders and decoders. If
// bktEncoder is nil, encoder will be used. If encoder is nil, EncoderCast is
// used. If bktDecoder is nil, decoder is used. If decoder is nil, DecoderCast
// is used.
func NewFactory(encoder, bktEncoder Encoder, decoder, bktDecoder Decoder) store.Factory {
	return newFactory(encoder, bktEncoder, decoder, bktDecoder)
}

func newFactory(encoder, bktEncoder Encoder, decoder, bktDecoder Decoder) *factory {
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
	return &factory{
		encoder:    encoder,
		bktEncoder: encoder,
		decoder:    decoder,
		bktDecoder: bktDecoder,
	}
}

func (d *dir) encode(key []byte) string {
	if d.encoder != nil {
		return d.encoder(key)
	}
	return string(key)

}

func (d *dir) bkt(path string) *dir {
	return &dir{
		path:    path,
		factory: d.factory,
	}
}

// Store fulfills store.Factory.
func (d *dir) Store(bkt []byte) (store.Store, error) {
	n := path.Join(d.path, d.bktEncoder(bkt))
	s, _ := os.Stat(n)
	if s != nil && !s.IsDir() {
		return nil, ErrValExists
	}
	os.MkdirAll(n, 0777)
	return d.bkt(n), nil
}

// Put a key, value pair. Fulfills store.Store.
func (d *dir) Put(key, value []byte) error {
	n := path.Join(d.path, d.encode(key))
	s, _ := os.Stat(n)
	if s != nil && s.IsDir() {
		return ErrBktExists
	}
	f, err := os.Create(n)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(value)
	return err
}

// Get a Record. Fulfills store.Store.
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

// Next key. Fulfills store.Store.
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

// Delete key. Fulfills store.Store.
func (d *dir) Delete(key []byte) error {
	n := path.Join(d.path, d.encode(key))
	os.RemoveAll(n)
	return nil
}
