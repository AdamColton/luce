package server

import (
	"strings"

	"github.com/adamcolton/luce/util/unixsocket"
)

// RunSocket for the admin interface. This is not invoked by ListenAndServe
// and needs to be run seperately.
func (s *Server) RunSocket() {
	sck := &unixsocket.Commands{
		Name: "luce-server",
		Addr: s.Socket,
		Commands: []unixsocket.Command{
			{
				Name:  "help",
				Usage: "show this menu",
				Action: func(ctx *unixsocket.Context) {
					ctx.WriteString(ctx.Socket.Help())
				},
			}, {
				Name:  "close",
				Usage: "the server (and client)",
				Action: func(ctx *unixsocket.Context) {
					ctx.WriteString("Closing Server. Goodbye\n")
					s.Close()
					ctx.Socket.Close()
				},
			}, {
				Name:  "exit",
				Usage: "exit the client",
				Action: func(ctx *unixsocket.Context) {
					ctx.WriteString("Closing Client. Goodbye\n")
					ctx.Close()
				},
			}, {
				Name:  "user",
				Usage: "create user",
				Action: func(ctx *unixsocket.Context) {
					var r struct {
						Name, Password string
					}
					ok := ctx.PopulateStruct("user", &r)
					if !ok {
						ctx.WriteString("  Operation Cancelled")
						return
					}

					_, err := s.Users.Create(r.Name, r.Password)
					if ctx.Error(err) {
						return
					}
					ctx.WriteString("  Created User")
				},
			}, {
				Name:  "users",
				Usage: "list users",
				Action: func(ctx *unixsocket.Context) {
					us := s.Users.List()
					for _, name := range us {
						u, err := s.Users.GetByName(name)
						if err != nil {
							ctx.Printf("  Error: %s", err.Error())
						}
						ctx.Printf("  %s", name)
						if len(u.Groups) > 0 {
							ctx.Printf(" (%s)", strings.Join(u.Groups, ", "))
						}
						ctx.Printf("\n")
					}
				},
			}, {
				Name:  "group",
				Usage: "create group",
				Action: func(ctx *unixsocket.Context) {
					var groupName string
					ok := ctx.Input("(group name) ", &groupName)
					if !ok {
						return
					}
					_, err := s.Users.Group(groupName)
					ctx.Error(err)
				},
			}, {
				Name:  "groups",
				Usage: "list groups",
				Action: func(ctx *unixsocket.Context) {
					ctx.Printf("  %s", strings.Join(s.Users.Groups(), "\n  "))
				},
			}, {
				Name:  "user-group",
				Usage: "add user to group",
				Action: func(ctx *unixsocket.Context) {
					var user, group string
					ok := ctx.Input("(group) ", &group)
					if !ok {
						ctx.WriteString("  Operation Cancelled")
						return
					}
					g := s.Users.HasGroup(group)
					if g == nil {
						ctx.WriteString("  Group not found")
					}

					ok = ctx.Input("(user) ", &user)
					if !ok {
						ctx.WriteString("  Operation Cancelled")
						return
					}
					u, err := s.Users.GetByName(user)
					if ctx.Error(err) {
						return
					}

					err = g.AddUser(u)
					if ctx.Error(err) {
						return
					}

					err = s.Users.Update(u)
					ctx.Error(err)
				},
			}, {
				Name:  "rm-user-group",
				Usage: "remove user from group",
				Action: func(ctx *unixsocket.Context) {
					var user, group string
					ok := ctx.Input("(group) ", &group)
					if !ok {
						ctx.WriteString("  Operation Cancelled")
						return
					}
					g := s.Users.HasGroup(group)
					if g == nil {
						ctx.WriteString("  Group not found")
					}

					ok = ctx.Input("(user) ", &user)
					if !ok {
						ctx.WriteString("  Operation Cancelled")
						return
					}
					u, err := s.Users.GetByName(user)
					if ctx.Error(err) {
						return
					}

					err = g.RemoveUser(u)
					if ctx.Error(err) {
						return
					}

					err = s.Users.Update(u)
					ctx.Error(err)
				},
			}, {
				Name:  "setport",
				Usage: "change port that server is running on",
				Action: func(ctx *unixsocket.Context) {
					var port string
					ok := ctx.Input("(port) ", &port)
					if !ok {
						ctx.WriteString("  Operation Cancelled")
						return
					}
					s.Close()
					s.Addr = port
					go s.ListenAndServe()
				},
			}, {
				Name:  "settings",
				Usage: "display server setttings",
				Action: func(ctx *unixsocket.Context) {
					ctx.Printf("  AdminLockUserCreation %t\n", s.Settings.AdminLockUserCreation)
				},
			}, {
				Name:  "adminLockUserCreation",
				Usage: "Sets if user creation requires a user to already be logged in",
				Action: func(ctx *unixsocket.Context) {
					var lock string
					ok := ctx.Input("(lock create user Y/N) ", &lock)
					if !ok {
						ctx.WriteString("  Operation Cancelled")
						return
					}
					s.Settings.AdminLockUserCreation = lock == "Y" || lock == "y"
				},
			}, {
				Name:  "routes",
				Usage: "View service routes",
				Action: func(ctx *unixsocket.Context) {
					for id, r := range s.serviceRoutes {
						if r.active {
							ctx.WriteString(id)
							ctx.WriteString("\n")
						}
					}
				},
			},
		},
	}
	sck.Run()
}
