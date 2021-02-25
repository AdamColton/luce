package bstore

import (
	"os"
	"testing"

	"github.com/adamcolton/luce/store/testsuite"
	"github.com/testify/assert"
)

func TestBasic(t *testing.T) {
	name := "test.db"
	f := Factory(name, 0777, nil)
	defer func() {
		os.Remove(name)
	}()
	testsuite.TestAll(t, f)

	s, err := f.Store(nil)
	assert.NoError(t, err)
	s.Put([]byte("foo"), []byte("bar"))
}
