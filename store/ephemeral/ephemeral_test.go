package ephemeral_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/store/ephemeral"
	"github.com/adamcolton/luce/store/testsuite"
)

func TestAll(t *testing.T) {
	f := ephemeral.Factory(bytebtree.New, 1)
	testsuite.TestAll(t, f)
}
