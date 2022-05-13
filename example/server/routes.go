package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/adamcolton/luce/lhttp/formdecoder"
	"github.com/adamcolton/luce/lhttp/midware"
	"github.com/adamcolton/luce/util/lusers/lusess"
)

func (s *Server) routes() {
	m := midware.New(
		s.Users,
		midware.NewDecoder(formdecoder.New(), "Form"),
	)
	r := s.Router

	r.HandleFunc("/", m.Handle(s.Home))

	r.HandleFunc("/login", m.Handle(s.GetLogin)).Methods("GET")
	r.HandleFunc("/login", m.Handle(s.PostLogin)).Methods("POST")

	r.HandleFunc("/user/create", m.Handle(s.GetCreateUser)).Methods("GET")
	r.HandleFunc("/user/create", m.Handle(s.PostCreateUser)).Methods("POST")

	r.HandleFunc("/user/grid", s.GetUserGrid).Methods("GET")

	ws := midware.NewWebSocket()
	r.HandleFunc("/websocket", s.GetWebsocket).Methods("GET")
	r.HandleFunc("/websocket/connect", ws.HandleSocketChans(s.WebSocket))

	r.HandleFunc("/token", s.tokens.Post).Methods("POST")

	r.HandleFunc("/logout", m.Handle(s.Logout)).Methods("GET")
}

func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusFound)
}

func (s *Server) Template(w io.Writer, name string, data interface{}) {
	tmp := bytes.NewBuffer(nil)
	err := s.Tmpls.ExecuteTemplate(tmp, name, data)
	if err == nil {
		w.Write(tmp.Bytes())
	} else {
		w.Write([]byte(err.Error()))
	}

}

func (s *Server) Home(w http.ResponseWriter, r *http.Request, data struct {
	Session *lusess.Session
}) {
	s.Template(w, "home.html", data.Session.User())
}

func (s *Server) GetWebsocket(w http.ResponseWriter, r *http.Request) {
	s.Template(w, "websocket.html", nil)
}

func (s *Server) Logout(w http.ResponseWriter, r *http.Request, data struct {
	Session *lusess.Session
}) {
	data.Session.Logout()
	data.Session.Save()
	redirect(w, r)
}

func (s *Server) GetLogin(w http.ResponseWriter, r *http.Request, data struct {
	Session *lusess.Session
}) {
	if data.Session.User() != nil {
		redirect(w, r)
	} else {
		s.login(w, r, "")
	}
}

func (s *Server) login(w http.ResponseWriter, r *http.Request, message string) {
	s.Template(w, "login.html", struct {
		Message string
	}{
		Message: message,
	})
}

func (s *Server) PostLogin(w http.ResponseWriter, r *http.Request, data struct {
	Session *lusess.Session
	Form    *lusess.Login
}) {
	_, err := data.Session.Login(data.Form)
	if err != nil {
		log.Print(err)
		return
	}
	data.Session.Save()

	redirect(w, r)
}

func (s *Server) GetCreateUser(w http.ResponseWriter, r *http.Request, data struct {
	Session *lusess.Session
}) {
	if s.Settings.AdminLockUserCreation {
		u := data.Session.User()
		if u == nil || !u.In("admin") {
			fmt.Fprint(w, "Must be logged in as admin to create user")
			return
		}
	}
	s.Template(w, "createUser.html", nil)
}

func (s *Server) PostCreateUser(w http.ResponseWriter, r *http.Request, data struct {
	Session *lusess.Session
	Form    *struct {
		Username, Password, Again string
	}
}) {
	currentUser := data.Session.User()
	if s.Settings.AdminLockUserCreation {
		if currentUser == nil || !currentUser.In("admin") {
			fmt.Fprint(w, "Must be logged in as admin to create user")
			return
		}
	}

	data.Form.Username = strings.TrimSpace(data.Form.Username)
	data.Form.Password = strings.TrimSpace(data.Form.Password)
	data.Form.Again = strings.TrimSpace(data.Form.Again)
	if len(data.Form.Username) < 4 {
		s.Template(w, "createUser.html", "Username must be at least 4 characters")
		return
	}
	if len(data.Form.Password) < 4 {
		s.Template(w, "createUser.html", "Password must be at least 4 characters")
		return
	}
	if data.Form.Password != data.Form.Again {
		s.Template(w, "createUser.html", "Passwords must match")
		return
	}

	newUser, err := s.Users.Create(data.Form.Username, data.Form.Password)
	if err != nil {
		log.Print(err)
		return
	}

	if currentUser == nil {
		data.Session.SetUser(newUser)
		err = data.Session.Save()
		if err != nil {
			log.Print(err)
			return
		}
	}

	redirect(w, r)
}

func (s *Server) GetUserGrid(w http.ResponseWriter, r *http.Request) {
	s.Users.List()
	type User struct {
		Name   string
		Groups []bool
	}
	var out struct {
		Groups []string
		Users  []User
	}
	out.Groups = s.Users.Groups()

	users := s.Users.List()
	out.Users = make([]User, len(users))
	for i, name := range users {
		out.Users[i].Name = name
		user, _ := s.Users.GetByName(name)
		out.Users[i].Groups = make([]bool, len(out.Groups))
		for j, g := range out.Groups {
			out.Users[i].Groups[j] = user.In(g)
		}
	}

	s.Template(w, "userGrid.html", out)
}
