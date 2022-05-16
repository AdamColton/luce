package filestore

import (
	"encoding/base64"
	"os"
	"reflect"
	"testing"

	"github.com/adamcolton/luce/store/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	name := []byte("test")
	decoder := func(s string) []byte {
		b, _ := base64.RawURLEncoding.DecodeString(s)
		return b
	}
	enc := base64.RawURLEncoding.EncodeToString
	f, err := NewFactory(enc, enc, decoder, decoder).Store(name)
	assert.NoError(t, err)
	defer func() {
		os.RemoveAll(string(name))
	}()
	testsuite.TestAll(t, f)
}

func TestEncoders(t *testing.T) {
	tt := map[string]struct {
		cases map[string]string
		Encoder
	}{
		"EncoderCast": {
			Encoder: EncoderCast,
			cases: map[string]string{
				"test": "test",
			},
		},
		"EncoderReplacer": {
			Encoder: EncoderReplacer("foo", "bar"),
			cases: map[string]string{
				"foo test foo": "bar test bar",
				"test":         "test",
				"fooo":         "baro",
			},
		},
		"EncoderExt": {
			Encoder: EncoderExt("txt"),
			cases: map[string]string{
				"test": "test.txt",
			},
		},
		"EncoderMany": {
			Encoder: EncoderMany(
				EncoderReplacer("foo", "bar"),
				EncoderExt("txt"),
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
		Decoder
	}{
		"DecoderCast": {
			Decoder: DecoderCast,
			cases: map[string]string{
				"test": "test",
			},
		},
		"DecoderRemoveExt": {
			Decoder: DecoderRemoveExt,
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

func TestNewFactory(t *testing.T) {
	ext := EncoderExt("txt")
	tt := map[string]struct {
		in, expected factory
	}{
		"nils": {
			expected: factory{
				encoder:    EncoderCast,
				bktEncoder: EncoderCast,
				decoder:    DecoderCast,
				bktDecoder: DecoderCast,
			},
		},
		"bkt-nils": {
			expected: factory{
				encoder:    ext,
				bktEncoder: ext,
				decoder:    DecoderRemoveExt,
				bktDecoder: DecoderRemoveExt,
			},
			in: factory{
				encoder: ext,
				decoder: DecoderRemoveExt,
			},
		},
	}

	var fnEq = func(fn1, fn2 any) {
		assert.Equal(t, reflect.ValueOf(fn1).Pointer(), reflect.ValueOf(fn2).Pointer())
	}

	for n, tc := range tt {
		t.Run(n, func(t *testing.T) {
			f := newFactory(tc.in.encoder, tc.in.bktEncoder, tc.in.decoder, tc.in.bktDecoder)

			fnEq(tc.expected.encoder, f.encoder)
			fnEq(tc.expected.bktDecoder, f.bktDecoder)
			fnEq(tc.expected.decoder, f.decoder)
			fnEq(tc.expected.bktDecoder, f.bktDecoder)

		})
	}
}
