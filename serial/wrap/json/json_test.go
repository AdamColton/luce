package json

import (
	"testing"

	"github.com/adamcolton/luce/serial/wrap/testutil"
)

func TestAll(t *testing.T) {
	testutil.SerialFuncsRoundTrip(t, Serialize, Deserialize)
}
