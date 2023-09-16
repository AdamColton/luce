// Command unixsocket is an example of a unix socket that serves a command line.
// Run it and connect with `nc -U /tmp/socket-demo.sock`, then enter "help".
// It writes plain text, so Client, which needs a socket made with CLISocket,
// can not connect to it.
package main

import (
	"fmt"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/unixsocket"
)

// User is the argument of the user command. The command asks for each field.
type User struct {
	Name, Password string
}

// main runs the socket until a client enters "close".
func main() {
	err := unixsocket.
		New("/tmp/socket-demo.sock", nil).
		Runner(setCommands).
		Run()
	lerr.Panic(err)
}

// setCommands gives each connection its prompt, start message and commands.
func setCommands(r *cli.Runner) {
	r.Prompt = "> "
	r.StartMessage = "Welcome to the commands demo\nenter 'help' for more\n"

	r.Commands = lerr.Must(handler.Cmds([]*handler.Command{
		{
			Name: "",
			Action: func() {
				r.WriteString("Unrecognized Command\nEnter 'help' to see commands\n")
			},
		}, {
			Name:  "help",
			Usage: "show this menu",
			Action: func() {
				r.ShowCommands(nil)
			},
		}, {
			Name:  "close",
			Usage: "the server (and client)",
			Action: func() {
				r.WriteString("Closing Server. Goodbye\n")
				r.Close = true
				r.Exit = true
			},
		}, {
			Name:  "exit",
			Usage: "exit the client",
			Action: func() {
				r.WriteString("Closing Client. Goodbye\n")
				r.Exit = true
			},
		}, {
			Name:  "sayHi",
			Usage: "says hi",
			Action: func() {
				r.WriteString("Hi!")
			},
		}, {
			Name:  "user",
			Usage: "create user",
			Action: func(u *User) {
				fmt.Printf("Create User\n\tName: %s\n\tPassword: %s", u.Name, u.Password)
				r.WriteString("  Created User")
			},
		},
	}))
}
