package luce

import (
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/unixsocket"
)

type UnixSocketClient struct {
	Disabled bool
	Ctx      cli.Context
}

type UnixsocketReq struct {
}
type UnixsocketResp struct {
	Err error
}

func (usc *UnixSocketClient) UnixsocketHandler(r *UnixsocketReq) *UnixsocketResp {
	return &UnixsocketResp{
		Err: unixsocket.Client(usc.Ctx),
	}
}

func (usc *UnixSocketClient) UnixsocketUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Alias:    "us",
		Disabled: usc.Disabled,
	}
}
