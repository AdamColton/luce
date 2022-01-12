package main

import (
	"fmt"
	"html/template"

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

	conn.Add(
		Admin,
		g.GetQuery("admin").RequireGroup("admin"),
	)

	conn.Run()
}

var hiTmpl = template.Must(template.New("sayHi").Parse(`<!DOCTYPE html>
<html>
	<head><title>Say Hi</title></head>
	<body>Hi, {{.Name}}</body>
</html>
`))

func SayHi(req *service.Request) *service.Response {
	n := struct {
		Name string
	}{
		Name: req.PathVars["name"],
	}
	if req.User != nil && req.User.Name != "" {
		n.Name = req.User.Name
	}
	return req.ResponseTemplate("", hiTmpl, n)
}

func Query(req *service.Request) *service.Response {
	r := req.Response(nil)
	fmt.Fprintf(r, "Query: %v Path: %v", req.Query, req.Path)
	return r
}

func Admin(req *service.Request) *service.Response {
	return req.ResponseString("You are admin")
}
