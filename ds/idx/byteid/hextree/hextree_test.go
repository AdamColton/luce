package hextree

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/testsuite"
)

func TestHexTree(t *testing.T) {
	testsuite.TestBasicInsertGet(t, New)
	testsuite.TestDeleteRecycle(t, New)
	testsuite.TestNext(t, New)
}
