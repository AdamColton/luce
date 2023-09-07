package gob

import (
	"testing"

	"github.com/adamcolton/luce/serial/wrap/testutil"
)

func TestAll(t *testing.T) {
	testutil.Suite(t, Serialize, Deserialize)
}
