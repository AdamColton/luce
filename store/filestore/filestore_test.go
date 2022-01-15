package filestore

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/adamcolton/luce/store/testsuite"
	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	name := "test"
	decoder := func(s string) []byte {
		b, _ := base64.RawURLEncoding.DecodeString(s)
		return b
	}
	enc := base64.RawURLEncoding.EncodeToString
	f, err := Factory(name, enc, enc, decoder, decoder)
	assert.NoError(t, err)
	defer func() {
		os.RemoveAll(name)
	}()
	testsuite.TestAll(t, f)
}
