package service_test

import (
	"net"
	"testing"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/unixsocket"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	sock := t.TempDir() + "/newservice.sock"
	srvr := unixsocket.New(sock, func(conn net.Conn) {
		lerr.Must(service.NewConn(conn))
	})
	go srvr.Run()
	assert.NoError(t, srvr.AwaitRunning())
	defer srvr.Close()

	c, err := service.NewService("widget", "widget-host", "/widget", sock)
	assert.NoError(t, err)
	assert.Equal(t, "widget", c.Service.Name)
	assert.Equal(t, "/widget", c.Service.Base)
	assert.Equal(t, "widget-host.{domain:.*}", c.Service.Host)

	c, err = service.NewService("widget", "", "/widget", sock)
	assert.NoError(t, err)
	assert.Equal(t, "", c.Service.Host)

	_, err = service.NewService("widget", "", "/widget", t.TempDir()+"/no-such-socket")
	assert.Error(t, err)
}

func TestNewClientDialError(t *testing.T) {
	_, err := service.NewClient(t.TempDir() + "/no-such-socket")
	assert.Error(t, err)
}
