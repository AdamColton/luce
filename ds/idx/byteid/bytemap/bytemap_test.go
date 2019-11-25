package bytemap_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/idx/byteid/bytemap"
	"github.com/adamcolton/luce/ds/idx/byteid/testsuite"
)

func TestSuite(t *testing.T) {
	testsuite.TestAll(t, bytemap.New)
}
