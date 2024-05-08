package gob_test

import (
	"testing"

	"github.com/adamcolton/luce/serial/wrap/gob"
	"github.com/adamcolton/luce/serial/wrap/testutil"
)

func TestAll(t *testing.T) {
	gob.Register((*testutil.Person)(nil))
	testutil.SerialFuncsRoundTrip(t, gob.Serialize, gob.Deserialize)
}

func TestInterfaces(t *testing.T) {
	testutil.SerialInterfacesRoundTrip(t, gob.Serializer{}, gob.Deserializer{})
}
