package filestore_test

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/adamcolton/luce/store"
	"github.com/adamcolton/luce/store/filestore"
	"github.com/adamcolton/luce/store/testsuite"
	"github.com/stretchr/testify/assert"
)

var (
	b64Enc = base64.RawURLEncoding.EncodeToString
	b64Dec = func(s string) []byte {
		b, err := base64.RawURLEncoding.DecodeString(s)
		if err != nil {
			return nil
		}
		return b
	}
)

// names lists the entries of the directory, and whether each is a directory.
func names(t *testing.T, dir string) map[string]bool {
	entries, err := os.ReadDir(dir)
	assert.NoError(t, err)
	out := make(map[string]bool)
	for _, e := range entries {
		out[e.Name()] = e.IsDir()
	}
	return out
}

func TestAll(t *testing.T) {
	// Random keys are not valid file names, so they are encoded.
	root := filepath.Join(t.TempDir(), "test")
	f, err := filestore.NewFactory(b64Enc, b64Enc, b64Dec, b64Dec).NestedStore([]byte(root))
	assert.NoError(t, err)
	testsuite.TestAll(t, f)
}

func TestAllSeparateEncoders(t *testing.T) {
	// Values are files that end in ".val". Sub-Stores are directories that don't.
	enc := filestore.EncoderMany(b64Enc, filestore.EncoderExt("val"))
	dec := func(s string) []byte {
		return b64Dec(string(filestore.DecoderRemoveExt(s)))
	}
	root := filepath.Join(t.TempDir(), "test")
	f, err := filestore.NewFactory(enc, b64Enc, dec, b64Dec).NestedStore([]byte(root))
	assert.NoError(t, err)
	testsuite.TestAll(t, f)
}

func TestSeparateEncoders(t *testing.T) {
	root := t.TempDir()
	f := filestore.NewFactory(filestore.EncoderExt("txt"), filestore.EncoderCast, filestore.DecoderRemoveExt, filestore.DecoderCast)
	s, err := f.NestedStore([]byte(root))
	assert.NoError(t, err)

	assert.NoError(t, s.Put([]byte("apple"), []byte("red")))
	fruit, err := s.NestedStore([]byte("fruit"))
	assert.NoError(t, err)
	assert.NoError(t, fruit.Put([]byte("pear"), []byte("green")))

	// Only the values get the extension.
	assert.Equal(t, map[string]bool{"apple.txt": false, "fruit": true}, names(t, root))
	assert.Equal(t, map[string]bool{"pear.txt": false}, names(t, filepath.Join(root, "fruit")))

	assert.Equal(t, store.Record{Found: true, Value: []byte("red")}, s.Get([]byte("apple")))
	r := s.Get([]byte("fruit"))
	assert.True(t, r.Found)
	assert.Equal(t, []byte("green"), r.Store.Get([]byte("pear")).Value)
	assert.Equal(t, 2, s.Len())

	// Values and Sub-Stores are in one key space.
	assert.Equal(t, filestore.ErrBktExists, s.Put([]byte("fruit"), []byte("value")))
	_, err = s.NestedStore([]byte("apple"))
	assert.Equal(t, filestore.ErrValExists, err)
	assert.Equal(t, []byte("apple"), s.Next(nil))
	assert.Equal(t, []byte("fruit"), s.Next([]byte("apple")))
	assert.Nil(t, s.Next([]byte("fruit")))

	assert.NoError(t, s.Delete([]byte("apple")))
	assert.NoError(t, s.Delete([]byte("fruit")))
	assert.NoError(t, s.Delete([]byte("missing")))
	assert.Empty(t, names(t, root))
}

func TestNextOrder(t *testing.T) {
	// The names of these keys sort in the opposite order to the keys.
	lo, hi := []byte{0x00}, []byte{0xd0}
	assert.Greater(t, b64Enc(lo), b64Enc(hi))

	s, err := filestore.NewFactory(b64Enc, nil, b64Dec, nil).NestedStore([]byte(t.TempDir()))
	assert.NoError(t, err)
	assert.Nil(t, s.Next(nil))
	assert.NoError(t, s.Put(hi, []byte("hi")))
	assert.NoError(t, s.Put(lo, []byte("lo")))

	assert.Equal(t, lo, s.Next(nil))
	assert.Equal(t, hi, s.Next(lo))
	assert.Equal(t, hi, s.Next([]byte{0x01}))
	assert.Nil(t, s.Next(hi))
}

func TestUndecodableName(t *testing.T) {
	root := t.TempDir()
	s, err := filestore.NewFactory(b64Enc, nil, b64Dec, nil).NestedStore([]byte(root))
	assert.NoError(t, err)
	assert.NoError(t, s.Put([]byte("key"), []byte("value")))

	// A file that isn't base64 can't be a key.
	assert.NoError(t, os.WriteFile(filepath.Join(root, "not base64!"), nil, 0600))
	assert.Equal(t, []byte("key"), s.Next(nil))
	assert.Nil(t, s.Next([]byte("key")))
}

func TestErrors(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file")
	assert.NoError(t, os.WriteFile(file, []byte("x"), 0600))
	f := filestore.NewFactory(nil, nil, nil, nil)

	// A directory can't be made where there is a file, or below one.
	s, err := f.NestedStore([]byte(file))
	assert.Error(t, err)
	assert.Nil(t, s)
	s, err = f.NestedStore([]byte(filepath.Join(file, "below")))
	assert.Error(t, err)
	assert.Nil(t, s)

	// The same goes for a Sub-Store. Its name is used as a path, so it can hold
	// a separator.
	s, err = f.NestedStore([]byte(root))
	assert.NoError(t, err)
	sub, err := s.NestedStore([]byte("file/below"))
	assert.Error(t, err)
	assert.Nil(t, sub)
}

func TestMissingDirectory(t *testing.T) {
	root := filepath.Join(t.TempDir(), "gone")
	s, err := filestore.NewFactory(nil, nil, nil, nil).NestedStore([]byte(root))
	assert.NoError(t, err)
	assert.NoError(t, s.Put([]byte("k"), []byte("v")))
	assert.NoError(t, os.RemoveAll(root))

	assert.Error(t, s.Put([]byte("k"), []byte("v")))
	assert.Equal(t, store.Record{}, s.Get([]byte("k")))
	assert.Nil(t, s.Next(nil))
	assert.NoError(t, s.Delete([]byte("k")))
	// Len can't return an error.
	assert.Panics(t, func() { s.Len() })
}

func TestEncoders(t *testing.T) {
	tt := map[string]struct {
		cases map[string]string
		filestore.Encoder
	}{
		"EncoderCast": {
			Encoder: filestore.EncoderCast,
			cases: map[string]string{
				"test": "test",
			},
		},
		"EncoderReplacer": {
			Encoder: filestore.EncoderReplacer("foo", "bar"),
			cases: map[string]string{
				"foo test foo": "bar test bar",
				"test":         "test",
				"fooo":         "baro",
			},
		},
		"EncoderExt": {
			Encoder: filestore.EncoderExt("txt"),
			cases: map[string]string{
				"test": "test.txt",
			},
		},
		"EncoderMany": {
			Encoder: filestore.EncoderMany(
				filestore.EncoderReplacer("foo", "bar"),
				filestore.EncoderExt("txt"),
			),
			cases: map[string]string{
				"foo": "bar.txt",
			},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			for in, out := range tc.cases {
				t.Run(in, func(t *testing.T) {
					assert.Equal(t, out, tc.Encoder([]byte(in)))
				})
			}
		})
	}
}

func TestDecoders(t *testing.T) {
	tt := map[string]struct {
		cases map[string]string
		filestore.Decoder
	}{
		"DecoderCast": {
			Decoder: filestore.DecoderCast,
			cases: map[string]string{
				"test": "test",
			},
		},
		"DecoderRemoveExt": {
			Decoder: filestore.DecoderRemoveExt,
			cases: map[string]string{
				"test.txt": "test",
				"foo":      "foo",
			},
		},
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			for in, out := range tc.cases {
				t.Run(in, func(t *testing.T) {
					assert.Equal(t, []byte(out), tc.Decoder(in))
				})
			}
		})
	}
}
