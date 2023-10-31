package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/lhttp/jsondecoder"
	"github.com/adamcolton/luce/lhttp/midware"
	"github.com/adamcolton/luce/lhttp/valuedecoder"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	s := mux.NewRouter()
	s.HandleFunc("/", home)

	m := midware.New(
		midware.NewDecoder(valuedecoder.Form(), "Form"),
		midware.NewDecoder(jsondecoder.New(), "JSON"),
		midware.URL("id", "ID"),
		midware.NewWebSocket().Initilizer("To", "From", "Socket"),
	)
	s.HandleFunc("/decode", getPerson).Methods("GET")
	s.HandleFunc("/decode/{id:[0-9]+}", m.Handle(postPerson)).Methods("POST")

	s.HandleFunc("/decode/json", getJsonPerson).Methods("GET")
	s.HandleFunc("/decode/json", m.Handle(postJsonPerson)).Methods("POST")

	s.HandleFunc("/socket", socketDemo).Methods("GET")
	s.HandleFunc("/socket/handler", m.Handle(socketHandler))
	s.HandleFunc("/socket/chan", m.Handle(chanHandler))

	lerr.Panic(http.ListenAndServe(":8081", s))
}

const (
	header = `<!DOCTYPE html>
<html>
	<head>
		<title>%s</title>
	</head>
	<body>`

	footer = `</body></html>`

	homepage = `<div>Luce Midware Demo</div>
	<div><a href="/decode">Decode Form Demo</a></div>
	<div><a href="/decode/json">Decode Json Demo</a></div>
	<div><a href="/socket?/socket/handler">Socket Handler Demo</a></div>
	<div><a href="/socket?/socket/chan">Socket Channel Demo</a></div>`

	personForm = `<form method="POST" action="/decode/31415">
	<div>
		First <input type="text" name="First" autofocus />
	</div>
	<div>
		Last <input type="text" name="Last" />
	</div>
	<div>
		<button>Go</button>
	</div>
</form>
<div><a href="/">Home</a></div>`

	decodedPerson = `<div>First: %s Last: %s ID:%s</div>
<div><a href="/">Home</a></div>`

	jsonPerson = `<div><a href="/">Home</a></div>
	<div>Look at network logs and server output</div>
	<script>
	var xhr = new XMLHttpRequest();
	xhr.open('POST', "/decode/json");
	xhr.send('{"First":"Brian", "Last":"kernighan"}');
</script>`

	socketBody = `<div><a href="/">Home</a></div>
<div id="log"></div>
<script>
	var conn = new WebSocket("ws:"+window.location.host+window.location.search.replace("?",""));
	conn.onmessage = function(msg){
		document.getElementById("log").innerHTML += msg.data+"<br>";
		console.log(msg);
	}
	var count = 0
	conn.onopen = function(){
		setInterval(function(){
			count++
			conn.send('client to server: '+count);
		},1000)
	}
	conn.onclose=function(){
		console.log("closed")
	}
</script>`
)

func render(w http.ResponseWriter, title, content string, args ...interface{}) {
	fmt.Fprintf(w, header, title)
	fmt.Fprintf(w, content, args...)
	fmt.Fprint(w, footer)
}

func home(w http.ResponseWriter, r *http.Request) {
	render(w, "Midware Demo", homepage)
}

type Person struct {
	First, Last string
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	render(w, "Decode Person", personForm)
}

func postPerson(w http.ResponseWriter, r *http.Request, data *struct {
	Form *Person
	ID   string
}) {
	render(w, "Decode Person", decodedPerson, data.Form.First, data.Form.Last, data.ID)
}

func getJsonPerson(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get json person")
	render(w, "Decode Json Person", jsonPerson)
}

func postJsonPerson(w http.ResponseWriter, r *http.Request, data *struct {
	JSON *Person
}) {
	fmt.Println("post json person")
	fmt.Println(data.JSON)
	fmt.Fprint(w, data.JSON)
}

func socketDemo(w http.ResponseWriter, r *http.Request) {
	render(w, "Decode Person", socketBody)
}

func socketHandler(w http.ResponseWriter, r *http.Request, data *struct {
	Socket *websocket.Conn
}) {
	go func() {
		var err error
		count := 0
		for err == nil {
			count++
			msg := fmt.Sprintf("Server to client: %d", count)
			err = data.Socket.WriteMessage(1, []byte(msg))
			time.Sleep(time.Second)
		}
	}()

	for {
		_, msg, err := data.Socket.ReadMessage()
		if err != nil {
			//lost connection
			break
		}
		fmt.Println(string(msg))
	}
	fmt.Println("Socket Closed")
}

func chanHandler(w http.ResponseWriter, r *http.Request, data *struct {
	To   chan<- []byte
	From <-chan []byte
}) {
	done := false
	go func() {
		for count := 1; !done; count++ {
			msg := fmt.Sprintf("server to client: %d", count)
			data.To <- []byte(msg)
			time.Sleep(time.Second)
		}
	}()

	for msg := range data.From {
		fmt.Println(string(msg))
	}
	fmt.Println("Socket Closed")
	done = true
}
