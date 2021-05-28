package main

import (
	"fmt"
	"net/http"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/serial/typestring"
	"github.com/adamcolton/luce/serial/wrap/json"

	"github.com/adamcolton/luce/ds/bus"
	"github.com/adamcolton/luce/ds/bus/serialbus"
	"github.com/adamcolton/luce/util/lusers"
	"github.com/adamcolton/luce/util/lusers/lusess"
)

type wsInstance struct {
	send *serialbus.Sender
	s    *Server
	u    *lusers.User
}

func (s *Server) WebSocket(to chan<- []byte, from <-chan []byte, r *http.Request) {
	tm := typestring.NewTypeMap(
		(*login)(nil),
		(*loginStatus)(nil),
	)

	wsi := &wsInstance{
		send: &serialbus.Sender{
			TypeSerializer: tm.WriterSerializer(json.Serialize),
			Chan:           to,
		},
		s: s,
	}

	l, err := bus.NewListener(&serialbus.Receiver{
		In:               from,
		TypeDeserializer: tm.ReaderDeserializer(json.Deserialize),
		TypeRegistrar:    tm,
	}, nil, nil)
	lerr.Panic(err)
	bus.RegisterHandlerType(l, wsi)

	l.Run()
	close(to)
}

type login struct {
	Username, Password string
}

func (*login) TypeIDString() string {
	return "login"
}

type loginStatus struct {
	Success bool
	Token   string
}

func (*loginStatus) TypeIDString() string {
	return "loginStatus"
}

func (wsi *wsInstance) HandleLogin(l *login) {
	ls := &loginStatus{}

	u, err := wsi.s.Users.UserStore.Login(l.Username, l.Password)
	ls.Success = err == nil
	if !ls.Success {
		wsi.send.Send(ls)
		return
	}

	wsi.u = u

	fn := func(w http.ResponseWriter, r *http.Request, sess *lusess.Session) {
		sess.SetUser(u)
		sess.Save()
	}
	h := wsi.s.Users.HandlerFunc(fn)
	ls.Token, _ = wsi.s.tokens.Register(h)
	fmt.Println("Sending token", ls.Token)

	wsi.send.Send(ls)
}
