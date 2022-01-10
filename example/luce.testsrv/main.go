package main

import (
	"fmt"

	"github.com/adamcolton/luce/tools/server/service"
)

func main() {
	conn := service.MustClient("/tmp/luceserver.service")

	g := service.RouteConfigGen{
		Base: "/testsrv",
	}

	conn.Add(
		SayHi,
		g.Get("sayHi/{name}").WithUser(),
	)

	conn.Add(
		Query,
		g.GetQuery("query").WithPrefix(),
	)

	conn.Run()
}

func SayHi(req *service.Request) *service.Response {
	name := req.PathVars["name"]
	if req.User != nil && req.User.Name != "" {
		name = req.User.Name
	}
	return req.ResponseString(
		fmt.Sprintf("Hi, %s!", name),
	)
}

func Query(req *service.Request) *service.Response {
	r := req.Response(nil)
	fmt.Fprintf(r, "Query: %v Path: %v", req.Query, req.Path)
	return r
}
