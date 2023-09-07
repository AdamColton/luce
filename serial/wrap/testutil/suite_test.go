package testutil_test

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/adamcolton/luce/serial"
	"github.com/adamcolton/luce/serial/wrap/testutil"
)

// The suites are run against encoding/json, which does a correct round trip.

func encode(v any, w io.Writer) error { return json.NewEncoder(w).Encode(v) }
func decode(v any, r io.Reader) error { return json.NewDecoder(r).Decode(v) }

func TestSerialFuncsRoundTrip(t *testing.T) {
	testutil.SerialFuncsRoundTrip(t, encode, decode)
}

func TestSerialInterfacesRoundTrip(t *testing.T) {
	testutil.SerialInterfacesRoundTrip(t, serial.WriterSerializer(encode), serial.ReaderDeserializer(decode))
}

func TestEncDec(t *testing.T) {
	testutil.EncDec(t,
		func(v any) []byte {
			b, _ := json.Marshal(v)
			return b
		},
		func(b []byte, v any) {
			json.Unmarshal(b, v)
		},
	)
}
