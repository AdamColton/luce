package gob

import (
	"testing"

	"github.com/adamcolton/luce/serial/wrap/testutil"
)

func TestAll(t *testing.T) {
	testutil.SerialFuncsRoundTrip(t, Serialize, Deserialize)
}

func TestInterfaces(t *testing.T) {
	testutil.SerialInterfacesRoundTrip(t, Serializer{}, Deserializer{})
}

func TestEncDec(t *testing.T) {
	testutil.EncDec(t, Enc, Dec)
}
