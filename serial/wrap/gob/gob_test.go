package gob_test

import (
	"bytes"
	"testing"

	"github.com/adamcolton/luce/serial/wrap/gob"
	"github.com/adamcolton/luce/serial/wrap/testutil"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	gob.Register((*testutil.Person)(nil))
	testutil.SerialFuncsRoundTrip(t, gob.Serialize, gob.Deserialize)
}

func TestErrors(t *testing.T) {
	// A func can't be encoded by gob.
	assert.Error(t, gob.Serialize(func() {}, &bytes.Buffer{}))

	// The data is not gob.
	var i int
	assert.Error(t, gob.Deserialize(&i, bytes.NewBufferString("not gob")))
}
