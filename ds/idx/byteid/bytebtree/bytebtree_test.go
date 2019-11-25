package bytebtree_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytebtree"
	"github.com/adamcolton/luce/ds/idx/byteid/testsuite"
)

func TestSuite(t *testing.T) {
	testsuite.TestAll(t, bytebtree.New)
}
