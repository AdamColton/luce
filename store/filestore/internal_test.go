package filestore

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFactory(t *testing.T) {
	ext, dirExt := EncoderExt("txt"), EncoderExt("dir")
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
		"all-set": {
			expected: factory{
				encoder:    ext,
				bktEncoder: dirExt,
				decoder:    DecoderRemoveExt,
				bktDecoder: DecoderCast,
			},
			in: factory{
				encoder:    ext,
				bktEncoder: dirExt,
				decoder:    DecoderRemoveExt,
				bktDecoder: DecoderCast,
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
			fnEq(tc.expected.bktEncoder, f.bktEncoder)
			fnEq(tc.expected.decoder, f.decoder)
			fnEq(tc.expected.bktDecoder, f.bktDecoder)
		})
	}
}
