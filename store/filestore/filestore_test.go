package filestore_test

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/adamcolton/luce/store/filestore"
	"github.com/adamcolton/luce/store/testsuite"
	"github.com/stretchr/testify/assert"
)

// TODO: move os mock logic before this and use that for testing
func TestAll(t *testing.T) {
	name := []byte("test")
	decoder := func(s string) []byte {
		b, _ := base64.RawURLEncoding.DecodeString(s)
		return b
	}
	enc := base64.RawURLEncoding.EncodeToString
	f, err := filestore.NewFactory(enc, enc, decoder, decoder).Store(name)
	assert.NoError(t, err)
	defer func() {
		os.RemoveAll(string(name))
	}()
	testsuite.TestAll(t, f)
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
