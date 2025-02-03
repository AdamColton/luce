package badgerstore_test

import (
	"os"
	"testing"

	"github.com/adamcolton/luce/store/badgerstore"
	"github.com/adamcolton/luce/store/testsuite"
)

func TestBasic(t *testing.T) {
	name := "testdb"
	f := badgerstore.Factory(name)
	defer func() {
		os.RemoveAll(name)
	}()
	testsuite.TestFlat(t, f)
}
