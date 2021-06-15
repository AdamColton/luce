package main

import (
	"fmt"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/unixsocket"
)

func main() {
	sck := &unixsocket.Socket{
		Name:         "socket-demo",
		Addr:         "/tmp/socket-demo.sock",
		StartMessage: "Welcome to the socket demo\nenter 'help' for more\n",
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
				Name:  "sayHi",
				Usage: "says hi",
				Action: func(ctx *unixsocket.Context) {
					ctx.WriteString("Hi!")
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

					fmt.Printf("Create User\n\tName: %s\n\tPassword: %s", r.Name, r.Password)
					ctx.WriteString("  Created User")
				},
			},
		},
	}
	lerr.Panic(sck.Run())
}
