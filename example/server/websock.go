package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adamcolton/luce/util/lusers"
	"github.com/adamcolton/luce/util/lusers/lusess"
)

type wsInstance struct {
	to chan<- []byte
	s  *Server
	u  *lusers.User
}

func (s *Server) WebSocket(to chan<- []byte, from <-chan []byte, r *http.Request) {
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("User closed connection")
				return
			}
		}
	}()

	wsi := wsInstance{
		to: to,
		s:  s,
	}

	for msg := range from {
		split := bytes.Index(msg, []byte{byte(' ')})
		if split == -1 {
			fmt.Println("Error: ", string(msg))
			continue
		}
		kind := string(msg[:split])
		data := msg[split+1:]

		switch kind {
		case "login":
			wsi.login(data)
		default:
			fmt.Println("Unknown: ", kind, string(data))
		}
	}
	stop <- true
}

func (wsi *wsInstance) login(data []byte) {
	var l struct {
		Username, Password string
	}
	json.Unmarshal(data, &l)

	u, err := wsi.s.Users.GetByName(l.Username)
	if err != nil || u == nil || u.CheckPassword(l.Password) != nil {
		wsi.to <- []byte("loginFailed")
		return
	}

	wsi.u = u

	fn := func(w http.ResponseWriter, r *http.Request, sess *lusess.Session) {
		sess.SetUser(u)
		sess.Save()
	}
	h := wsi.s.Users.HandlerFunc(fn)
	token, _ := wsi.s.tokens.Register(h)
	fmt.Println("Sending token", token)

	wsi.to <- []byte(fmt.Sprintf(`loginSuccess %s`, token))
}
