package bytemap

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/testsuite"
)

func TestSuiteBasicInsert(t *testing.T) {
	testsuite.TestAll(t, New)
}
