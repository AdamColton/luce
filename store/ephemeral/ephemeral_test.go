package ephemeral

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/store/testsuite"
)

func TestBasic(t *testing.T) {
	f := Factory(bytebtree.New, 1)
	testsuite.TestAll(t, f)
}
