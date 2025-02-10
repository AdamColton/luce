package service_test

import (
	"net"
	"testing"
	"time"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/type32"
	"github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/reflector"
	"github.com/adamcolton/luce/util/unixsocket"
	"github.com/stretchr/testify/assert"
)

type testMessage struct {
	ID   uint32
	Body string
}

func (testMessage) TypeID32() uint32 {
	return 3141592653
}

func TestRoundTrip(t *testing.T) {
	expected := "this is a test"
	handleServiceSocket := func(netConn net.Conn) {
		conn := lerr.Must(service.NewConn(netConn))
		msg := testMessage{
			ID:   31415,
			Body: expected,
		}

		// if testMessage is not registered, we get an error
		err := conn.Sender.Send(msg)
		expectErr := type32.ErrTypeNotFound{reflector.Type[testMessage]()}
		assert.Equal(t, expectErr, err)

		// after calling service.Register, it works
		service.Register(testMessage{})
		err = conn.Sender.Send(msg)
		assert.NoError(t, err)
	}

	sock := "/tmp/luceserver.sock"
	srvr := unixsocket.New(sock, handleServiceSocket)
	go func() {
		time.Sleep(time.Second)
		srvr.Close()
	}()

	var got string
	go func() {
		time.Sleep(time.Millisecond * 10)
		clnt := lerr.Must(service.NewClient(sock))
		err := clnt.Listener.RegisterInterface(func(msg testMessage) {
			got = msg.Body
			srvr.Close()
		})
		assert.NoError(t, err)
		clnt.Run()
	}()

	err := srvr.Run()
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
