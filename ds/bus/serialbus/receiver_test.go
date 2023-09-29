package serialbus_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/serial/wrap/json"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// note is a message that these tests send.
type note struct {
	Text string
}

func (*note) TypeID32() uint32 {
	return 9001
}

// serialized is the bytes of a message on a serial bus.
func serialized(t *testing.T, tm type32.TypeMap, text string) []byte {
	b, err := tm.WriterSerializer(json.Serialize).SerializeType(&note{Text: text}, nil)
	assert.NoError(t, err)
	return b
}

func TestReceiver(t *testing.T) {
	in := make(chan []byte)
	out := make(chan any, 1)
	errCh := make(chan error, 1)
	tm := type32.NewTypeMap()
	r := &serialbus.Receiver{
		In:               in,
		TypeDeserializer: tm.ReaderDeserializer(json.Deserialize),
		TypeRegistrar:    tm,
	}
	r.SetOut(out)
	assert.NoError(t, r.SetErrorHandler(errCh))
	assert.NoError(t, r.RegisterType(&note{}))
	done := timeout.Run(r.Run)

	assert.NoError(t, timeout.After(1000, func() {
		in <- serialized(t, tm, "hello")
		assert.Equal(t, &note{Text: "hello"}, <-out)

		// A message that can't be deserialized goes to the error handler.
		in <- []byte("not a message")
		assert.Error(t, <-errCh)
	}))

	close(in)
	assert.NoError(t, timeout.After(1000, done))
}

func TestReceiverErrorHandler(t *testing.T) {
	r := &serialbus.Receiver{}

	// A handler that can't be used leaves the current one.
	errCh := make(chan error, 1)
	assert.NoError(t, r.SetErrorHandler(errCh))
	assert.Equal(t, lerr.ErrHandlerFunc, r.SetErrorHandler("not a handler"))

	tm := type32.NewTypeMap()
	r.TypeDeserializer = tm.ReaderDeserializer(json.Deserialize)
	in := make(chan []byte, 1)
	in <- []byte("not a message")
	close(in)
	r.In = in
	r.SetOut(make(chan any))
	r.Run()
	assert.Error(t, <-errCh)

	// nil ignores errors.
	assert.NoError(t, r.SetErrorHandler(nil))
}

func TestReceiverNoRegistrar(t *testing.T) {
	assert.NoError(t, (&serialbus.Receiver{}).RegisterType(&note{}))
}

func TestNewListener(t *testing.T) {
	in := make(chan []byte)
	tm := type32.NewTypeMap()
	notes := make(chan string)
	errCh := make(chan error, 1)

	l, err := serialbus.NewListener(in, tm.ReaderDeserializer(json.Deserialize), tm, errCh, func(n *note) {
		notes <- n.Text
	})
	assert.NoError(t, err)
	done := timeout.Run(l.Run)

	assert.NoError(t, timeout.After(1000, func() {
		in <- serialized(t, tm, "hello")
		assert.Equal(t, "hello", <-notes)

		// An error while receiving goes to the error handler.
		in <- []byte("not a message")
		assert.Error(t, <-errCh)
	}))
	close(in)
	assert.NoError(t, timeout.After(1000, done))

	// A handler that can't be registered.
	l, err = serialbus.NewListener(in, tm.ReaderDeserializer(json.Deserialize), tm, nil, "not a handler")
	assert.Nil(t, l)
	assert.Error(t, err)
}
