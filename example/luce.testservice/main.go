// This demonstrates how to make a service for luce.server.
package main

import (
	"fmt"
	"html/template"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/tools/server/service"
)

func main() {
	conn := lerr.Must(service.NewClient("/tmp/luceserver.service"))

	srv := conn.Service
	srv.Base = "/testsrv"
	srv.Name = "Test Service"

	// Need to update hosts file for this to work
	//conn.Service.Host = "somehost.{domain:.*}"

	// ??? https://somehost.adamcolton.local:6060/testsrv/sayHi/adam
	// /testsrv/sayHi/{name}
	// "somehost.{domain:.*}"

	service.NewServiceRoute("home").
		Handle(conn, Home)
	conn.Service.AddLink("Home", "", "home")

	service.NewServiceRoute("sayHi/{name}").
		WithUser().
		Handle(conn, SayHi)
	conn.Service.AddLink("Hi, Adam", "", "sayHi/Adam")

	service.NewServiceRoute("query").
		WithQuery().
		Handle(conn, Query)
	conn.Service.AddLink("Query", "", "query?foo=bar")

	service.NewServiceRoute("admin").
		WithQuery().
		RequireGroup("admin").
		Handle(conn, Admin)
	conn.Service.AddLink("Admin Only", "", "admin")

	conn.RunService()
}

var home = `<!DOCTYPE html>
<html>
	<head><title>Example Test Service</title></head>
	<body>This is the example test service</body>
</html>
`

func Home(req *service.Request) *service.Response {
	return req.ResponseString(home)
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

func Host(req *service.Request) *service.Response {
	return req.ResponseString("DOMAIN:" + req.PathVars["domain"])
}
