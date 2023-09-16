package filestore

import (
	"bytes"
	"os"
	"path/filepath"

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
// Every key needs its own name, and the name is used as a single path element,
// so it must not contain a path separator or be "." or "..".
type Encoder func([]byte) string

// Decoder is used to convert a directory or file path to a key. It should be
// the inverse of the Encoder, and return nil for a name that can't be decoded.
type Decoder func(string) []byte

// dir is a store.NestedStore held in a directory. A value is a file, named by
// the encoder, and a Sub-Store is a sub-directory, named by the bktEncoder.
type dir struct {
	path string
	*factory
}

// factory holds the encoders and decoders shared by every dir it creates.
type factory struct {
	encoder, bktEncoder Encoder
	decoder, bktDecoder Decoder
}

// NestedStore returns the store for the directory path, creating the directory
// if needed. The path is used as it is, so a relative path starts from the
// working directory.
func (f *factory) NestedStore(path []byte) (store.NestedStore, error) {
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

// NewFactory creates a factory with the given encoders and decoders. The
// encoder and decoder name and read values, and the bktEncoder and bktDecoder
// name and read Sub-Stores. If bktEncoder is nil, encoder will be used. If
// encoder is nil, EncoderCast is used. If bktDecoder is nil, decoder is used.
// If decoder is nil, DecoderCast is used. Each Decoder should invert its
// Encoder.
func NewFactory(encoder, bktEncoder Encoder, decoder, bktDecoder Decoder) store.NestedFactory {
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
		bktEncoder: bktEncoder,
		decoder:    decoder,
		bktDecoder: bktDecoder,
	}
}

func (d *dir) bkt(path string) *dir {
	return &dir{
		path:    path,
		factory: d.factory,
	}
}

// stat returns nil if there is nothing at the path.
func stat(path string) os.FileInfo {
	info, _ := os.Stat(path)
	return info
}

// lookup finds what is stored at key. A value is a file named by the encoder
// and a Sub-Store is a directory named by the bktEncoder. It returns the path
// and whether it is a Sub-Store.
func (d *dir) lookup(key []byte) (n string, isBkt, found bool) {
	n = filepath.Join(d.path, d.encoder(key))
	if info := stat(n); info != nil && !info.IsDir() {
		return n, false, true
	}
	n = filepath.Join(d.path, d.bktEncoder(key))
	if info := stat(n); info != nil && info.IsDir() {
		return n, true, true
	}
	return "", false, false
}

// Len counts the values and Sub-Stores in the directory. Store.Len can't return
// an error, so it panics if the directory can't be read.
func (d *dir) Len() int {
	files, e := os.ReadDir(d.path)
	if e != nil {
		panic(e)
	}
	return len(files)
}

// NestedStore gets the Sub-Store at bkt, creating its directory if needed. It
// returns ErrValExists if bkt is the key of a value.
func (d *dir) NestedStore(bkt []byte) (store.NestedStore, error) {
	if info := stat(filepath.Join(d.path, d.encoder(bkt))); info != nil && !info.IsDir() {
		return nil, ErrValExists
	}
	n := filepath.Join(d.path, d.bktEncoder(bkt))
	if err := os.MkdirAll(n, 0777); err != nil {
		return nil, err
	}
	return d.bkt(n), nil
}

// Put a key, value pair. It creates the file or replaces its contents, and
// returns ErrBktExists if key is a Sub-Store. Fulfills store.Store.
func (d *dir) Put(key, value []byte) error {
	if info := stat(filepath.Join(d.path, d.bktEncoder(key))); info != nil && info.IsDir() {
		return ErrBktExists
	}
	return os.WriteFile(filepath.Join(d.path, d.encoder(key)), value, 0666)
}

// Get a Record. The Record holds the Sub-Store if key is one. A file that can't
// be read is not found. Fulfills store.Store.
func (d *dir) Get(key []byte) store.Record {
	n, isBkt, found := d.lookup(key)
	if !found {
		return store.Record{}
	}
	if isBkt {
		return store.Record{Found: true, Store: d.bkt(n)}
	}
	v, err := os.ReadFile(n)
	return store.Record{Found: err == nil, Value: v}
}

// Next returns the smallest key greater than key, or the lowest key if key is
// nil. Every name in the directory is decoded to compare the keys, so it takes
// O(n), and the order is by key, not by name. Names that decode to nil are
// skipped. It returns nil if there is no such key or the directory can't be
// read. Fulfills store.Store.
func (d *dir) Next(key []byte) []byte {
	files, err := os.ReadDir(d.path)
	if err != nil {
		return nil
	}
	var next []byte
	for _, file := range files {
		var k []byte
		if file.IsDir() {
			k = d.bktDecoder(file.Name())
		} else {
			k = d.decoder(file.Name())
		}
		if k == nil || (key != nil && bytes.Compare(k, key) <= 0) {
			continue
		}
		if next == nil || bytes.Compare(k, next) < 0 {
			next = k
		}
	}
	return next
}

// Delete key, or the Sub-Store it names along with everything in it. Deleting a
// key that is not in the store is not an error. Fulfills store.Store.
func (d *dir) Delete(key []byte) error {
	n, _, found := d.lookup(key)
	if !found {
		return nil
	}
	return os.RemoveAll(n)
}
