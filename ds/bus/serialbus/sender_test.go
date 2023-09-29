package serialbus_test

import (
	"errors"
	"testing"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/serial/wrap/json"
	"github.com/stretchr/testify/assert"
)

// failSerializer is a TypeSerializer that fails.
type failSerializer struct{}

func (failSerializer) SerializeType(any, []byte) ([]byte, error) {
	return nil, errors.New("can't serialize")
}

func TestSender(t *testing.T) {
	ch := make(chan []byte, 1)
	tm := type32.NewTypeMap()
	assert.NoError(t, tm.RegisterType(&note{}))
	var s bus.Sender = &serialbus.Sender{
		TypeSerializer: tm.WriterSerializer(json.Serialize),
		Chan:           ch,
	}
	assert.NoError(t, s.Send(&note{Text: "hello"}))
	assert.Len(t, ch, 1)

	// What was sent can be received.
	got, err := tm.ReaderDeserializer(json.Deserialize).DeserializeType(<-ch)
	assert.NoError(t, err)
	assert.Equal(t, &note{Text: "hello"}, got)
}

func TestSenderError(t *testing.T) {
	ch := make(chan []byte, 1)
	s := &serialbus.Sender{TypeSerializer: failSerializer{}, Chan: ch}
	assert.EqualError(t, s.Send(&note{}), "can't serialize")
	assert.Len(t, ch, 0)
}
