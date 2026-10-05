// This demonstrates how to make a service for luce.server.
package main

import (
	"fmt"
	"html/template"
	"strconv"
	"time"

	"github.com/adamcolton/luce/ds/lmap"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/tools/server/service"
	"github.com/adamcolton/luce/util/lusers"
)

var conn *service.Client

func main() {
	conn = lerr.Must(service.NewClient("/tmp/luceserver.service"))

	srv := conn.Service
	srv.Base = "/testsrv"
	srv.Name = "Test Service"

	// Need to update hosts file for this to work
	//conn.Service.Host = "somehost.{domain:.*}"

	// ??? https://somehost.adamcolton.local:6060/testsrv/sayHi/adam
	// /testsrv/sayHi/{name}
	// "somehost.{domain:.*}"

	service.NewRoute("home").
		Handle(conn, Home)
	conn.Service.AddLink("Home", "", "home")

	service.NewRoute("sayHi/{name}").
		WithUser().
		Handle(conn, SayHi)
	conn.Service.AddLink("Hi, Adam", "", "sayHi/Adam")

	service.NewRoute("query").
		WithQuery().
		Handle(conn, Query)
	conn.Service.AddLink("Query", "", "query?foo=bar")

	service.NewRoute("admin").
		WithQuery().
		RequireGroup("admin").
		Handle(conn, Admin)
	conn.Service.AddLink("Admin Only", "", "admin")

	service.NewRoute("ws").
		WithSocket().
		Handle(conn, OpenWebSocket)

	service.NewRoute("websocket").
		Handle(conn, WebSocket)
	conn.Service.AddLink("WebSocket Demo", "", "websocket")

	lerr.Panic(conn.RegisterHandlers(SocketMessage, SocketClose))

	conn.Run()
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

var sockets = lmap.NewSafe[uint32, *lusers.User](nil)

func OpenWebSocket(req *service.Request) *service.Response {
	id := req.WebSocket
	sockets.Set(id, req.User)
	go func() {
		for i := range 10 {
			if !sockets.Contains(id) {
				return
			}
			msg := strconv.Itoa(i)
			conn.SocketMessage(id, []byte(msg))
			time.Sleep(time.Second * 2)
		}
		conn.SocketClose(id)
	}()
	return req.ResponseString("this does nothing")
}

func SocketMessage(msg service.SocketMessage) {
	fmt.Println("got msg: ", string(msg.Body))
}

func SocketClose(c service.SocketClose) {
	fmt.Println("socket closed")
	sockets.Delete(c.ID)
}

var wsTmpl = template.Must(template.New("websocket").Parse(`<!DOCTYPE html>
<html>
	<head><title>WebSocket Demo</title></head>
	<body>
		<div id="buttons">
			<button id="send">Send</button>
			<button id="close">Close</button>
		</div>
		<script>
			var counter = 0;
			var wsScheme = location.protocol === "https:" ? "wss:" : "ws:";
			var ws = new WebSocket(wsScheme + "//" + location.host + "{{.WSPath}}");
			ws.onclose = function() {
				console.log("the service closed the websocket");
			};
			ws.onmessage = function(e) {
				console.log("message from service:", e.data);
			};
			document.getElementById("send").onclick = function() {
				counter++;
				ws.send(String(counter));
			};
			document.getElementById("close").onclick = function() {
				ws.close();
				document.getElementById("buttons").innerHTML = "Done";
			};
		</script>
	</body>
</html>
`))

func WebSocket(req *service.Request) *service.Response {
	n := struct {
		WSPath string
	}{
		WSPath: conn.Service.Base + "/ws",
	}
	return req.ResponseTemplate("", wsTmpl, n)
}
