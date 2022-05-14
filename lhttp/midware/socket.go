package midware

import (
	"net/http"
	"reflect"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/lhttp"
	"github.com/adamcolton/luce/util/filter"
	"github.com/gorilla/websocket"
)

// WebSocket can perform two different midware operations. It can attach a
// *websocket.Conn directly or it can attach a set of []byte channels connected
// to a websocket.
type WebSocket struct {
	Upgrader       websocket.Upgrader
	ToBuf, FromBuf int
	lhttp.ErrHandler
}

// NewWebSocket with default buffer sizes.
func NewWebSocket() WebSocket {
	return WebSocket{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (ws WebSocket) Handler(handler lhttp.SocketHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		socket, err := ws.Upgrader.Upgrade(w, r, nil)
		if !ws.Check(w, r, lerr.Wrap(err, "while_upgrading_socket")) {
			handler(socket, r)
		}
	}
}

// HandleSocketChans abstracts the websocket as a pair of channels. The handler
// must close the to channel when it is done.
func (ws WebSocket) HandleSocketChans(handler lhttp.ChanHandler) http.HandlerFunc {
	return ws.Handler(func(socket *websocket.Conn, r *http.Request) {
		to := make(chan []byte, ws.ToBuf)
		from := make(chan []byte, ws.FromBuf)

		socket.SetCloseHandler(func(code int, text string) error {
			close(to)
			return nil
		})

		sw := lhttp.NewSocket(socket)

		go sw.RunReader(from)
		go handler(to, from, r)
		sw.RunSender(to)
		socket.Close()
	})
}

// ChannelInitilizer creates an Initilizer that will abstract the websocket to
// a pair of []byte channels.
func (ws WebSocket) Initilizer(to, from, conn string) Initilizer {
	return webSocketInitilizer{
		WebSocket: ws,
		to:        to,
		from:      from,
		conn:      conn,
	}
}

type webSocketInitilizer struct {
	WebSocket
	to, from, conn string
}

var (
	checkSend = filter.
			IsNilRef((*chan<- []byte)(nil)).
			Check(typeErr("Invalid Websocket 'To' type: "))

	checkRecv = filter.
			IsNilRef((*<-chan []byte)(nil)).
			Check(typeErr("Invalid Websocket 'From' type: "))

	checkConn = filter.TypeCheck(
		filter.IsType((*websocket.Conn)(nil)),
		typeErr("Invalid Websocket 'Conn' type: "),
	)
)

// Initilize checks the "to" and "from" fields to validate the names and types.
// It will panic if they are not.
func (i webSocketInitilizer) Initilize(fieldType reflect.Type) DataInserter {
	var hasToField, hasFromField, hasConnField bool
	var toField, fromField, connField reflect.StructField

	if fieldName(i.to) {
		toField, hasToField = fieldType.FieldByName(i.to)
	}
	if fieldName(i.from) {
		fromField, hasFromField = fieldType.FieldByName(i.from)
	}
	if fieldName(i.conn) {
		connField, hasConnField = fieldType.FieldByName(i.conn)
	}

	if !hasToField && !hasFromField && !hasConnField {
		return nil
	}
	di := webSocketDataInserter{
		WebSocket: i.WebSocket,
	}
	if hasToField {
		checkSend.Panic(toField.Type)
		di.to = toField.Index
	}
	if hasFromField {
		checkRecv.Panic(fromField.Type)
		di.from = fromField.Index
	}
	if hasConnField {
		checkConn.Panic(connField.Type)
		di.conn = connField.Index
	}
	return di
}

type webSocketDataInserter struct {
	WebSocket
	to, from, conn []int
}

// Insert sets the "to" and "from" channels on the data field.
func (di webSocketDataInserter) Insert(w http.ResponseWriter, r *http.Request, data reflect.Value) (func(), error) {
	fbi := data.Elem().FieldByIndex

	conn, err := di.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	if di.conn != nil {
		fbi(di.conn).Set(reflect.ValueOf(conn))
	}

	var to chan []byte
	if di.to != nil || di.from != nil {
		sw := lhttp.NewSocket(conn)
		if di.to != nil {
			to = make(chan []byte, di.ToBuf)
			fbi(di.to).Set(reflect.ValueOf(to))
			go sw.RunSender(to)
		}
		if di.from != nil {
			from := make(chan []byte, di.FromBuf)
			fbi(di.from).Set(reflect.ValueOf(from))
			go sw.RunReader(from)
		}
	}

	callbackClose := func() {
		conn.Close()
		defer func() {
			// ignore double close error
			recover()
		}()
		if to != nil {
			close(to)
		}
	}
	return callbackClose, nil
}
