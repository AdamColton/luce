package bytetree

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/testsuite"
)

func TestInsert(t *testing.T) {
	testsuite.TestBasicInsertGet(t, New)
}
