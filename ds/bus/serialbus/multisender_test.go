package serialbus_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/serial/wrap/json"
	"github.com/stretchr/testify/assert"
)

func TestMultiSenderChannels(t *testing.T) {
	tm := type32.NewTypeMap()
	assert.NoError(t, tm.RegisterType(&note{}))
	var ms bus.MultiSender = serialbus.NewMultiSender(tm.WriterSerializer(json.Serialize))

	a, b, c, d := make(chan []byte, 3), make(chan []byte, 3), make(chan []byte, 3), make(chan []byte, 3)
	assert.NoError(t, ms.Add("a", a))
	assert.NoError(t, ms.Add("b", (chan<- []byte)(b)))
	assert.NoError(t, ms.Add("c", c))
	ms.(*serialbus.MultiSender).AddCh("d", d)

	// Only the keys that are given, and a key that isn't there is skipped.
	assert.NoError(t, ms.Send(&note{Text: "some"}, "b", "d", "missing"))
	assert.Len(t, a, 0)
	assert.Len(t, b, 1)
	assert.Len(t, c, 0)
	assert.Len(t, d, 1)

	// All of them, with the same bytes.
	assert.NoError(t, ms.Send(&note{Text: "all"}))
	assert.Len(t, a, 1)
	assert.Len(t, b, 2)
	assert.Len(t, c, 1)
	assert.Len(t, d, 2)
	assert.Equal(t, <-a, <-c)

	// A deleted channel is not sent to.
	ms.Delete("a")
	ms.Delete("not there")
	assert.NoError(t, ms.Send(&note{Text: "after delete"}))
	assert.Len(t, a, 0)
	assert.Len(t, c, 1)
}

func TestMultiSenderErrors(t *testing.T) {
	ms := serialbus.NewMultiSender(failSerializer{})

	// Only a channel of []byte can be added.
	assert.EqualError(t, ms.Add("a", make(chan string)), "Expected chan<- []byte")
	assert.EqualError(t, ms.Add("a", "not a channel"), "Expected chan<- []byte")
	assert.Empty(t, ms.Chans)

	ch := make(chan []byte, 1)
	assert.NoError(t, ms.Add("a", ch))
	assert.EqualError(t, ms.Send(&note{}), "can't serialize")
	assert.Len(t, ch, 0)
}
