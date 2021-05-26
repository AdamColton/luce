package main

import (
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/unixsocket"
)

type UnixsocketReq struct {
}
type UnixsocketResp struct {
	Err error
}

func (m *Modes) UnixsocketHandler(r *UnixsocketReq) *UnixsocketResp {
	return &UnixsocketResp{
		Err: unixsocket.Client(m),
	}
}

func (*Modes) UnixsocketUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Alias: "us",
	}
}
