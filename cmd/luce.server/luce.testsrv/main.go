package main

import (
	"fmt"

	"github.com/adamcolton/luce/tools/server/service"
)

func main() {
	conn := service.MustClient("/tmp/luceserver.service")

	g := service.RouteConfigGen{}

	conn.Add(
		SayHi,
		g.Get("/sayHi/{name}"),
	)

	conn.Add(
		SayHiUser,
		g.Get("/sayHi/{name}").WithUser(),
	)

	conn.Run()
}

func SayHi(req service.Request) service.Response {
	return req.ResponseString(
		fmt.Sprintf("Hi, %s!", req.PathVars["name"]),
	)
}

func SayHiUser(req service.Request) service.Response {
	name := "Anon"
	if req.User != nil && req.User.Name != "" {
		name = req.User.Name
	}
	return req.ResponseString(
		fmt.Sprintf("Hi, %s!", name),
	)
}
